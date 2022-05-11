/*
En el fichero space se pueden ver todas las estructuras y constantes utilizadas en DAC.

El flujo de datos en resumen es este
spaceErrors ->  *pointer

	lDAC ->  *pointer

		space ->  *pointer

			spaceFile ->  *pointer

				WBuffer -> *pointer
				RBuffer -> *pointer

Este flujo optimiza la memoria, y la reutilización de en los siguientes espacios, en caso de
tener que actualizar un espacio es inevitable crear un nuevo objeto, que cree dicho espacio
esto es debido a que el flujo de datos ya esta optimizado.

Ejemplo teorico:
Si quisieras un nuevo space reutilizarias spaceErrors y lDAC.
Si quisieras un nuevo spaceFile reutilizarias spaceErrors , lDAC y space.
Si quisieras un nuevo WBuffer o RBuffer reutilizarias spaceErrors , lDAC, space y spaceFile.

Ejemplo practico:
1 -> spaceErrors, 1 -> lDAC, -> 1 space -> 50 spaceFile -> 100  WBuffer
														   1000 rBuffer
En este caso estamos usando 1 space que genera 50 spaceFile que generan 100 WBuufer y
1000 rBuffer cada uno.

Otro ejemplo practico:
1 -> spaceErrors, 1 -> lDAC, -> 100 space -> 50 spaceFile -> 100  WBuffer
														     1000 rBuffer
En este caso estamos usando 100 space que generan 50 spaceFile cada uno que generan a su vez 100 WBuufer y
1000 rBuffer cada uno.

Como test he creado 1000 000 de archivos que permanecen abiertos en un
linux ubuntu y equivalen a 700 mb de ram. Todo son punteros a archivos que van a
permanecer abiertos disminuyendo muchisimo los tiempos de escritura y lectura.

El siguiente paso es comprender como funcionan los archivos hay tres tipos por ahora,
disk deferDisk y permDisk.

Los archivos disk se cierran cada x tiempo.
Los archivos deferDisk permanecen en memoria hasta que se supera un limite
de archivos abiertos.
Los permDisk permanecen siempre abiertos como los logs de la aplicación.

¿Como funciona DAC?

En resumen DAC almacena los datos y lee los datos muy rapido a coste de desperdiciar memoria en
disco.

¿Por qué desperdiciar esa memoria?
Lo que perdemos en memoria lo ganamos en procesamiento, hasta ahora el coste de procesamiento
siempre ha sido mas caro que el coste en memoria.

¿Qué ocurre si se invierte ese coste?
No importa DAC optimiza el tiempo de ejecucion al maximo simplemente comprime los archivos
con un algoritmo de compresion exacto la mayoria estan llenos de bytes nulos.

¿Cómo funciona exacatamente?

Los archivos se dividen en fields y lineas, a su vez las lineas se dividen en columnas.

Los fields ocupan un espacio fijo al inicio del archivo.

Las lineas ocupan un espacio fijo segun se van creando tengan datos utiles o no.

Ejemplo:
Cuando un usuario te pide datos como un email ese usuario tendra una id que es el numero de linea
 o el id suyo en la aplicación y apunta a un numero de linea, leerias los datos directamente
 de esa linea para ese usuario.

 Usuario: Franky, Id: 101,
 Guardar email pues escribiria en el archivo de emails en la linea 101.
 Leer email pues leeria en el archivo email en la linea 101.

¿Porqué complicarse la vida frente a sistema de bases de datos?

La primera razon DAC es divertido, resolveras problemas que nuncan pensaste, valoraras las
soluciones de bases de datos que lo resuelven todo.

La segunda razon es que actualmente es la solucion mas rapida de lectura y escritura de
almacenamiento de datos, no lo he medido, la razon es que no hay otra forma de hacerlo
mas rapido mediante el software teoricamente.

Ademas optimiza el uso de la memoria para todo tipo de archivos video, audio, fotos, html
jss gracias a la escritura por rangos en fields, no como otras soluciones de bases de datos
que no estan optimizadas para guardar archivos grandes.

La tercera es muy seguro, no hay inyeccion sql porque no es un lenguaje sql.

La cuarta optimizado para la ejecuion en paralelo y la concurrencia, puedes escribir todos los
datos con go write si quieres y funciona.

La quinta y la mejor almacenamiento de datos NATIVO en golang, ya sabes amo Golang y lo
quiero hacer todo en golang y en cuanto menos herramientas de tercero no optimizadas para
golang mejor.

La sexta muy facil de usar, poco a poco ire añadiendo porcelana de las principales operaciones
que se harian en lenguaje sql en DAC.

Diviertete y no utilices DAC para un proyecto serio ya que le quedan años de desarrollo todavia.

Todas las funciones que se pueden utlizar estan en los archivos que empiezan por public.

*/

package dac

import (

	"os"
	"sync"
	"time"
)

/**********************************************************************************************/
/* Constantes */
/**********************************************************************************************/

const(welcome string = `Bienvend@ a DAC, Errores: ACTIVADOS, Documentacion: https://github.com/FrankyGolang/AccessDataGolang`)

