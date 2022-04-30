package bd

import (
	"log"
	"runtime"
	"strconv"
	"strings"
	"time"
)




/*
var GlobalError = &spaceErrors{
	logFatalErrors:   false,
	logConsoleErrors: true,

	logFileError: true,

	logTimeUse :   true,   
	logFileTimeUse: true,
	logTimeOpenFile: true,


	logMemoryUse:  true,
	logFileMemoryUse: true,

	separatorLog: " ",
	levelsUrl: 8,
}
*/

/*
var errorLog = &space{
	FileNativeType: PermDisk,
	Dir: autoload_options.Patch + "/errors/dac/",
	extension: DacByte,
	indexSizeFieldsArray: []spaceLen{
		{"exceptionCount", 64,},
		{"warningCount", 64,},
		{"fatalCount", 64,},
	},
	indexSizeColumnsArray: []spaceLen{
		{"date",      20,},
		{"typeError", 20,},
		{"fileName",  20,},
		{"funcion",   35,},
		{"line",       6,},
		{"endLine1",   1,},

		{"message",   64,},
		{"endLine2",   1,},

		{"DAC", 64,},
		{"endLine3",   2,},
	},
}

var memoryLog = &space{
	FileNativeType: PermDisk,
	Dir: autoload_options.Patch + "/errors/dac/",
	extension: DacByte,
	indexSizeColumnsArray: []spaceLen{
		{"date",      20,},
		{"fileName",  20,},
		{"funcion",   35,},
		{"line",       6,},
		{"endLine1",   1,},

		{"DAC", 64,},
		{"endLine2",   1,},

		{"Alloc",   20,},
		{"totalAlloc", 20,},
		{"memSys", 20,},
		{"objLargerMemory", 20,},
		{"countGC", 20,},
		{"endLine3",   2,},
	},
}

var timeLog = &space{
	FileNativeType: PermDisk,
	Dir: autoload_options.Patch + "/errors/dac/",
	extension: DacByte,
	indexSizeColumnsArray: []spaceLen{
		{"date",      20,},
		{"fileName",  20,},
		{"funcion",   35,},
		{"line",       6,},
		{"endLine1",   1,},
		
		{"DAC", 64,},
		{"nanosecond",   15,},
		{"endLine2",   2,},
	},
}
*/

func negritaTerminal(str string)string{
	
	return strings.Join( []string{"\u001b[01m" ,ColorWhite, "\u001b[40m",   str , Reset} , "")
}

func (EDAC *ErrorsDac ) uint64ToString(uintData uint64)string {

	

	var uintFormat string
	uintStr := strconv.FormatUint( uintData, 10) 

	divuintStr  := len(uintStr) / 3 
	restuintStr := len(uintStr) % 3 

	if restuintStr > 0 {

		uintFormat = strings.Join([]string{ uintStr[:restuintStr] , EDAC.separatorLog } , "")
		uintStr = uintStr[restuintStr:]
	
	}

	for x := 0; x < divuintStr ; x++{
		
		uintFormat += strings.Join([]string{ uintStr[:3] , EDAC.separatorLog } , "")
		uintStr = uintStr[3:]

	}

	return uintFormat
}

func errorsLog(){

//	errorLog.OSpaceInit()
//	timeLog.OSpaceInit()
//	memoryLog.OSpaceInit()


}

type ErrorsDac struct  {
	*spaceErrors
	fileName string
	typeError errorDac
	Url string
	messageLog string

	levelsUrl int
	separatorLog string

	timeNow *time.Time
}



func (sP *space )ErrorSpaceDefault(typeError errorDac, messageLog string){

	EDAC := &ErrorsDac {
		spaceErrors: sP.spaceErrors,
		fileName: "",
		typeError: typeError,
		Url: sP.dir,
		messageLog: messageLog,
		levelsUrl:    sP.levelsUrl ,
		separatorLog:  sP.separatorLog,
		timeNow: nil,
	}
	EDAC.LogNewError()
}

