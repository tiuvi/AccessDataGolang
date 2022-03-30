package bd

import(
	"log"
	"time"
)


//Timer que cierra cuando hay mas de 10 000 archivos abiertos
//Maximos archivos Ubuntu 1 048 576 Files  --  700 MB Ram 
func dacTimerCloserDeferFile(){

	//Numero de archivos maximos
	fileOpen := 10000
	//Especifica cada cuanto se ejecuta la funcion
	tikectDb := time.Tick(time.Duration(30) * time.Second)
	for range tikectDb {

		if len(deferSpace.DeferFile) == 0 {
			log.Println("0 DeferFile Map")
			continue
		}

		if len(deferSpace.Info) == 0{
			log.Println("0 DeferFile Slice")
			continue
		}

		//Bloqueamos el mapa para operaciones en paralelo
		deferSpace.Lock()
		//Recorremos el slice con la informacion
		for len(deferSpace.Info) > fileOpen  {

			//Buscamos el primer elemento del slice añadido
			value, found := deferSpace.DeferFile[deferSpace.Info[0].Name]
			//Si encontramos el elemento
			if found {
				//Cerramos el descriptor del archivo
				err := value.File.Close()
				if err != nil {
					log.Println(err)
				}
			}
	
			//Borrado del elemento del mapa y de la variable
			delete(deferSpace.DeferFile  ,  deferSpace.Info[0].Name)
			//Borramos el primer elemento del slice que hemos borrado
			deferSpace.Info = deferSpace.Info[1:]
		}
		//Quitamos el candado al mapa
		deferSpace.Unlock()
	}
}


//Cierra todos los archivos abiertos cada X segundos , minutos, horas, dias....
func dacTimerCloserDiskFile(){

	//Cada cuanto se ejecuta la funcion
	tikectDiskFile := time.Tick(time.Duration(30) * time.Second)
	for range tikectDiskFile {
	
		if len(diskSpace.DiskFile) == 0 {
			log.Println("0 DiskFile")
			continue
		}

		//Abrimos el candado para operaciones de mapas concurrentes
		diskSpace.Lock()
		
		//Recorremos todo el mapa 
		for ind, value := range diskSpace.DiskFile   {
			//Cerramos todos los archivos
			err := value.File.Close()
			if err != nil {
				log.Println(err)
			}
			
			//Borrado del elemento del mapa y de la variable
			delete(diskSpace.DiskFile  ,  ind)
			
		}
		//Quitamos el candado al mapa
		diskSpace.Unlock()

	}

}