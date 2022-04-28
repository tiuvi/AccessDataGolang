package bd

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
)


//Diferenciar archivos de Bit de archivos de Byte
type FileCoding int
const(
	Bit FileCoding = iota + 1
	Byte
)



const(
	Preformat  = "preformat"
	Postformat = "postformat"
)



const(
	//Archivo multicolumna especial para guardar varios valores
	DacByte string  = "dacByte"
	//Lista de bit con dos estados posibles verdadero y falso
	DacBit string   = "dacBit"
)

var extensionFile = map[string]string{
	DacByte:     "Archivo que incluye fields y columnas de bytes",
	DacBit:      "Archivo que incluye fields y una lista de bits",
}



type spaceLen struct {
	name string
	len  int64
}
type Space struct  {

	//TE dice si chequear los archivos y en caso de que si se usa la superglobal SpaceErrors
	Check bool
	*SpaceErrors
	//Indica el estado del archivo en la aplicacion
	FileNativeType FileNativeType
	//Propiedades comunes a todos los archivos
	Dir string

	//Cambiar Extension
	Extension string

	SizeLine int64

	//Indice de fields
	IndexSizeFieldsArray []spaceLen
	IndexSizeFields map[string][2]int64
	lenFields int64

	//Indice de columnas y tamaño de columna
	IndexSizeColumnsArray []spaceLen
	IndexSizeColumns map[string][2]int64
	lenColumns int64
	//Formateadores antes y despues
	Hooker map[string]func(*[]byte)
	
	FileCoding   FileCoding

	compilation bool

}


type spaceFile struct {
	*Space
	File *os.File
	Url string
	//Numero de lineas ATOMICO de un archivo
	SizeFileLine *int64
	//Mutex que sirve para leer y escribir , actualizar mapas , actualizar arrays
	sync.RWMutex
}


type spaceDisk struct{
	DiskFile map[string]*spaceFile
	sync.RWMutex
}

type deferFileInfo struct {
	Name string
	Time time.Time
}

type spaceDeferDisk struct {
	DeferFile map[string]*spaceFile
	Info []deferFileInfo
	sync.RWMutex
}


type spacePermDisk struct{
	PermDisk map[string]*spaceFile
	sync.RWMutex
}

func NewDac(){

	go dacTimerCloserDeferFile()
	
	go dacTimerCloserDiskFile()

	go errorsLog()

}

func (obj *Space ) OSpaceInit()bool  {

	return obj.ospaceCompilationFile()

}

func (obj *Space ) OSpace( name string ,folder... string)*spaceFile  {


	if !obj.compilation {

		log.Fatalln("Este espacio no se ha copilado, copilar mejora la seguridad.", obj.Dir)

	}


	if len(name) == 0 {

		log.Fatalln("Nombre de archivo vacio en: ", obj.Dir, obj.Extension)

	}


	if CheckFileNativeType(obj.FileNativeType, Disk ){

		return obj.ospaceDisk(name, folder)

	} 

	if CheckFileNativeType(obj.FileNativeType, DeferDisk ){

		return obj.ospaceDeferDisk(name, folder)

	} 

	if CheckFileNativeType(obj.FileNativeType, PermDisk ){

		return obj.ospacePermDisk(name, folder)

	}

	log.Fatalln("Es obligatorio definir el FileNativeType de la estructura. ", obj.Dir,obj.Extension)
	return nil
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
