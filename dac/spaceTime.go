package dac

import (
	"log"
	"time"
)

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
			_, found := deferSpace.deferFile[deferSpace.info[0].name]

			//Si encontramos el elemento
			if found {

				//Borrado del elemento del mapa y de la variable
				delete(deferSpace.deferFile, deferSpace.info[0].name)

				//Borramos el primer elemento del slice que hemos borrado
				deferSpace.info = deferSpace.info[1:]
			}

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
		for urlDiskFile := range diskSpace.diskFile {

			//Borrado del elemento del mapa y de la variable
			delete(diskSpace.diskFile, urlDiskFile)

		}
		
		//Quitamos el candado al mapa
		diskSpace.Unlock()

	}

}



func dacTimerCloserGlobalCache(){

	//Cada cuanto se ejecuta la funcion
	tikecGlobalCache := time.Tick(24 * time.Hour)
	for range tikecGlobalCache {

		if len(globalCache.cache) == 0 {
		
			continue
		}

		//Abrimos el candado para operaciones de mapas concurrentes
		globalCache.Lock()

		//Recorremos todo el mapa
		for directory := range globalCache.cache {

			//Borrado del elemento del mapa y de la variable
			delete(globalCache.cache, directory)

		}
		log.Println(globalCache.cache)
		//Quitamos el candado al mapa
		globalCache.Unlock()

	}
}