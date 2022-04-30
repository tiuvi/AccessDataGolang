package bd

import(
	"log"
	"os"
	"strings"
	"sync/atomic"
)


func (obj *space )newSpaceFile(folderString string, name string)*spaceFile{

	url := strings.Join([]string{ obj.dir , folderString , name , "." , obj.extension }, "") 

	var err error
	//Creamos una nueva referencia a spaceFile
	spacef := new(spaceFile)

	//Abrimos el archivo
	spacef.file, err = os.OpenFile(url , os.O_RDWR | os.O_CREATE, 0666)
	if err != nil {
		//Migrar los errores de archivo a un log de archivo
		log.Println("Error al abrir o crear el archivo.", err)

		if  os.IsNotExist(err) {
			
			err = os.MkdirAll(obj.dir + folderString , 0666)
			if err != nil {

				log.Println("Error al crear la carpeta para ese archivo. ", err)
			} else {

				// obj.LogNewError(Message , "Nueva ruta de archivo en: ", obj.dir + folderString )
				
			}
			spacef.file, err = os.OpenFile(url , os.O_RDWR | os.O_CREATE, 0666)
			if err != nil {
			
				log.Println("Error al abrir el archivo, Error al crear la carpeta para ese archivo.", err)
			}
		}
	}

	spacef.space = obj
	//Url pasada como valor dir + name + extension -> name dinamico
	spacef.url = url
	//Iniciamos un puntero a SizeFileLine manejado atomicamente por
	//un contador atomico
	spacef.sizeFileLine = new(int64)
	//Iniciamos el contador atomico
	atomic.StoreInt64(spacef.sizeFileLine, spacef.ospaceAtomicUpdateSizeFileLine())
	//Guardamos nuestro puntero de estructura space en el mapa global DiskFile
	


return spacef
}


func (obj *spaceFile ) ospaceAtomicUpdateSizeFileLine()int64 {

var line int64
size, err := obj.file.Seek(0, 2)
if err != nil {

	log.Println(err)

}

if size > 0 {

	size -= obj.lenFields 

}

if size  % obj.sizeLine == 0 {

	line = (size / obj.sizeLine)

}


if size  % obj.sizeLine != 0 {

	line = (size / obj.sizeLine) + 1
	
}

if checkFileCoding(obj.fileCoding , bit) {

	if line > 0 {

		line *= 8 

	}


}

return line -1
}



