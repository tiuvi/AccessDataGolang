package bd

	import "log"



func (buf *WBuffer) Wspace()*int64{


	if checkFileCoding(buf.fileCoding , bit){

			return buf.writeBitSpace()
		
	}

	if checkFileCoding(buf.fileCoding , bytes){

		return buf.writeByteSpace()
		
	}

	log.Fatalln("Error Grave, Uspace.go ; Funcion: Wspace ;" +
	"No Hubo coincidencias fileCoding")
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

	log.Fatalln("Error Grave, Uspace.go ; Funcion: Rspace ;" +
	"No Hubo coincidencias fileCoding")
	return
}
