package bd

import(
	"log"
	"os"
	"time"
	"sync/atomic"
)


func (obj *Space ) ospaceCompilationFile()bool {

	if len(obj.Dir) == 0 {

		log.Fatalln("Ruta directorio vacio en: ", obj.Dir, obj.Extension)

	}
	
	if len(obj.Extension) == 0 {

		log.Fatalln("Extension vacia en: ", obj.Dir, obj.Extension )

	}

	if _ , found := extensionFile[obj.Extension];  !found{

		log.Fatalln("Extension no valida en: ", obj.Dir, obj.Extension)

	}



	obj.lenColumns = int64(len(obj.IndexSizeColumns))
	obj.lenFields  = int64(len(obj.IndexSizeFields))

	if obj.lenColumns == 0 && obj.lenFields == 0{
		
		log.Fatalln("Iniciaste un archivo de acceso a datos sin columnas y sin campos.", obj.Dir)
	}

	if _ , found := obj.IndexSizeColumns["buffer"]; found {

		log.Fatalln("La palabra buffer esta reservada para el programa.", obj.Dir, obj.IndexSizeColumns)

	}

	checkMap := make(map[string]bool)

	if obj.lenFields != 0 {

		obj.lenFields       = 0
		var checkSizeFields int64 = 0

		for name, val := range obj.IndexSizeFields{

			checkMap[name] = true

			calcSizeLine := (val[1] - val[0])
			if calcSizeLine <= 0 {
	
				log.Fatalln("Los fields no pueden tener un tamaño inferior a cero.",
				obj.Dir,obj.IndexSizeColumns)
			}

			obj.lenFields += calcSizeLine

			if val[1] >= checkSizeFields {

				checkSizeFields = val[1]
			}
		}

		if checkSizeFields != int64(obj.lenFields){
			log.Fatalln("Los campos estan mal escritos, Ejemplo: field1: 0,20; field2:20,30;", obj.Dir)
		}

	}



	//Actualizamos el valor del ancho de la linea
	if obj.lenColumns != 0 {

		var checkSizeColumns int64 = 0

		for name , val := range obj.IndexSizeColumns{

			if obj.IndexSizeFields != nil {

				if found := checkMap[name]; found{

					log.Fatalln("El campo: " + name +" coincide con la columna: " + name + " en: ",
					obj.Dir)
				
				}
			}
	
			calcSizeLine := (val[1] - val[0])
			if calcSizeLine <= 0 {

				log.Fatalln("Las columnas no pueden tener un tamaño inferior a cero.",
				obj.Dir,obj.IndexSizeColumns)
			}

			obj.SizeLine += calcSizeLine

			if val[1] >= checkSizeColumns {

				checkSizeColumns = val[1]
			}
			
		}
		
		if checkSizeColumns != obj.SizeLine {

			log.Fatalln("Las columnas estan mal escritos, Ejemplo: column1: 0,20; column2:20,30;", obj.Dir)
		
		}

	}
	

	
	//Lectura de archivos monocolumna
	if obj.lenColumns == 1 && obj.Extension == "odac" {

		if obj.lenColumns > 1 {

			log.Fatalln("As iniciado un archivo de una columna con multiples columnas", obj.Dir,obj.IndexSizeColumns)

		}

		obj.ospaceCompilationFileUpdateColumn(obj.lenColumns)
		obj.compilation = true
		return true
	}

	//Si el archivo es multicolumna, contamos las columnas.
	if obj.lenColumns > 1 && obj.Extension == "mdac"{

		if obj.lenColumns == 1 {

			log.Fatalln("As iniciado un archivo de un multicolumna con una columna", obj.Dir,obj.IndexSizeColumns)

		}

		obj.ospaceCompilationFileUpdateColumn(obj.lenColumns)
		obj.compilation = true
		return true
	}
	

	//sram: fichero que sincroniza un archivo de n lineas
	//con un mapa en la memoria ram asociando valor linea -> n linea
	if obj.Extension == "sram"{

		obj.ospaceCompilationFileUpdateColumn(obj.lenColumns)
		obj.FileNativeType = RamSearch | obj.FileNativeType
		obj.compilation = true
		return true
	}
	

	//iram: archivo que sincroniza un archivo de n lineas 
	//con un array en la memoria asociando n lineas -> n valores
	if obj.Extension == "iram"{
		
		obj.ospaceCompilationFileUpdateColumn(obj.lenColumns)
		obj.FileNativeType |= RamIndex
		obj.compilation = true
		return true
	}

	if obj.Extension == "bram" {

		obj.ospaceCompilationFileUpdateColumn(obj.lenColumns)
		obj.FileNativeType |= RamIndex | RamSearch
		obj.compilation = true
		return true
	}


	//bdisk: Lista de bit en un archivo disk
	if obj.Extension == "bitlist" {

		obj.SizeLine        = 1
		obj.FileCoding      = Bit
		obj.FileTipeBit     = ListBit
		obj.compilation = true
		return true
	}

	log.Fatalln("Error de copilacion, Ospace.go; funcion: ospaceCompilationFile ; Linea:193", obj.Dir)
	return false
}



//Variable Global de ospaceDisk
var diskSpace = &spaceDisk{
	DiskFile: make(map[string]*spaceFile),
}

