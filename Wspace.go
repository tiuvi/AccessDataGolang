package dac

import (
	"fmt"
	"sync/atomic"
	"time"
)

func (WB *WBuffer) writeByteSpace()*int64{


	if EDAC &&
	WB.logTimeWriteFile && !WB.isErrorFile {
		defer WB.NewLogDeferTimeMemorySF("writeBytes", time.Now())
	}

	//Buffer de bytes
	if checkFileTypeBuffer(WB.typeBuff , buffBytes ){


		if WB.indexSizeFields != nil && WB.wRangues != nil {


			size, found := WB.indexSizeFields[WB.columnName]
			if found {

				WB.WriteIndexSizeField(WB.columnName,size,*WB.wRangues,WB.buffer )
				
				return nil
			}

		}
		
		if WB.indexSizeColumns != nil  && WB.wLines != nil { 

			size , found := WB.indexSizeColumns[WB.columnName]
			if found {

				line := WB.line
				if line == -1 {
			
					line = atomic.AddInt64(WB.sizeFileLine, 1)
			
				}
			
				if line > *WB.sizeFileLine {
			
					atomic.AddInt64(WB.sizeFileLine, line - *WB.sizeFileLine )
					
				}

				if WB.hooker != nil {

					WB.hookerPreFormatPointer(WB.buffer,WB.columnName)

				}
				
				WB.spacePaddingPointer(WB.buffer , size)

				_ , err := WB.file.WriteAt(*WB.buffer , WB.lenFields + (WB.sizeLine * line) + size[0])
				if err != nil && EDAC && 
				WB.ECSFD( true,"Error al escribir en el archivo \n\r" + fmt.Sprintln(err)){}
				
				return &line
			}
		}
		return nil
	}


	if checkFileTypeBuffer(WB.typeBuff , buffMap ){

		if WB.indexSizeFields != nil && WB.wRangues != nil {

			for columnName, fieldBuffer := range WB.bufferMap {
		
				size, found := WB.indexSizeFields[columnName]
				if found {
				
					WB.WriteIndexSizeField(columnName,size,*WB.wRangues,&fieldBuffer )
					
				}
			}
		}

		if WB.indexSizeColumns != nil  && WB.wLines != nil { 

			line := WB.line
			//Añade una linea.
			if line == -1 {
		
				line = atomic.AddInt64(WB.sizeFileLine, 1)
		
			}
			//Añade lineas hasta llegar a la linea indicada.
			if line > *WB.sizeFileLine {
		
				atomic.AddInt64(WB.sizeFileLine, line - *WB.sizeFileLine )
				
			}
			
		
			//ind -> index val -> valor
			for colName , bufBytes := range WB.bufferMap {

				//value-> valor found -> Encontrado en el mapa
				size, found := WB.indexSizeColumns[colName]
				//Si no encontramos la columna seguimos con el ciclo for
				if found {
					
					if WB.hooker != nil {
			
						WB.hookerPreFormatPointer(&bufBytes,colName)
					}
					
					WB.spacePaddingPointer(&bufBytes , size)
					
					_ , err := WB.file.WriteAt(bufBytes, WB.lenFields + (WB.sizeLine * line) + size[0])
					if err != nil && EDAC && 
					WB.ECSFD( true,"Error al escribir en el archivo \n\r" + fmt.Sprintln(err)){}
					
				}
			}
			return &line
		}
	}




	if checkFileTypeBuffer(WB.typeBuff , buffChan ){

		for CHAN := range WB.channel {


			if WB.indexSizeFields != nil && CHAN.wRangues != nil {

				size, found := WB.indexSizeFields[CHAN.colName]
				if found {
					
					WB.WriteIndexSizeField(CHAN.colName,size,*CHAN.wRangues, &CHAN.buffer)
					continue
				}
				
			}

			if WB.indexSizeColumns != nil  && CHAN.wLines != nil { 

				size, found := WB.indexSizeColumns[CHAN.colName]
				if found {
				
					if WB.hooker != nil {
						
						WB.hookerPreFormatPointer(&CHAN.buffer, CHAN.colName)
					
					}
			
					WB.spacePaddingPointer(&CHAN.buffer , size)
						
		
					_, err := WB.file.WriteAt(CHAN.buffer, WB.lenFields + (WB.sizeLine * CHAN.line) + size[0])
					if err != nil && EDAC && 
					WB.ECSFD( true,"Error al escribir en el archivo \n\r" + fmt.Sprintln(err)){}
				
					continue
				}
			}
		}
	}

	return nil
}


