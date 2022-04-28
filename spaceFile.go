package bd

import(
	"log"
	"os"
	"sync/atomic"
)


func (obj *Space )newSpaceFile(folderString string, name string)*spaceFile{

	url := obj.Dir + folderString + name + "." + obj.Extension

	var err error
	//Creamos una nueva referencia a spaceFile
	spacef := new(spaceFile)

	//Abrimos el archivo
	spacef.File, err = os.OpenFile(url , os.O_RDWR | os.O_CREATE, 0666)
	if err != nil {
		//Migrar los errores de archivo a un log de archivo
		log.Println("Error al abrir o crear el archivo.", err)

		if  os.IsNotExist(err) {
			
			err = os.MkdirAll(obj.Dir + folderString , 0666)
			if err != nil {

				log.Println("Error al crear la carpeta para ese archivo. ", err)
			} else {

				// obj.LogNewError(Message , "Nueva ruta de archivo en: ", obj.Dir + folderString )
				
			}
			spacef.File, err = os.OpenFile(url , os.O_RDWR | os.O_CREATE, 0666)
			if err != nil {
			
				log.Println("Error al abrir el archivo, Error al crear la carpeta para ese archivo.", err)
			}
		}
	}

	spacef.Space = obj
	//Url pasada como valor dir + name + extension -> name dinamico
	spacef.Url = url
	//Iniciamos un puntero a SizeFileLine manejado atomicamente por
	//un contador atomico
	spacef.SizeFileLine = new(int64)
	//Iniciamos el contador atomico
	atomic.StoreInt64(spacef.SizeFileLine, spacef.ospaceAtomicUpdateSizeFileLine())
	//Guardamos nuestro puntero de estructura space en el mapa global DiskFile
	


return spacef
}


func (obj *spaceFile ) ospaceAtomicUpdateSizeFileLine()int64 {

var line int64
size, err := obj.File.Seek(0, 2)
if err != nil {

	log.Println(err)

}

if size > 0 {

	size -= obj.lenFields 

}

if size  % obj.SizeLine == 0 {

	line = (size / obj.SizeLine)

}


if size  % obj.SizeLine != 0 {

	line = (size / obj.SizeLine) + 1
	
}

if CheckFileCoding(obj.FileCoding , Bit) {
	if line > 0 {

		line *= 8 

	}


}

return line -1
}