//Abrimos el espacio en disco
func (obj *Space ) ospaceDisk(name string)*spaceFile {

	url := obj.Dir + name + "." + obj.Extension

	spacef , found := diskSpace.DiskFile[url]
	if !found {

		//Creamos una nueva referencia a spaceFile
		spacef = obj.newSpaceFile(url)
		//Unico bloqueo cuando se abre el archivo para mantener la atomicidad
		diskSpace.Lock()
		//Guardamos nuestra referencia al archivo en el mapa
		diskSpace.DiskFile[url] = spacef
		//Quitamos el cerrojo de la estructura diskSpace
		diskSpace.Unlock()

	}
	
	return spacef
}

var deferSpace = &spaceDeferDisk{
	DeferFile: make(map[string]*spaceFile),
	Info: make([]deferFileInfo,0),
}

func (obj *Space ) ospaceDeferDisk(name string)*spaceFile{

	url := obj.Dir + name + "." + obj.Extension

	
	spacef , found := deferSpace.DeferFile[url]
	if !found {

		//Creamos una nueva referencia a spaceFile
		spacef = obj.newSpaceFile(url)
		//Unico bloqueo cuando se abre el archivo para mantener la atomicidad
		deferSpace.Lock()
		//Creamos una nueva referencia a spaceFile
		deferSpace.DeferFile[url] = spacef
		//Añadimos un elemento a la cola de array para su posterior
		//eleiminacion del mapa en orden
		deferSpace.Info = append( deferSpace.Info , deferFileInfo{url,time.Now()})
		//Quitamos el cerrojo de la estructura diskSpace
		deferSpace.Unlock()

	}
	return spacef
}

//Variable Global de ospaceDisk
var permSpace = &spacePermDisk{
	PermDisk: make(map[string]*spaceFile),
}

func (obj *Space ) ospacePermDisk(name string)*spaceFile{

	url := obj.Dir + name + "." + obj.Extension

	spacef , found := permSpace.PermDisk[url]
	if !found {

		//Creamos una nueva referencia a spaceFile
		spacef = obj.newSpaceFile(url)
		//Unico bloqueo cuando se abre el archivo para mantener la atomicidad
		permSpace.Lock()
		//Creamos una nueva referencia a spaceFile
		permSpace.PermDisk[url] = spacef
		//Quitamos el cerrojo de la estructura diskSpace
		permSpace.Unlock()
	

	}

	return spacef

}





func (obj *Space ) ospaceDirectory(name string){


	url := obj.Dir + name + "." + obj.Extension

	_ , err := os.Stat(url)
	if err != nil {

		log.Println("Crea una funcion de creacion de directorios.")

	}

	return
/*
	if infoDir.IsDir() {

		obj.FileTypeDir = EmptyDir
		return
	}
	*/

}











func (obj *Space ) ospaceCompilationFileUpdateColumn(LenIndexSizeColumns int64) {

		//Lectura de archivos monocolumna
		if LenIndexSizeColumns == 1 {

			obj.FileCoding = Byte
			obj.FileTipeByte = OneColumn
			return
		}
	
		//Si el archivo es multicolumna, contamos las columnas.
		if LenIndexSizeColumns > 1 {
	
			obj.FileCoding = Byte
			obj.FileTipeByte = MultiColumn
			return
		}

}



func (obj *Space )newSpaceFile(url string)*spaceFile{

		var err error
		//Creamos una nueva referencia a spaceFile
		spacef := new(spaceFile)

		//Abrimos el archivo
		spacef.File, err = os.OpenFile(url , os.O_RDWR | os.O_CREATE, 0666)
		if err != nil {
			//Migrar los errores de archivo a un log de archivo
			log.Println("Error al abrir o crear el archivo.", err)
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

	line = (size / obj.SizeLine)

	if CheckFileCoding(obj.FileCoding , Bit) {
		if line > 0 {

			line *= 8 

		}
	

	}

	return line -1
}

func (sF *spaceFile) ospaceCompilationFileRamSearch() {

	sF.FileNativeType |= RamSearch

	var field string
	for val, ind := range sF.IndexSizeColumns {

		if ind[0] == 0 {

			field = val
			break
		}
	}

	mapColumn := sF.BRspace( BuffMap, 0, *sF.SizeFileLine  , field)
	mapColumn.Rspace()

	sF.Search = make(map[string]int64)

	var x int64
	for x = 0 ; x <= *sF.SizeFileLine; x++{
		
		sF.Search[ string( mapColumn.BufferMap[field][x] ) ] = x
		
		
	}

}




func (sF *spaceFile ) ospaceCompilationFileRamIndex() {

	sF.FileNativeType |= RamIndex

	var field string

	for val, ind := range sF.IndexSizeColumns {

		if ind[0] == 0 {

			field = val
			break
		}
	}

	mapColumn := sF.BRspace(BuffMap,0, *sF.SizeFileLine, field)
	mapColumn.Rspace()

	sF.Index = make([]string ,0)
	
	var x int64
	for x = 0 ; x <= *sF.SizeFileLine; x++{
		
		sF.Index = append(sF.Index, string( mapColumn.BufferMap[field][x] ))

	}

}