func (WB *WBuffer) writeBitSpace()*int64{

	if EDAC &&
	WB.logTimeWriteFile && !WB.isErrorFile {
		defer WB.NewLogDeferTimeMemory("writeBits", time.Now())
	}
	
	if checkFileTypeBuffer(WB.typeBuff , buffMap ){


		if WB.indexSizeFields != nil && WB.wRangues != nil {

			for columnName, fieldBuffer := range WB.bufferMap {
		
				size, found := WB.indexSizeFields[columnName]
				if found {
				
					WB.WriteIndexSizeField(columnName,size,*WB.wRangues,&fieldBuffer )
					
				}
			}
		}

		

			if WB.indexSizeColumns != nil  && WB.wLines != nil { 

				line := WB.line
				//Añade una linea.
				if line == -1 {
			
					line = atomic.AddInt64(WB.sizeFileLine, 1)
			
				}
				//Añade lineas hasta llegar a la linea indicada.
				if line > *WB.sizeFileLine {
			
					atomic.AddInt64(WB.sizeFileLine, line - *WB.sizeFileLine )
					
				}

				var byteLine int64 =  line / 8
				var bitLine  int64  =  line % 8 
				bufferBit := make([]byte , 1 )
				WB.Lock()
				defer WB.Unlock()


				for colName, buffer := range WB.bufferMap {

					size , found := WB.indexSizeColumns[colName]
					if found {

						_ , err := WB.file.ReadAt(bufferBit , WB.lenFields + (byteLine * WB.sizeLine) + size[0])
						if err != nil{

							bufferBit = []byte{0}
							err = nil
						}
					

						//log.Printf("Binary Before: %b", bufferBit ) 

						switch string(buffer) {

						case "on":
							bufferBit = writeBit(bitLine ,true , bufferBit )
						case "off":
							bufferBit = writeBit(bitLine ,false , bufferBit )
						}

						//log.Printf("Binary After: %b", bufferBit ) 

						_ , err = WB.file.WriteAt(bufferBit, WB.lenFields + (byteLine * WB.sizeLine) + size[0] )	
						if err != nil && EDAC && 
						WB.ECSFD( true,"Error al escribir en el archivo de bits \n\r" + fmt.Sprintln(err)){}
				
					}
				}
				return &line
			}
		return nil
	}



	//Buffer de bytes
	if checkFileTypeBuffer(WB.typeBuff , buffBytes ){

		if WB.indexSizeFields != nil && WB.wRangues != nil {


			size, found := WB.indexSizeFields[WB.columnName]
			if found {

				WB.WriteIndexSizeField(WB.columnName,size,*WB.wRangues,WB.buffer )
				
				return nil
			}

		}

		if WB.indexSizeColumns != nil  && WB.wLines != nil { 

			size , found := WB.indexSizeColumns[WB.columnName]
			if found {

				line := WB.line
				if line == -1 {
			
					line = atomic.AddInt64(WB.sizeFileLine, 1)
			
				}
			
				if line > *WB.sizeFileLine {
			
					atomic.AddInt64(WB.sizeFileLine, line - *WB.sizeFileLine )
					
				}

				var byteLine int64 =  line / 8
				var bitLine  int64  =  line % 8 
				bufferBit := make([]byte , 1 )

				WB.Lock()
				defer WB.Unlock()
				
				_ , err := WB.file.ReadAt(bufferBit , WB.lenFields + (byteLine * WB.sizeLine) + size[0])
				if err != nil{

					bufferBit = []byte{0}
					err = nil
				}
			

				//log.Printf("Binary Before: %b", bufferBit ) 

				switch string(*WB.buffer) {

				case "on":
					bufferBit = writeBit(bitLine ,true , bufferBit )
				case "off":
					bufferBit = writeBit(bitLine ,false , bufferBit )
				}

			//	log.Printf("Binary After: %b", bufferBit ) 

				_ , err = WB.file.WriteAt(bufferBit, WB.lenFields + (byteLine * WB.sizeLine) + size[0] )	
				if err != nil && EDAC && 
				WB.ECSFD( true,"Error al escribir en el archivo de bits \n\r" + fmt.Sprintln(err)){}
				
				return &line
			}
		}
		return nil
	}

	//Buffer de bytes
	if checkFileTypeBuffer(WB.typeBuff , buffChan ){
	
		for CHAN := range WB.channel {

		
			if WB.indexSizeFields != nil && CHAN.wRangues != nil {

				size, found := WB.indexSizeFields[CHAN.colName]
				if found {
					
					WB.WriteIndexSizeField(CHAN.colName,size,*CHAN.wRangues, &CHAN.buffer)
					continue
				}
				
			}

			if WB.indexSizeColumns != nil  && CHAN.wLines != nil { 

				size , found := WB.indexSizeColumns[CHAN.colName]
				if found {
				
					var byteLine int64 =  CHAN.line / 8
					var bitLine  int64  =  CHAN.line % 8 
					bufferBit := make([]byte , 1 )
					WB.Lock()
					
					_ , err := WB.file.ReadAt(bufferBit , WB.lenFields + (byteLine * WB.sizeLine) + size[0])
					if err != nil{
	
						bufferBit = []byte{0}
						err = nil
					}

					switch string(CHAN.buffer) {

						case "on":
							bufferBit = writeBit(bitLine ,true , bufferBit )
						case "off":
							bufferBit = writeBit(bitLine ,false , bufferBit )
					}

					_, err = WB.file.WriteAt(bufferBit , WB.lenFields + (WB.sizeLine * byteLine) + size[0])
					if err != nil && EDAC && 
					WB.ECSFD( true,"Error al escribir en el archivo de bits \n\r" + fmt.Sprintln(err)){}
				

					WB.Unlock()

					continue
				}
			}
		}

	}

	return nil
}
