package bd

import(
	"time"
	"strings"
)


func (obj *Space ) ospaceCompilationFile()bool {


	if obj.Check {

		obj.SpaceErrors = GlobalError
		obj.LogNewError(MessageCopilation ,obj.Dir, `Las alertas en ospace son fatales con el tipo: MessageCopilation.`)
	}

	if len(obj.Dir) == 0 {

		obj.LogNewError(MessageCopilation ,obj.Dir, `Variable Dir vacia.`)

	}
	
	if len(obj.Extension) == 0 {

		obj.LogNewError(MessageCopilation ,obj.Dir, `Extension vacia.`)

	}

	if _ , found := extensionFile[obj.Extension]; !found{

		obj.LogNewError(MessageCopilation ,obj.Dir, `Extension no valida.`)

	}



	if len(obj.IndexSizeFieldsArray) > 0 {

		if obj.IndexSizeFields == nil {

			obj.IndexSizeFields = make(map[string][2]int64)

		}

		var afterValue int64
		for ind , value := range obj.IndexSizeFieldsArray {
			
			name := value.name
			len  := value.len

			if ind == 0 {

				obj.IndexSizeFields[name] = [2]int64{0 , len }

			}else {

				obj.IndexSizeFields[name] = [2]int64{ afterValue ,len + afterValue }

			}

			afterValue += len
		}
	}
	

	if len(obj.IndexSizeColumnsArray) > 0 {

		if obj.IndexSizeColumns == nil {

			obj.IndexSizeColumns = make(map[string][2]int64)

		}

		var afterValue int64
		for ind , value := range obj.IndexSizeColumnsArray {

			name := value.name
			len  := value.len
		
			if ind == 0 {

				obj.IndexSizeColumns[name] = [2]int64{0 , len }

			}else {

				obj.IndexSizeColumns[name] = [2]int64{ afterValue ,len + afterValue}

			}

			afterValue += len
		}
	}



	obj.lenColumns = int64(len(obj.IndexSizeColumns))
	obj.lenFields  = int64(len(obj.IndexSizeFields))

	if obj.lenColumns == 0 && obj.lenFields == 0 {

		obj.LogNewError(MessageCopilation ,obj.Dir, `Iniciaste un espacio sin columnas y sin campos.`)
	
	}

	if _ , found := obj.IndexSizeColumns["buffer"]; found {

		obj.LogNewError(MessageCopilation ,obj.Dir, `La palabra buffer en columnas esta reservada para el uso del programa.`)

	}

	if _ , found := obj.IndexSizeFields["buffer"]; found {

		obj.LogNewError(MessageCopilation ,obj.Dir, `La palabra buffer en campos esta reservada para el uso del programa.`)

	}

	checkMap := make(map[string]bool)

	if obj.lenFields != 0 {

		obj.lenFields       = 0
		var checkSizeFields int64 = 0

		for name, val := range obj.IndexSizeFields{

			if found := checkMap[name]; found{

				obj.LogNewError(MessageCopilation ,obj.Dir, "El campo: " + name +" coincide con el campo: " + name)

			}

			checkMap[name] = true

			calcSizeLine := (val[1] - val[0])
			if calcSizeLine <= 0 {

				obj.LogNewError(MessageCopilation ,obj.Dir, `Los fields no pueden tener un tamaño igual o inferior a cero.`)

			}

			obj.lenFields += calcSizeLine

			if val[1] >= checkSizeFields {

				checkSizeFields = val[1]
			}
		}

		if checkSizeFields != int64(obj.lenFields){
		
			obj.LogNewError(MessageCopilation ,obj.Dir, `Los campos estan mal escritos, Ejemplo: field1: 0,20; field2:20,30`)

		}

	}



	//Actualizamos el valor del ancho de la linea
	if obj.lenColumns != 0 {

		var checkSizeColumns int64 = 0

		for name , val := range obj.IndexSizeColumns {

			if obj.IndexSizeFields != nil {

				if found := checkMap[name]; found{

					obj.LogNewError(MessageCopilation ,obj.Dir, "El campo: " + name +" coincide con la columna: " + name)

				}
			}
	
			calcSizeLine := (val[1] - val[0])
			if calcSizeLine <= 0 {

				obj.LogNewError(MessageCopilation ,obj.Dir, "Las columnas no pueden tener un tamaño igual o inferior a cero." )

			}

			obj.SizeLine += calcSizeLine

			if val[1] >= checkSizeColumns {

				checkSizeColumns = val[1]
			}
			
		}
		
		if checkSizeColumns != obj.SizeLine {
			
			obj.LogNewError(MessageCopilation ,obj.Dir, "Las columnas estan mal escritas, Ejemplo: column1: 0,20; column2:20,30;" )
		
		}

	}
	

	
	//Lectura de archivos de byte
	if  obj.Extension == DacByte {

		obj.FileCoding = Byte
		obj.compilation = true
		return true
	}


	//bdisk: Lista de bit en un archivo disk
	if obj.Extension == DacBit {

		obj.FileCoding  = Bit
		obj.compilation = true
		return true
	}

	obj.LogNewError(MessageCopilation ,obj.Dir,
	"No se han encontrado coincidencias con las extensiones de archivo predeterminadas." )

	return false
}