func (sP *space ) NewErrorSpace(fileName string, typeError errorDac, messageLog string){

	EDAC := &ErrorsDac {
		spaceErrors: sP.spaceErrors,
		fileName: fileName,
		typeError: typeError,
		Url: sP.dir,
		messageLog: messageLog,
		levelsUrl:    sP.levelsUrl ,
		separatorLog:  sP.separatorLog,
		timeNow: nil,
	}
	EDAC.LogNewError()
}

func (sF *spaceFile ) ErrorSpaceFileDefault(typeError errorDac, messageLog string){

	EDAC := &ErrorsDac {
		spaceErrors: sF.spaceErrors,
		fileName: "",
		typeError: typeError,
		Url: sF.url,
		messageLog: messageLog,
		levelsUrl:    sF.levelsUrl ,
		separatorLog:  sF.separatorLog,
		timeNow: nil,
	}
	EDAC.LogNewError()
}

func (sF *spaceFile ) NewErrorSpace(fileName string, typeError errorDac, messageLog string){

	EDAC := &ErrorsDac {
		spaceErrors: sF.spaceErrors,
		fileName: fileName,
		typeError: typeError,
		Url: sF.url,
		messageLog: messageLog,
		levelsUrl:    sF.levelsUrl ,
		separatorLog:  sF.separatorLog,
		timeNow: nil,
	}
	EDAC.LogNewError()
}

func (sP *space ) LogDeferTimeMemoryDefault(timeNow time.Time){

	EDAC := &ErrorsDac {
		spaceErrors: sP.spaceErrors,
		fileName: "",
		typeError: TimeMemory,
		Url: sP.dir,
		messageLog: "",
		levelsUrl:    sP.levelsUrl ,
		separatorLog:  sP.separatorLog,
		timeNow: &timeNow,
	}
	EDAC.LogNewError()

}

func (sP *space ) NewLogDeferTimeMemory(fileName string, timeNow time.Time){

	EDAC := &ErrorsDac {
		spaceErrors: sP.spaceErrors,
		fileName: fileName,
		typeError: TimeMemory,
		Url: sP.dir,
		messageLog: "",
		levelsUrl:    sP.levelsUrl ,
		separatorLog:  sP.separatorLog,
		timeNow: &timeNow,
	}
	EDAC.LogNewError()

}

func (sF *spaceFile ) LogDeferTimeMemorySF(timeNow time.Time){

	EDAC := &ErrorsDac {
		spaceErrors: sF.spaceErrors,
		fileName: "",
		typeError: TimeMemory,
		Url: sF.url,
		messageLog: "",
		levelsUrl:    sF.levelsUrl ,
		separatorLog:  sF.separatorLog,
		timeNow: &timeNow,
	}
	EDAC.LogNewError()

}



func (sF *spaceFile ) NewLogDeferTimeMemorySF(fileName string, timeNow time.Time){

	EDAC := &ErrorsDac {
		spaceErrors: sF.spaceErrors,
		fileName: fileName,
		typeError: TimeMemory,
		Url: sF.url,
		messageLog: "",
		levelsUrl:    sF.levelsUrl ,
		separatorLog:  sF.separatorLog,
		timeNow: &timeNow,
	}
	EDAC.LogNewError()

}



