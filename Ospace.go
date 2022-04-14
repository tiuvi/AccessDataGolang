package bd

import(
	"log"
	"time"
)


func (obj *Space ) ospaceCompilationFile()bool {

	if len(obj.Dir) == 0 {

		log.Fatalln("Ruta directorio vacio en: ", obj.Dir, obj.Extension)

	}
	
	if len(obj.Extension) == 0 {

		log.Fatalln("Extension vacia en: ", obj.Dir, obj.Extension )

	}

	if _ , found := extensionFile[obj.Extension];  !found{

		log.Fatalln("Extension no valida en: ", obj.Dir, obj.Extension)

	}



	obj.lenColumns = int64(len(obj.IndexSizeColumns))
	obj.lenFields  = int64(len(obj.IndexSizeFields))

	if obj.lenColumns == 0 && obj.lenFields == 0{
		
		log.Fatalln("Iniciaste un archivo de acceso a datos sin columnas y sin campos.", obj.Dir)
	}

	if _ , found := obj.IndexSizeColumns["buffer"]; found {

		log.Fatalln("La palabra buffer esta reservada para el programa.", obj.Dir, obj.IndexSizeColumns)

	}

	if _ , found := obj.IndexSizeFields["buffer"]; found {

		log.Fatalln("La palabra buffer esta reservada para el programa.", obj.Dir, obj.IndexSizeColumns)

	}

	checkMap := make(map[string]bool)

	if obj.lenFields != 0 {

		obj.lenFields       = 0
		var checkSizeFields int64 = 0

		for name, val := range obj.IndexSizeFields{

			checkMap[name] = true

			calcSizeLine := (val[1] - val[0])
			if calcSizeLine <= 0 {
	
				log.Fatalln("Los fields no pueden tener un tamaño inferior a cero.",
				obj.Dir,obj.IndexSizeColumns)
			}

			obj.lenFields += calcSizeLine

			if val[1] >= checkSizeFields {

				checkSizeFields = val[1]
			}
		}

		if checkSizeFields != int64(obj.lenFields){
			log.Fatalln("Los campos estan mal escritos, Ejemplo: field1: 0,20; field2:20,30;", obj.Dir)
		}

	}



	//Actualizamos el valor del ancho de la linea
	if obj.lenColumns != 0 {

		var checkSizeColumns int64 = 0

		for name , val := range obj.IndexSizeColumns{

			if obj.IndexSizeFields != nil {

				if found := checkMap[name]; found{

					log.Fatalln("El campo: " + name +" coincide con la columna: " + name + " en: ",
					obj.Dir)
				
				}
			}
	
			calcSizeLine := (val[1] - val[0])
			if calcSizeLine <= 0 {

				log.Fatalln("Las columnas no pueden tener un tamaño inferior a cero.",
				obj.Dir,obj.IndexSizeColumns)
			}

			obj.SizeLine += calcSizeLine

			if val[1] >= checkSizeColumns {

				checkSizeColumns = val[1]
			}
			
		}
		
		if checkSizeColumns != obj.SizeLine {

			log.Fatalln("Las columnas estan mal escritos, Ejemplo: column1: 0,20; column2:20,30;", obj.Dir)
		
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

		obj.SizeLine        = 1
		obj.FileCoding      = Bit
		obj.compilation = true
		return true
	}


	log.Fatalln("Archivo: Ospace.go; Funcion: ospaceCompilationFile ; Linea:150 ;" +
	"No se han encontrado coincidencias de archivos con esas extensiones.", obj.Dir)
	return false
}



//Variable Global de ospaceDisk
var diskSpace = &spaceDisk{
	DiskFile: make(map[string]*spaceFile),
}

//Abrimos el espacio en disco
func (obj *Space ) ospaceDisk(name string)*spaceFile {

	url := obj.Dir + name + "." + obj.Extension

	spacef , found := diskSpace.DiskFile[url]
	if !found {

		//Creamos una nueva referencia a spaceFile
		spacef = obj.newSpaceFile(url)
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

func (obj *Space ) ospaceDeferDisk(name string)*spaceFile{

	url := obj.Dir + name + "." + obj.Extension

	
	spacef , found := deferSpace.DeferFile[url]
	if !found {

		//Creamos una nueva referencia a spaceFile
		spacef = obj.newSpaceFile(url)
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

func (obj *Space ) ospacePermDisk(name string)*spaceFile{

	url := obj.Dir + name + "." + obj.Extension

	spacef , found := permSpace.PermDisk[url]
	if !found {

		//Creamos una nueva referencia a spaceFile
		spacef = obj.newSpaceFile(url)
		//Unico bloqueo cuando se abre el archivo para mantener la atomicidad
		permSpace.Lock()
		//Creamos una nueva referencia a spaceFile
		permSpace.PermDisk[url] = spacef
		//Quitamos el cerrojo de la estructura diskSpace
		permSpace.Unlock()
	

	}

	return spacef

}


