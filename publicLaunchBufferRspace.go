package bd

import (
	"log"
)

//Crea un nuevo buffer de bytes de lectura
func (sF *spaceFile) NewReaderBytes() (buf *RBuffer) {

	return &RBuffer{
		spaceFile: sF,
		typeBuff:  buffBytes,
	}
}

//Crea un nuevo canal de lectura
func (sF *spaceFile) NewReaderChan() (buf *RBuffer) {

	return &RBuffer{
		spaceFile: sF,
		typeBuff:  buffChan,
	}
}

//Crea un nuevo mapa[string]bytes de lectura
func (sF *spaceFile) NewReaderMapBytes() (buf *RBuffer) {

	return &RBuffer{
		spaceFile: sF,
		typeBuff:  buffMap,
	}
}

//Activa el postformateado de los datos.
func (RB *RBuffer) PostFormatBRspace() {

	RB.postFormat = true
}

//Lee una sola linea de buffer. (Columnas)
func (RB *RBuffer) OneLineBRspace(line int64) {

	if RB.spaceErrors != nil && *RB.sizeFileLine < line {

		log.Fatalln("Error de buffer archivo: ", RB.url, " Linea final: ", line, " Numero de lineas del archivo: ", *RB.sizeFileLine)

	}

	if RB.rLines != nil {

		log.Fatalln("Accion incompatible en paralelo, crea un nuevo RBuffer para mayor seguridad.")

	}

	RB.rLines = &rLines{
		startLine: line,
		endLine:   line + 1,
	}

}

//Lee multiples lineas de buffer.(Columnas)
func (RB *RBuffer) MultiLineBRspace(startLine int64, endLine int64) {

	if RB.spaceErrors != nil && *RB.sizeFileLine < endLine {

		log.Fatalln("Error de buffer archivo: ", RB.url, " Linea final: ", endLine, " Numero de lineas del archivo: ", *RB.sizeFileLine)

	}

	if RB.typeBuff == buffBytes {

		log.Fatalln("Incompatible con buffer de bytes")

	}

	if RB.rLines != nil {

		log.Fatalln("Accion incompatible en paralelo, crea un nuevo RBuffer para mayor seguridad.")

	}

	RB.rLines = &rLines{
		startLine: startLine,
		endLine:   endLine + 1,
	}

}

//Lee todas las lineas de un archivo (Columnas)
func (RB *RBuffer) AllLinesBRspace() {

	if RB.typeBuff == buffBytes {

		log.Fatalln("Incompatible con buffer de bytes")

	}

	if RB.rLines != nil {

		log.Fatalln("Accion incompatible en paralelo, crea un nuevo RBuffer para mayor seguridad.")

	}

	RB.rLines = &rLines{
		startLine: 0,
		endLine:   *RB.sizeFileLine + 1,
	}

}

//Lee un field sin rangos.(Campos)
func (RB *RBuffer) NoRangeFieldsBRspace() {

	if RB.rRangues != nil {

		log.Fatalln("Accion incompatible en paralelo, crea un nuevo RBuffer para mayor seguridad.")

	}

	RB.rRangues = &rRangues{
		rangeBytes:  0,
		totalRangue: 1,
		rangue:      0,
	}
}

//Lee un fields por rangos, de manera dinamica.(Campos)
func (RB *RBuffer) RangeFieldsBRspace(RangeBytes int64) {

	if RB.rRangues != nil {

		log.Fatalln("Accion incompatible en paralelo, crea un nuevo RBuffer para mayor seguridad.")

	}

	RB.rRangues = &rRangues{
		rangeBytes:  RangeBytes,
		totalRangue: 1,
		rangue:      0,
	}

}

//CAlcula el numero de rangos de un field o nil si no existe.
func (RB *RBuffer) CalcRangeFieldBRspace(field string, RangeBytes int64) *int64 {

	if RB.indexSizeFields != nil {

		size, found := RB.indexSizeFields[field]
		if found {

			sizeTotal := size[1] - size[0]
			if RangeBytes < sizeTotal && RangeBytes > 0 {

				TotalRangue := sizeTotal / RB.rangeBytes
				restoRangue := sizeTotal % RB.rangeBytes
				if restoRangue != 0 {

					TotalRangue += 1
				}
				return &TotalRangue
			}
		}
	}
	return nil
}



//BRspace: crea un buffer de lectura, se puede elegir si añadir el postformat de los campos.
//Las lineas estan precalculadas inicio 0 - fin - 0 equivale a la linea 0, 0 - 1 equivale a la linea 0 y 1.
//data son los fields y las columnas que se desean.
func (RB *RBuffer) BRspace(data ...string) {

	if RB.spaceErrors != nil {

		if int64(len(data)) > RB.lenColumns+RB.lenFields {

			log.Fatalln("Error el archivo solo tiene: ", RB.lenColumns+RB.lenFields, "columnas y campos", RB.url)

		}

		if len(data) == 0 {

			log.Fatalln("No se puede enviar un buffer vacio en:", RB.url)

		}

		if checkFileTypeBuffer(RB.typeBuff, buffBytes) {

			if RB.spaceErrors != nil && len(data) > 1 {

				log.Fatalln("El Buffer de Bytes solo es compatible con un unico campo.", RB.url)

			}
		}

		for _, colname := range data {

			if RB.spaceErrors != nil {

				RB.checkColFil(colname, "Archivo: BRspace.go ; Funcion: BRspace")
			}
		}

	}

	RB.colName = &data

	//Buffer de bytes
	if checkFileTypeBuffer(RB.typeBuff, buffBytes) {

		if RB.indexSizeFields != nil && RB.rRangues != nil {

			_, found := RB.indexSizeFields[(*RB.colName)[0]]
			if found {

				RB.FieldBuffer = new([]byte)
			}
		}

		if RB.indexSizeColumns != nil && RB.rLines != nil {

			size, found := RB.indexSizeColumns[(*RB.colName)[0]]
			if found {

				RB.Buffer = new([]byte)
				*RB.Buffer = make([]byte, size[1]-size[0])
			}
		}
		return
	}

	if checkFileTypeBuffer(RB.typeBuff, buffChan) {

		if RB.indexSizeFields != nil && RB.rRangues != nil {

			RB.FieldBuffer = new([]byte)
		}

		if RB.indexSizeColumns != nil && RB.rLines != nil {

			RB.Buffer = new([]byte)
			*RB.Buffer = make([]byte, RB.sizeLine)
			RB.Channel = make(chan RChanBuf, 1)
		}
		return
	}

	if checkFileTypeBuffer(RB.typeBuff, buffMap) {

		RB.BufferMap = make(map[string][][]byte)

		if RB.indexSizeFields != nil && RB.rRangues != nil {

			RB.FieldBuffer = new([]byte)
		}

		if RB.indexSizeColumns != nil && RB.rLines != nil {

			RB.BufferMap["buffer"] = make([][]byte, 1)
			RB.BufferMap["buffer"][0] = make([]byte, RB.sizeLine*(RB.endLine-RB.startLine))
		}
		return
	}
}
