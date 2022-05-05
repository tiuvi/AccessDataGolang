package dac

import (
	"log"
)

/*
* Funciones de buffer de bytes lectura
*
 */

func (sF *spaceFile) GetOneField(col string) (RBuf *RBuffer) {

	RBuf = sF.NewReaderBytes()
	RBuf.NoRangeFieldsBRspace()
	RBuf.BRspace(col)
	RBuf.Rspace()

	return RBuf
}

func (sF *spaceFile) GetOneLine(col string, line int64) *RBuffer {

	RBuf := sF.NewReaderBytes()
	RBuf.OneLineBRspace(line)
	RBuf.BRspace(col)
	RBuf.Rspace()

	return RBuf
}

/*
* Funciones de buffer de bytes escritura
*
 */
 func (SF *spaceFile) SetOneField(col string, bufferBytes *[]byte)*int64 {

	WBuffBytes := SF.NewWriterBytes()
	
	WBuffBytes.NewNoRangeWBspace()
	
	WBuffBytes.SendBWspace(col , bufferBytes )
	
	return WBuffBytes.Wspace()
}


func (sF *spaceFile) SetOneLine(col string, line int64, bufferBytes *[]byte) *int64 {

	RBuf := sF.NewWriterBytes()
	RBuf.UpdateLineWBspace(line)
	RBuf.SendBWspace(col, bufferBytes)

	return RBuf.Wspace()
}

func (sF *spaceFile) NewOneLine(col string, bufferBytes *[]byte) *int64 {

	RBuf := sF.NewWriterBytes()
	RBuf.NewLineWBspace()
	RBuf.SendBWspace(col, bufferBytes)

	return RBuf.Wspace()
}

/*
* Funciones de buffer de mapas lectura
*
 */
func (sF *spaceFile) GetAllLines(col string) *RBuffer {

	RBuf := sF.NewReaderMapBytes()
	RBuf.AllLinesBRspace()
	RBuf.BRspace(col)
	RBuf.Rspace()

	return RBuf
}

/*
* Funciones de buffer de mapas escritura
*
 */






/*
*
* Funcion de inicio dac con errores
*
*/
func NewDac(path string,fopenDefer int, secondsEvent int64 ){

		// Inicio de DAC
		lDAC := NewLaunchDac()
		lDAC.OnAllErrors()
		lDAC.ConfCloserDiskFile(secondsEvent)
		lDAC.ConfCloserDeferFile(fopenDefer ,secondsEvent )
		//lDAC.OffEvalErrors()
		lDAC.OnCreateDACFolder()
		lDAC.OnGoDACFolder()
		lDAC.SetGlobalDACFolder(path)
}





/*
* Funciones de nuevo espacio global
*
 */

func NewSfPermBytes(fields map[string]int64, columns map[string]int64, dirName ...string) *spaceFile {

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

	return space.OSpace(dirName[len(dirName)-1])
}
