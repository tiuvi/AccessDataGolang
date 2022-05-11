package http

import (
	"net/http"
	"strconv"
	"strings"
	."dac"
	"log"
)

const (
	_           = iota // ignore first value by assigning to blank identifier
	KiloByte  = 1 << (10 * iota)
	MegaByte
	GigaByte
	TeraByte
	PetaByte
	ExaByte
	ZettaByte
	YottaByte
)

type httpSpeakerOption struct {
	bandwidth  int64
	isCache    bool
	isProtect  bool
}


type httpSpeaker struct {
	*httpSpeakerOption
	response   http.ResponseWriter
	request    *http.Request
	file       *PublicSpaceFile
	fileSize   int64
	contenType string
	isRange    bool

}

func SetHttpSpeakerOptions(bandwidth uint32, isProtect bool, isCache bool) *httpSpeakerOption{

	return &httpSpeakerOption{
		bandwidth:  int64(bandwidth),
		isCache: isCache,
		isProtect:isProtect,
	}
}


func  (HSO *httpSpeakerOption) NewHttpSpeaker(response http.ResponseWriter, request *http.Request, file *PublicSpaceFile) {


	H := new(httpSpeaker)
	H.httpSpeakerOption = HSO
	H.response =  response
	H.request  =  request
	H.file     =  file
	H.contenType =  file.GetExtension()
	H.fileSize   =  file.CalcSizeField(H.contenType)
	H.isRange = 	 IsRagesExtension(H.contenType)

	if !H.isCache {
		H.handlerCache()
	}

	if H.isProtect {
		H.handlerSecurity()
	}

	H.handlerContentType()

	H.writeRanges()

}




func (H *httpSpeaker) handlerContentType() {

	if value, found := IsExtensionContent(H.contenType); found {

		H.response.Header().Set("Content-Type", value)
	

		if H.contenType == "glb" {
			H.response.Header().Set("Access-Control-Allow-Origin", "*")
		}

	}
}

func (H *httpSpeaker) onRange() {

	H.isRange = true
}

func (H *httpSpeaker) writeRanges() {

	//Si el tamño no supera al ancho de banda permitido entonces sirvelo completo
	if H.fileSize <= H.bandwidth && !H.isRange{

		H.HandlerContentLength(H.fileSize)
		RBuffer := H.file.GetOneFieldChan(H.contenType, H.bandwidth)
		for chanBufer := range RBuffer {

			H.response.Write(chanBufer.Buffer)

		}

	}

	//Si el tamaño supera al ancho de banda permitido entonces sirvelo por rangos
	if H.fileSize > H.bandwidth  && H.isRange{

		//Enviamos los encabezados de rangos y recibimos el encabezado Range
		if startRange, rango := H.handlerHeaderRanges(); startRange != nil {

			startRange := *startRange
			rango      := *rango

			//Expandir el rango de metadatos.
			buffer := H.file.GetOneFieldRanges(H.contenType, H.bandwidth, rango)

			if startRange > rango*H.bandwidth {

				*buffer.FieldBuffer = (*buffer.FieldBuffer)[startRange-rango*H.bandwidth:]
				H.response.Write(*buffer.FieldBuffer)

			} else {

				H.response.Write(*buffer.FieldBuffer)
			}

		}

	}
}

func (H *httpSpeaker) handlerHeaderRanges() (*int64, *int64) {

	var startRange  int64
	var rango int64
	var err error

	//Obteniendo el rango y formateandolo
	rangeHeader := H.request.Header.Get("Range")

	// StatusPartialContent 206
	H.response.WriteHeader(http.StatusPartialContent)

	if rangeHeader == "" {

		H.HandlerContentLength(H.fileSize)

		return nil, nil
	}

	rangeHeader = strings.Replace(rangeHeader, "bytes=", "", -1)
	ranges     := strings.SplitN(rangeHeader, "-", 2)

	if len(ranges) < 1 {

		return nil, nil
	}

	startRange, err = strconv.ParseInt(ranges[0], 10, 64)
	if err != nil {

		return nil, nil
	}

	if startRange > H.fileSize || startRange < 0 {

		return nil, nil
	}

	//Calculo del rango
	rango = startRange / H.bandwidth

	//Calculando el rango final y el tamaño del rango
	var calcFinalRango int64 = (rango + 1) * H.bandwidth
	if calcFinalRango == 0 {

		calcFinalRango = H.bandwidth
	} else if calcFinalRango > H.fileSize {

		calcFinalRango = H.fileSize
	}

	H.HandlerContentLength(calcFinalRango - startRange)

	//Emitiendo el header de rangos
	H.response.Header().Set("Content-Range",
		strings.Join([]string{"bytes ",
			ranges[0],
			"-",
			strconv.FormatInt(calcFinalRango-1, 10),
			"/",
			strconv.FormatInt(H.fileSize, 10)}, ""))

	return &startRange, &rango
}

func (H *httpSpeaker) handlerCache() {

	H.response.Header().Set("Cache-Control", "no-store, max-age=0")

}

func (H *httpSpeaker) handlerSecurity() {

	H.response.Header().Set("strict-transport-security", "max-age=31536000; includeSubDomains; preload")
	H.response.Header().Set("x-frame-options", "SAMEORIGIN")
	H.response.Header().Set("x-xss-protection", "1; mode=block")

}

func (H *httpSpeaker) HandlerContentLength(length int64) {

	H.response.Header().Set("Content-Length", strconv.FormatInt(length, 10))

}



func (HSO *httpSpeakerOption) NewContentRoute (url string,extension []string , dirName string ){

	//Validar extensiones compatibles

	testWrite := "testing"
	texto := "creating file system"

	for _, extName := range extension {
		
		if content := NewContentWrite(extName, int64(len(texto)), dirName , extName ,testWrite ); content != nil {
			
			if !content.CheckDirSF() { 
				log.Println("CheckDirSF")
				content.SetOneFieldString(extName, texto)
				content.DeleteFile()
			}
		

			http.HandleFunc( url + "/" + extName + "/",
			func(response http.ResponseWriter, request *http.Request) {

				path, extName := SanitizeUrl(0, 4, dirName + request.URL.Path)


				if file := NewContentRead(extName, path...); file != nil {

					//Abrebiamos una estructura para el response y el request
					HSO.NewHttpSpeaker(response, request , file )

				}
			})
		}
	}

}