func (EDAC *ErrorsDac ) LogNewError(){

	date := time.Now()
	dateString := date.Format("02/01/2006 15:04:05")


	ptr , _, _, _ := runtime.Caller(2)
	firstFrame := runtime.CallersFrames([]uintptr{ptr})
	frame, _ := firstFrame.Next()
	
	fileString     := frame.File
	funcNameString := frame.Function
	urlNameString  := EDAC.Url
	lineStr        := strconv.Itoa(frame.Line)

	fileDir    := strings.Split(frame.File, "/")
	if len(fileDir) > 2 {

		fileString = strings.Join(fileDir[len(fileDir)-2:] , "/")

	}
	
	funcName       := strings.Split(frame.Function , "/")
	if len(funcName) > 1 {

		funcNameString = strings.Join(funcName[len(funcName)-1:]  , "/")
	}
	

	urlName       := strings.Split(EDAC.Url , "/")
	if len(urlName) > EDAC.levelsUrl {

		urlNameString = strings.Join(urlName[len(urlName)-EDAC.levelsUrl:] , "/")
	}

	if EDAC.logFileError && EDAC.typeError != MessageCopilation && EDAC.timeNow == nil{
	
		if len(EDAC.fileName) == 0{

			go WriteNewError("message",dateString, EDAC.typeError ,fileString ,funcNameString , lineStr , EDAC.messageLog , urlNameString )

		}else{

			go WriteNewError(EDAC.fileName + "Message" ,dateString, EDAC.typeError ,fileString ,funcNameString , lineStr , EDAC.messageLog , urlNameString )

		}
	
	}


	var logString string
	var nT func(string)string = negritaTerminal

	if (EDAC.logFatalErrors || EDAC.logConsoleErrors) && EDAC.timeNow == nil {

		var color string
		switch EDAC.typeError {
	
		case Fatal,MessageCopilation:
			color = ColorRed
		case Exception:
			color = ColorGreen
		case Warning:
			color = ColorYellow
		case Message ,NewRouteFolder:
			color = ColorBlue
		}

		logString = strings.Join([]string{ 
		nT("Fecha:") ," ", dateString ," - ",nT("Tipo de error:"), " " ,color , string(EDAC.typeError) ,Reset, "\n\r", 
		nT("Archivo:"), " "  , fileString ," - ",nT("Funcion:"), " " , funcNameString ," - ",nT("Nº Linea:"), " "  , lineStr ,"\n\r",
		nT("Ruta DAC:"), " " , urlNameString, "\n\r",
		nT("Mensaje:"), " " , EDAC.messageLog, "\n\r" } , "")
	
		if EDAC.logFatalErrors  {

			log.Fatalln(logString )
	
		} 
	
		if EDAC.logConsoleErrors {
		
			go log.Println(logString)
		
		}
	}


	var memoryStats runtime.MemStats
	var allocString string
	var totalAlloc string
	var memSys string
	var bigObjSize string
	var gcCount string
	if (EDAC.logMemoryUse || EDAC.logFileMemoryUse) && EDAC.timeNow != nil {

		runtime.ReadMemStats(&memoryStats)
		allocString = EDAC.uint64ToString(memoryStats.Alloc)
		totalAlloc  = EDAC.uint64ToString(memoryStats.TotalAlloc)
		memSys      = EDAC.uint64ToString(memoryStats.Sys)
		bigObjSize  = EDAC.uint64ToString(uint64(memoryStats.BySize[60].Size))
		gcCount     = EDAC.uint64ToString(uint64(memoryStats.NumGC))
	}

	//Log de memoria y archivo de memoria
	if EDAC.logFileMemoryUse && EDAC.typeError != MessageCopilation && EDAC.timeNow != nil {

		if len(EDAC.fileName) == 0{

			go WriteNewMemoryUse("memory",dateString, fileString ,funcNameString , lineStr , urlNameString, 
			allocString, totalAlloc,memSys , bigObjSize, gcCount )
		}else{

			go WriteNewMemoryUse(EDAC.fileName + "memory",dateString, fileString ,funcNameString , lineStr , urlNameString, 
			allocString, totalAlloc,memSys , bigObjSize, gcCount )
		}
	
	}

	// For info on each, see: https://golang.org/pkg/runtime/#MemStats
	if EDAC.logMemoryUse && EDAC.timeNow != nil {


		logString := strings.Join([]string{ 
		nT("Fecha:")               ," ", dateString ," - ",nT("Tipo de error:"), " " , ColorBlue , "memoryStats" ,Reset, "\n\r", 
		nT("Archivo:"), " "  , fileString ," - ",nT("Funcion:"), " " , funcNameString ," - ",nT("Nº Linea:"), " "  , lineStr ,"\n\r",
		nT("Ruta DAC:"), " "          , urlNameString , "\n\r",
		nT("Memoria en uso:")         ," ", allocString ," Bytes, - ",
		nT("Memoria total usada:")    ," ", totalAlloc  ," Bytes, - ",
		nT("Memoria sistema:")        ," ", memSys      ," Bytes, - ",
		nT("Objeto de mayor tamaño:") ," ", bigObjSize  ," Bytes, - ",
		nT("Contador GC:")            ," ", gcCount     ," veces.\n\r" } , "")
		
		go log.Println(logString)

	}


	//Variables para iniciar tanto el log como el archivo.
	var elapsed int64
	if (EDAC.logFileTimeUse || EDAC.logTimeUse) && EDAC.timeNow != nil  {
		
		elapsed = time.Since(*EDAC.timeNow).Nanoseconds()

	}

	//Log de tiempo transcurrido y archivo de tiempo transcurrido
	if EDAC.logFileTimeUse && EDAC.typeError != MessageCopilation && EDAC.timeNow != nil {

		if len(EDAC.fileName) == 0{

			go WriteNewTimeUse("time",dateString, fileString ,funcNameString , lineStr , urlNameString, EDAC.uint64ToString(uint64(elapsed)) )

		}else{

			go WriteNewTimeUse(EDAC.fileName + "Time",dateString, fileString ,funcNameString , lineStr , urlNameString, EDAC.uint64ToString(uint64(elapsed)) )

		}
	}

	if EDAC.logTimeUse && EDAC.timeNow != nil {

		logString := strings.Join([]string{ 
			nT("Fecha:")               ," ", dateString ," - ",nT("Tipo de error:"), " " , ColorBlue , "timeStats" ,Reset, "\n\r", 
			nT("Archivo:"), " "  , fileString ," - ",nT("Funcion:"), " " , funcNameString ," - ",nT("Nº Linea:"), " "  , lineStr ,"\n\r",
			nT("Ruta DAC:"), " " , urlNameString, "\n\r",
			nT("Tiempo transcurrido:")     ," ", EDAC.uint64ToString(uint64(elapsed)) ," NanoSegundos.\n\r" } , "")
			go log.Println(logString)
	}
	
}



