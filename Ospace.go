package bd

import(
	"log"
	"os"
	"time"
	"sync/atomic"
)


func (obj *Space ) ospaceDisk(){

	var err error

	_ , found := diskSpace.DiskFile[obj.Url]
	if !found {
		obj.File, err = os.OpenFile(obj.Name + "." + obj.Extension , os.O_RDWR | os.O_CREATE, 0666)
		if err != nil {
			log.Fatalln("Error al abrir o crear el archivo.", err)
		}
		diskSpace.Lock()
		diskSpace.DiskFile[obj.Url] = obj.File
	
		diskSpace.Unlock()
	}

	
}

func (obj *Space ) ospaceDeferDisk(){

	var err error

	_ , found := mSpace.Descriptor[obj.Url]
	if !found {

		obj.File, err = os.OpenFile(obj.Url , os.O_RDWR | os.O_CREATE, 0666)
		if err != nil {
			log.Fatalln("Error al abrir o crear el archivo.", err)
		}
		
		mSpace.Lock()
		mSpace.Descriptor[obj.Url] = obj.File
		mSpace.Information = append( mSpace.Information , DescriptorInformation{obj.Url,time.Now()}  )
		mSpace.Unlock()

	}
	
}

func (obj *Space ) ospacePermDisk(){

	var err error
	if obj.compilation {
		return
	}

	obj.File, err = os.OpenFile(obj.Url , os.O_RDWR | os.O_CREATE, 0666)
	if err != nil {
		log.Fatalln("Error al abrir o crear el archivo.", err)
	}
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



func (obj *Space ) ospaceCompilationFile() {

	if obj.compilation {
		return
	}
	
	LenIndexSizeColumns := len(obj.IndexSizeColumns)
	if LenIndexSizeColumns == 0 {
		
		log.Fatalln("Iniciaste un archivo de acceso a datos sin columnas", obj.Url)
	}

	//Actualizamos el valor del ancho de la linea
	for _, val := range obj.IndexSizeColumns{
		obj.SizeLine += (val[1] - val[0])
	}
	log.Println("Actualizacion SizeLine: " , obj.SizeLine)



	_ ,found := extensionFile[obj.Extension]
	if !found {
		log.Fatalln("Extension de archivo no soportado. ", obj.Url )
	}

	//Lectura de archivos monocolumna
	if LenIndexSizeColumns == 1 && obj.Extension == "odac" {

		obj.ospaceCompilationFileUpdateSizeFile()
		obj.FileCoding = Byte
		obj.FileTipeByte = OneColumn
		obj.compilation = true
		return
	}

	//Si el archivo es multicolumna, contamos las columnas.
	if LenIndexSizeColumns > 1 && obj.Extension == "mdac"{

		obj.ospaceCompilationFileUpdateSizeFile()
		obj.FileCoding = Byte
		obj.FileTipeByte = MultiColumn
		obj.compilation = true
		return
	}
	

	//sram: fichero que sincroniza un archivo de n lineas
	//con un mapa en la memoria ram asociando valor linea -> n linea
	
	if obj.Extension == "sram"{

		obj.FileNativeType |= RamSearch
		obj.ospaceCompilationFileUpdateSizeFile()
		obj.ospaceCompilationFileUpdateColumn(LenIndexSizeColumns)
		obj.ospaceCompilationFileRamSearch()
		obj.compilation = true
		return
	}
	

	//iram: archivo que sincroniza un archivo de n lineas 
	//con un array en la memoria asociando n lineas -> n valores
	
	if obj.Extension == "iram"{
		
		obj.FileNativeType |= RamIndex
		obj.ospaceCompilationFileUpdateSizeFile()
		obj.ospaceCompilationFileUpdateColumn(LenIndexSizeColumns)
		obj.ospaceCompilationFileRamIndex()
		obj.compilation = true
		return
	}

	if obj.Extension == "bram" {

		obj.FileNativeType |= RamIndex | RamSearch
		obj.ospaceCompilationFileUpdateSizeFile()
		obj.ospaceCompilationFileUpdateColumn(LenIndexSizeColumns)
		obj.ospaceCompilationFileRamSearch()
		obj.ospaceCompilationFileRamIndex()
		obj.compilation = true
		return
	}


	//bdisk: Lista de bit en un archivo disk
	if obj.Extension == "bitlist" {

		obj.SizeLine        = 1
		obj.ospaceCompilationFileUpdateSizeFile()
		obj.FileCoding      = Bit
		obj.FileTipeBit     = ListBit
		obj.compilation = true
		return
	}

	

}



func (obj *Space ) ospaceCompilationFileRamSearch() {

	obj.FileNativeType |= RamSearch

	var field string
	for val, ind := range obj.IndexSizeColumns {

		if ind[0] == 0 {

			field = val
			break
		}
	}

	mapColumn := *obj.NewSearchSpace(0, obj.SizeFileLine  , field)
	obj.Rspace(mapColumn)

	obj.Search = make(map[string]int64)

	var x int64
	for x = 0 ; x <= obj.SizeFileLine; x++{
		
		obj.Search[ string( mapColumn.Buffer[field][x] ) ] = x
		
		
	}

}

func (obj *Space ) ospaceCompilationFileRamIndex() {

	obj.FileNativeType |= RamIndex

	var field string

	for val, ind := range obj.IndexSizeColumns {

		if ind[0] == 0 {

			field = val
			break
		}
	}

	mapColumn := *obj.NewSearchSpace(0, obj.SizeFileLine, field)
	obj.Rspace(mapColumn)

	obj.Index = make([]string ,0)
	
	var x int64
	for x = 0 ; x <= obj.SizeFileLine; x++{
		
		obj.Index = append(obj.Index, string( mapColumn.Buffer[field][x] ))

	}

}


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


func (obj *Space ) ospaceCompilationFileUpdateSizeFile() {

	info, err := obj.File.Stat()
	if err != nil {
		log.Println("File.Stat error: ",err)
	}

	//Calculamos el numero de lineas del fichero
	//obj.SizeFileLine = (info.Size() / obj.SizeLine) - 1
	atomic.AddInt64(&obj.SizeFileLine, (info.Size() / obj.SizeLine) - 1 )
}