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



const(
	//Archivo mono cololmna especial para guardar un solo valor
	Odac  = "odac"
	//Archivo multicolumna especial para guardar varios valores
	Mdac  = "mdac"
	//Archivo con un indice o array permanente en la ram
	Iram  = "iram"
	//Archivo con un mapa permanente en la ram
	Sram  = "sram"
	//Archivo con un mapa y un indice permanente en la ram
	Bram  = "bram"
	//Lista de bit con dos estados posibles verdadero y falso
	BitList  = "bitlist"
)

var extensionFile = map[string]string{
	//Archivos
	Odac:     "Archivo mono cololmna especial para guardar un solo valor",
	Mdac:     "Archivo multicolumna especial para guardar varios valores",
	Iram:     "Archivo con un indice o array permanente en la ram",
	Sram:     "Archivo con un mapa permanente en la ram",
	Bram:     "Archivo con un mapa y un indice permanente en la ram",
	BitList:  "Lista de bit con dos estados posibles verdadero y falso",
}

type Space struct  {

	//Indica el estado del archivo en la aplicacion
	FileNativeType FileNativeType
	//Propiedades comunes a todos los archivos
	Dir string
	Name string
	//Cambiar Extension
	Extension string
	url string

	SizeLine int64

	//Indice de columnas y tamaño de columna
	IndexSizeColumns map[string][2]int64
	LenColumns int
	//Formateadores antes y despues
	Hooker map[string]func([]byte)[]byte
	
	FileCoding   FileCoding
	FileTipeByte FileTipeByte
	FileTipeBit  FileTipeBit
	FileTypeDir  FileTypeDir
	compilation bool

}

type spaceFile struct {
	*Space
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

func (obj *Space ) OSpace()*spaceFile  {


	if !obj.compilation {

		obj.ospaceCompilationFile()

	}

	if len(obj.url) == 0 {

		if len(obj.Name) == 0 {

			log.Fatalln("Nombre de archivo vacio en: ", obj.Dir, obj.Extension)
	
		}
		obj.url = obj.Dir + obj.Name + "." + obj.Extension
		
	}


	if CheckFileNativeType(obj.FileNativeType, Disk ){

		return obj.ospaceDisk()

	} else 	if CheckFileNativeType(obj.FileNativeType, DeferDisk ){

		return obj.ospaceDeferDisk()

	} else if CheckFileNativeType(obj.FileNativeType, PermDisk ){

		return obj.ospacePermDisk()

	} else if CheckFileNativeType(obj.FileNativeType, Directory ){

		obj.ospaceDirectory()
		return nil
	}

	log.Fatalln("Es obligatorio definir el FileNativeType de la estructura. ", obj.Dir,obj.Name,obj.Extension)
	return nil
}


func (obj *Space ) ReNameSpace(name string)*Space {

	if len(name) == 0 {

		log.Fatalln("Nombre de archivo vacio en: ", obj.Dir, obj.Extension)

	}

	NewFile := *obj
	NewFile.Name = name
	NewFile.url = NewFile.Dir + NewFile.Name + "." + NewFile.Extension

	return &NewFile
}






func CheckFileNativeType(base FileNativeType, compare FileNativeType)(bool){

	if (base & compare) != 0 {

		return true

	}
	return false
}

func CheckFileCoding(base FileCoding, compare FileCoding)(bool){

	if (base & compare) != 0 {

		return true

	}
	return false
}

func CheckFileTipeBit(base FileTipeBit, compare FileTipeBit)(bool){

	if (base & compare) != 0 {

		return true

	}
	return false
}

func CheckFileTipeByte(base FileTipeByte, compare FileTipeByte)(bool){

	if (base & compare) != 0 {

		return true

	}
	return false
}

func CheckFileTypeDir(base FileTypeDir, compare FileTypeDir)(bool){

	if (base & compare) != 0 {

		return true

	}
	return false
}