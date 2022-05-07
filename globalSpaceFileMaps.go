package dac

import (
	"strings"
	"time"
	"os"
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


func(SF *spaceFile) DeleteFile(){

	if SF.fileNativeType == disk {

		spaceFile , found := diskSpace.diskFile[SF.url]
		if found {
			
			diskSpace.Lock()
			spaceFile.file.Close()
			delete(diskSpace.diskFile, SF.url)
			diskSpace.Unlock()
	
		}

	}

	if SF.fileNativeType == deferDisk {

		spaceFile , found := deferSpace.deferFile[SF.url]
		if found {
			
			deferSpace.Lock()
			spaceFile.file.Close()
			delete(deferSpace.deferFile , SF.url)
			deferSpace.Unlock()

		}
	}

	if SF.fileNativeType == permDisk {

		spaceFile , found := permSpace.permDisk[SF.url]
		if found {
			
			permSpace.Lock()
			spaceFile.file.Close()
			delete(permSpace.permDisk , SF.url)
			permSpace.Unlock()
	
		}
	}

	err := os.Remove(SF.url)
	if EDAC &&
	SF.ECSFD( err != nil ,"Error al borrar el archivo"){}
	
}

func getFileDisk(url string)*spaceFile{

	spacef, found := diskSpace.diskFile[url]
	if !found {
		return nil
	}

	return spacef
}

func deleteFileDisk(url string){

	spaceFile , found := diskSpace.diskFile[url]
	if found {
		
		diskSpace.Lock()
		spaceFile.file.Close()
		delete(diskSpace.diskFile, url)
		diskSpace.Unlock()

	}
}
