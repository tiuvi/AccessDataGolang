package bd

import (
		"log"
		"runtime"
		"strconv"
		"strings"
		"goblue/core/autoload_options"
)

type errorDac string
const(
	message    errorDac   =  "message"
	exception  errorDac   =  "exception"
	warning    errorDac   =  "warning"
	fatal      errorDac   =  "fatal"
)

var GlobalError = &SpaceErrors{
	LogFatalErrors: false,
	LogConsoleErrors: true,
	LogFileError: true,
	LevelsUrl: 3 ,
	LogSpeeds : true,   

}

var errorLog = &Space{
	Check: true,
	FileNativeType: PermDisk,
	Dir: autoload_options.Patch + "/errors/",
	Extension: DacByte,
	IndexSizeFields: map[string][2]int64{
		"exceptionCount": {0,64},
		"warningCount": {64,128},
		"fatalCount": {128,192},
	},
	
	IndexSizeColumns: map[string][2]int64{
		"typeError": {0,10},
		"fileName": {10,30},
		"funcion": {30,65},
		"line": {65,70},
		"message": {70,326},
		"url": {326,582},
		
	},
}
var ErrorLogFile *spaceFile

func errorsLog(){

	log.SetFlags(log.LstdFlags | log.Lshortfile)

	errorLog.OSpaceInit()
	ErrorLogFile := errorLog.OSpace("errors")
	_ = ErrorLogFile



}

func (sP *Space ) LogNewError(typeError errorDac, messageLog string, url string ){

	ptr , _, _, _ := runtime.Caller(1)
	firstFrame := runtime.CallersFrames([]uintptr{ptr})
	frame, _ := firstFrame.Next()
	
	fileString     := frame.File
	funcNameString := frame.Function
	urlNameString  := url

	fileDir    := strings.Split(frame.File, "/")
	if len(fileDir) > 2 {

		fileString = strings.Join(fileDir[len(fileDir)-2:] , "/")

	}
	
	funcName       := strings.Split(frame.Function , "/")

	if len(funcName) > 1 {

		funcNameString = strings.Join(funcName[len(funcName)-1:]  , "/")
	}
	

	urlName       := strings.Split(url , "/")
	if len(urlName) > sP.LevelsUrl {

		urlNameString = strings.Join(urlName[len(urlName)-sP.LevelsUrl:] , "/")
	}
	

	lineStr := strconv.Itoa(frame.Line)

	if sP.LogFatalErrors || typeError == fatal {

		log.Fatalln(typeError , fileString , funcNameString , lineStr , messageLog , urlNameString )

	} else {

		if sP.LogConsoleErrors {
		
			log.Println(typeError , fileString , funcNameString , lineStr, messageLog , urlNameString )
		
		}
		
		if sP.LogFileError {

			WriteNewError(typeError ,fileString ,funcNameString , lineStr , messageLog , urlNameString )

		}
	}
}

func WriteNewError(typeError errorDac,fileName string,funcion string, line string, messageLog string, url string ) {

	if ErrorLogFile != nil {

		typeErrorByte	:= []byte(typeError)
		fileNameByte	:= []byte(fileName)
		funcionByte  	:= []byte(funcion)
		lineByte  	    := []byte(line)
		messageLogByte	:= []byte(messageLog)
		urlByte         := []byte(url) 

		bufferW := ErrorLogFile.NewWBspace(BuffMap)
		bufferW.NewLineWBspace()

		bufferW.SendBWspace("fileName"  , &fileNameByte )
		bufferW.SendBWspace("funcion"   , &funcionByte )
		bufferW.SendBWspace("line"      , &lineByte )

		bufferW.SendBWspace("message"   , &messageLogByte )
		bufferW.SendBWspace("url"       , &urlByte )
		bufferW.SendBWspace("typeError" , &typeErrorByte )
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