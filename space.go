package bd

import (
	"os"
	"sync"
	"time"
)

/**********************************************************************************************/
/* Constantes */
/**********************************************************************************************/

type fileNativeType int64
const (
	disk fileNativeType = 1 << iota
	deferDisk
	permDisk
)

//Diferenciar archivos de Bit de archivos de Byte
type fileCoding int
const (
	bit fileCoding = iota + 1
	bytes
)

const (
	preformat  = "preformat"
	postformat = "postformat"
)

const (
	//Archivo multicolumna especial para guardar varios valores
	dacByte string = "dacByte"
	//Lista de bit con dos estados posibles verdadero y falso
	dacBit string = "dacBit"
)

var extensionFile = map[string]string{
	dacByte: "Archivo que incluye fields y columnas de bytes",
	dacBit:  "Archivo que incluye fields y una lista de bits",
}

const (
	ColorWhite  = "\u001b[37m"
	ColorBlack  = "\u001b[30m"
	ColorRed    = "\u001b[31m"
	ColorGreen  = "\u001b[32m"
	ColorYellow = "\u001b[33m"
	ColorBlue   = "\u001b[34m"
	Reset       = "\u001b[0m"
)

//Tipos de buffer para las funciones de lectura y escritura.
type fileTypeBuffer int64
const(
	buffMap fileTypeBuffer = 1 << iota
	buffChan
	buffBytes
)




/**********************************************************************************************/
/* 
spaceErrors -> Gestion de errores de la aplicación, test de velocidad y test de memoria.
 */
/**********************************************************************************************/

type errorDac string

const (
	NewRouteFolder errorDac = "newRouteFolder"

	Message           errorDac = "message"
	MessageCopilation errorDac = "messageCopilation"

	Exception errorDac = "exception"
	Warning   errorDac = "warning"
	Fatal     errorDac = "fatal"

	TimeMemory errorDac = "TimeMemoryStat"
)

//Errores globales dac y seguimiento mediante archivo o log.
//En DAC no hay excepciones, ni errores que manejar.
type spaceErrors struct {

	//Log fatales que impiden la propagacion de errores (Recomendado)
	logFatalErrors bool

	//Visualiza todo en consola
	logConsoleErrors bool

	//Activar fichero de errores.
	//Predeterminado: DAC/errors
	logFileError bool

	//Activar log cronometro
	logTimeUse bool

	//Activar fichero de cronometro.
	//Predeterminado: DAC/timer
	logFileTimeUse bool
	//Activar cronometros de apertura de archivos
	logTimeOpenFile bool

	//Activar log de memoria
	logMemoryUse bool

	//Activar fichero de memoria.
	//Predeterminado: DAC/memory
	logFileMemoryUse bool

	//Separador para los resultados en nanosegundos: coma, espacio , punto...
	separatorLog string

	//Niveles de la url DAC a mostrar.
	levelsUrl int
}

//space globales
var timeLog *space
var memoryLog *space
var errorLog *space





/**********************************************************************************************/
/* 
spaceErrors -> Gestion de errores de la aplicación, test de velocidad y test de memoria.
lDAC -> Inicio de la aplicacion global e instancia del DAC central, compartido con el 
resto de espacios, equivalente a una base de datos mysql
*/
/**********************************************************************************************/

//Inicio del path DAC
type lDAC struct {
	*spaceErrors
	newDACFolder    bool
	goDACFolder     bool
	globalDACFolder string
}

//Variable global de DAC
var globalLaunchDac *lDAC





/**********************************************************************************************/
/*
spaceErrors -> Gestion de errores de la aplicación, test de velocidad y test de memoria.
lDAC -> Inicio de la aplicacion global e instancia del DAC central, compartido con el 
resto de espacios, equivalente a una base de datos mysql
space -> Genera un espacio que seria igual a una tabla mysql, se divide en dos tipos de campos
fields campos unicos y columnas campos dinamicos, los espacios son compatibles y recomendados
 para todo tipo de datos tablas de nombres, emails, video, audio, html, jss, css...
*/
/**********************************************************************************************/

type spaceLen struct {
	name string
	len  int64
}
type space struct {
	*lDAC
	//TE dice si chequear los archivos y en caso de que si se usa la superglobal SpaceErrors

	//Indica el estado del archivo en la aplicacion
	fileNativeType fileNativeType
	//Propiedades comunes a todos los archivos
	dir string

	//Cambiar Extension
	extension string

	sizeLine int64

	//Indice de fields
	indexSizeFieldsArray []spaceLen
	indexSizeFields      map[string][2]int64
	lenFields            int64

	//Indice de columnas y tamaño de columna
	indexSizeColumnsArray []spaceLen
	indexSizeColumns      map[string][2]int64
	lenColumns            int64
	//Formateadores antes y despues
	hooker map[string]func(*[]byte)

	fileCoding fileCoding

	compilation bool
}





/**********************************************************************************************/
/* 
spaceErrors -> Gestion de errores de la aplicación, test de velocidad y test de memoria.

lDAC -> Inicio de la aplicacion global e instancia del DAC central, compartido con el 
resto de espacios, equivalente a una base de datos mysql

space -> Genera un espacio que seria igual a una tabla mysql, se divide en dos tipos de campos
fields campos unicos y columnas campos dinamicos, los espacios son compatibles y recomendados
 para todo tipo de datos tablas de nombres, emails, video, audio, html, jss, css...

spaceFile -> Apunta a un archivo y se guardan actualmente en tres diferentes mapas donde
permanecen abiertos hasta que un timer los cierra, funcionan con locs de lectura y escritura.
*/
/**********************************************************************************************/

