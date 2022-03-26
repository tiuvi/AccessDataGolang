package bd

//Actualizacion de archivos, api rest
//Diferenciar el acceso entre datos privados y publicos
//id , user , token
//Siguiente paso imprimir una respuesta

import (
	"log"
	"os"
	"strings"
	"sync"
	"time"
)

type FileNativeType int64
const(
	disk FileNativeType = 1 << iota
	tdisk
	RamSearch
	RamIndex
)


//Diferenciar archivos de Bit de archivos de Byte
type FileCoding int
const(
	Bit FileCoding = iota + 1
	Byte
	Directory
)

//Diferenciar archivos de una sola columna, de los archivos multicolumna
type FileTipeByte int
const(
	OneColumn FileTipeByte = iota + 1
	MultiColumn
	FullFile
	MultiColumnDefer
)

//Actualmente solo para lista de bits
type FileTipeBit int
const(
	ListBit FileTipeBit = iota + 1
)

type FileTypeDir int
const(
	EmptyDir FileTypeDir = iota + 1
)

type Hook string 
const(
	Preformat  = "preformat"
	Postformat = "postformat"
)





type Space struct  {

	//Propiedades comunes a todos los archivos
	Url string
	Name string
	Exist bool
	Extension string

	File *os.File
	Spacef spaceDescriptors

	Size_line int64
	SizeFileLine int64
	NumberLines int64

	sync.RWMutex
	

	//Indice de columnas y tamaño de columna
	IndexSizeColumns map[string][2]int64

	//Tipo de archivo nativo
	FileNativeType FileNativeType

	//Casos de uso para las funciones writes y read
	FileCoding   FileCoding
	FileTipeByte FileTipeByte
	FileTipeBit  FileTipeBit
	FileTypeDir  FileTypeDir
	//Formateadores antes y despues
	Hooker map[string]func([]byte)[]byte

	//Mapa del archivo en memoria
	Search map[string]int64 

	//Array del archivo en memoria
	Index []string

	err error
}

type DescriptorInformation struct {
	Name string
	Descriptor int64
	Time time.Time
}
type DescriptorInformations []DescriptorInformation

type spaceDescriptors struct {
	Descriptor map[string]int64
	Information []DescriptorInformation
	sync.RWMutex
}

var mSpace = &spaceDescriptors{
	Descriptor: make(map[string]int64,0),
	Information: make(DescriptorInformations,0),
}


func NewDac(){

	//Archivos abiertos
	fileOpen := 1000
	//tikectDb := time.Tick(time.Duration(1) * time.Hour)
	tikectDb := time.Tick(time.Duration(30) * time.Second)
	for range tikectDb {
		/*
		log.Println("Nuevo ticket")
		log.Println("-----","     ","-----")
		log.Println(len(mSpace.Information))
		log.Println(len(mSpace.Descriptor))
		log.Println("-----","     ","-----")
		*/

		if len(mSpace.Descriptor) == 0 {
			log.Println("Descriptor parada ")
			continue
		}

		if len(mSpace.Information) == 0{
			log.Println("Information parada ")
			continue
		}
		log.Println("Lock parada")
		mSpace.Lock()
		
	
		//log.Println("Nuevo Bucle: ", "","")
		//log.Println("-----","     ","-----")
		//log.Println("-----","     ","-----")
		for value := range mSpace.Descriptor{
			log.Println(value)
		}
		//log.Println("-----","     ","-----")
		//log.Println("-----","     ","-----")
		for ind := range mSpace.Information{
			log.Println(mSpace.Information[ind])
		}

		
		//for ind := range mSpace.Information {
		ind :=0
		for len(mSpace.Information) > fileOpen  {
		//	log.Println("-----","     ","-----")
			//log.Println(ind)
			//log.Println("-----","     ","-----")

			/*
			if fileOpen > len(mSpace.Information) {
				break
			}
			*/
			//Cierre del descriptor
			descriptor := os.NewFile(uintptr(mSpace.Information[ind].Descriptor), mSpace.Information[ind].Name)
			err := descriptor.Close()
			if err != nil {
				log.Println(err)
			}
			log.Println("-----","     ","-----")
			log.Println(len(mSpace.Information))
			log.Println(mSpace.Information[ind].Name)
			log.Println("-----","     ","-----")
			//Borrado
			delete(mSpace.Descriptor  ,  mSpace.Information[ind].Name)
			mSpace.Information = append(mSpace.Information[:ind], mSpace.Information[ind+1:]... )
		
				
		}
		//mSpace.Information = mSpace.Information[5:]
			
		/*
		log.Println("-----","     ","-----")
		log.Println("-----","     ","-----")
		for value := range mSpace.Descriptor{
			log.Println(value)
		}
		log.Println("-----","     ","-----")
		log.Println("-----","     ","-----")
		for ind := range mSpace.Information{
			log.Println(mSpace.Information[ind])
		}
		log.Println("-----","     ","-----")
		log.Println("-----","     ","-----")
		log.Println("-----"," END ","-----")
		*/
		mSpace.Unlock()
	}
	



}









