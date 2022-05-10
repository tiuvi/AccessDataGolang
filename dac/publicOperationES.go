package dac



/*
* Funciones de buffer de bytes lectura para imagenes
*
 */

 func (sF *spaceFile) GetOneFieldRanges(col string, RangeBytes int64 , Rangue int64) (RBuf *RBuffer) {

	RBuf = sF.NewReaderBytes()
	RBuf.RangeFieldsBRspace(RangeBytes ,Rangue )

	RBuf.BRspace(col)
	RBuf.Rspace()

	return RBuf
}

func (sF *spaceFile) GetOneFieldChan(col string, RangeBytes int64 ) (<-chan RChanBuf) {

	RBuf := sF.NewReaderChan()
	RBuf.RangeFieldsBRspace(RangeBytes ,0 )
	RBuf.BRspace(col)
	go RBuf.Rspace()

	return  RBuf.Channel
}

/*
* Lectura con postFormat
*
**/
func (sF *spaceFile) GetOneField(col string) (RBuf *RBuffer) {

	RBuf = sF.NewReaderBytes()
	RBuf.OnPostFormat()
	RBuf.NoRangeFieldsBRspace()
	RBuf.BRspace(col)
	RBuf.Rspace()

	return RBuf
}

func (sF *spaceFile) GetOneLine(col string, line int64) *RBuffer {

	RBuf := sF.NewReaderBytes()
	RBuf.OnPostFormat()
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
	RBuf.OnPostFormat()
	RBuf.AllLinesBRspace()
	RBuf.BRspace(col)
	RBuf.Rspace()

	return RBuf
}

/*
* Funciones de buffer de mapas escritura
*
 */




