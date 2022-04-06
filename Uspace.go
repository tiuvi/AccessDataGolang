package bd

	import "log"


//Esta funcion escribe en el numero de linea requerido
//Pasandole un string en forma de bytes
func (obj *Space ) Wspace(buf *WBuffer)int64{


	spaceFile := obj.OSpace()


	if CheckFileCoding(obj.FileCoding , Bit){

		if CheckFileTipeBit(obj.FileTipeBit,ListBit ){
			//Funcion ListBit interfaz
			
			return spaceFile.WriteListBitSpace(buf)
		}

		log.Fatalln("Error Grave, Uspace.go ; Funcion: Wspace ;" +
		"No Hubo coincidencias FileTipeBit")
		return -1
	}

	if CheckFileCoding(obj.FileCoding , Byte){

		if CheckFileTipeByte(obj.FileTipeByte,OneColumn){

			return spaceFile.WriteColumnSpace(buf)
		}

		if CheckFileTipeByte(obj.FileTipeByte,MultiColumn){

			return spaceFile.WriteColumnSpace(buf)
		}

		if  CheckFileTipeByte(obj.FileTipeByte,FullFile){

			log.Fatalln("Error Grave, Uspace.go ; Funcion: Wspace ;" +
			"Funcion writeFullFile en desarrollo...")
			return -1
		}

		log.Fatalln("Error Grave, Uspace.go ; Funcion: Wspace ;" +
		"No Hubo coincidencias FileTipeByte")
		return -1
	}


	if CheckFileCoding(obj.FileCoding , Dir){

		if CheckFileTypeDir(obj.FileTypeDir, EmptyDir){

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



func (obj *Space) Rspace (column *RBuffer){



	spaceFile := obj.OSpace()

	
	if CheckFileCoding(obj.FileCoding , Bit){

		if CheckFileTipeBit(obj.FileTipeBit,ListBit ){
			//Funcion ListBit interfaz
			spaceFile.ListBitSpace(column)
			return
		}

		log.Fatalln("Error Grave, Uspace.go ; Funcion: Rspace ;" +
		"No Hubo coincidencias FileTipeBit")
		return
	}

	if CheckFileCoding(obj.FileCoding , Byte){

		if CheckFileTipeByte(obj.FileTipeByte,OneColumn){
			//Funcion OneColumn interfaz
			spaceFile.OneColumnSpace(column)
			return
		}

		if CheckFileTipeByte(obj.FileTipeByte,MultiColumn){
			//Funcion Multicolumn interfaz
			spaceFile.MultiColumnSpace(column)
			return
		}

		if  CheckFileTipeByte(obj.FileTipeByte,FullFile){
			//Añadir modificacion de url desde aqui
			spaceFile.FullFileSpace(column)
		}

		log.Fatalln("Error Grave, Uspace.go ; Funcion: Rspace ;" +
		"No Hubo coincidencias FileTipeByte")
		return
	}


	if CheckFileCoding(obj.FileCoding , Dir){

		if CheckFileTypeDir(obj.FileTypeDir, EmptyDir){

			//Añadir modificacion de url desde aqui
			spaceFile.ReadEmptyDirSpace(column)
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
	

	//Preformat global
	function, exist := obj.Hooker[Preformat]
	if exist {

		str_write = string(function([]byte(str_write)))

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