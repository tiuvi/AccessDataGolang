package bd

import (
	"goblue/core/autoload_options"
	"log"
	"runtime"
	"strconv"
	"strings"
	"time"
)


const (
	ColorWhite        = "\u001b[37m"
	ColorBlack        = "\u001b[30m"
	ColorRed          = "\u001b[31m"
	ColorGreen        = "\u001b[32m"
	ColorYellow       = "\u001b[33m"
	ColorBlue         = "\u001b[34m"
	Reset             = "\u001b[0m"
)

type errorDac string
const(
	
	NewRouteFolder errorDac       =   "newRouteFolder" 

	Message    errorDac           =  "message"
	MessageCopilation  errorDac   = "messageCopilation"

	Exception  errorDac   =  "exception"
	Warning    errorDac   =  "warning"
	Fatal      errorDac   =  "fatal"

	TimeMemory errorDac = "TimeMemoryStat"
)

type SpaceErrors struct {
	//Si fatal close program si no log normales.
	LogFatalErrors   bool
	LogConsoleErrors bool

	LevelsUrl int
	LogFileError     bool

	LogTimeUse        bool
	LogFileTimeUse   bool
	LogTimeOpenFile bool

	LogMemoryUse     bool
	LogFileMemoryUse     bool

	SeparatorLog string
}

var GlobalError = &SpaceErrors{
	LogFatalErrors:   false,
	LogConsoleErrors: true,

	LogFileError: true,

	LogTimeUse :   true,   
	LogFileTimeUse: true,
	LogTimeOpenFile: true,


	LogMemoryUse:  true,
	LogFileMemoryUse: true,

	SeparatorLog: " ",
	LevelsUrl: 8,
}

