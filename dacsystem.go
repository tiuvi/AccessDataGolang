package bd

import(
	"log"
	"time"
)

func dacTimerCloserDeferFile(){

		//Archivos abiertos
	//1 048 576 maximos
	fileOpen := 10000
	//tikectDb := time.Tick(time.Duration(1) * time.Hour)
	tikectDb := time.Tick(time.Duration(24) * time.Second)
	for range tikectDb {

		log.Println("-----","     ","-----")
		log.Println("Nuevo Bucle tikectDb")
		log.Println("-----","     ","-----")

		if len(mSpace.Descriptor) == 0 {
			log.Println("Descriptor parada ")
			continue
		}

		if len(mSpace.Information) == 0{
			log.Println("Information parada ")
			continue
		}

		mSpace.Lock()
		

		for len(mSpace.Information) > fileOpen  {

			value, found := mSpace.Descriptor[mSpace.Information[0].Name]
			if found {
				err := value.Close()
				if err != nil {
					log.Println(err)
				}
			}
	
			//Borrado del elemento del mapa y de la variable
			delete(mSpace.Descriptor  ,  mSpace.Information[0].Name)
			//mSpace.Information = append(mSpace.Information[:ind], mSpace.Information[ind+1:]... )
			mSpace.Information = mSpace.Information[1:]
		}
		mSpace.Unlock()
	}
}


func dacTimerCloserDiskFile(){

		//Archivos de disco
		tikectDiskFile := time.Tick(time.Duration(10) * time.Second)
		for range tikectDiskFile {
	
			log.Println("-----","     ","-----")
			log.Println("Nuevo Bucle tikectDiskFile")
			log.Println("-----","     ","-----")
	
			if len(diskSpace.DiskFile) == 0 {
				log.Println("DiskFile parada")
				continue
			}
	
			diskSpace.Lock()
			
				for ind, value := range diskSpace.DiskFile   {
	
					err := value.Close()
					if err != nil {
						log.Println(err)
					}
					
					//Borrado del elemento del mapa y de la variable
					delete(diskSpace.DiskFile  ,  ind)
	
				}
			
			diskSpace.Unlock()
	
		}
		
}