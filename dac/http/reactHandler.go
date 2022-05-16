package http

import (
	. "dac"
	"log"
	"net/http"
	"time"
)

//Esta funcion crea una ruta para manejar una aplicación con React
func (HSO *httpSpeakerOption) NewReactApp(url string) {

	//Validar extensiones compatibles

	testWrite := "testing3141592"
	texto := "creating file system"
	extName := "txt"

	if content := NewContentWrite(extName, int64(len(texto)), HSO.dirName, testWrite); content != nil {

		content.DeleteFile()

		http.HandleFunc(url,
			func(response http.ResponseWriter, request *http.Request) {

				path, extName := SanitizeUrlRGP( 0, 5, HSO.dirName + request.URL.Path)
				if len(path) == 0 || len(extName) == 0 {
					path = []string{HSO.dirName , "index"}
					extName = "html"
				}

				if file := NewContentRead(extName, path...); file != nil {

					//Abrebiamos una estructura para el response y el request
					HSO.NewHttpSpeaker(response, request, file)

				}
			})
	}
}


//Esta funcion crea muchas rutas de usuarios para que cada usuario pueda tener
//su propia pwa basada en una sola aplicación react.
func (HSO *httpSpeakerOption) NewReactAppPWA(url string ) {

	//Validar extensiones compatibles
	testWrite := "testing3141592"
	texto := "creating file system"
	extName := "txt"

	if content := NewContentWrite(extName, int64(len(texto)), HSO.dirName , testWrite); content != nil {

		content.DeleteFile()

		http.HandleFunc(url, HSO.NewAppRoute)

		if url[len(url)-1:] == "/" {
			//Ranguear usuarios
			urlUsers := url[:len(url)-1]
			http.HandleFunc(urlUsers, HSO.NewAppRoute)
			http.HandleFunc(urlUsers+"franky", HSO.NewAppRoute)
			http.HandleFunc(urlUsers+"maria", HSO.NewAppRoute)
			http.HandleFunc(urlUsers+"estefania", HSO.NewAppRoute)
			http.HandleFunc(urlUsers+"happycode", HSO.NewAppRoute)
		}

	}

}

func (HSO *httpSpeakerOption) NewAppRoute(response http.ResponseWriter, request *http.Request) {

	tiempoPagina := time.Now()
//	log.Println("Url: ", request.URL.Path)
	path, extName := SanitizeUrlRGP( 0, 5, HSO.dirName +request.URL.Path)
	
//	log.Println("Path: ", path, extName)

	if len(path) == 0 || len(extName) == 0 {
		path = []string{HSO.dirName , "index"}
		extName = "html"
	}

//	log.Println("Path After If: ", path, extName)

	if file := NewContentRead(extName, path...); file != nil {

		//Abrebiamos una estructura para el response y el request
		HSO.NewHttpSpeaker(response, request, file)

	}
	go log.Println("Tiempo en pagina",time.Since(tiempoPagina).Nanoseconds())
}
