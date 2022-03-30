package bd

import(
	"log"
	"os"
	"time"
	"sync/atomic"
)


func (obj *Space ) ospaceCompilationFile() {

	
	LenIndexSizeColumns := len(obj.IndexSizeColumns)
	if LenIndexSizeColumns == 0 {
		
		log.Fatalln("Iniciaste un archivo de acceso a datos sin columnas", obj.url)
	}

	//Actualizamos el valor del ancho de la linea
	for _, val := range obj.IndexSizeColumns{

		obj.SizeLine += (val[1] - val[0])
	}




	_ ,found := extensionFile[obj.Extension]
	if !found {

		log.Fatalln("Extension de archivo no soportado. ", obj.url )

	}



	//Lectura de archivos monocolumna
	if LenIndexSizeColumns == 1 && obj.Extension == "odac" {

		obj.ospaceCompilationFileUpdateColumn(LenIndexSizeColumns)
		obj.compilation = true
		return
	}

	//Si el archivo es multicolumna, contamos las columnas.
	if LenIndexSizeColumns > 1 && obj.Extension == "mdac"{

		obj.ospaceCompilationFileUpdateColumn(LenIndexSizeColumns)
		obj.compilation = true
		return
	}
	

	//sram: fichero que sincroniza un archivo de n lineas
	//con un mapa en la memoria ram asociando valor linea -> n linea
	if obj.Extension == "sram"{

		obj.ospaceCompilationFileUpdateColumn(LenIndexSizeColumns)
		obj.FileNativeType |= RamSearch
		obj.compilation = true
		return
	}
	

	//iram: archivo que sincroniza un archivo de n lineas 
	//con un array en la memoria asociando n lineas -> n valores
	if obj.Extension == "iram"{
		
		obj.ospaceCompilationFileUpdateColumn(LenIndexSizeColumns)
		obj.FileNativeType |= RamIndex
		obj.compilation = true
		return
	}

	if obj.Extension == "bram" {

		obj.ospaceCompilationFileUpdateColumn(LenIndexSizeColumns)
		obj.FileNativeType |= RamIndex | RamSearch
		obj.compilation = true
		return
	}


	//bdisk: Lista de bit en un archivo disk
	if obj.Extension == "bitlist" {

		obj.SizeLine        = 1
		obj.FileCoding      = Bit
		obj.FileTipeBit     = ListBit
		obj.compilation = true
		return
	}

	

}



//Variable Global de ospaceDisk
var diskSpace = &spaceDisk{
	DiskFile: make(map[string]*spaceFile),
}

//Abrimos el espacio en disco
func (obj *Space ) ospaceDisk()*spaceFile {

	spacef , found := diskSpace.DiskFile[obj.url]
	if !found {

		//Creamos una nueva referencia a spaceFile
		spacef = obj.newSpaceFile()
		//Unico bloqueo cuando se abre el archivo para mantener la atomicidad
		diskSpace.Lock()
		//Guardamos nuestra referencia al archivo en el mapa
		diskSpace.DiskFile[obj.url] = spacef
		//Quitamos el cerrojo de la estructura diskSpace
		diskSpace.Unlock()

	}

	return spacef
}

var deferSpace = &spaceDeferDisk{
	DeferFile: make(map[string]*spaceFile),
	Info: make([]deferFileInfo,0),
}

func (obj *Space ) ospaceDeferDisk()*spaceFile{

	spacef , found := deferSpace.DeferFile[obj.url]
	if !found {

		//Creamos una nueva referencia a spaceFile
		spacef = obj.newSpaceFile()
		//Unico bloqueo cuando se abre el archivo para mantener la atomicidad
		deferSpace.Lock()
		//Creamos una nueva referencia a spaceFile
		deferSpace.DeferFile[obj.url] = spacef
		//Añadimos un elemento a la cola de array para su posterior
		//eleiminacion del mapa en orden
		deferSpace.Info = append( deferSpace.Info , deferFileInfo{obj.url,time.Now()})
		//Quitamos el cerrojo de la estructura diskSpace
		deferSpace.Unlock()

	}
	return spacef
}

//Variable Global de ospaceDisk
var permSpace = &spacePermDisk{
	PermDisk: make(map[string]*spaceFile),
}

