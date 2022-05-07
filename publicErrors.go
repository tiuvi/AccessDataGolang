package dac

import (
	"time"
)

/*

ELDACF
if EDAC && 
obj.ELDACF( true," "){}

Evaluacion global de excepciones
if EDAC && 
obj.ECSD( true," "){}

Error mas evaluacion global
if err != nil && EDAC && 
obj.ECSD( true," \n\r" + fmt.Sprintln(err)){}

*/


var runCallerGlobal = 3

//Coge la url space.dir para el error
func (LDAC *lDAC) ErrorLDACDefault(typeError errorDac, messageLog string)bool {

	EDAC := &errorsDac{
		spaceErrors:  LDAC.spaceErrors,
		fileName:     "",
		typeError:    typeError,
		url:          LDAC.globalDACFolder,
		messageLog:   messageLog,
		levelsUrl:    99999,
		separatorLog: "",
		runCaller: 3,
		timeNow:      nil,
	}
	EDAC.logNewError()
	return true
}

func (LDAC *lDAC) ELDAC (conditional bool, msg string)bool {
	
	if LDAC.spaceErrors != nil {

		if conditional {

			return	LDAC.ErrorLDACDefault(MessageCopilation, msg)
		}	
	}
	return false
}

func (LDAC *lDAC) ELDACF (conditional bool, msg string)bool {
	
	if LDAC.spaceErrors != nil {

		if conditional {

			return	LDAC.ErrorLDACDefault(Fatal, msg)
		}	
	}
	return false
}

//Coge la url space.dir para el error
func (sP *space) ErrorSpaceDefault(typeError errorDac, messageLog string)bool {

	EDAC := &errorsDac{
		spaceErrors:  sP.spaceErrors,
		fileName:     "",
		typeError:    typeError,
		url:          sP.dir,
		messageLog:   messageLog,
		levelsUrl:    sP.levelsUrl,
		separatorLog: sP.separatorLog,
		runCaller: runCallerGlobal ,
		timeNow:      nil,
	}
	EDAC.logNewError()
	return true
}

//#ErrorCheckSpaceDefault -> ECSD
func (sP *space) ECSD(conditional bool, messageLog string)bool {

	if sP.spaceErrors != nil {
	
		if conditional {
			return sP.ErrorSpaceDefault(Fatal , messageLog)
		}
	}
	return false
}


//Coge la url space.dir para el error
func (sP *space) NewErrorSpace(fileName string, typeError errorDac, messageLog string)bool {

	EDAC := &errorsDac{
		spaceErrors:  sP.spaceErrors,
		fileName:     fileName,
		typeError:    typeError,
		url:          sP.dir,
		messageLog:   messageLog,
		levelsUrl:    sP.levelsUrl,
		separatorLog: sP.separatorLog,
		runCaller:runCallerGlobal,
		timeNow:      nil,
	}
	EDAC.logNewError()
	return true
}

//Coge la url space.dir para el error
func (sP *space) NewRouteErrorSpace(typeError errorDac, messageLog string, fileName string, fileFolder ...string)bool {

	EDAC := &errorsDac{
		spaceErrors:  sP.spaceErrors,
		fileName:     fileName,
		fileFolder: fileFolder,
		typeError:    typeError,
		url:          sP.dir,
		messageLog:   messageLog,
		levelsUrl:    sP.levelsUrl,
		separatorLog: sP.separatorLog,
		runCaller:runCallerGlobal,
		timeNow:      nil,
	}
	EDAC.logNewError()
	return true
}


//Coge la url spaceFile.url para el error
func (sF *spaceFile) ErrorSpaceFileDefault(typeError errorDac, messageLog string)bool {

	EDAC := &errorsDac{
		spaceErrors:  sF.spaceErrors,
		fileName:     "",
		typeError:    typeError,
		url:          sF.url,
		messageLog:   messageLog,
		levelsUrl:    sF.levelsUrl,
		separatorLog: sF.separatorLog,
		runCaller:runCallerGlobal,
		timeNow:      nil,
	}
	EDAC.logNewError()
	return true
}

//#ErrorCheckSpaceFileDefault -> ECSFD
func (sF *spaceFile) ECSFD(conditional bool, messageLog string)bool {

	if sF.spaceErrors != nil {
	
		if conditional {
			return sF.ErrorSpaceDefault(Fatal , messageLog)
		}
	}
	return false
}

//Coge la url spaceFile.url para el error
func (sF *spaceFile) NewErrorSpaceFile(fileName string, typeError errorDac, messageLog string)bool {

	EDAC := &errorsDac{
		spaceErrors:  sF.spaceErrors,
		fileName:     fileName,
		typeError:    typeError,
		url:          sF.url,
		messageLog:   messageLog,
		levelsUrl:    sF.levelsUrl,
		separatorLog: sF.separatorLog,
		runCaller:runCallerGlobal,
		timeNow:      nil,
	}
	EDAC.logNewError()
	return true
}