func WriteNewTimeUse(spaceName string, date string,fileName string,funcion string, line string, url string , nanosecond string) {

	var bufferW *WBuffer
	var file *spaceFile


	file = timeLog.OSpace(spaceName)


	if file != nil {

		bufferW = file.NewWBspace(buffMap)
		
	}

	if bufferW != nil {

		bufferW.NewLineWBspace()
		dateByte	    := []byte(date)
		fileNameByte	:= []byte(fileName)
		funcionByte  	:= []byte(funcion)
		lineByte  	    := []byte(line)
		urlByte         := []byte(url) 
		nanosecondByte	:= []byte(nanosecond)
		endLineByte1     := []byte("\n")
		endLineByte2     := []byte("\n\r")

		bufferW.SendBWspace("date"      , &dateByte )
		bufferW.SendBWspace("fileName"  , &fileNameByte )
		bufferW.SendBWspace("funcion"   , &funcionByte )
		bufferW.SendBWspace("line"      , &lineByte )
		bufferW.SendBWspace("endLine1"       , &endLineByte1 )

		bufferW.SendBWspace("DAC"       , &urlByte )
		bufferW.SendBWspace("nanosecond"       , &nanosecondByte )
		bufferW.SendBWspace("endLine2"       , &endLineByte2 )
		bufferW.Wspace();

	}

}

