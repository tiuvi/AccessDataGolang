package dac

import (
	"log"
	"runtime"
	"strconv"
	"strings"
	"time"
)

func (LDAC *lDAC) onErrorsLog() {

	//errorLog
	space := globalDac.NewSpace()
	space.OnErrorFile()
	space.NewTimeFileDisk()
	space.NewDacByte()
	space.SetSubDir("errors", "dac", "errors")
	space.NewField("exceptionCount", 64)
	space.NewField("warningCount", 64)
	space.NewField("fatalCount", 64)

	space.NewColumnByte("date", 20)
	space.NewColumnByte("typeError", 20)
	space.NewColumnByte("fileName", 20)
	space.NewColumnByte("funcion", 35)
	space.NewColumnByte("line", 6)
	space.NewColumnByte("endLine1", 1)

	space.NewColumnByte("message", 64)
	space.NewColumnByte("endLine2", 1)

	space.NewColumnByte("DAC", 64)
	space.NewColumnByte("endLine3", 2)
	space.OSpaceInit()
	errorLog = space.SetPublicSpace()

	//memoryLog
	space = globalDac.NewSpace()
	space.OnErrorFile()
	space.NewTimeFileDisk()
	space.NewDacByte()
	space.SetSubDir("errors", "dac", "memory")

	space.NewColumnByte("date", 20)
	space.NewColumnByte("fileName", 20)
	space.NewColumnByte("funcion", 35)
	space.NewColumnByte("line", 6)
	space.NewColumnByte("endLine1", 1)

	space.NewColumnByte("DAC", 64)
	space.NewColumnByte("endLine2", 1)

	space.NewColumnByte("Alloc", 20)
	space.NewColumnByte("totalAlloc", 20)
	space.NewColumnByte("memSys", 20)
	space.NewColumnByte("objLargerMemory", 20)
	space.NewColumnByte("countGC", 20)
	space.NewColumnByte("endLine3", 2)
	space.OSpaceInit()
	memoryLog = space.SetPublicSpace()

	//memoryLog
	space = globalDac.NewSpace()
	space.OnErrorFile()
	space.NewTimeFileDisk()
	space.NewDacByte()
	space.SetSubDir("errors", "dac", "timers")

	space.NewColumnByte("date", 20)
	space.NewColumnByte("fileName", 20)
	space.NewColumnByte("funcion", 35)
	space.NewColumnByte("line", 6)
	space.NewColumnByte("endLine1", 1)

	space.NewColumnByte("DAC", 64)
	space.NewColumnByte("nanosecond", 15)
	space.NewColumnByte("endLine2", 2)

	space.OSpaceInit()
	timeLog = space.SetPublicSpace()

}

func BoldConsole(str string) string {

	return strings.Join([]string{"\u001b[01m", ColorWhite, "\u001b[40m", str, Reset}, "")
}

func (errorDac errorDac )PrintConsole()string{

	var color string
	switch errorDac {

		case Fatal, MessageCopilation:
			color = ColorRed
		case Exception:
			color = ColorGreen
		case Warning:
			color = ColorYellow
		case Message, NewRouteFolder, TimeMemory:
			color = ColorBlue
	}
	return strings.Join([]string{"\u001b[01m", color,  string(errorDac) , Reset}, "")
}

func Uint64ToStringSep(uintData uint64, sep string) string {

	return counterJoin(strconv.FormatUint(uintData, 10), sep)
}

func counterJoin(stringJoin string, sep string) (stringJoinFormat string) {

	divuintStr := len(stringJoin) / 3
	restuintStr := len(stringJoin) % 3

	if restuintStr > 0 {

		stringJoinFormat = strings.Join([]string{stringJoin[:restuintStr], sep}, "")
		stringJoin = stringJoin[restuintStr:]

	}

	for x := 0; x < divuintStr; x++ {

		stringJoinFormat += strings.Join([]string{stringJoin[:3], sep}, "")
		stringJoin = stringJoin[3:]

	}

	return stringJoinFormat
}