//Iniciamos el objeto Ospace

//Todos los datos se organizan por numero de linea, este numero de linea
// es inmutable, en caso de querer borrar simplemente se deja vacio y 
//se rellenara mas adelante con otro dato.

//Tenemos dos tipos de organizacion de datos una unica columna asignada
//a una fila y multiples columnas asociadas a una fila

// Tenemos 4 tipos de archivos Ramsearch, Ramindex, Tdisk y disk
// Dependiendo de si el archivo permanece abierto y si necesita
//sincronizacion con ram.

//Retocar funcion para una ejecucion diferente en tiempo de copilacion y tiempo de
//ejecucion

func (obj *Space ) Ospace(){

	//Tiempo de ejecucion
	if obj.Exist {

		switch obj.FileCoding {

			case Directory:
				return

			default:
				obj.File, obj.err = os.OpenFile(obj.Name + "." + obj.Extension , os.O_RDWR | os.O_CREATE, 0666)
				return
		}
	}

	
	//Obtenemos la extension del fichero
	fileName := strings.SplitN(obj.Url,".",2)
	if len(fileName) == 1 {
		obj.Name = fileName[0]
	}
	if len(fileName) == 2 {
		obj.Name = fileName[0]
		obj.Extension =  fileName[1]
	}


	if obj.Extension == "dir" {

		obj.Extension = ""

		infoDir , err := os.Stat(obj.Name)
		if err != nil {
	
			log.Println("Crea una funcion de creacion de directorios.")
	
		}
	
		if infoDir.IsDir() {
	
			obj.FileCoding  = Directory
			obj.FileTypeDir = EmptyDir
			obj.Exist = true
			return
		}
	}



	//Tiempo de compilacion
	obj.File, obj.err = os.OpenFile(obj.Name + "." + obj.Extension , os.O_RDWR | os.O_CREATE, 0666)

	if obj.err != nil {

		log.Println(obj.err)
		obj.Exist = false
		return
	
	}

	obj.Exist = true

	info, err := obj.File.Stat()
	if err != nil {
		
		log.Println(err)
		obj.Exist = false
	
	}

	obj.Exist = true
	

	obj.SizeFileLine = (info.Size() / obj.Size_line) - 1





	//Si el archivo es multicolumna, contamos las columnas.
	if len(obj.IndexSizeColumns) > 1 {

		obj.FileCoding = Byte
		obj.FileTipeByte = MultiColumn

	}
	
	if len(obj.IndexSizeColumns) == 1  {

		obj.FileCoding = Byte
		obj.FileTipeByte = OneColumn

	}
		

	//sram: fichero que sincroniza un archivo de n lineas
	//con un mapa en la memoria ram asociando valor linea -> n linea
	
	if obj.Extension == "sram" || obj.Extension == "bram"{
		
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
	

	//iram: archivo que sincroniza un archivo de n lineas 
	//con un array en la memoria asociando n lineas -> n valores
	
	if obj.Extension == "iram" || obj.Extension == "bram"{
			
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


	//tdisk son archivos que se abren en tiempo de copilacion y 
	//permanecen abiertos en tiempo de ejecucion
	if obj.Extension == "tdisk" {

		obj.FileNativeType  |= tdisk

	}

	//disk: Son archivos normales que se abren y se cierran en
	//tiempo de ejecucion
	if obj.Extension == "disk" {
		
		obj.FileNativeType  |= disk

	}

		//disk: Son archivos normales que se abren y se cierran en
	//tiempo de ejecucion
	if obj.Extension == "defer" {
		
		obj.FileTipeByte = MultiColumnDefer
		mSpace.Lock()
		mSpace.Descriptor[obj.Url] = int64(obj.File.Fd())
		mSpace.Information = append( mSpace.Information , DescriptorInformation{obj.Url, int64(obj.File.Fd()) , time.Now()}  )
		mSpace.Unlock()

		//log.Println(mSpace.Information[len(mSpace.Information)-1])

	}

	

	//bdisk: Lista de bit en un archivo disk
	if obj.Extension == "bdisk" {

		//La primera propiedad indica si se cierra el archivo
		//obj.FileNativeType  |= disk
		obj.FileCoding      = Bit
		obj.FileTipeBit     = ListBit
		defer obj.File.Close()
		
	
	}

}