//Variable Global de ospaceDisk
var diskSpace = &spaceDisk{
	DiskFile: make(map[string]*spaceFile),
}



//Abrimos el espacio en disco
func (obj *Space ) ospaceDisk( name string, folder []string)*spaceFile {
	
	var folderString string
	if len(folder) > 0 {

		folderString = strings.Join(folder, "/")
		folderString += "/"
	}
	url := obj.Dir + folderString + name + "." + obj.Extension


	spacef , found := diskSpace.DiskFile[url]
	if !found {

		//Creamos una nueva referencia a spaceFile
		spacef = obj.newSpaceFile(folderString, name)
		//Unico bloqueo cuando se abre el archivo para mantener la atomicidad
		diskSpace.Lock()
		//Guardamos nuestra referencia al archivo en el mapa
		diskSpace.DiskFile[url] = spacef
		//Quitamos el cerrojo de la estructura diskSpace
		diskSpace.Unlock()

	}
	
	return spacef
}

var deferSpace = &spaceDeferDisk{
	DeferFile: make(map[string]*spaceFile),
	Info: make([]deferFileInfo,0),
}

func (obj *Space ) ospaceDeferDisk(name string, folder []string)*spaceFile{

	var folderString string
	if len(folder) > 0 {

		folderString = strings.Join(folder, "/")
		folderString += "/"
	}
	url := obj.Dir + folderString + name + "." + obj.Extension

	
	spacef , found := deferSpace.DeferFile[url]
	if !found {

		//Creamos una nueva referencia a spaceFile
		spacef = obj.newSpaceFile(folderString, name)
		//Unico bloqueo cuando se abre el archivo para mantener la atomicidad
		deferSpace.Lock()
		//Creamos una nueva referencia a spaceFile
		deferSpace.DeferFile[url] = spacef
		//Añadimos un elemento a la cola de array para su posterior
		//eleiminacion del mapa en orden
		deferSpace.Info = append( deferSpace.Info , deferFileInfo{url,time.Now()})
		//Quitamos el cerrojo de la estructura diskSpace
		deferSpace.Unlock()

	}
	return spacef
}

//Variable Global de ospaceDisk
var permSpace = &spacePermDisk{
	PermDisk: make(map[string]*spaceFile),
}

func (obj *Space ) ospacePermDisk( name string,folder []string)*spaceFile{

	var folderString string
	if len(folder) > 0 {

		folderString = strings.Join(folder, "/")
		folderString += "/"
	}
	url := obj.Dir + folderString + name + "." + obj.Extension

	spacef , found := permSpace.PermDisk[url]
	if !found {

		//Creamos una nueva referencia a spaceFile
		spacef = obj.newSpaceFile(folderString, name)
		//Unico bloqueo cuando se abre el archivo para mantener la atomicidad
		permSpace.Lock()
		//Creamos una nueva referencia a spaceFile
		permSpace.PermDisk[url] = spacef
		//Quitamos el cerrojo de la estructura diskSpace
		permSpace.Unlock()
	

	}

	return spacef

}