func CallerSystemInfo(urlDac string , numCaller int ,levelsUrl int)(dateNow string , urlFile string ,nameFunction string, lineFunction string,urlDacFormat string , ok bool) {

	/*
	Por si falla runtime.callers

	count := runtime.Callers(0, make([]uintptr,15))

	if count == 9 {
		EDAC.runCaller = EDAC.runCaller - 1
	}
	if count == 10 {
		EDAC.runCaller = EDAC.runCaller - 1
	}
	if count == 11 {
		EDAC.runCaller = EDAC.runCaller - 1
	}

	log.Println("contador: ", count)
*/

	//Inicio de los datos de Archivo, funcion, linea y url del dac desde donde se llama.
	ptr, _, _, ok := runtime.Caller(numCaller)
	if !ok {

		return "","","","","",ok

	}

	dateNow     = time.Now().Format("02/01/2006 15:04:05")

	firstFrame := runtime.CallersFrames([]uintptr{ptr})
	frame, _ := firstFrame.Next()




	urlFile = frame.File
	uFileFormat := strings.Split(urlFile , "/")
	if len(uFileFormat) > 2 {

		urlFile = strings.Join(uFileFormat[len(uFileFormat)-2:], "/")

	}

	nameFunction = frame.Function
	fNameFormat := strings.Split(nameFunction , "/")
	if len(fNameFormat) > 1 {

		nameFunction = strings.Join(fNameFormat[len(fNameFormat)-1:], "/")
	}


	lineFunction  = strconv.Itoa(frame.Line)


	urlDacFormat = urlDac
	uDacFormat  := strings.Split(urlDac, "/")
	if len(uDacFormat) > levelsUrl {

		urlDacFormat = strings.Join(uDacFormat[len(uDacFormat)-levelsUrl:], "/")
	}

	return
}

func LogMemory(sep string) (allocString string, totalAlloc string, memSys string, bigObjSize string, gcCount string) {

	var memoryStats runtime.MemStats

	runtime.ReadMemStats(&memoryStats)
	allocString = Uint64ToStringSep(memoryStats.Alloc, sep)
	totalAlloc = Uint64ToStringSep(memoryStats.TotalAlloc, sep)
	memSys = Uint64ToStringSep(memoryStats.Sys, sep)
	bigObjSize = Uint64ToStringSep(uint64(memoryStats.BySize[60].Size), sep)
	gcCount = Uint64ToStringSep(uint64(memoryStats.NumGC), sep)

	return
}

