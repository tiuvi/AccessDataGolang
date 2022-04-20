package bd

import (
//	"log"
)

/*
* Funciones de buffer de bytes
*
 */
func (sF *spaceFile) GetOneLine(col string, line int64) *RBuffer {

	RBuf := sF.NewBRspace(BuffBytes)
	RBuf.OneLineBRspace(line)
	RBuf.BRspace(col)
	RBuf.Rspace()

	return RBuf
}

/*
* Funciones de buffer de mapas
*
 */
func (sF *spaceFile) GetAllLines(col string) *RBuffer {

	RBuf := sF.NewBRspace(BuffMap)
	RBuf.AllLinesBRspace()
	RBuf.BRspace(col)
	RBuf.Rspace()

	return RBuf
}
