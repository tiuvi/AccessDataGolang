package dac



func (SF *spaceFile) GetOneFieldString(col string)string {

	buffer := SF.GetOneField(col)

	return string(*buffer.FieldBuffer)

}

func (SF *spaceFile) GetOneFieldBytes(col string)[]byte {

	buffer := SF.GetOneField(col)

	return *buffer.FieldBuffer

}






func (SF *spaceFile) SetOneFieldString(col string,str string)*int64 {

	buffer := []byte(str)

	return SF.SetOneField(col , &buffer)

}