func (EDAC *errorsDac) logNewError() {

	var UTSS func(uint64, string) string = Uint64ToStringSep
	/*
	* Variables para iniciar tanto el log como el archivo.
	*
	* Fin del tiempo cronometrado
	* Nombre de archivo, funcion, linea y url dac
	* Leer memoria actual en uso
	*
	 */

	//Finalizacion del log de time
	var elapsed int64
	if (EDAC.logFileTimeUse || EDAC.logTimeUse) && EDAC.timeNow != nil {

		elapsed = time.Since(*EDAC.timeNow).Nanoseconds()

	}



	dateNow , urlFile ,nameFunction,lineFunction,urlDacFormat,ok := CallerSystemInfo(EDAC.url , EDAC.runCaller ,EDAC.levelsUrl)
	if !ok {

		return
	}

	EDAC.runCaller = EDAC.runCaller + 1


	//Inicio de la memoria.

	var allocString string
	var totalAlloc string
	var memSys string
	var bigObjSize string
	var gcCount string
	if (EDAC.logMemoryUse || EDAC.logFileMemoryUse) && EDAC.timeNow != nil {

		allocString, totalAlloc, memSys, bigObjSize, gcCount = LogMemory(EDAC.separatorLog)

	}

	/*
	* Funciones de log de errores asincronas excepto el log fatal.
	* Log Fatal
	* Log consola
	* Log memory
	* Log time
	 */

	var logString string
	var nT func(string) string = BoldConsole

	if (EDAC.logFatalErrors || EDAC.logConsoleErrors) && EDAC.timeNow == nil {


		logString = strings.Join([]string{
			nT("Fecha:"), " ", dateNow, " - ", nT("Tipo de error:"), " ", EDAC.typeError.PrintConsole() , "\n\r",
			nT("Archivo:"), " ", urlFile, " - ", nT("Funcion:"), " ", nameFunction, " - ", nT("Nº Linea:"), " ", lineFunction, "\n\r",
			nT("Ruta DAC:"), " ", urlDacFormat, "\n\r",
			nT("Mensaje:"), " ", EDAC.messageLog, "\n\r"}, "")

		if EDAC.logFatalErrors && (EDAC.typeError == Fatal || EDAC.typeError == MessageCopilation) {

			log.Fatalln(logString)

		}

		if EDAC.logConsoleErrors {

			go log.Println(logString)

		}
	}

	// For info on each, see: https://golang.org/pkg/runtime/#MemStats
	if EDAC.logMemoryUse && EDAC.timeNow != nil {

		logString := strings.Join([]string{
			nT("Fecha:"), " ", dateNow, " - ", nT("Tipo de error:"), " ", EDAC.typeError.PrintConsole(), "\n\r",
			nT("Archivo:"), " ", urlFile, " - ", nT("Funcion:"), " ", nameFunction, " - ", nT("Nº Linea:"), " ", lineFunction, "\n\r",
			nT("Ruta DAC:"), " ", urlDacFormat, "\n\r",
			nT("Memoria en uso:"), " ", allocString, " Bytes, - ",
			nT("Memoria total usada:"), " ", totalAlloc, " Bytes, - ",
			nT("Memoria sistema:"), " ", memSys, " Bytes, - ",
			nT("Objeto de mayor tamaño:"), " ", bigObjSize, " Bytes, - ",
			nT("Contador GC:"), " ", gcCount, " veces.\n\r"}, "")

		go log.Println(logString)

	}

	if EDAC.logTimeUse && EDAC.timeNow != nil {

		logString := strings.Join([]string{
			nT("Fecha:"), " ", dateNow, " - ", nT("Tipo de error:"), " ", EDAC.typeError.PrintConsole(), "\n\r",
			nT("Archivo:"), " ", urlFile, " - ", nT("Funcion:"), " ", nameFunction, " - ", nT("Nº Linea:"), " ", lineFunction, "\n\r",
			nT("Ruta DAC:"), " ", urlDacFormat, "\n\r",
			nT("Tiempo transcurrido:"), " ", UTSS(uint64(elapsed), EDAC.separatorLog), " NanoSegundos.\n\r"}, "")
		go log.Println(logString)
	}

	/*
	* Funciones de escritura de archivos asincronas
	* - writeNewError
	* - writeNewTimeUse
	* - writeNewMemoryUse
	 */

	//Log de errores
	if EDAC.logFileError && EDAC.typeError != MessageCopilation && EDAC.timeNow == nil {

		go EDAC.writeNewError(dateNow, EDAC.typeError, urlFile, nameFunction, lineFunction, EDAC.messageLog, urlDacFormat)

	}

	//Log de tiempo transcurrido y archivo de tiempo transcurrido
	if EDAC.logFileTimeUse && EDAC.typeError != MessageCopilation && EDAC.timeNow != nil {

		go EDAC.writeNewTimeUse(dateNow, urlFile, nameFunction, lineFunction, urlDacFormat, UTSS(uint64(elapsed), EDAC.separatorLog))

	}

	//Log de memoria y archivo de memoria
	if EDAC.logFileMemoryUse && EDAC.typeError != MessageCopilation && EDAC.timeNow != nil {

		go EDAC.writeNewMemoryUse(dateNow, urlFile, nameFunction, lineFunction, urlDacFormat,
			allocString, totalAlloc, memSys, bigObjSize, gcCount)

	}
}

