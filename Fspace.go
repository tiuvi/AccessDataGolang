package bd

import (
//	"log"
)

/*
* Funciones de buffer de bytes lectura
*
*/
func (sF *spaceFile)GetOneLine(col string, line int64) *RBuffer {

	RBuf := sF.NewBRspace(BuffBytes)
	RBuf.OneLineBRspace(line)
	RBuf.BRspace(col)
	RBuf.Rspace()

	return RBuf
}

/*
* Funciones de buffer de bytes escritura
*
*/
func (sF *spaceFile)SetOneLine(col string, line int64,bufferBytes *[]byte )*int64 {

	RBuf := sF.NewWBspace(BuffBytes)
	RBuf.UpdateLineWBspace(line)
	RBuf.SendBWspace(col,bufferBytes)
	
	return RBuf.Wspace()
}

func (sF *spaceFile)NewOneLine(col string, bufferBytes *[]byte )*int64 {

	RBuf := sF.NewWBspace(BuffBytes)
	RBuf.NewLineWBspace()
	RBuf.SendBWspace(col,bufferBytes)
	
	return RBuf.Wspace()
}


/*
* Funciones de buffer de mapas lectura
*
*/
func (sF *spaceFile)GetAllLines(col string) *RBuffer {

	RBuf := sF.NewBRspace(BuffMap)
	RBuf.AllLinesBRspace()
	RBuf.BRspace(col)
	RBuf.Rspace()

	return RBuf
}

/*
* Funciones de buffer de mapas escritura
*
*/

