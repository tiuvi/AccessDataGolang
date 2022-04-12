package bd

	import "log"


//Esta funcion escribe en el numero de linea requerido
//Pasandole un string en forma de bytes
func (buf *WBuffer) Wspace()int64{


	//spaceFile := obj.OSpace()


	if CheckFileCoding(buf.FileCoding , Bit){

		if CheckFileTipeBit(buf.FileTipeBit,ListBit ){
			//Funcion ListBit interfaz
			
			return buf.WriteListBitSpace()
		}

		log.Fatalln("Error Grave, Uspace.go ; Funcion: Wspace ;" +
		"No Hubo coincidencias FileTipeBit")
		return -1
	}

	if CheckFileCoding(buf.FileCoding , Byte){

		if CheckFileTipeByte(buf.FileTipeByte,OneColumn){

			return buf.WriteColumnSpace()
		}

		if CheckFileTipeByte(buf.FileTipeByte,MultiColumn){

			return buf.WriteColumnSpace()
		}

		if  CheckFileTipeByte(buf.FileTipeByte,FullFile){

			log.Fatalln("Error Grave, Uspace.go ; Funcion: Wspace ;" +
			"Funcion writeFullFile en desarrollo...")
			return -1
		}

		log.Fatalln("Error Grave, Uspace.go ; Funcion: Wspace ;" +
		"No Hubo coincidencias FileTipeByte")
		return -1
	}


	if CheckFileCoding(buf.FileCoding , Dir){

		if CheckFileTypeDir(buf.FileTypeDir, EmptyDir){

			log.Fatalln("Error Grave, Uspace.go ; Funcion: Wspace ;" +
			"Funcion WriteEmptyDirSpace en desarrollo...")
			//Añadir modificacion de url desde aqui
			//spaceFile.WriteEmptyDirSpace(buf)
			return -1
		}

		log.Fatalln("Error Grave, Uspace.go ; Funcion: Wspace ;" +
		"No Hubo coincidencias FileTypeDir")
		return -1
	}

	log.Fatalln("Error Grave, Uspace.go ; Funcion: Wspace ;" +
	"No Hubo coincidencias FileCoding")
	return -1
}



func (buf *RBuffer)Rspace (){



	//spaceFile := obj.OSpace()

	
	if CheckFileCoding(buf.FileCoding , Bit){

		if CheckFileTipeBit(buf.FileTipeBit,ListBit ){
			//Funcion ListBit interfaz
			buf.ListBitSpace()
			return
		}

		log.Fatalln("Error Grave, Uspace.go ; Funcion: Rspace ;" +
		"No Hubo coincidencias FileTipeBit")
		return
	}

	if CheckFileCoding(buf.FileCoding , Byte){

		if CheckFileTipeByte(buf.FileTipeByte,OneColumn){
			//Funcion OneColumn interfaz
			buf.OneColumnSpace()
			return
		}

		if CheckFileTipeByte(buf.FileTipeByte,MultiColumn){
			//Funcion Multicolumn interfaz
			buf.MultiColumnSpace()
			return
		}

		if  CheckFileTipeByte(buf.FileTipeByte,FullFile){
			//Añadir modificacion de url desde aqui
			buf.FullFileSpace()
		}

		log.Fatalln("Error Grave, Uspace.go ; Funcion: Rspace ;" +
		"No Hubo coincidencias FileTipeByte")
		return
	}


	if CheckFileCoding(buf.FileCoding , Dir){

		if CheckFileTypeDir(buf.FileTypeDir, EmptyDir){

			//Añadir modificacion de url desde aqui
			buf.ReadEmptyDirSpace()
			return
		}

		log.Fatalln("Error Grave, Uspace.go ; Funcion: Rspace ;" +
		"No Hubo coincidencias FileTypeDir")
		return
	}

	log.Fatalln("Error Grave, Uspace.go ; Funcion: Rspace ;" +
	"No Hubo coincidencias FileCoding")
	return
}






func (obj *spaceFile ) Rmapspace(str_write string)(value int64, found bool){

	if int64(len(str_write)) > obj.SizeLine{

		str_write = str_write[:obj.SizeLine]

	}
	
	if obj.Hooker != nil {

		bufferByte := []byte(str_write)
		obj.hookerPreFormatPointer(&bufferByte, Preformat)
		str_write = string(bufferByte)
	}



	

	obj.RLock()
	value, found = obj.Search[str_write]
	obj.RUnlock()
	return
}


func (obj *spaceFile ) Rindexspace(line int64)(value string, found bool){

	if  int64(len(obj.Index) ) > line {

		value, found = obj.Index[line] , true

	} else {

		value, found = "", false
	}

	return
}