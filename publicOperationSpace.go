package dac

import (
	"log"
)

/*
*
* Funcion de inicio dac con errores
*
 */

 func NewBasicDac(path string) {

	// Inicio de DAC
	lDAC := NewLaunchDac()
	lDAC.OnBasicErrors()

	lDAC.ConfCloserDiskFile(86400)
	lDAC.ConfCloserDeferFile(612, 86400)

	lDAC.OnCreateDACFolder()
	lDAC.OnGoDACFolder()
	lDAC.SetGlobalDACFolder(path)
}

func NewDac(path string, fopenDefer int, secondsEvent int64) {

	// Inicio de DAC
	lDAC := NewLaunchDac()
	lDAC.OnAllErrors()
	lDAC.ConfCloserDiskFile(secondsEvent)
	lDAC.ConfCloserDeferFile(fopenDefer, secondsEvent)
	lDAC.OnCreateDACFolder()
	lDAC.OnGoDACFolder()
	lDAC.SetGlobalDACFolder(path)
}



/*
* Funciones de nuevo espacio global
*
 */

func NewSfPermBytes(fields map[string]int64, columns map[string]int64, dirName ...string) *PublicSpaceFile {

	DAC := GetGlobalDac()
	if DAC == nil {
		log.Fatal("No se inicio DAC")
	}

	if EDAC &&
		DAC.ELDACF(len(dirName) == 0, "No se puede iniciar un espacio sin nombres") ||
		DAC.ELDACF(len(dirName) == 1, "No es recomdable guardar archivos en la ruta principal.") ||
		DAC.ELDACF(len(fields) == 0 && len(columns) == 0, "Tratas de iniciar un espacio sin fiedls y sin columnas.") {
	}

	//Creacion de un espacio
	space := DAC.NewSpace()
	space.NewTimeFilePermDisk()
	space.NewDacByte()
	space.SetSubDir(dirName[:len(dirName)-1]...)

	if len(fields) > 0 {
		for name, size := range fields {

			space.NewField(name, size)
		}
	}

	if len(columns) > 0 {

		for name, size := range columns {

			space.NewColumnByte(name, size)
		}
	}

	space.OSpaceInit()

	return &PublicSpaceFile{
		spaceFile: space.OSpace(dirName[len(dirName)-1]),
	}
}
