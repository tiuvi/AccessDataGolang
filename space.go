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
	"io/fs"
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
	Bit FileCoding = iota
	Byte
)

//Diferenciar archivos de una sola columna, de los archivos multicolumna
type FileTipeByte int
const(
	OneColumn FileTipeByte = iota
	MultiColumn
	FullFile
)

//Actualmente solo para lista de bits
type FileTipeBit int
const(
	ListBit FileTipeBit = iota
)

type Hook string 
const(
	Preformat  = "preformat"
	Postformat = "postformat"
)



type Space struct  {

	//Propiedades comunes a todos los archivos
	Url string
	Exetnsion string
	File *os.File
	Size_line int64
	Size_file int64
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
	
	//Formateadores antes y despues
	Hooker map[string]func([]byte)[]byte

	//Mapa del archivo en memoria
	Search map[string]int64 

	//Array del archivo en memoria
	Index []string

	err error
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

func (obj *Space ) Ospace(){

	//Abrimos el archivo y que permanezca abierto
	obj.File, obj.err = os.OpenFile(obj.Url , os.O_RDWR | os.O_CREATE, 0666)

	//Si error lo pintamos en consola
	if obj.err != nil {

		log.Println(obj.err)

	}

	//Obtenemos filstat en tiempo de copilacion
	var info fs.FileInfo
	info , obj.err = obj.File.Stat()
	
	//Si error lo pintamos en consola
	if obj.err != nil {

		log.Println(obj.err)

	}

	//Obtenemos el tamaño en tiempo de copilacion y lo usamos en tiempo de
	//ejecucion
	obj.Size_file = info.Size()

	//Si el archivo es multicolumna, contamos las columnas.
	if len(obj.IndexSizeColumns) > 1 {

		obj.FileCoding = Byte
		obj.FileTipeByte = MultiColumn

	}
	
	if len(obj.IndexSizeColumns) == 1  {

		obj.FileCoding = Byte
		obj.FileTipeByte = OneColumn

	}
		
	//Obtenemos la extension del fichero
	fileName := strings.SplitN(obj.Url,".",2)

	//sram: fichero que sincroniza un archivo de n lineas
	//con un mapa en la memoria ram asociando valor linea -> n linea
	
	if fileName[1] == "sram" || fileName[1] == "bram"{
		
		obj.FileNativeType |= RamSearch

		var field string
		for val, ind := range obj.IndexSizeColumns {

			if ind[0] == 0 {

				field = val
				break
			}
		}

		mapColumn := *obj.NewSearchSpace(0,(obj.Size_file / obj.Size_line) - 1, field)
		obj.Rspace(mapColumn)
	
		obj.Search = make(map[string]int64)

		var x int64
		for x = 0 ; x < (obj.Size_file / obj.Size_line); x++{
			
			obj.Search[ string( mapColumn.Buffer[field][x] ) ] = x
			
			
		}
	
		
	}
	

	//iram: archivo que sincroniza un archivo de n lineas 
	//con un array en la memoria asociando n lineas -> n valores
	
	if fileName[1] == "iram" || fileName[1] == "bram"{
			
		obj.FileNativeType |= RamIndex

		var field string

		for val, ind := range obj.IndexSizeColumns {

			if ind[0] == 0 {

				field = val
				break
			}
		}

		mapColumn := *obj.NewSearchSpace(0,(obj.Size_file / obj.Size_line) - 1, field)
		obj.Rspace(mapColumn)

		obj.Index = make([]string ,0)
		
		var x int64
		for x = 0 ; x < (obj.Size_file / obj.Size_line); x++{
			
			obj.Index = append(obj.Index, string( mapColumn.Buffer[field][x] ))

		}
		
	}


	//tdisk son archivos que se abren en tiempo de copilacion y 
	//permanecen abiertos en tiempo de ejecucion
	if fileName[1] == "tdisk" {

		obj.FileNativeType  |= tdisk

	}

	//disk: Son archivos normales que se abren y se cierran en
	//tiempo de ejecucion
	if fileName[1] == "disk" {
		
		obj.FileNativeType  |= disk

	}

	//bdisk: Lista de bit en un archivo disk
	if fileName[1] == "bdisk" {

		//La primera propiedad indica si se cierra el archivo
		obj.FileNativeType  |= disk
		obj.FileCoding      = Bit
		obj.FileTipeBit     = ListBit
	
	}

}
