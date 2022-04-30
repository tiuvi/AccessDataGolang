package bd

import(
	"time"
	"strings"
)


func (obj *space ) ospaceCompilationFile()bool {


	if obj.spaceErrors != nil {

	
		obj.ErrorSpaceDefault(MessageCopilation , `Las alertas en ospace son fatales con el tipo: MessageCopilation.`)

		if len(obj.dir) == 0 {

			obj.ErrorSpaceDefault(MessageCopilation , `Variable dir vacia.`)
	
		}
		
		if  len(obj.extension) == 0 {
	
			obj.ErrorSpaceDefault(MessageCopilation , `extension vacia.`)
	
		}
	
		if _ , found := extensionFile[obj.extension]; !found{
	
			obj.ErrorSpaceDefault(MessageCopilation , `extension no valida.`)
	
		}
	}


	




	if len(obj.indexSizeFieldsArray) > 0 {

		if obj.indexSizeFields == nil {

			obj.indexSizeFields = make(map[string][2]int64)

		}

		var afterValue int64
		for ind , value := range obj.indexSizeFieldsArray {
			
			name := value.name
			len  := value.len

			if ind == 0 {

				obj.indexSizeFields[name] = [2]int64{0 , len }

			}else {

				obj.indexSizeFields[name] = [2]int64{ afterValue ,len + afterValue }

			}

			afterValue += len
		}
	}
	

	if len(obj.indexSizeColumnsArray) > 0 {

		if obj.indexSizeColumns == nil {

			obj.indexSizeColumns = make(map[string][2]int64)

		}

		var afterValue int64
		for ind , value := range obj.indexSizeColumnsArray {

			name := value.name
			len  := value.len
		
			if ind == 0 {

				obj.indexSizeColumns[name] = [2]int64{0 , len }

			}else {

				obj.indexSizeColumns[name] = [2]int64{ afterValue ,len + afterValue}

			}

			afterValue += len
		}
	}



	obj.lenColumns = int64(len(obj.indexSizeColumns))
	obj.lenFields  = int64(len(obj.indexSizeFields))
	if obj.spaceErrors != nil {
		
		if obj.lenColumns == 0 && obj.lenFields == 0 {

			obj.ErrorSpaceDefault(MessageCopilation , `Iniciaste un espacio sin columnas y sin campos.`)
		
		}
	
		if _ , found := obj.indexSizeColumns["buffer"]; found {
	
			obj.ErrorSpaceDefault(MessageCopilation , `La palabra buffer en columnas esta reservada para el uso del programa.`)
	
		}
	
		if _ , found := obj.indexSizeFields["buffer"]; found {
	
			obj.ErrorSpaceDefault(MessageCopilation , `La palabra buffer en campos esta reservada para el uso del programa.`)
	
		}
	}


	checkMap := make(map[string]bool)

	if obj.lenFields != 0 {

		obj.lenFields       = 0
		var checkSizeFields int64 = 0

		for name, val := range obj.indexSizeFields{

			if obj.spaceErrors != nil {
				if found := checkMap[name]; found{

					obj.ErrorSpaceDefault(MessageCopilation , "El campo: " + name +" coincide con el campo: " + name)

				}
			}

			checkMap[name] = true

			calcSizeLine := (val[1] - val[0])
			if obj.spaceErrors != nil && calcSizeLine <= 0 {

				obj.ErrorSpaceDefault(MessageCopilation , `Los fields no pueden tener un tamaño igual o inferior a cero.`)

			}

			obj.lenFields += calcSizeLine

			if val[1] >= checkSizeFields {

				checkSizeFields = val[1]
			}
		}

		if obj.spaceErrors != nil && checkSizeFields != int64(obj.lenFields){
		
			obj.ErrorSpaceDefault(MessageCopilation , `Los campos estan mal escritos, Ejemplo: field1: 0,20; field2:20,30`)

		}

	}



	//Actualizamos el valor del ancho de la linea
	if obj.lenColumns != 0 {

		var checkSizeColumns int64 = 0

		for name , val := range obj.indexSizeColumns {

			if obj.indexSizeFields != nil {

				if obj.spaceErrors != nil {

					if found := checkMap[name]; found{

						obj.ErrorSpaceDefault(MessageCopilation , "El campo: " + name +" coincide con la columna: " + name)
	
					}
				}
			}
	
			calcSizeLine := (val[1] - val[0])
			if obj.spaceErrors != nil  && calcSizeLine <= 0 {

				obj.ErrorSpaceDefault(MessageCopilation , "Las columnas no pueden tener un tamaño igual o inferior a cero." )

			}

			obj.sizeLine += calcSizeLine

			if val[1] >= checkSizeColumns {

				checkSizeColumns = val[1]
			}
			
		}
		
		if obj.spaceErrors != nil  && checkSizeColumns != obj.sizeLine {
			
			obj.ErrorSpaceDefault(MessageCopilation , "Las columnas estan mal escritas, Ejemplo: column1: 0,20; column2:20,30;" )
		
		}

	}
	

	
	//Lectura de archivos de byte
	if  obj.extension == dacByte {

		obj.fileCoding = bytes
		obj.compilation = true
		return true
	}


	//bdisk: Lista de bit en un archivo disk
	if obj.extension == dacBit {

		obj.fileCoding  = bit
		obj.compilation = true
		return true
	}

	if obj.spaceErrors != nil {
		obj.ErrorSpaceDefault(MessageCopilation ,
			"No se han encontrado coincidencias con las extensiones de archivo predeterminadas." )
		
	}
	return false
}











func (obj *space ) ospaceDeferDisk(name string, folder []string)*spaceFile{

	var folderString string
	if len(folder) > 0 {

		folderString = strings.Join(folder, "/")
		folderString += "/"
	}
	url := strings.Join([]string{ obj.dir , folderString , name , "." , obj.extension }, "") 
	
	
	spacef , found := deferSpace.deferFile[url]
	if !found {

		//Creamos una nueva referencia a spaceFile
		spacef = obj.newSpaceFile(folderString, name)
		//Unico bloqueo cuando se abre el archivo para mantener la atomicidad
		deferSpace.Lock()
		//Creamos una nueva referencia a spaceFile
		deferSpace.deferFile[url] = spacef
		//Añadimos un elemento a la cola de array para su posterior
		//eleiminacion del mapa en orden
		deferSpace.info = append( deferSpace.info , deferFileInfo{url,time.Now()})
		//Quitamos el cerrojo de la estructura diskSpace
		deferSpace.Unlock()

	}
	return spacef
}



func (obj *space ) ospacePermDisk( name string,folder []string)*spaceFile{

	var folderString string
	if len(folder) > 0 {

		folderString = strings.Join(folder, "/")
		folderString += "/"
	}
	url := strings.Join([]string{ obj.dir , folderString , name , "." , obj.extension }, "") 
	

	spacef , found := permSpace.permDisk[url]
	if !found {

		//Creamos una nueva referencia a spaceFile
		spacef = obj.newSpaceFile(folderString, name)
		//Unico bloqueo cuando se abre el archivo para mantener la atomicidad
		permSpace.Lock()
		//Creamos una nueva referencia a spaceFile
		permSpace.permDisk[url] = spacef
		//Quitamos el cerrojo de la estructura diskSpace
		permSpace.Unlock()
	

	}

	return spacef

}








