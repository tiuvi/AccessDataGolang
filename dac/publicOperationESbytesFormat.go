package dac


/*
* Lectura de fields bytes
*/

func (SF *spaceFile) GetOneFieldString(col string)string {

	buffer := SF.GetOneField(col)

	return string(*buffer.FieldBuffer)

}

func (SF *spaceFile) GetOneFieldBytes(col string)[]byte {

	buffer := SF.GetOneField(col)

	return *buffer.FieldBuffer

}

func (SF *spaceFile) GetOneFieldBytesRaw(col string)[]byte {

	buffer := SF.GetOneFieldRaw(col)

	return *buffer.FieldBuffer

}


/*
* Escritura de fields bytes
*/

func (SF *spaceFile) SetOneFieldString(col string,str string)*int64 {

	buffer := []byte(str)

	return SF.SetOneField(col , &buffer)

}


//Obtiene una linea en un archivo y lo formatea directamente a string
func (SF *spaceFile) GetOneLineString(col string, line int64)string {

	RBuffer := SF.GetOneLine(col, line)
	return string(*RBuffer.Buffer)
}

//Crea una nueva linea en un archivo desde un string en la columna especificada.
func (SF *spaceFile) NewOneLineString(col string, str string)*int64{

	WBuffer := []byte(str)
	return SF.NewOneLine(col , &WBuffer )
}

//Actualiza una nueva linea en un archivo desde un string en la columna especificada.
func (SF *spaceFile) SetOneLineString(col string,line int64, str string)*int64{

	WBuffer := []byte(str)
	return SF.SetOneLine(col , line , &WBuffer )
}

