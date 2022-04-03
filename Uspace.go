package bd

	import "log"


//Esta funcion escribe en el numero de linea requerido
//Pasandole un string en forma de bytes
func (obj *Space ) Wspace(line int64, column map[string][]byte)int64{


	spaceFile := obj.OSpace()
	log.Println("Write space: " , *spaceFile.SizeFileLine)

	//Difenciar archivos de bit de archivos de byte
	switch obj.FileCoding {

		//Primer caso archivo de bit
		case Bit:

			//Segun el tipo de archivo aplicamos una funcion diferente
			switch obj.FileTipeBit {
				//Caso lista de bits
				case ListBit:	
				spaceFile.WriteListBitSpace(line, column)
				break
			}

		//Segundo caso archivos de byte
		case Byte:
	
			//Segun el tipo de archivo aplicamos una funcion diferente
			switch obj.FileTipeByte {
				
				case OneColumn,MultiColumn:
					//Funcion OneColumn interfaz
					return spaceFile.WriteColumnSpace(line, column)
			}

		case Dir:
			switch obj.FileTypeDir {

				case EmptyDir:
					//Añadir modificacion de url desde aqui
					spaceFile.WriteEmptyDirSpace(line, column)
					break

			}
	}
	return -1
}



func (obj *Space) Rspace (column *Buffer){



	spaceFile := obj.OSpace()

	
	//Difenciar archivos de bit de archivos de byte
	switch obj.FileCoding {

		//Primer caso archivo de bit
		case Bit:

			//Segun el tipo de archivo aplicamos una funcion diferente
			switch obj.FileTipeBit {
			
				case ListBit:
					//Funcion ListBit interfaz
					spaceFile.ListBitSpace(column)
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
					spaceFile.OneColumnSpace(column)
					break
				case MultiColumn:
					//Funcion Multicolumn interfaz
					spaceFile.MultiColumnSpace(column)
					break
				case FullFile:
					//Añadir modificacion de url desde aqui
					spaceFile.FullFileSpace(column)
					break

				default: 
					return
			}
		// Rompemos switch byte
		break

		case Dir:
			switch obj.FileTypeDir {

				case EmptyDir:
					//Añadir modificacion de url desde aqui
					spaceFile.ReadEmptyDirSpace(column)
					break

				default: 
					return
			}
		break

		//Defalt si no es bytes o bit
		default: 
				return
	}


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