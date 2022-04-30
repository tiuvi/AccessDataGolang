package bd

import(
	"log"
	"time"
	"strings"
)


//Abrimos el espacio en disco
func (obj *space ) ospaceDisk( name string, folder []string)*spaceFile {
	
	if obj.spaceErrors != nil && obj.logTimeOpenFile{
		defer obj.LogDeferTimeMemoryDefault(time.Now())

		//defer obj.LogDeferTimeMemory(obj.dir,time.Now())

	}
	var folderString string
	if len(folder) > 0 {

		folderString = strings.Join(folder, "/")
		folderString += "/"
	}
	url := obj.dir + folderString + name + "." + obj.extension


	spacef , found := diskSpace.diskFile[url]
	if !found {

		//Creamos una nueva referencia a spaceFile
		spacef = obj.newSpaceFile(folderString, name)
		//Unico bloqueo cuando se abre el archivo para mantener la atomicidad
		diskSpace.Lock()
		//Guardamos nuestra referencia al archivo en el mapa
		diskSpace.diskFile[url] = spacef
		//Quitamos el cerrojo de la estructura diskSpace
		diskSpace.Unlock()

	}
	
	return spacef
}


//Timer que cierra cuando hay mas de 10 000 archivos abiertos
//Maximos archivos Ubuntu 1 048 576 Files  --  700 MB Ram 
func dacTimerCloserDeferFile(){

	//Numero de archivos maximos
	fileOpen := 10000
	//Especifica cada cuanto se ejecuta la funcion
	tikectDb := time.Tick(time.Duration(300) * time.Second)
	for range tikectDb {

		if len(deferSpace.deferFile) == 0 {
			log.Println("0 DeferFile Map")
			continue
		}

		if len(deferSpace.info) == 0{
			log.Println("0 DeferFile Slice")
			continue
		}

		//Bloqueamos el mapa para operaciones en paralelo
		deferSpace.Lock()
		//Recorremos el slice con la informacion
		for len(deferSpace.info) > fileOpen  {

			//Buscamos el primer elemento del slice añadido
			value, found := deferSpace.deferFile[deferSpace.info[0].name]
			//Si encontramos el elemento
			if found {
				//Cerramos el descriptor del archivo
				err := value.file.Close()
				if err != nil {
					log.Println(err)
				}
			}
	
			//Borrado del elemento del mapa y de la variable
			delete(deferSpace.deferFile  ,  deferSpace.info[0].name)
			//Borramos el primer elemento del slice que hemos borrado
			deferSpace.info = deferSpace.info[1:]
		}
		//Quitamos el candado al mapa
		deferSpace.Unlock()
	}
}


//Cierra todos los archivos abiertos cada X segundos , minutos, horas, dias....
func dacTimerCloserDiskFile(){

	//Cada cuanto se ejecuta la funcion
	tikectDiskFile := time.Tick(time.Duration(300) * time.Second)
	for range tikectDiskFile {
	
		if len(diskSpace.diskFile) == 0 {
			log.Println("0 diskFile")
			continue
		}

		//Abrimos el candado para operaciones de mapas concurrentes
		diskSpace.Lock()
		
		//Recorremos todo el mapa 
		for ind, value := range diskSpace.diskFile   {
			//Cerramos todos los archivos
			err := value.file.Close()
			if err != nil {
				log.Println(err)
			}
			
			//Borrado del elemento del mapa y de la variable
			delete(diskSpace.diskFile  ,  ind)
			
		}
		//Quitamos el candado al mapa
		diskSpace.Unlock()

	}

}