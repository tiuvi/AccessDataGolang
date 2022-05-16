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
func NewSf(fileNativeType fileNativeType, fileCoding fileCoding, extension string, fields map[string]int64, columns map[string]int64, dirName ...string) *PublicSpaceFile {

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

	switch fileNativeType {

	case permDisk:
		space.NewTimeFilePermDisk()
	case deferDisk:
		space.NewTimeFileDeferDisk()
	case disk:
		space.NewTimeFileDisk()
	case openFile:
		space.NewTimeFileOpenFile()
	}

	if len(fields) > 0 {
		for name, size := range fields {

			space.NewField(name, size)
		}
	}

	switch extension {

	case dacBit:
		space.NewDacBit()
	case dacByte:
		space.NewDacByte()
	default:
		space.SetExtension(extension)
		
	}

	switch fileCoding {
	case bit:
		space.SetFileCodgingBit()
		if len(columns) > 0 {

			for name := range columns {

				space.NewColumnBit(name)
			}
		}
	case bytes:
		space.SetFileCodgingByte()
		if len(columns) > 0 {

			for name, size := range columns {

				space.NewColumnByte(name, size)
			}
		}
	}

	space.SetSubDir(dirName[:len(dirName)-1]...)

	space.OSpaceInit()

	return &PublicSpaceFile{
		spaceFile: space.OSpace(dirName[len(dirName)-1]),
	}
}

func NewSfPermBytes(fields map[string]int64, columns map[string]int64, dirName ...string) *PublicSpaceFile {

	return NewSf(permDisk, bytes, dacByte, fields, columns, dirName...)

}

func NewSfDeferDiskBytes(fields map[string]int64, columns map[string]int64, dirName ...string) *PublicSpaceFile {

	return NewSf(deferDisk, bytes, dacByte, fields, columns, dirName...)

}

func NewSfDiskBytes(fields map[string]int64, columns map[string]int64, dirName ...string) *PublicSpaceFile {

	return NewSf(disk, bytes, dacByte, fields, columns, dirName...)

}

func NewSfopenFile(fields map[string]int64, columns map[string]int64, dirName ...string) *PublicSpaceFile {

	return NewSf(openFile, bytes, dacByte, fields, columns, dirName...)

}
