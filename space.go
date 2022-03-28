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
	Disk FileNativeType = 1 << iota
	DeferDisk 
	PermDisk
	Directory
	RamSearch
	RamIndex
)


//Diferenciar archivos de Bit de archivos de Byte
type FileCoding int
const(
	Bit FileCoding = iota + 1
	Byte
	Dir
)

//Diferenciar archivos de una sola columna, de los archivos multicolumna
type FileTipeByte int
const(
	OneColumn FileTipeByte = iota + 1
	MultiColumn
	FullFile
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

	
	//Indica el estado del archivo en la aplicacion
	FileNativeType FileNativeType
	//Propiedades comunes a todos los archivos
	Url string
	//Indice de columnas y tamaño de columna
	IndexSizeColumns map[string][2]int64

	Name string
	Extension string

	File *os.File
	Spacef spaceDescriptors
	spaceDisk spaceDisk

	compilation bool
	SizeLine int64
	SizeFileLine int64

	//Casos de uso para las funciones writes y read
	FileCoding   FileCoding
	FileTipeByte FileTipeByte
	FileTipeBit  FileTipeBit
	FileTypeDir  FileTypeDir

	//Formateadores antes y despues
	Hooker map[string]func([]byte)[]byte
	
	sync.RWMutex
	//Mapa del archivo en memoria
	Search map[string]int64 

	//Array del archivo en memoria
	Index []string

}



type DescriptorInformation struct {
	Name string
	Time time.Time
}
type DescriptorInformations []DescriptorInformation

type spaceDescriptors struct {
	Descriptor map[string] *os.File
	Information []DescriptorInformation
	sync.RWMutex

}

var mSpace = &spaceDescriptors{
	Descriptor: make(map[string]*os.File,0),
	Information: make(DescriptorInformations,0),
}

type spaceDisk struct{
	DiskFile map[string] *os.File
	sync.RWMutex
}

var diskSpace = &spaceDisk{
	DiskFile: make(map[string]*os.File,0),
}

var extensionFile = map[string]string{
	//Archivos
	"odac":"Archivo mono cololmna especial para guardar un solo valor",
	"mdac":"Archivo multicolumna especial para guardar varios valores",
	"iram":"Archivo con un indice o array permanente en la ram",
	"sram":"Archivo con un mapa permanente en la ram",
	"bram":"Archivo con un mapa y un indice permanente en la ram",
	"bitlist":"Lista de bit con dos estados posibles verdadero y falso",
}

var extensionDir = map[string]string{
	//Directorios
	"dir":"Gestiona una carpeta como si fuera un unico archivo",
}

func NewDac(){

	log.Println("-----","     ","-----")
	log.Println("New DAc")
	log.Println("-----","     ","-----")

	go dacTimerCloserDeferFile()
	
	go dacTimerCloserDiskFile()

		
}









//Iniciamos el objeto Ospace

//Todos los datos se organizan por numero de linea, este numero de linea
// es inmutable, en caso de querer borrar simplemente se deja vacio y 
//se rellenara mas adelante con otro dato.

//Tenemos dos tipos de organizacion de datos una unica columna asignada
//a una fila y multiples columnas asociadas a una fila

// Tenemos 4 tipos de archivos Ramsearch, Ramindex, PermDisk, DeferDisk y Disk
// Dependiendo de si el archivo permanece abierto y si necesita
//sincronizacion con ram.

//Retocar funcion para una ejecucion diferente en tiempo de copilacion y tiempo de
//ejecucion

func (obj *Space ) Ospace(){


	//Url obligatoria, siempre tenemos que tener definida a que parte de 
	//los archivos se va acceder.
	//En caso de ruta vacia se entendera como un error grave de seguridad.
	 if obj.Url == ""{
		log.Fatalln("Ruta de archivo no definida.")
	 }

	//Tipado de los archivos, la extenion declara el tipo de archivo explicitamente.
	//En caso de dos puntos se entendera como fallo de seguridad y se cerrara el programa
	if obj.Name == "" || obj.Extension == "" {

		fileName := strings.SplitN(obj.Url,".",2)
		if len(fileName) == 1 {
			log.Fatalln("Extension de archivo no definida.")
	
		}
		if len(fileName) == 2 {
			obj.Name = fileName[0]
			obj.Extension =  fileName[1]
		}
		if len(fileName) > 2 {
			log.Fatalln("Solo se permite un punto para declarar la extension.")
		}
	}


	//Funciones que abren archivos
	//Tipan el acceso a los archivos o directorios
	switch obj.FileNativeType {
		
		case obj.FileNativeType & Disk:
			obj.ospaceDisk()
			obj.ospaceCompilationFile()
			break
		case obj.FileNativeType & DeferDisk:
			obj.ospaceDeferDisk()
			obj.ospaceCompilationFile()
			break
		case obj.FileNativeType & PermDisk:
			obj.ospacePermDisk()
			obj.ospaceCompilationFile()
			break
		case obj.FileNativeType & Directory:
			obj.ospaceDirectory()
			break
		default:
			log.Fatalln("Es obligatorio definir el FileNativeType de la estructura. ", obj.Url)
	}

}