var errorLog = &Space{
	Check: true,
	FileNativeType: PermDisk,
	Dir: autoload_options.Patch + "/errors/dac/",
	Extension: DacByte,
	IndexSizeFieldsArray: []spaceLen{
		{"exceptionCount", 64,},
		{"warningCount", 64,},
		{"fatalCount", 64,},
	},
	IndexSizeColumnsArray: []spaceLen{
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

var memoryLog = &Space{
	Check: true,
	FileNativeType: PermDisk,
	Dir: autoload_options.Patch + "/errors/dac/",
	Extension: DacByte,
	IndexSizeColumnsArray: []spaceLen{
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

var timeLog = &Space{
	Check: true,
	FileNativeType: PermDisk,
	Dir: autoload_options.Patch + "/errors/dac/",
	Extension: DacByte,
	IndexSizeColumnsArray: []spaceLen{
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


func negritaTerminal(str string)string{
	
	return strings.Join( []string{"\u001b[01m" ,ColorWhite, "\u001b[40m",   str , Reset} , "")
}

func (EDAC *ErrorsDac ) uint64ToString(uintData uint64)string {

	

	var uintFormat string
	uintStr := strconv.FormatUint( uintData, 10) 

	divuintStr  := len(uintStr) / 3 
	restuintStr := len(uintStr) % 3 

	if restuintStr > 0 {

		uintFormat = strings.Join([]string{ uintStr[:restuintStr] , EDAC.SeparatorLog } , "")
		uintStr = uintStr[restuintStr:]
	
	}

	for x := 0; x < divuintStr ; x++{
		
		uintFormat += strings.Join([]string{ uintStr[:3] , EDAC.SeparatorLog } , "")
		uintStr = uintStr[3:]

	}

	return uintFormat
}

func errorsLog(){

	errorLog.OSpaceInit()
	timeLog.OSpaceInit()
	memoryLog.OSpaceInit()
}

type ErrorsDac struct  {
	*SpaceErrors
	FileName string
	TypeError errorDac
	Url string
	MessageLog string

	LevelsUrl int
	SeparatorLog string

	TimeNow *time.Time
}



func (sP *Space )ErrorSpaceDefault(typeError errorDac, MessageLog string){

	EDAC := &ErrorsDac {
		SpaceErrors: sP.SpaceErrors,
		FileName: "",
		TypeError: typeError,
		Url: sP.Dir,
		MessageLog: MessageLog,
		LevelsUrl:    sP.LevelsUrl ,
		SeparatorLog:  sP.SeparatorLog,
		TimeNow: nil,
	}
	EDAC.LogNewError()
}

func (sP *Space ) NewErrorSpace(FileName string, typeError errorDac, MessageLog string){

	EDAC := &ErrorsDac {
		SpaceErrors: sP.SpaceErrors,
		FileName: FileName,
		TypeError: typeError,
		Url: sP.Dir,
		MessageLog: MessageLog,
		LevelsUrl:    sP.LevelsUrl ,
		SeparatorLog:  sP.SeparatorLog,
		TimeNow: nil,
	}
	EDAC.LogNewError()
}

func (sF *spaceFile ) ErrorSpaceFileDefault(typeError errorDac, MessageLog string){

	EDAC := &ErrorsDac {
		SpaceErrors: sF.SpaceErrors,
		FileName: "",
		TypeError: typeError,
		Url: sF.Url,
		MessageLog: MessageLog,
		LevelsUrl:    sF.LevelsUrl ,
		SeparatorLog:  sF.SeparatorLog,
		TimeNow: nil,
	}
	EDAC.LogNewError()
}

func (sF *spaceFile ) NewErrorSpace(FileName string, typeError errorDac, MessageLog string){

	EDAC := &ErrorsDac {
		SpaceErrors: sF.SpaceErrors,
		FileName: FileName,
		TypeError: typeError,
		Url: sF.Url,
		MessageLog: MessageLog,
		LevelsUrl:    sF.LevelsUrl ,
		SeparatorLog:  sF.SeparatorLog,
		TimeNow: nil,
	}
	EDAC.LogNewError()
}

func (sP *Space ) LogDeferTimeMemoryDefault(timeNow time.Time){

	EDAC := &ErrorsDac {
		SpaceErrors: sP.SpaceErrors,
		FileName: "",
		TypeError: TimeMemory,
		Url: sP.Dir,
		MessageLog: "",
		LevelsUrl:    sP.LevelsUrl ,
		SeparatorLog:  sP.SeparatorLog,
		TimeNow: &timeNow,
	}
	EDAC.LogNewError()

}

func (sP *Space ) NewLogDeferTimeMemory(FileName string, timeNow time.Time){

	EDAC := &ErrorsDac {
		SpaceErrors: sP.SpaceErrors,
		FileName: FileName,
		TypeError: TimeMemory,
		Url: sP.Dir,
		MessageLog: "",
		LevelsUrl:    sP.LevelsUrl ,
		SeparatorLog:  sP.SeparatorLog,
		TimeNow: &timeNow,
	}
	EDAC.LogNewError()

}

func (sF *spaceFile ) LogDeferTimeMemorySF(timeNow time.Time){

	EDAC := &ErrorsDac {
		SpaceErrors: sF.SpaceErrors,
		FileName: "",
		TypeError: TimeMemory,
		Url: sF.Url,
		MessageLog: "",
		LevelsUrl:    sF.LevelsUrl ,
		SeparatorLog:  sF.SeparatorLog,
		TimeNow: &timeNow,
	}
	EDAC.LogNewError()

}



func (sF *spaceFile ) NewLogDeferTimeMemorySF(FileName string, timeNow time.Time){

	EDAC := &ErrorsDac {
		SpaceErrors: sF.SpaceErrors,
		FileName: FileName,
		TypeError: TimeMemory,
		Url: sF.Url,
		MessageLog: "",
		LevelsUrl:    sF.LevelsUrl ,
		SeparatorLog:  sF.SeparatorLog,
		TimeNow: &timeNow,
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
	if len(urlName) > EDAC.LevelsUrl {

		urlNameString = strings.Join(urlName[len(urlName)-EDAC.LevelsUrl:] , "/")
	}

	if EDAC.LogFileError && EDAC.TypeError != MessageCopilation && EDAC.TimeNow == nil{
	
		if len(EDAC.FileName) == 0{

			go WriteNewError("message",dateString, EDAC.TypeError ,fileString ,funcNameString , lineStr , EDAC.MessageLog , urlNameString )

		}else{

			go WriteNewError(EDAC.FileName + "Message" ,dateString, EDAC.TypeError ,fileString ,funcNameString , lineStr , EDAC.MessageLog , urlNameString )

		}
	
	}


	var logString string
	var nT func(string)string = negritaTerminal

	if (EDAC.LogFatalErrors || EDAC.LogConsoleErrors) && EDAC.TimeNow == nil {

		var color string
		switch EDAC.TypeError {
	
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
		nT("Fecha:") ," ", dateString ," - ",nT("Tipo de error:"), " " ,color , string(EDAC.TypeError) ,Reset, "\n\r", 
		nT("Archivo:"), " "  , fileString ," - ",nT("Funcion:"), " " , funcNameString ," - ",nT("Nº Linea:"), " "  , lineStr ,"\n\r",
		nT("Ruta DAC:"), " " , urlNameString, "\n\r",
		nT("Mensaje:"), " " , EDAC.MessageLog, "\n\r" } , "")
	
		if EDAC.LogFatalErrors  {

			log.Fatalln(logString )
	
		} 
	
		if EDAC.LogConsoleErrors {
		
			go log.Println(logString)
		
		}
	}


	var memoryStats runtime.MemStats
	var allocString string
	var totalAlloc string
	var memSys string
	var bigObjSize string
	var gcCount string
	if (EDAC.LogMemoryUse || EDAC.LogFileMemoryUse) && EDAC.TimeNow != nil {

		runtime.ReadMemStats(&memoryStats)
		allocString = EDAC.uint64ToString(memoryStats.Alloc)
		totalAlloc  = EDAC.uint64ToString(memoryStats.TotalAlloc)
		memSys      = EDAC.uint64ToString(memoryStats.Sys)
		bigObjSize  = EDAC.uint64ToString(uint64(memoryStats.BySize[60].Size))
		gcCount     = EDAC.uint64ToString(uint64(memoryStats.NumGC))
	}

	//Log de memoria y archivo de memoria
	if EDAC.LogFileMemoryUse && EDAC.TypeError != MessageCopilation && EDAC.TimeNow != nil {

		if len(EDAC.FileName) == 0{

			go WriteNewMemoryUse("memory",dateString, fileString ,funcNameString , lineStr , urlNameString, 
			allocString, totalAlloc,memSys , bigObjSize, gcCount )
		}else{

			go WriteNewMemoryUse(EDAC.FileName + "memory",dateString, fileString ,funcNameString , lineStr , urlNameString, 
			allocString, totalAlloc,memSys , bigObjSize, gcCount )
		}
	
	}

	// For info on each, see: https://golang.org/pkg/runtime/#MemStats
	if EDAC.LogMemoryUse && EDAC.TimeNow != nil {


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
	if (EDAC.LogFileTimeUse || EDAC.LogTimeUse) && EDAC.TimeNow != nil  {
		
		elapsed = time.Since(*EDAC.TimeNow).Nanoseconds()

	}

	//Log de tiempo transcurrido y archivo de tiempo transcurrido
	if EDAC.LogFileTimeUse && EDAC.TypeError != MessageCopilation && EDAC.TimeNow != nil {

		if len(EDAC.FileName) == 0{

			go WriteNewTimeUse("time",dateString, fileString ,funcNameString , lineStr , urlNameString, EDAC.uint64ToString(uint64(elapsed)) )

		}else{

			go WriteNewTimeUse(EDAC.FileName + "Time",dateString, fileString ,funcNameString , lineStr , urlNameString, EDAC.uint64ToString(uint64(elapsed)) )

		}
	}

	if EDAC.LogTimeUse && EDAC.TimeNow != nil {

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

		bufferW = file.NewWBspace(BuffMap)
		
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

		bufferW = file.NewWBspace(BuffMap)
		
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

		bufferW = file.NewWBspace(BuffMap)
		
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
func (sP *Space)checkColFil(name string, err string){

	mensaje := err + " ; La columna o field: " + name + " no existe en el archivo ; " + sP.Dir

	if sP.IndexSizeColumns != nil {

		_, found := sP.IndexSizeColumns[name]
		if !found {
		
			if sP.IndexSizeFields != nil {

				_, found := sP.IndexSizeFields[name]
				if !found {
		
					log.Fatalln(mensaje)
				
				}
				return
			}
		}
	}

	if sP.IndexSizeColumns == nil {

		if sP.IndexSizeFields != nil {

			_, found := sP.IndexSizeFields[name]
			if !found {
	
				log.Fatalln(mensaje)
			
			}
			return
		}
	}

	if sP.IndexSizeColumns == nil && sP.IndexSizeFields == nil {

		log.Fatalln("El archivo no tiene columnas o campos que sincronizar." + sP.Dir)

	}

}