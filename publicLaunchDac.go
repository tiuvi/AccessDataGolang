package dac

import (
	"fmt"
	"os"
)


//Inicia DAC y crea la carpeta para la aplicación
func NewLaunchDac()(ldac *lDAC) {

	 ldac = new(lDAC)
	 ldac.timersFile = &timersFile{
		fileOpenDeferFile:  612, // Linux permite 1024 descriptores por proceso.
		timeEventDeferFile: 21600, //6 horas

		timeEventDiskFile: 21600, //6 horas
	 }

	 return ldac
}


/**********************************************************************************************/
/* Funciones: spaceErrors  */
/**********************************************************************************************/

//Inicia la activacion de errores.
func (LDAC *lDAC) OnErrors() {

	LDAC.spaceErrors = new(spaceErrors)
	LDAC.levelsUrl = 8
	LDAC.separatorLog = " "
	EDAC = true
}

//Desactiva los errores.
func (LDAC *lDAC) OffErrors() {

	LDAC.spaceErrors = nil
	EDAC = false
}
//Inicia la activacion de errores.
func (LDAC *lDAC) OnEvalErrors() {

	if LDAC.spaceErrors == nil {

		LDAC.OnErrors()
	}

	EDAC = true
}

//Desactiva la evaluacion de los errores.
func (LDAC *lDAC) OffEvalErrors() {

	EDAC = false
}

//Activa logConsoleErrors
func (LDAC *lDAC) OnLogConsoleErrors() {

	LDAC.logConsoleErrors = true

}

//Desactiva logConsoleErrors
func (LDAC *lDAC) OffLogConsoleErrors() {

	LDAC.logConsoleErrors = false

}

//Activa logFatalErrors
func (LDAC *lDAC) OnLogFatalErrors() {

	LDAC.logFatalErrors = true

}

//Desactiva logFatalErrors
func (LDAC *lDAC) OffLogFatalErrors() {

	LDAC.logFatalErrors = false

}

//Activa logFileError
func (LDAC *lDAC) OnLogFileError() {

	LDAC.logFileError = true

}

//Desactiva logFileError
func (LDAC *lDAC) OffLogFileError() {

	LDAC.logFileError = false

}

//Activa logTimeUse
func (LDAC *lDAC) OnLogTimeUse() {

	LDAC.logTimeUse = true

}

//Desactiva logTimeUse
func (LDAC *lDAC) OffLogTimeUse() {

	LDAC.logTimeUse = false

}

//Activa logFileTimeUse
func (LDAC *lDAC) OnLogFileTimeUse() {

	LDAC.logFileTimeUse = true

}

//Desactiva logFileTimeUse
func (LDAC *lDAC) OffLogFileTimeUse() {

	LDAC.logFileTimeUse = false

}

//Activa logTimeOpenFile
func (LDAC *lDAC) OnLogTimeOpenFile() {

	LDAC.logTimeOpenFile = true

}

//Desactiva logTimeOpenFile
func (LDAC *lDAC) OffLogTimeOpenFile() {

	LDAC.logTimeOpenFile = false

}

//Activa logMemoryUse
func (LDAC *lDAC) OnLogMemoryUse() {

	LDAC.logMemoryUse = true

}

//Desactiva logMemoryUse
func (LDAC *lDAC) OffLogMemoryUse() {

	LDAC.logMemoryUse = false

}

//Activa LogFileMemoryUse
func (LDAC *lDAC) OnLogFileMemoryUse() {

	LDAC.logFileMemoryUse = true

}

//Desactiva LogFileMemoryUse
func (LDAC *lDAC) OffLogFileMemoryUse() {

	LDAC.logFileMemoryUse = false

}

//Activa logTimeReadFile
func (LDAC *lDAC) OnLogTimeReadFile() {

	LDAC.logTimeReadFile = true

}

//Desactiva logTimeReadFile
func (LDAC *lDAC) OffLogTimeReadFile() {

	LDAC.logTimeReadFile = false

}

//Activa logTimeReadFile
func (LDAC *lDAC) OnLogTimeWriteFile() {

	LDAC.logTimeWriteFile = true

}

//Desactiva logTimeReadFile
func (LDAC *lDAC) OffLogTimeWriteFile() {

	LDAC.logTimeWriteFile = false

}



//Elije el separador para los nanosegundos
func (LDAC *lDAC) SetSeparatorLog(separator string) {

	LDAC.separatorLog = separator

}

//Elimina el separador para los nanosegundos
func (LDAC *lDAC) UnSetSeparatorLog() {

	LDAC.separatorLog = ""

}

//Elije el separador para los levelsUrl
func (LDAC *lDAC) SetLevelsUrl(level int) {

	if level > 0 {
		LDAC.levelsUrl = level
	}
}

//Elije el separador para los nanosegundos
func (LDAC *lDAC) UnSetLevelsUrl() {

	LDAC.levelsUrl = 0

}

