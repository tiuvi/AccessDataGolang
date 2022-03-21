package bd

import (

	"log"
	//"time"
)


//Esta funcion escribe en el numero de linea requerido
//Pasandole un string en forma de bytes
func (obj *Space ) Wspace(line int64, column map[string][]byte){

	//Revisar apertura disk
		if (obj.FileNativeType & disk) != 0 {

			obj.Ospace()

		}


	//Difenciar archivos de bit de archivos de byte
	switch obj.FileCoding {

		//Primer caso archivo de bit
		case Bit:

			//Segun el tipo de archivo aplicamos una funcion diferente
			switch obj.FileTipeBit {
				//Caso lista de bits
				case ListBit:	
					obj.WriteListBitSpace(line, column)
					break
				default:
					return
			}
		// Rompemos switch bit
		break

		//Segundo caso archivos de byte
		case Byte:

			//Segun el tipo de archivo aplicamos una funcion diferente
			switch obj.FileTipeByte {
			
				case OneColumn,MultiColumn:
					//Funcion OneColumn interfaz
					obj.WriteColumnSpace(line, column)
					break

				default: 
					return
			}
		// Rompemos switch byte
		break

		//Defalt si no es bytes o bit
		default: 
				return
	}

		

	
	if (obj.FileNativeType & tdisk) != 0 {

	}

	if (obj.FileNativeType & disk) != 0 {

		defer obj.File.Close()

	}
	
}



//Esta funcion lee la linea simplemente le pasamos el numero de linea
//Que queremos leer y devuelve un buffer

//func (obj *Space) Rspace (line int64)(buffer []byte){
func (obj *Space) Rspace (column Buffer){


	//var start time.Time
	//start = time.Now()
	//log.Println( "Buffer: ", time.Since(start).Nanoseconds())

	//Iniciando archivos cerrados
	if (obj.FileNativeType & disk) != 0 {

		obj.Ospace()

	}




	//Difenciar archivos de bit de archivos de byte
	switch obj.FileCoding {

		//Primer caso archivo de bit
		case Bit:

			//Segun el tipo de archivo aplicamos una funcion diferente
			switch obj.FileTipeBit {
			
				case ListBit:
					//Funcion ListBit interfaz
					obj.ListBitSpace(column)
					break
				default:
					return
			}
		// Rompemos switch bit
		break

		//Segundo caso archivos de byte
		case Byte:

			//Segun el tipo de archivo aplicamos una funcion diferente
			switch obj.FileTipeByte {
			
				case OneColumn:
					//Funcion OneColumn interfaz
					obj.OneColumnSpace(column)
					break
				case MultiColumn:
					//Funcion Multicolumn interfaz
					obj.MultiColumnSpace(column)
					break
				case FullFile:
					//Funcion FullFileSpace interfaz
					obj.FullFileSpace(column)
					break

				default: 
					return
			}
		// Rompemos switch byte
		break

		//Defalt si no es bytes o bit
		default: 
				return
	}


	if obj.err != nil {
		log.Println(obj.err)
	}

	if (obj.FileNativeType & disk) != 0 {

		defer obj.File.Close()

	}

	return
}






func (obj *Space ) Rmapspace(str_write string)(value int64, found bool){

	if int64(len(str_write)) > obj.Size_line{

		str_write = str_write[:obj.Size_line]

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


func (obj *Space ) Rindexspace(line int64)(value string, found bool){

	if  int64(len(obj.Index) ) > line {

		value, found = obj.Index[line] , true

	} else {

		value, found = "", false
	}

	return
}