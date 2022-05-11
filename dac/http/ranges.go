package http

import (
	"net/http"
	"strconv"
	"strings"
	."dac"
	
)

type httpSpeaker struct {
	response   http.ResponseWriter
	request    *http.Request
	file       *PublicSpaceFile
	fileSize   int64
	contenType string
	bandwidth  int64
	isRange    bool
}

func NewHttpSpeaker(response http.ResponseWriter, request *http.Request, file *PublicSpaceFile, bandwidth uint32, protect bool, cache bool) {

	H := &httpSpeaker{
		response:   response,
		request:    request,
		file:       file,
		contenType: file.GetExtension(),
		fileSize:   file.CalcSizeField(file.GetExtension()),
		bandwidth:  int64(bandwidth),
	}

	if !cache {
		H.handlerCache()
	}

	if protect {
		H.handlerSecurity()
	}
	H.handlerContentType()

	H.writeRanges()

}

func (H *httpSpeaker) handlerContentType() {

	if IsExtension(H.contenType) {

		switch H.contenType {

		//Extensiones de contenido web
		case "html":
			H.response.Header().Set("Content-Type", "text/html; charset=UTF-8")
		case "json":
			H.response.Header().Set("Content-Type", "application/json; charset=utf-8")
		case "js":
			H.response.Header().Set("Content-Type", "text/javascript; charset=UTF-8")
		case "css":
			H.response.Header().Set("Content-Type", "text/css; charset=UTF-8")

		//Extensiones de contenido imagen
		case "jpg":
			H.response.Header().Set("Content-Type", "image/jpeg")
		case "png":
			H.response.Header().Set("Content-Type", "image/png")
		case "webp":
			H.response.Header().Set("Content-Type", "image/webp")
		case "bmp":
			H.response.Header().Set("Content-Type", "image/bmp")
		case "svg":
			H.response.Header().Set("Content-Type", "image/svg+xml")
		case "gif":
			H.response.Header().Set("Content-Type", "image/gif")
		case "glb":
			H.response.Header().Set("Content-Type", "model/gltf-binary")
			H.response.Header().Set("Access-Control-Allow-Origin", "*")

		//Extensiones de audio
		case "mp3":
			H.response.Header().Set("Content-Type", "audio/mpeg")
			H.onRange()
		//Extensiones de contenido video
		case "mp4":
			H.response.Header().Set("Content-Type", "video/mp4")
			H.onRange()

		//Extensiones de contenido de documentos
		case "pdf":
			H.response.Header().Set("Content-Type", "application/pdf")
		case "txt":
			H.response.Header().Set("Content-Type", "text/plain; charset=UTF-8")

		}
	}
}

func (H *httpSpeaker) onRange() {

	H.isRange = true
}

func (H *httpSpeaker) writeRanges() {

	//Si el tamño no supera al ancho de banda permitido entonces sirvelo completo
	if H.fileSize <= H.bandwidth && !H.isRange {

		H.HandlerContentLength(H.fileSize)
		RBuffer := H.file.GetOneFieldChan(H.contenType, H.bandwidth)
		for chanBufer := range RBuffer {

			H.response.Write(chanBufer.Buffer)

		}

	}

	//Si el tamaño supera al ancho de banda permitido entonces sirvelo por rangos
	if H.fileSize > H.bandwidth && H.isRange {

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
