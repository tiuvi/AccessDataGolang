package dac



//Crea un nuevo buffer de bytes de lectura
func (sF *spaceFile) NewReaderBytes() (buf *RBuffer) {

	return &RBuffer{
		spaceFile: sF,
		typeBuff:  buffBytes,
		postFormat: false,
	}
}

//Crea un nuevo canal de lectura
func (sF *spaceFile) NewReaderChan() (buf *RBuffer) {

	return &RBuffer{
		spaceFile: sF,
		typeBuff:  buffChan,
		postFormat: false,
	}
}

//Crea un nuevo mapa[string]bytes de lectura
func (sF *spaceFile) NewReaderMapBytes() (buf *RBuffer) {

	return &RBuffer{
		spaceFile: sF,
		typeBuff:  buffMap,
		postFormat: false,
	}
}

//Activa el postformateado de los datos.
func (RB *RBuffer) OnPostFormat() {

	if EDAC && 
	RB.ECSFD(RB.postFormat == true, "Ya se activo el postformateo."){}

	RB.postFormat = true
}

//Lee una sola linea de buffer. (Columnas)
func (RB *RBuffer) OneLineBRspace(line int64) {

	if EDAC && 
	RB.ECSFD(line < 0, "Se trato de leer una linea inferior a 0") ||
	RB.ECSFD(*RB.sizeFileLine < line, "Se trato de leer una linea que no existe.") ||
	RB.ECSFD(RB.rLines != nil, "Llamada dos veces a la misma funcion, Accion incompatible en paralelo, crea un nuevo buffer."){}

	RB.rLines = &rLines{
		startLine: line,
		endLine:   line + 1,
	}

}

//Lee dos lineas de buffer. (Columnas)
func (RB *RBuffer) TwoLineBRspace(line int64) {

	if EDAC && 
	RB.ECSFD(line < 0, "Se trato de leer una linea inferior a 0") ||
	RB.ECSFD(RB.typeBuff == buffBytes, "Funcion Incompatible con buffer de bytes") ||
	RB.ECSFD(*RB.sizeFileLine < line, "Se trato de leer una linea que no existe.") ||
	RB.ECSFD(RB.rLines != nil, "Llamada dos veces a la misma funcion, Accion incompatible en paralelo, crea un nuevo buffer."){}

	RB.rLines = &rLines{
		startLine: line,
		endLine:   line + 2,
	}

}


//Lee tres lineas de buffer. (Columnas)
func (RB *RBuffer) ThreeLineBRspace(line int64) {

	if EDAC && 
	RB.ECSFD(line < 0, "Se trato de leer una linea inferior a 0") ||
	RB.ECSFD(RB.typeBuff == buffBytes, "Funcion Incompatible con buffer de bytes") ||
	RB.ECSFD(*RB.sizeFileLine < line, "Se trato de leer una linea que no existe.") ||
	RB.ECSFD(RB.rLines != nil, "Llamada dos veces a la misma funcion, Accion incompatible en paralelo, crea un nuevo buffer."){}

	RB.rLines = &rLines{
		startLine: line,
		endLine:   line + 3,
	}
}

//Lee multiples lineas de buffer.(Columnas)
func (RB *RBuffer) MultiLineBRspace(startLine int64, endLine int64) {

	if EDAC && 
	RB.ECSFD(startLine < 0, "Se trato de leer una linea inferior a 0") ||
	RB.ECSFD(startLine > endLine, "El principio de linea no puede ser mayor que el final de linea.") ||
	RB.ECSFD(*RB.sizeFileLine < endLine, "Se trato de leer una linea que no existe.") ||
	RB.ECSFD(RB.typeBuff == buffBytes, "Funcion Incompatible con buffer de bytes") ||
	RB.ECSFD(RB.rLines != nil, "Llamada dos veces a la misma funcion, Accion incompatible en paralelo, crea un nuevo buffer."){}

	RB.rLines = &rLines{
		startLine: startLine,
		endLine:   endLine + 1,
	}

}

//Lee todas las lineas de un archivo (Columnas)
func (RB *RBuffer) AllLinesBRspace() {

	if EDAC && 
	RB.ECSFD(RB.typeBuff == buffBytes, "Funcion Incompatible con buffer de bytes") ||
	RB.ECSFD(RB.rLines != nil, "Llamada dos veces a la misma funcion, Accion incompatible en paralelo, crea un nuevo buffer."){}

	RB.rLines = &rLines{
		startLine: 0,
		endLine:   *RB.sizeFileLine + 1,
	}

}