func (EDAC *errorsDac) writeNewTimeUse(date string, fileName string, funcion string, line string, url string, nanosecond string) {

	var bufferW *WBuffer
	var file *spaceFile

	file = timeLog.OSpace(EDAC.fileName+"Time", EDAC.fileFolder...)

	if file != nil {

		bufferW = file.NewWriterMapBytes()

	}

	if bufferW != nil {

		bufferW.NewLineWBspace()
		dateByte := []byte(date)
		fileNameByte := []byte(fileName)
		funcionByte := []byte(funcion)
		lineByte := []byte(line)
		urlByte := []byte(url)
		nanosecondByte := []byte(nanosecond)
		endLineByte1 := []byte("\n")
		endLineByte2 := []byte("\n\r")

		bufferW.SendBWspace("date", &dateByte)
		bufferW.SendBWspace("fileName", &fileNameByte)
		bufferW.SendBWspace("funcion", &funcionByte)
		bufferW.SendBWspace("line", &lineByte)
		bufferW.SendBWspace("endLine1", &endLineByte1)

		bufferW.SendBWspace("DAC", &urlByte)
		bufferW.SendBWspace("nanosecond", &nanosecondByte)
		bufferW.SendBWspace("endLine2", &endLineByte2)
		bufferW.Wspace()

	}

}

func (EDAC *errorsDac) writeNewMemoryUse(date string, fileName string, funcion string, line string, url string,
	alloc string, totalAlloc string, memSys string, bigObjSize string, gcCount string) {

	var bufferW *WBuffer
	var file *spaceFile

	file = memoryLog.OSpace(EDAC.fileName+"Memory", EDAC.fileFolder...)

	if file != nil {

		bufferW = file.NewWriterMapBytes()

	}

	if bufferW != nil {

		bufferW.NewLineWBspace()
		dateByte := []byte(date)
		fileNameByte := []byte(fileName)
		funcionByte := []byte(funcion)
		lineByte := []byte(line)
		urlByte := []byte(url)

		allocByte := []byte(alloc)
		totalAllocByte := []byte(totalAlloc)
		memSysByte := []byte(memSys)
		bigObjSizeByte := []byte(bigObjSize)
		gcCountByte := []byte(gcCount)

		endLineByte1 := []byte("\n")
		endLineByte2 := []byte("\n\r")

		bufferW.SendBWspace("date", &dateByte)
		bufferW.SendBWspace("fileName", &fileNameByte)
		bufferW.SendBWspace("funcion", &funcionByte)
		bufferW.SendBWspace("line", &lineByte)
		bufferW.SendBWspace("endLine1", &endLineByte1)

		bufferW.SendBWspace("DAC", &urlByte)
		bufferW.SendBWspace("endLine2", &endLineByte1)

		bufferW.SendBWspace("Alloc", &allocByte)
		bufferW.SendBWspace("totalAlloc", &totalAllocByte)
		bufferW.SendBWspace("memSys", &memSysByte)
		bufferW.SendBWspace("objLargerMemory", &bigObjSizeByte)
		bufferW.SendBWspace("countGC", &gcCountByte)
		bufferW.SendBWspace("endLine3", &endLineByte2)

		bufferW.Wspace()

	}

}

func (EDAC *errorsDac) writeNewError(date string, typeError errorDac, fileName string, funcion string, line string, messageLog string, url string) {

	var bufferW *WBuffer
	var file *spaceFile

	file = errorLog.OSpace(EDAC.fileName+"Errors", EDAC.fileFolder...)

	if file != nil {

		bufferW = file.NewWriterMapBytes()

	}

	if bufferW != nil {

		bufferW.NewLineWBspace()
		dateByte := []byte(date)
		typeErrorByte := []byte(typeError)
		fileNameByte := []byte(fileName)
		funcionByte := []byte(funcion)
		lineByte := []byte(line)
		messageLogByte := []byte(messageLog)
		urlByte := []byte(url)

		endLineByte1 := []byte("\n")
		endLineByte2 := []byte("\n\r")

		bufferW.SendBWspace("date", &dateByte)
		bufferW.SendBWspace("typeError", &typeErrorByte)
		bufferW.SendBWspace("fileName", &fileNameByte)
		bufferW.SendBWspace("funcion", &funcionByte)
		bufferW.SendBWspace("line", &lineByte)
		bufferW.SendBWspace("endLine1", &endLineByte1)
		bufferW.SendBWspace("message", &messageLogByte)
		bufferW.SendBWspace("endLine2", &endLineByte1)
		bufferW.SendBWspace("DAC", &urlByte)
		bufferW.SendBWspace("endLine3", &endLineByte2)
		bufferW.Wspace()

	}
}
