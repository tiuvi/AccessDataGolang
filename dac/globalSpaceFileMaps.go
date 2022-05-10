package dac

import (

	"os"
	"strings"
	"log"
	"time"
)

//Abrimos el archivo
func (obj *space) ospaceOpenFile(name string, folder []string) (spacef *spaceFile) {

	if EDAC &&
		obj.logTimeOpenFile && !obj.isErrorFile {
		defer obj.NewLogDeferTimeMemory("diskFile", time.Now())
	}

	var folderString string
	if len(folder) > 0 {

		folderString = strings.Join(folder, "/")
		folderString += "/"
	}


	//Creamos una nueva referencia a spaceFile
	return obj.newSpaceFile(folderString, name)

}

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

	url := strings.Join([]string{obj.dir, folderString, name, ".", obj.extension}, "")
	
	diskSpace.RLock()
	spacef, found := diskSpace.diskFile[url]
	diskSpace.RUnlock()

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

	deferSpace.RLock()
	spacef, found := deferSpace.deferFile[url]
	deferSpace.RUnlock()
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
	
	permSpace.RLock()
	spacef, found := permSpace.permDisk[url]
	permSpace.RUnlock()
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

		diskSpace.RLock()
		_ , found := diskSpace.diskFile[SF.url]
		diskSpace.RUnlock()

		if found {
			
			diskSpace.Lock()
			delete(diskSpace.diskFile, SF.url)
			diskSpace.Unlock()
	
		}

	}

	if SF.fileNativeType == deferDisk {
		
		deferSpace.RLock()
		_ , found := deferSpace.deferFile[SF.url]
		deferSpace.RUnlock()
		if found {
			
			deferSpace.Lock()
			delete(deferSpace.deferFile , SF.url)
			deferSpace.Unlock()

		}
	}

	if SF.fileNativeType == permDisk {

		permSpace.RLock()
		_ , found := permSpace.permDisk[SF.url]
		permSpace.RUnlock()
		if found {
			
			permSpace.Lock()
			delete(permSpace.permDisk , SF.url)
			permSpace.Unlock()
	
		}
	}

	err := os.Remove(SF.url)
	if EDAC &&
	SF.ECSFD( err != nil ,"Error al borrar el archivo"){}
	
}



func(PSF *PublicSpaceFile) DeleteFile(){


	if PSF.PublicSpaceCache != nil {

		PSF.PublicSpaceCache.DeleteFileCache(PSF)

	}
	err := os.Remove(PSF.url)
	if EDAC &&
	PSF.ECSFD( err != nil ,"Error al borrar el archivo"){}
}










/**************************************************************************/
//Public space file maps disk
/**************************************************************************/

func SeeGlobalCache(){

	log.Println(globalCache.cache)
}

func GetCache(directory string)*PublicSpaceCache{

	globalCache.RLock()
	PublicSpaceCache, found := globalCache.cache[directory]
	globalCache.RUnlock()
	if !found {
		return nil
	}
	return PublicSpaceCache
}

func  InsertCache(directory string)*PublicSpaceCache{

	globalCache.RLock()
	cache, found := globalCache.cache[directory]
	globalCache.RUnlock()
	if found {

		return cache
	}

	if !found {

		newCache := NewCache()
		globalCache.Lock()
		globalCache.cache[directory] = newCache
		globalCache.Unlock()
		return newCache
	}

	return nil
}

func DeleteCache(directory string){

	globalCache.RLock()
	_ , found := globalCache.cache[directory]
	globalCache.RUnlock()
	if found {
		
		globalCache.Lock()
		delete(globalCache.cache, directory)
		globalCache.Unlock()

	}
	
}


func NewCache()*PublicSpaceCache{

	return &PublicSpaceCache{
			cache: make(map[string]*PublicSpaceFile),	
	}

}

func (PSD *PublicSpaceCache) GetFileCache(url string)*PublicSpaceFile{

	PSD.RLock()
	PublicSpaceFile, found := PSD.cache[url]
	PSD.RUnlock()
	if !found {
		return nil
	}

	return PublicSpaceFile
}

func (PSD *PublicSpaceCache) InsertFileCache(PSF *PublicSpaceFile){

	PSD.RLock()
	_, found := PSD.cache[PSF.url]
	PSD.RUnlock()
	if !found {

		PSD.Lock()
		PSD.cache[PSF.url] = PSF
		PSD.Unlock()
	}
}

func (PSD *PublicSpaceCache) UpdateFileCache(PSF *PublicSpaceFile){

	PSD.RLock()
	_, found := PSD.cache[PSF.url]
	PSD.RUnlock()
	if found {

		PSD.Lock()
		PSD.cache[PSF.url] = PSF
		PSD.Unlock()
		
	}
}


func (PSD *PublicSpaceCache) DeleteFileCache(PSF *PublicSpaceFile){

	PSD.RLock()
	_ , found := PSD.cache[PSF.url]
	PSD.RUnlock()
	if found {
		
		PSD.Lock()
		delete(PSD.cache, PSF.url)
		PSD.Unlock()

	}
	
}