//Coge la url spaceFile.url para el error
func (sF *spaceFile) NewRouteErrorSpaceFile(typeError errorDac, messageLog string, fileName string, fileFolder ...string)bool {

	EDAC := &errorsDac{
		spaceErrors:  sF.spaceErrors,
		fileName:     fileName,
		fileFolder:   fileFolder,
		typeError:    typeError,
		url:          sF.url,
		messageLog:   messageLog,
		levelsUrl:    sF.levelsUrl,
		separatorLog: sF.separatorLog,
		runCaller:runCallerGlobal,
		timeNow:      nil,
	}
	EDAC.logNewError()
	return true
}

//NewRouteErrorSpaceFile message
func (SF *spaceFile) NRESM(conditional bool, messageLog string, fileName string, fileFolder ...string)bool{
	
	if SF.spaceErrors != nil {

		if conditional {

			return	SF.NewRouteErrorSpaceFile(Message , messageLog , fileName , fileFolder...)
		}	
	}
	return false
}







//Coge la url space.dir para el error
func  LogDeferTimeMemoryDefaultDac(timeNow time.Time)bool {

	EDAC := &errorsDac{
		spaceErrors:  globalDac.spaceErrors,
		fileName:     "",
		typeError:    TimeMemory,
		url:          globalDac.globalDACFolder,
		messageLog:   "",
		levelsUrl:    globalDac.levelsUrl,
		separatorLog: globalDac.separatorLog,
		runCaller:    runCallerGlobal ,
		timeNow:      &timeNow,
	}
	EDAC.logNewError()
	return true
}

//Coge la url space.dir para el error
func (sP *space) LogDeferTimeMemoryDefault(timeNow time.Time)bool {

	EDAC := &errorsDac{
		spaceErrors:  sP.spaceErrors,
		fileName:     "",
		typeError:    TimeMemory,
		url:          sP.dir,
		messageLog:   "",
		levelsUrl:    sP.levelsUrl,
		separatorLog: sP.separatorLog,
		runCaller:    runCallerGlobal ,
		timeNow:      &timeNow,
	}
	EDAC.logNewError()
	return true
}

//Coge la url space.dir para el error
func (sP *space) NewLogDeferTimeMemory(fileName string, timeNow time.Time)bool {

	EDAC := &errorsDac{
		spaceErrors:  sP.spaceErrors,
		fileName:     fileName,
		typeError:    TimeMemory,
		url:          sP.dir,
		messageLog:   "",
		levelsUrl:    sP.levelsUrl,
		separatorLog: sP.separatorLog,
		runCaller: runCallerGlobal ,
		timeNow:      &timeNow,
	}
	EDAC.logNewError()
	return true
}

//Coge la url space.dir para el error
func (sP *space) NewRouteLogDeferTimeMemory(timeNow time.Time, fileName string, fileFolder ...string )bool {

	EDAC := &errorsDac{
		spaceErrors:  sP.spaceErrors,
		fileName:     fileName,
		fileFolder:   fileFolder,
		typeError:    TimeMemory,
		url:          sP.dir,
		messageLog:   "",
		levelsUrl:    sP.levelsUrl,
		separatorLog: sP.separatorLog,
		runCaller:    runCallerGlobal ,
		timeNow:      &timeNow,
	}
	EDAC.logNewError()
	return true
}


//Coge la url spaceFile.url para el error
func (sF *spaceFile) LogDeferTimeMemorySF(timeNow time.Time)bool {

	EDAC := &errorsDac{
		spaceErrors:  sF.spaceErrors,
		fileName:     "",
		typeError:    TimeMemory,
		url:          sF.url,
		messageLog:   "",
		levelsUrl:    sF.levelsUrl,
		runCaller:    runCallerGlobal ,
		separatorLog: sF.separatorLog,
		timeNow:      &timeNow,
	}
	EDAC.logNewError()
	return true
}

//Coge la url spaceFile.url para el error
func (sF *spaceFile) NewLogDeferTimeMemorySF(fileName string, timeNow time.Time)bool {

	EDAC := &errorsDac{
		spaceErrors:  sF.spaceErrors,
		fileName:     fileName,
		typeError:    TimeMemory,
		url:          sF.url,
		messageLog:   "",
		levelsUrl:    sF.levelsUrl,
		separatorLog: sF.separatorLog,
		runCaller:    runCallerGlobal ,
		timeNow:      &timeNow,
	}
	EDAC.logNewError()
	return true
}

//Coge la url spaceFile.url para el error
func (sF *spaceFile) NewRouteLogDeferTimeMemorySF(timeNow time.Time, fileName string, fileFolder ...string)bool {

	EDAC := &errorsDac{
		spaceErrors:  sF.spaceErrors,
		fileName:     fileName,
		fileFolder:   fileFolder,
		typeError:    TimeMemory,
		url:          sF.url,
		messageLog:   "",
		levelsUrl:    sF.levelsUrl,
		separatorLog: sF.separatorLog,
		runCaller:    runCallerGlobal ,
		timeNow:      &timeNow,
	}
	EDAC.logNewError()
	return true
}