func WriteNewMemoryUse(spaceName string,date string, fileName string,funcion string, line string, url string, 
alloc string, totalAlloc string,memSys string, bigObjSize string, gcCount string){

	var bufferW *WBuffer
	var file *spaceFile

	file = memoryLog.OSpace(spaceName)

	if file != nil {

		bufferW = file.NewWBspace(buffMap)
		
	}

	if bufferW != nil {

		bufferW.NewLineWBspace()
		dateByte	    := []byte(date)
		fileNameByte	:= []byte(fileName)
		funcionByte  	:= []byte(funcion)
		lineByte  	    := []byte(line)
		urlByte         := []byte(url) 

		allocByte	    := []byte(alloc)
		totalAllocByte	:= []byte(totalAlloc)
		memSysByte	    := []byte(memSys)
		bigObjSizeByte	:= []byte(bigObjSize)
		gcCountByte	    := []byte(gcCount)

		endLineByte1     := []byte("\n")
		endLineByte2     := []byte("\n\r")

		bufferW.SendBWspace("date"      , &dateByte )
		bufferW.SendBWspace("fileName"  , &fileNameByte )
		bufferW.SendBWspace("funcion"   , &funcionByte )
		bufferW.SendBWspace("line"      , &lineByte )
		bufferW.SendBWspace("endLine1"  , &endLineByte1 )

		bufferW.SendBWspace("DAC"       , &urlByte )
		bufferW.SendBWspace("endLine2"  , &endLineByte1 )

		bufferW.SendBWspace("Alloc"       , &allocByte )
		bufferW.SendBWspace("totalAlloc"  , &totalAllocByte )
		bufferW.SendBWspace("memSys"           , &memSysByte )
		bufferW.SendBWspace("objLargerMemory"  , &bigObjSizeByte )
		bufferW.SendBWspace("countGC"          , &gcCountByte )
		bufferW.SendBWspace("endLine3"         , &endLineByte2 )


		bufferW.Wspace();

	}

}

func WriteNewError(spaceName string, date string, typeError errorDac,fileName string,funcion string, line string, messageLog string, url string ) {


	var bufferW *WBuffer
	var file *spaceFile

	file = errorLog.OSpace(spaceName)

	if file != nil {

		bufferW = file.NewWBspace(buffMap)
		
	}

	if bufferW != nil {

		bufferW.NewLineWBspace()
		dateByte	    := []byte(date)
		typeErrorByte	:= []byte(typeError)
		fileNameByte	:= []byte(fileName)
		funcionByte  	:= []byte(funcion)
		lineByte  	    := []byte(line)
		messageLogByte	:= []byte(messageLog)
		urlByte         := []byte(url) 

		endLineByte1     := []byte("\n")
		endLineByte2     := []byte("\n\r")

		bufferW.SendBWspace("date"      , &dateByte )
		bufferW.SendBWspace("typeError" , &typeErrorByte )
		bufferW.SendBWspace("fileName"  , &fileNameByte )
		bufferW.SendBWspace("funcion"   , &funcionByte )
		bufferW.SendBWspace("line"      , &lineByte )
		bufferW.SendBWspace("endLine1"  , &endLineByte1 )
		bufferW.SendBWspace("message"   , &messageLogByte )
		bufferW.SendBWspace("endLine2"  , &endLineByte1 )
		bufferW.SendBWspace("DAC"       , &urlByte )
		bufferW.SendBWspace("endLine3"  , &endLineByte2 )
		bufferW.Wspace();

	}
}



//check: revisa las columnas y los fields haber si existen como columna si no existe da error fatal.
//#core.go/check
func (sP *space)checkColFil(name string, err string){

	mensaje := err + " ; La columna o field: " + name + " no existe en el archivo ; " + sP.dir

	if sP.indexSizeColumns != nil {

		_, found := sP.indexSizeColumns[name]
		if !found {
		
			if sP.indexSizeFields != nil {

				_, found := sP.indexSizeFields[name]
				if !found {
		
					log.Fatalln(mensaje)
				
				}
				return
			}
		}
	}

	if sP.indexSizeColumns == nil {

		if sP.indexSizeFields != nil {

			_, found := sP.indexSizeFields[name]
			if !found {
	
				log.Fatalln(mensaje)
			
			}
			return
		}
	}

	if sP.indexSizeColumns == nil && sP.indexSizeFields == nil {

		log.Fatalln("El archivo no tiene columnas o campos que sincronizar." + sP.dir)

	}

}