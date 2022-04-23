package bd

	import "log"



func (buf *WBuffer) Wspace()*int64{


	if CheckFileCoding(buf.FileCoding , Bit){

			return buf.writeBitSpace()
		
	}

	if CheckFileCoding(buf.FileCoding , Byte){

		return buf.writeByteSpace()
		
	}

	log.Fatalln("Error Grave, Uspace.go ; Funcion: Wspace ;" +
	"No Hubo coincidencias FileCoding")
	return nil
}



func (buf *RBuffer)Rspace (){

	if CheckFileCoding(buf.FileCoding , Bit){

		buf.readBitSpace()
		return
	}

	if CheckFileCoding(buf.FileCoding , Byte){

		 buf.readByteSpace()
		 return
	}

	log.Fatalln("Error Grave, Uspace.go ; Funcion: Rspace ;" +
	"No Hubo coincidencias FileCoding")
	return
}
