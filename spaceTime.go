package dac

import (
	"fmt"
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