func (LDAC *lDAC) OnAllErrors() {

	LDAC.OnErrors()
	LDAC.OnLogFatalErrors()
	LDAC.OnLogFileError()

	LDAC.OnLogTimeUse()
	LDAC.OnLogFileTimeUse()
	
	LDAC.OnLogMemoryUse()
	LDAC.OnLogFileMemoryUse()

	LDAC.OnLogTimeOpenFile()
	LDAC.OnLogTimeReadFile()
	LDAC.OnLogTimeWriteFile()

	EDAC = true
}

func (LDAC *lDAC) OffAllErrors() {

	LDAC.OffLogFatalErrors()
	LDAC.OffLogConsoleErrors()
	LDAC.OffLogFileError()

	LDAC.OffLogTimeUse()
	LDAC.OffLogFileTimeUse()

	LDAC.OffLogMemoryUse()
	LDAC.OffLogFileMemoryUse()

	LDAC.OffLogTimeOpenFile()
	LDAC.OffLogTimeReadFile()
	LDAC.OffLogTimeWriteFile()

}





/**********************************************************************************************/
/* Funciones: lDAC  */
/**********************************************************************************************/


func (LDAC *lDAC) OnCreateDACFolder() {

	LDAC.newDACFolder = true

}

func (LDAC *lDAC) OffCreateDACFolder() {

	LDAC.newDACFolder = false

}

func (LDAC *lDAC) OnGoDACFolder() {

	LDAC.goDACFolder = true

}

func (LDAC *lDAC) OffGoDACFolder() {

	LDAC.goDACFolder = false

}


/**********************************************************************************************/
/* Configuracion eventos de fichero */
/**********************************************************************************************/

func (LDAC *lDAC) ConfCloserDiskFile(seconds int64){

	if EDAC && 
	LDAC.ELDAC(seconds < 600 ,"El valor no puede ser inferior a 600 segundos"){}

	LDAC.timersFile.timeEventDiskFile = seconds
}

func (LDAC *lDAC) ConfCloserDeferFile(fileopen int, seconds int64){

	if EDAC && 
	LDAC.ELDAC(fileopen < 612 , "El valor de archivos abiertos no puede ser inferior a 612.") ||
	LDAC.ELDAC(seconds < 600 , "El valor no puede ser inferior a 600 segundos"){}

	LDAC.timersFile.fileOpenDeferFile  = fileopen
	LDAC.timersFile.timeEventDeferFile = seconds
}


func (LDAC *lDAC) SetGlobalDACFolder(path string) {

	if EDAC && 
	LDAC.ELDAC(len(LDAC.globalDACFolder) > 0 ,"La ruta DAC ya esta establecida.") ||
	LDAC.ELDAC(len(path) == 0 ,"patch vacio"){}

	//Filtrado solo se permite letras mayusculas, letras minusculas, numeros y barras.
	path = regexPathGlobal(path)

	//Añadimos / barra al final si no la lleva
	if path[len(path)-1:] != "/" {

		path += "/"

	}


	//Verificamos la ruta antes de crear la carpeta
	if EDAC && 
	LDAC.ELDAC(!LDAC.goDACFolder ,"¿Correcto?: " + path + " |Activa - bd.OnGoDACFolder()"){}


	//Verificamos si el archivo existe y si no existe lo creamos.
	_, err := os.Stat(path)
	if err != nil {

		if EDAC && 
		LDAC.ELDAC(os.IsNotExist(err) && !LDAC.newDACFolder ,
		"No existe la ruta a la carpeta DAC, o has introducido una ruta incorrecta." + 
		"\r\n" +
		"¿Crear nueva ruta? , |Activa onCreateDACFolder() para crear la carpeta en la ruta seleccionada."){}

		
		if os.IsNotExist(err) {

			err = os.MkdirAll(path, os.ModeDir | os.ModePerm)
			if err != nil && EDAC && 
			LDAC.ELDAC( true,"Error al crear la carpeta DAC. \n\r" + fmt.Sprintln(err)){}
			
		}

	}


	LDAC.globalDACFolder = path

	if globalDac == nil {

		//Variable SuperGlobal de DAC
		globalDac = LDAC

		go LDAC.dacTimerCloserDeferFile()

		go LDAC.dacTimerCloserDiskFile()

		Space = make(map[string]*space)

		//Activacion de archivos de errores de DAC
		globalDac.onErrorsLog()
	}

}

func (LDAC *lDAC) GetGlobalDACFolder() string {
	
	if EDAC && 
	LDAC.ELDAC(len(LDAC.globalDACFolder) == 0,"La ruta DAC no esta establecida."){}
	
	return LDAC.globalDACFolder

}

func GetGlobalDac()*lDAC{
	
	if globalDac != nil {

		if len(globalDac.globalDACFolder) != 0 {

			return globalDac
		}
	} 
	return nil
}


