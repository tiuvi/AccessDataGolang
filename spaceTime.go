package bd

import (
	"fmt"
	"strings"
	"time"
)

//Abrimos el espacio en disco
func (obj *space) ospaceDisk(name string, folder []string) *spaceFile {

	if EDAC &&
		obj.logTimeOpenFile && !obj.isErrorFile {
		defer obj.NewLogDeferTimeMemory("diskFile", time.Now())
	}

	var folderString string
	if len(folder) > 0 {

		folderString = strings.Join(folder, "/")
		folderString += "/"
	}
	url := obj.dir + folderString + name + "." + obj.extension

	spacef, found := diskSpace.diskFile[url]
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

func (obj *space) ospaceDeferDisk(name string, folder []string) *spaceFile {

	if EDAC &&
		obj.logTimeOpenFile && !obj.isErrorFile {
		defer obj.NewLogDeferTimeMemory("deferDiskFile", time.Now())
	}

	var folderString string
	if len(folder) > 0 {

		folderString = strings.Join(folder, "/")
		folderString += "/"
	}
	url := strings.Join([]string{obj.dir, folderString, name, ".", obj.extension}, "")

	spacef, found := deferSpace.deferFile[url]
	if !found {

		//Creamos una nueva referencia a spaceFile
		spacef = obj.newSpaceFile(folderString, name)
		//Unico bloqueo cuando se abre el archivo para mantener la atomicidad
		deferSpace.Lock()
		//Creamos una nueva referencia a spaceFile
		deferSpace.deferFile[url] = spacef
		//Añadimos un elemento a la cola de array para su posterior
		//eleiminacion del mapa en orden
		deferSpace.info = append(deferSpace.info, deferFileInfo{url, time.Now()})
		//Quitamos el cerrojo de la estructura diskSpace
		deferSpace.Unlock()

	}
	return spacef
}

func (obj *space) ospacePermDisk(name string, folder []string) *spaceFile {

	if EDAC &&
		obj.logTimeOpenFile && !obj.isErrorFile {
		defer obj.NewLogDeferTimeMemory("permDiskFile", time.Now())
	}

	var folderString string
	if len(folder) > 0 {

		folderString = strings.Join(folder, "/")
		folderString += "/"
	}
	url := strings.Join([]string{obj.dir, folderString, name, ".", obj.extension}, "")

	spacef, found := permSpace.permDisk[url]
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

//Timer que cierra cuando hay mas de 10 000 archivos abiertos
//Maximos archivos Ubuntu 1 048 576 Files  --  700 MB Ram
func (LDAC *lDAC) dacTimerCloserDeferFile() {

	//Especifica cada cuanto se ejecuta la funcion
	tikectDb := time.Tick(time.Duration(LDAC.timeEventDeferFile) * time.Second)
	for range tikectDb {

		if len(deferSpace.deferFile) == 0 || len(deferSpace.info) == 0 {

			continue
		}


		//Bloqueamos el mapa para operaciones en paralelo
		deferSpace.Lock()

		//Recorremos el slice con la informacion
		for len(deferSpace.info) > LDAC.fileOpenDeferFile {

			//Buscamos el primer elemento del slice añadido
			value, found := deferSpace.deferFile[deferSpace.info[0].name]

			//Si encontramos el elemento
			if found {

				//Cerramos el descriptor del archivo
				err := value.file.Close()
				if err != nil && EDAC && 
				LDAC.ELDACF( true,"Error al cerrar el archivo en el mapa defer global. \n\r" + fmt.Sprintln(err)){}

			}

			//Borrado del elemento del mapa y de la variable
			delete(deferSpace.deferFile, deferSpace.info[0].name)

			//Borramos el primer elemento del slice que hemos borrado
			deferSpace.info = deferSpace.info[1:]

		}

		//Quitamos el candado al mapa
		deferSpace.Unlock()
	}
}



//Cierra todos los archivos abiertos cada X segundos , minutos, horas, dias....
func (LDAC *lDAC) dacTimerCloserDiskFile() {

	//Cada cuanto se ejecuta la funcion
	tikectDiskFile := time.Tick(time.Duration(LDAC.timeEventDiskFile) * time.Second)
	for range tikectDiskFile {

		if len(diskSpace.diskFile) == 0 {
		
			continue
		}

		//Abrimos el candado para operaciones de mapas concurrentes
		diskSpace.Lock()

		//Recorremos todo el mapa
		for ind, value := range diskSpace.diskFile {

			//Cerramos todos los archivos
			err := value.file.Close()
			if err != nil && EDAC && 
			LDAC.ELDACF( true,"Error al cerrar el archivo en el mapa defer global. \n\r" + fmt.Sprintln(err)){}

			//Borrado del elemento del mapa y de la variable
			delete(diskSpace.diskFile, ind)

		}
		
		//Quitamos el candado al mapa
		diskSpace.Unlock()

	}

}