func (RB *RBuffer) FirstLineBRspace() {

	if EDAC && 
	RB.ECSFD(RB.rLines != nil, "Llamada dos veces a la misma funcion, Accion incompatible en paralelo, crea un nuevo buffer."){}

	RB.rLines = &rLines{
		startLine: 0,
		endLine:   1,
	}
}

//Lee la ultima linea del archivo, necesario bloqueo para garantizar que no se escriba
//mientras se esta pidiendo la ultima linea (Columnas)
func (RB *RBuffer) LastLineBRspace() {

	if EDAC && 
	RB.ECSFD(RB.rLines != nil, "Llamada dos veces a la misma funcion, Accion incompatible en paralelo, crea un nuevo buffer."){}

	RB.rLines = &rLines{
		startLine: *RB.sizeFileLine ,
		endLine:   *RB.sizeFileLine + 1,
	}
}

//Lee un field sin rangos.(Campos)
func (RB *RBuffer) NoRangeFieldsBRspace() {

	if EDAC && 
	RB.ECSFD( RB.rRangues != nil, "Llamada dos veces a la misma funcion, Accion incompatible en paralelo, crea un nuevo buffer."){}

	RB.rRangues = &rRangues{
		rangeBytes:  0,
		totalRangue: 1,
		rangue:      0,
	}
}

//Lee un fields por rangos, de manera dinamica.(Campos)
func (RB *RBuffer) RangeFieldsBRspace(RangeBytes int64) {

	if EDAC && 
	RB.ECSFD( RB.rRangues != nil, "Llamada dos veces a la misma funcion, Accion incompatible en paralelo, crea un nuevo buffer.") ||
	RB.ECSFD( RangeBytes < 0, "Rango inferior a 0."){}

	RB.rRangues = &rRangues{
		rangeBytes:  RangeBytes,
		totalRangue: 1,
		rangue:      0,
	}

}





//BRspace: crea un buffer de lectura, se puede elegir si añadir el postformat de los campos.
//Las lineas estan precalculadas inicio 0 - fin - 0 equivale a la linea 0, 0 - 1 equivale a la linea 0 y 1.
//data son los fields y las columnas que se desean.
func (RB *RBuffer) BRspace(data ...string) {

	if EDAC && 
	RB.ECSFD(int64(len(data)) > (RB.lenColumns+RB.lenFields) , "El spacio no tiene tantos fields y Columnas.") ||
	RB.ECSFD(len(data) == 0 , "No se puede enviar un buffer vacio.") ||
	RB.ECSFD(RB.IsNotColFil(data...) , "Se ha iniciado un buffer de lectura con una columna o fields que no existe.") ||
	RB.ECSFD(RB.rRangues == nil && RB.rLines == nil,"Iniciaste un buffer vacio."){}
	

	RB.colName = &data
	var fieldsOk  bool = RB.indexSizeFields  != nil && RB.rRangues != nil
	var columnsOk bool = RB.indexSizeColumns != nil && RB.rLines   != nil

	//Buffer de bytes
	if checkFileTypeBuffer(RB.typeBuff, buffBytes) {

		if EDAC && 
		RB.ECSFD(len(data) > 1  , "El Buffer de Bytes solo es compatible con un unico campo."){}
	
		if fieldsOk {

			_, found := RB.indexSizeFields[(*RB.colName)[0]]
			if found {

				RB.FieldBuffer = new([]byte)
			}
		}

		if columnsOk {

			size, found := RB.indexSizeColumns[(*RB.colName)[0]]
			if found {

				RB.Buffer = new([]byte)
				*RB.Buffer = make([]byte, size[1]-size[0])
			}
		}
		return
	}

	if checkFileTypeBuffer(RB.typeBuff, buffChan) {

		if fieldsOk {

			RB.FieldBuffer = new([]byte)
		}

		if columnsOk {

			RB.Buffer = new([]byte)
			*RB.Buffer = make([]byte, RB.sizeLine)
			RB.Channel = make(chan RChanBuf, 1)
		}
		return
	}

	if checkFileTypeBuffer(RB.typeBuff, buffMap) {

		RB.BufferMap = make(map[string][][]byte)

		if fieldsOk {

			RB.FieldBuffer = new([]byte)
		}

		if columnsOk {

			RB.BufferMap["buffer"] = make([][]byte, 1)
			RB.BufferMap["buffer"][0] = make([]byte, RB.sizeLine*(RB.endLine-RB.startLine))
		}
		return
	}
	
	if EDAC && 
	RB.ECSFD(true , "Comportamiento inesperado, sin coincidencias en buffer de lectura."){}
}
