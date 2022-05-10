package dac


func (buf *WBuffer) Wspace()*int64{


	if checkFileCoding(buf.fileCoding , bit){

		return buf.writeBitSpace()
		
	}

	if checkFileCoding(buf.fileCoding , bytes){

		return buf.writeByteSpace()
		
	}

	if EDAC && 
	buf.ECSD( true,"Error grave, sin coincidencias de filecoding"){}

	return nil
}



func (buf *RBuffer)Rspace (){

	if checkFileCoding(buf.fileCoding , bit){

		buf.readBitSpace()
		return
	}

	if checkFileCoding(buf.fileCoding , bytes){

		 buf.readByteSpace()
		 return
	}

	if EDAC && 
	buf.ECSD( true,"Error grave, sin coincidencias de filecoding"){}

	return
}
