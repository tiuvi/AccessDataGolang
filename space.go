package bd

//Actualizacion de archivos, api rest
//Diferenciar el acceso entre datos privados y publicos
//id , user , token
//Siguiente paso imprimir una respuesta

import (
	"log"
	"os"
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
	Dir string
	Name string
	Extension string
	url string
	SizeLine int64

	//Indice de columnas y tamaño de columna
	IndexSizeColumns map[string][2]int64

	//Formateadores antes y despues
	Hooker map[string]func([]byte)[]byte
	
	FileCoding   FileCoding
	FileTipeByte FileTipeByte
	FileTipeBit  FileTipeBit
	FileTypeDir  FileTypeDir
	compilation bool

}

type spaceFile struct {

	File *os.File
	Url string
	IndexSizeColumns map[string][2]int64
	Hooker map[string]func([]byte)[]byte
	SizeLine int64
	SizeFileLine *int64

	sync.RWMutex
	//Mapa del archivo en memoria
	Search map[string]int64 
	//Array del archivo en memoria
	Index []string
}

type spaceDisk struct{
	DiskFile map[string]*spaceFile
	sync.RWMutex
}

type spaceDeferDisk struct {
	DeferFile map[string]*spaceFile
	Info []deferFileInfo
	sync.RWMutex
}

type deferFileInfo struct {
	Name string
	Time time.Time
}

type spacePermDisk struct{
	PermDisk map[string]*spaceFile
	sync.RWMutex
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

	log.Println("-----",  "New DAc"  ,"-----")

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

func (obj *Space ) Ospace()*spaceFile  {


	obj.url = obj.Dir + obj.Name + "." + obj.Extension

	
	switch obj.FileNativeType {
		
		case obj.FileNativeType & Disk:

			if !obj.compilation {

				obj.ospaceCompilationFile()
				log.Println("Copilacion exitosa")
			}
			return obj.ospaceDisk()
		
		case obj.FileNativeType & DeferDisk:
			if !obj.compilation {

				obj.ospaceCompilationFile()
		
			}
			return obj.ospaceDeferDisk()

			
		case obj.FileNativeType & PermDisk:
			if !obj.compilation {

				obj.ospaceCompilationFile()
		
			}
			return obj.ospacePermDisk()

			
		case obj.FileNativeType & Directory:
			obj.ospaceDirectory()
			break

		default:
			log.Fatalln("Es obligatorio definir el FileNativeType de la estructura. ", obj.Dir,obj.Name,obj.Extension)
	}

	return nil
}



type UnrealSpace[]*Space
func ( US UnrealSpace ) OpenDataAccess(){


}