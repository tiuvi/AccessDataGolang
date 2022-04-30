package bd

import (
	"log"
	"os"
	"regexp"
)



//Inicia DAC y crea la carpeta para la aplicación
func NewLaunchDac() *lDAC {

	if globalLaunchDac == nil {

		//Variable SuperGlobal de DAC
		globalLaunchDac = new(lDAC)

		go dacTimerCloserDeferFile()
	
		go dacTimerCloserDiskFile()
	
		go errorsLog()
	}

	return new(lDAC)

}

func regexPathGlobal(path string) string {

	//Filtrado solo se permite letras mayusculas, letras minusculas, numeros y barras.
	return regexp.MustCompile(`[^a-zA-Z0-9/]`).ReplaceAllString(path, "")

}

/**********************************************************************************************/
/* Funciones: spaceErrors  */
/**********************************************************************************************/

//Inicia la activacion de errores.
func (LDAC *lDAC) OnErrors() {

	LDAC.spaceErrors = new(spaceErrors)
	LDAC.levelsUrl = 8
	LDAC.separatorLog = " "
}

//Desactiva los errores.
func (LDAC *lDAC) OffErrors() {

	LDAC.spaceErrors = nil
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
	LDAC.OnLogConsoleErrors()
	LDAC.OnLogFileError()
	LDAC.OnLogTimeUse()
	LDAC.OnLogFileTimeUse()
	LDAC.OnLogTimeOpenFile()
	LDAC.OnLogMemoryUse()
	LDAC.OnLogFileMemoryUse()
}

func (LDAC *lDAC) OffAllErrors() {

	LDAC.OffLogFatalErrors()
	LDAC.OffLogConsoleErrors()
	LDAC.OffLogFileError()
	LDAC.OffLogTimeUse()
	LDAC.OffLogFileTimeUse()
	LDAC.OffLogTimeOpenFile()
	LDAC.OffLogMemoryUse()
	LDAC.OffLogFileMemoryUse()
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

func (LDAC *lDAC) SetGlobalDACFolder(path string) {

	if len(LDAC.globalDACFolder) > 0 {

		log.Fatal("La ruta DAC ya esta establecida en: ", "\r\n",
			LDAC.globalDACFolder, "\r\n",
			"¿Correcto? Revisa - bd.SetGlobalDACFolder")
	}

	if len(path) == 0 {

		log.Fatal("Estas enviando una cadena vacia.", "\r\n",
			"¿Correcto? Revisa - bd.SetGlobalDACFolder")

	}

	//Filtrado solo se permite letras mayusculas, letras minusculas, numeros y barras.
	path = regexPathGlobal(path)

	//Añadimos / barra al final si no la lleva
	if path[len(path)-1:] != "/" {

		path += "/"

	}

	//Añadimos el nombre del archivo DAC
	path += "DAC/"

	//Verificamos la ruta antes de crear la carpeta
	if !LDAC.goDACFolder {

		log.Fatal("Chequea la ruta final antes de iniciar DAC: ", "\r\n",
			path, "\r\n",
			"¿Correcto? Activa - bd.OnGoDACFolder()")

	}

	//Verificamos si el archivo existe y si no existe lo creamos.
	_, err := os.Stat(path)
	if err != nil {

		if os.IsNotExist(err) && !LDAC.newDACFolder {

			log.Fatal("No existe la ruta a la carpeta DAC, o has introducido una ruta incorrecta.", "\r\n",
				"¿Si te gustaria crear la carpeta DAC? , Usa onCreateDACFolder() para crear la carpeta en la ruta seleccionada.")

		}

		if os.IsNotExist(err) {

			err = os.MkdirAll(path, 0666)
			if err != nil {

				log.Fatal("Error al crear la carpeta DAC ", err)

			}

			log.Println("Genial se a creado la carpeta dac en: ", "\r\n",
				path)

		}

	}

	LDAC.globalDACFolder = path

}

func (LDAC *lDAC) GetGlobalDACFolder() string {

	if len(LDAC.globalDACFolder) == 0 {

		log.Fatal("La ruta DAC no esta establecida en: ", "\r\n",
			LDAC.globalDACFolder, "\r\n",
			"¿Correcto? Activa - bd.SetGlobalDACFolder")
	}

	return LDAC.globalDACFolder

}