type fileNativeType int64
const (
	disk fileNativeType = 1 << iota
	deferDisk
	permDisk
	openFile
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
	//Extensiones de contenido web
	html = "html"
	json = "json"
	js   = "js"
	css  = "css"

	//Extensiones de contenido imagen
	glb  = "glb"
	gif  = "gif"
	svg  = "svg"
	png  = "png"
	jpg  = "jpg"
	webp = "webp"
	bmp  = "bmp"
	
	//Extensiones de audio
	mp3 = "mp3"

	//Extensiones de contenido video
	mp4 = "mp4"

	//Extensiones de contenido de documentos
	pdf = "pdf"
	txt = "txt"
)

var extensionFile = map[string]string{
	dacByte: "Archivo que incluye fields y columnas de bytes",
	dacBit:  "Archivo que incluye fields y una lista de bits",

	html:  "text/html; charset=UTF-8",
	json : "application/json; charset=utf-8",
	js   : "text/javascript; charset=UTF-8",
	css  : "text/css; charset=UTF-8",

	//Extensiones de contenido imagen
	glb  : "model/gltf-binary",
	gif  : "image/gif",
	svg  : "image/svg+xml",
	png  : "image/png",
	jpg  : "image/jpeg",
	webp : "image/webp",
	bmp :  "image/bmp",

	//Extensiones de audio
	mp3 : "audio/mpeg",
	
	//Extensiones de contenido video
	mp4 :  "video/mp4",

	//Extensiones de contenido de documentos
	pdf : "application/pdf",
	txt : "text/plain; charset=UTF-8",
}

var allowedRange = map[string]bool{
	//Extensiones de audio
	mp3 : true,
//Extensiones de contenido video
	mp4 :  true,
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


	//Activar log de memoria
	logMemoryUse bool

	//Activar fichero de memoria.
	//Predeterminado: DAC/memory
	logFileMemoryUse bool

	//Separador para los resultados en nanosegundos: coma, espacio , punto...
	separatorLog string

	//Niveles de la url DAC a mostrar.
	levelsUrl int

	//Activar cronometros de apertura de archivos
	logTimeOpenFile bool
	//Cronometro de lectura de archivos
	logTimeReadFile bool
	//Cronometro de lectura de archivos
	logTimeWriteFile bool
}

//space globales
var timeLog *PublicSpace
var memoryLog *PublicSpace
var errorLog *PublicSpace





/**********************************************************************************************/
/* 
spaceErrors -> *pointer

lDAC -> Inicio de la aplicacion global e instancia del DAC central, compartido con el 
resto de espacios, equivalente a una base de datos mysql
*/
/**********************************************************************************************/

type timersFile struct {

	fileOpenDeferFile int
	timeEventDeferFile int64
	timeEventDiskFile int64

}
//Inicio del path DAC
type lDAC struct {
	*spaceErrors
	*timersFile
	newDACFolder    bool
	goDACFolder     bool
	globalDACFolder string
}

//Variable global de DAC
var globalDac *lDAC





/**********************************************************************************************/
/*
spaceErrors ->  *pointer

lDAC ->  *pointer

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
	isErrorFile bool
	//Indica el estado del archivo en la aplicacion
	fileNativeType fileNativeType
	//Propiedades comunes a todos los archivos
	dir string

	//Cambiar Extension
	extension string

	sizeLine int64
	sizeField int64

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

type PublicSpace struct  {
	*space
}



var Space map[string]*space



/**********************************************************************************************/
/* 
spaceErrors ->  *pointer

lDAC ->  *pointer

space ->  *pointer

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



//Archivos
type PublicSpaceFile struct  {
	*spaceFile
	*PublicSpaceCache
}

//Rutas de archivo
type PublicSpaceCache struct {
	cache map[string]*PublicSpaceFile
	sync.RWMutex
}

//Directorios
type spaceGlobalCache struct {
	cache map[string]*PublicSpaceCache
	sync.RWMutex
}

//Iniciando variable global de cache de rutas de directorios
var globalCache = &spaceGlobalCache{
	cache: make(map[string]*PublicSpaceCache),
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
space 
spaceFile 
spaceErrors -> pointer
 Compatible con space y spaceFile, hereda url del dac space o spaceFile
*/
/**********************************************************************************************/

type errorsDac struct  {
	*spaceErrors
	fileName string
	fileFolder []string

	typeError errorDac
	url string
	messageLog string

	levelsUrl int
	separatorLog string
	//Cuantas funciones ascienden para localizar la funcion de llamada
	runCaller int

	timeNow *time.Time
}
var EDAC bool

/**********************************************************************************************/
/* 
spaceErrors ->  *pointer

lDAC ->  *pointer

space ->  *pointer
 
spaceFile ->  *pointer

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
spaceErrors ->  *pointer

lDAC ->  *pointer

space ->  *pointer
 
spaceFile ->  *pointer

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
	columnName string
	typeBuff fileTypeBuffer
	preFormat bool
	*wLines
	*wRangues

	buffer *[]byte
	bufferMap map[string][]byte
	channel chan wChanBuf
}