func (obj *Space ) ospacePermDisk()*spaceFile{

	spacef , found := permSpace.PermDisk[obj.url]
	if !found {

		//Creamos una nueva referencia a spaceFile
		spacef = obj.newSpaceFile()
		//Unico bloqueo cuando se abre el archivo para mantener la atomicidad
		permSpace.Lock()
		//Creamos una nueva referencia a spaceFile
		permSpace.PermDisk[obj.url] = spacef
		//Quitamos el cerrojo de la estructura diskSpace
		permSpace.Unlock()

	}

	return spacef

}





func (obj *Space ) ospaceDirectory(){


	infoDir , err := os.Stat(obj.Name)
	if err != nil {

		log.Println("Crea una funcion de creacion de directorios.")

	}

	if infoDir.IsDir() {

		obj.FileTypeDir = EmptyDir
		return
	}

}





/*
func (obj *spaceFile ) ospaceCompilationFileRamSearch() {

	obj.FileNativeType |= RamSearch

	var field string
	for val, ind := range obj.IndexSizeColumns {

		if ind[0] == 0 {

			field = val
			break
		}
	}

	mapColumn := *obj.NewSearchSpace(0, *obj.SizeFileLine  , field)
	obj.Rspace(mapColumn)

	obj.Search = make(map[string]int64)

	var x int64
	for x = 0 ; x <= obj.SizeFileLine; x++{
		
		obj.Search[ string( mapColumn.Buffer[field][x] ) ] = x
		
		
	}

}

func (obj *spaceFile ) ospaceCompilationFileRamIndex() {

	obj.FileNativeType |= RamIndex

	var field string

	for val, ind := range obj.IndexSizeColumns {

		if ind[0] == 0 {

			field = val
			break
		}
	}

	mapColumn := *obj.NewSearchSpace(0, *obj.SizeFileLine, field)
	obj.Rspace(mapColumn)

	obj.Index = make([]string ,0)
	
	var x int64
	for x = 0 ; x <= obj.SizeFileLine; x++{
		
		obj.Index = append(obj.Index, string( mapColumn.Buffer[field][x] ))

	}

}
*/

func (obj *Space ) ospaceCompilationFileUpdateColumn(LenIndexSizeColumns int) {

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

/*
func (obj *Space ) ospaceCompilationFileUpdateSizeFile() {
	
	
	test := time.Now()
	info, err := obj.File.Stat()
	log.Println("tiempo de escritura info: ", time.Since(test).Nanoseconds())

	if err != nil {
		log.Println("File.Stat error: ",err)
	}

	//Calculamos el numero de lineas del fichero
	//obj.SizeFileLine = (info.Size() / obj.SizeLine) - 1
	atomic.AddInt64(&obj.SizeFileLine, (info.Size() / obj.SizeLine) - 1 )
	

	size, err := obj.File.Seek(0, 2)
	if err != nil {
		log.Println(err)
	}
	atomic.AddInt64(&obj.SizeFileLine, (size / obj.SizeLine) - 1 )
}
*/

func (obj *Space )newSpaceFile()*spaceFile{

		var err error
		//Creamos una nueva referencia a spaceFile
		spacef := new(spaceFile)

		//Abrimos el archivo
		spacef.File, err = os.OpenFile(obj.url , os.O_RDWR | os.O_CREATE, 0666)
		if err != nil {
			//Migrar los errores de archivo a un log de archivo
			log.Println("Error al abrir o crear el archivo.", err)
		}

		//Url pasada como valor dir + name + extension -> name dinamico
		spacef.Url = obj.url
		//Columnas pasadas por referencia al spacio
		spacef.IndexSizeColumns = obj.IndexSizeColumns 
		//Hookers pasados por referencia para filtrar datos
		spacef.Hooker = obj.Hooker
		//Size line pasado por valor (Valor total de la linea)
		spacef.SizeLine     = obj.SizeLine
		//Iniciamos un puntero a SizeFileLine manejado atomicamente por
		//un contador atomico
		spacef.SizeFileLine = new(int64)
		//Iniciamos el contador atomico
		atomic.StoreInt64(spacef.SizeFileLine, spacef.ospaceAtomicUpdateSizeFileLine())
		//Guardamos nuestro puntero de estructura space en el mapa global DiskFile
		


	return spacef
}

func (obj *spaceFile ) ospaceAtomicUpdateSizeFileLine()int64 {
	
	size, err := obj.File.Seek(0, 2)
	if err != nil {
		log.Println(err)
	}
	return (size / obj.SizeLine) -1
}

