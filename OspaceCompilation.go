package bd



func (obj *space) ospaceCompilationFile()bool {

	//eMSGoCF: Error Mensaje ospaceCompilationFile
	//#Error MessageCopilation -> ErrorSpaceDefault
	var eMSGoCF func(c bool,b string)bool = func (conditional bool, msg string)bool {
	
		if obj.spaceErrors != nil {
	
			if conditional {
	
				return obj.ErrorSpaceDefault(MessageCopilation, msg)
			}	
		}
		return false
	}

	eMSGoCF(obj.fileNativeType == 0 ,`Variable fileNativeType no definida.`)

	eMSGoCF(obj.fileNativeType == 0,`Variable fileNativeType no definida.`)

	eMSGoCF(len(obj.dir) == 0 ,`Variable dir vacia.`)

	eMSGoCF(len(obj.extension) == 0, `Extension vacia.`)


	eMSGoCF(obj.IsNotExtension(obj.extension),`Extension no valida.`)

	
	

	if len(obj.indexSizeFieldsArray) > 0 {

		if obj.indexSizeFields == nil {

			obj.indexSizeFields = make(map[string][2]int64)

		}

		var afterValue int64
		for ind, value := range obj.indexSizeFieldsArray {

			
			eMSGoCF(value.len < 0,`indexSizeFieldsArray valor negativo.`)
			
			eMSGoCF(obj.IsField(value.name),`Nombre repetido en indexSizeFieldsArray.`)

			

			if ind == 0 {

				obj.indexSizeFields[value.name] = [2]int64{0, value.len}

			} else {

				obj.indexSizeFields[value.name] = [2]int64{afterValue, value.len + afterValue}

			}

			afterValue += value.len
		}
	}





	if len(obj.indexSizeColumnsArray) > 0 {

		if obj.indexSizeColumns == nil {

			obj.indexSizeColumns = make(map[string][2]int64)

		}

		var afterValue int64
		for ind, value := range obj.indexSizeColumnsArray {

			eMSGoCF(value.len < 0,`indexSizeFieldsArray valor negativo.`)
			
			eMSGoCF(obj.IsColumn(value.name),`Nombre repetido en indexSizeColumnsArray.`)

			if ind == 0 {

				obj.indexSizeColumns[value.name] = [2]int64{0, value.len}

			} else {

				obj.indexSizeColumns[value.name] = [2]int64{afterValue, value.len + afterValue}

			}

			afterValue += value.len
		}
	}



	if obj.indexSizeFields != nil {

		obj.lenFields  = int64(len(obj.indexSizeFields))
	
		eMSGoCF(obj.lenFields == 0,`Iniciaste un mapa de campos, sin campos`)


		var checkSizeFields int64 = 0

		for name, val := range obj.indexSizeFields {

			eMSGoCF(obj.IsColumn(name),`Los nombres de fields y columnas tiene que ser unicos.`)

			calcSize := (val[1] - val[0])
		
			eMSGoCF(calcSize <= 0, `El field` + name + ` no puede tener un tamaño igual o inferior a cero.`)
			
			obj.sizeField += calcSize

			if val[1] >= checkSizeFields {

				checkSizeFields = val[1]
			}
		}

		eMSGoCF(checkSizeFields != obj.sizeField,`Los campos estan mal escritos, Ejemplo: field1: 0,20; field2:20,30`)

		eMSGoCF(obj.IsField("buffer"),`La palabra buffer en campos esta reservada para el uso del programa.`)

	}


	//Actualizamos el valor del ancho de la linea
	if obj.indexSizeColumns != nil {

		obj.lenColumns = int64(len(obj.indexSizeColumns))
		eMSGoCF(obj.lenColumns == 0,`Iniciaste un mapa de columnas, sin columnas`)
	
		var checkSizeColumns int64 = 0

		for name , val := range obj.indexSizeColumns {

			eMSGoCF(obj.IsField(name),`Los nombres de fields y columnas tiene que ser unicos.`)


			calcSizeLine := (val[1] - val[0])
		

			eMSGoCF(calcSizeLine <= 0 , "Las columnas no pueden tener un tamaño igual o inferior a cero.")

			

			obj.sizeLine += calcSizeLine

			if val[1] >= checkSizeColumns {

				checkSizeColumns = val[1]
			}

		}

		eMSGoCF(checkSizeColumns != obj.sizeLine , "Las columnas estan mal escritas, Ejemplo: column1: 0,20; column2:20,30;")

		eMSGoCF(obj.IsColumn("buffer"),`La palabra buffer en columnas esta reservada para el uso del programa.`)

	}

	eMSGoCF(obj.indexSizeColumns == nil  && obj.indexSizeFields == nil ,`Iniciaste un espacio sin columnas y sin campos.`)



	//Lectura de archivos de byte
	if obj.extension == dacByte {

		obj.fileCoding = bytes
		obj.compilation = true
		return true
	}

	//bdisk: Lista de bit en un archivo disk
	if obj.extension == dacBit {

		obj.fileCoding = bit
		obj.compilation = true
		return true
	}


	eMSGoCF(true,"No se han encontrado coincidencias con las extensiones de archivo predeterminadas.")

	return false
}