//Variable Global diskSpace: Estos archivos son temporales disco y se cierran cada x * time
var diskSpace = &spaceDisk{
	diskFile: make(map[string]*spaceFile),
}

//Variable Global deferSpace: Estos archivos son semipermanentes en el disco y se cierran
//cuando se sobrepasa cierto limite de archivos.
var deferSpace = &spaceDeferDisk{
	deferFile: make(map[string]*spaceFile),
	info:      make([]deferFileInfo, 0),
}

//Variable Global permSpace: Estos archivos permanecen abiertos durante toda la aplicación.
var permSpace = &spacePermDisk{
	permDisk: make(map[string]*spaceFile),
}

type spaceFile struct {
	*space
	file *os.File
	url  string
	//Numero de lineas ATOMICO de un archivo
	sizeFileLine *int64
	//Mutex que sirve para leer y escribir , actualizar mapas , actualizar arrays
	sync.RWMutex
}

type spaceDisk struct {
	diskFile map[string]*spaceFile
	sync.RWMutex
}

type deferFileInfo struct {
	name string
	time time.Time
}

type spaceDeferDisk struct {
	deferFile map[string]*spaceFile
	info      []deferFileInfo
	sync.RWMutex
}

type spacePermDisk struct {
	permDisk map[string]*spaceFile
	sync.RWMutex
}






/**********************************************************************************************/
/* 
spaceErrors -> Gestion de errores de la aplicación, test de velocidad y test de memoria.

lDAC -> Inicio de la aplicacion global e instancia del DAC central, compartido con el 
resto de espacios, equivalente a una base de datos mysql

space -> Genera un espacio que seria igual a una tabla mysql, se divide en dos tipos de campos
fields campos unicos y columnas campos dinamicos, los espacios son compatibles y recomendados
 para todo tipo de datos tablas de nombres, emails, video, audio, html, jss, css...
 
spaceFile -> Apunta a un archivo y se guardan actualmente en tres diferentes mapas donde
permanecen abiertos hasta que un timer los cierra, funcionan con locs de lectura y escritura.

RBuffer -> Inicia el buffer de lectura para leer un archivo, compatible con un solo campo,
una sola linea, multiples campos, multiples lineas, (multiples campos y multiples lineas).
Ademas compatible con canales para enviar datos segun se van leyendo (Costo de memoria minimo).
+Ademas compatible con lectura por rangos en campos (Costo de memoria el minimo que tu quieras),
perfecto para videos.
*/
/**********************************************************************************************/

//canal del buffer de lectura.
type RChanBuf struct {
	Line    int64
	ColName string
	Buffer  []byte
}

//Buffer de lectura compatible con multiples lineas y fieldas.
//El buffer solo puede leer una columna o fields tampoco es compatible con multples lineas.
//Buffer map es compatible con multiples lineas, fields y columnas
//El canal de buffer es compatible con mutiples lineas fields y columnas.
type rLines struct {
	startLine int64
	endLine   int64
}

type rRangues struct {
	rangue      int64
	totalRangue int64
	rangeBytes  int64
}

type RBuffer struct {
	*spaceFile
	typeBuff   fileTypeBuffer
	postFormat bool
	colName    *[]string

	*rLines
	*rRangues

	FieldBuffer *[]byte
	Buffer      *[]byte
	BufferMap   map[string][][]byte
	Channel     chan RChanBuf
}





/**********************************************************************************************/
/*
spaceErrors -> Gestion de errores de la aplicación, test de velocidad y test de memoria.

lDAC -> Inicio de la aplicacion global e instancia del DAC central, compartido con el 
resto de espacios, equivalente a una base de datos mysql

space -> Genera un espacio que seria igual a una tabla mysql, se divide en dos tipos de campos
fields campos unicos y columnas campos dinamicos, los espacios son compatibles y recomendados
 para todo tipo de datos tablas de nombres, emails, video, audio, html, jss, css...
 
spaceFile -> Apunta a un archivo y se guardan actualmente en tres diferentes mapas donde
permanecen abiertos hasta que un timer los cierra, funcionan con locs de lectura y escritura.

WBuffer -> Inicia un buffer de escritura para escribir en un archivo compatible con un solo campo,
una sola linea, multiples campos, multiples columnas, (multiples campos y multiples columnas),
(No compatible con multilinea).
Ademas escritura en tiempo real gracias a los canales, donde puedes escribir en un canal que 
no se cierra nunca.
Ademas escritura por rangos en fields perfecto para guardar un video , un audio o archivos grandes 
por rangos.
*/
/**********************************************************************************************/


type wLines struct{
	line  int64
}

type wRangues struct {
	rangue      int64 
	rangeBytes  int64
}

//El canal de escritura linea , nombre de columna y el buffer.
type wChanBuf struct{
	*wLines
	*wRangues
	colName string
	buffer 	[]byte
}

//Buffer de escritura con tres tipos de buffer
//Tipo buffer unicamente puede escribir en una columna o un field.
//Tipo mapaBuffer puede escribir simultaneametne en columnas y fields.
//Abre un canal que puede actualizar tantno columnas como fields.
type WBuffer struct {
	*spaceFile
	ColumnName string
	typeBuff fileTypeBuffer
	preFormat bool
	*wLines
	*wRangues

	buffer *[]byte
	bufferMap map[string][]byte
	channel chan wChanBuf
}