package bd

import (
	"log"
	"sync/atomic"
)

func (WB *WBuffer) writeByteSpace()*int64{


	//Buffer de bytes
	if CheckFileTypeBuffer(WB.typeBuff , BuffBytes ){


		if WB.IndexSizeFields != nil && WB.WRangues != nil {


			size, found := WB.IndexSizeFields[WB.ColumnName]
			if found {

				WB.WriteIndexSizeField(WB.ColumnName,size,*WB.WRangues,WB.Buffer )
				
				return nil
			}

		}
		
		if WB.IndexSizeColumns != nil  && WB.WLines != nil { 

			size , found := WB.IndexSizeColumns[WB.ColumnName]
			if found {

				line := WB.Line
				if line == -1 {
			
					line = atomic.AddInt64(WB.SizeFileLine, 1)
			
				}
			
				if line > *WB.SizeFileLine {
			
					atomic.AddInt64(WB.SizeFileLine, line - *WB.SizeFileLine )
					
				}

				if WB.Hooker != nil {

					WB.hookerPreFormatPointer(WB.Buffer,WB.ColumnName)

				}
				
				WB.spacePaddingPointer(WB.Buffer , size)

				_ , err := WB.File.WriteAt(*WB.Buffer , WB.lenFields + (WB.SizeLine * line) + size[0])
				if err != nil{

					log.Println("Error de buffer en Wspace")

				}
				return &line
			}
		}
		return nil
	}


	if CheckFileTypeBuffer(WB.typeBuff , BuffMap ){

		if WB.IndexSizeFields != nil && WB.WRangues != nil {

			for columnName, fieldBuffer := range WB.BufferMap {
		
				size, found := WB.IndexSizeFields[columnName]
				if found {
				
					WB.WriteIndexSizeField(columnName,size,*WB.WRangues,&fieldBuffer )
					
				}
			}
		}

		if WB.IndexSizeColumns != nil  && WB.WLines != nil { 

			line := WB.Line
			//Añade una linea.
			if line == -1 {
		
				line = atomic.AddInt64(WB.SizeFileLine, 1)
		
			}
			//Añade lineas hasta llegar a la linea indicada.
			if line > *WB.SizeFileLine {
		
				atomic.AddInt64(WB.SizeFileLine, line - *WB.SizeFileLine )
				
			}
			
		
			//ind -> index val -> valor
			for colName , bufBytes := range WB.BufferMap {

				//value-> valor found -> Encontrado en el mapa
				size, found := WB.IndexSizeColumns[colName]
				//Si no encontramos la columna seguimos con el ciclo for
				if found {
					
					if WB.Hooker != nil {
			
						WB.hookerPreFormatPointer(&bufBytes,colName)
					}
					
					WB.spacePaddingPointer(&bufBytes , size)
					
					_ , err := WB.File.WriteAt(bufBytes, WB.lenFields + (WB.SizeLine * line) + size[0])
					if err != nil{

						log.Println("Error de buffer en Wspace")
	
					}
				}
			}
			return &line
		}
	}




	if CheckFileTypeBuffer(WB.typeBuff , BuffChan ){

		for CHAN := range WB.Channel {


			if WB.IndexSizeFields != nil && CHAN.WRangues != nil {

				size, found := WB.IndexSizeFields[CHAN.ColName]
				if found {
					
					WB.WriteIndexSizeField(CHAN.ColName,size,*CHAN.WRangues, &CHAN.Buffer)
					continue
				}
				
			}

			if WB.IndexSizeColumns != nil  && CHAN.WLines != nil { 

				size, found := WB.IndexSizeColumns[CHAN.ColName]
				if found {
				
					if WB.Hooker != nil {
						
						WB.hookerPreFormatPointer(&CHAN.Buffer, CHAN.ColName)
					
					}
			
					WB.spacePaddingPointer(&CHAN.Buffer , size)
						
		
					_, err := WB.File.WriteAt(CHAN.Buffer, WB.lenFields + (WB.SizeLine * CHAN.Line) + size[0])
					if err != nil{

						log.Println("Error de buffer en Wspace")

					}
					continue
				}
			}
		}
	}

	return nil
}


func (WB *WBuffer) writeBitSpace()*int64{
	
	if CheckFileTypeBuffer(WB.typeBuff , BuffMap ){


		if WB.IndexSizeFields != nil && WB.WRangues != nil {

			for columnName, fieldBuffer := range WB.BufferMap {
		
				size, found := WB.IndexSizeFields[columnName]
				if found {
				
					WB.WriteIndexSizeField(columnName,size,*WB.WRangues,&fieldBuffer )
					
				}
			}
		}

		

			if WB.IndexSizeColumns != nil  && WB.WLines != nil { 

				line := WB.Line
				//Añade una linea.
				if line == -1 {
			
					line = atomic.AddInt64(WB.SizeFileLine, 1)
			
				}
				//Añade lineas hasta llegar a la linea indicada.
				if line > *WB.SizeFileLine {
			
					atomic.AddInt64(WB.SizeFileLine, line - *WB.SizeFileLine )
					
				}

				var byteLine int64 =  line / 8
				var bitLine  int64  =  line % 8 
				bufferBit := make([]byte , 1 )
				WB.Lock()
				defer WB.Unlock()


				for colName, buffer := range WB.BufferMap {

					size , found := WB.IndexSizeColumns[colName]
					if found {

						_ , err := WB.File.ReadAt(bufferBit , WB.lenFields + (byteLine * WB.SizeLine) + size[0])
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

						_ , err = WB.File.WriteAt(bufferBit, WB.lenFields + (byteLine * WB.SizeLine) + size[0] )	
						if err != nil{

							log.Println("Error de buffer en Wspace")

						}
					}
				}
				return &line
			}
		return nil
	}



	//Buffer de bytes
	if CheckFileTypeBuffer(WB.typeBuff , BuffBytes ){

		if WB.IndexSizeFields != nil && WB.WRangues != nil {


			size, found := WB.IndexSizeFields[WB.ColumnName]
			if found {

				WB.WriteIndexSizeField(WB.ColumnName,size,*WB.WRangues,WB.Buffer )
				
				return nil
			}

		}

		if WB.IndexSizeColumns != nil  && WB.WLines != nil { 

			size , found := WB.IndexSizeColumns[WB.ColumnName]
			if found {

				line := WB.Line
				if line == -1 {
			
					line = atomic.AddInt64(WB.SizeFileLine, 1)
			
				}
			
				if line > *WB.SizeFileLine {
			
					atomic.AddInt64(WB.SizeFileLine, line - *WB.SizeFileLine )
					
				}

				var byteLine int64 =  line / 8
				var bitLine  int64  =  line % 8 
				bufferBit := make([]byte , 1 )

				WB.Lock()
				defer WB.Unlock()
				
				_ , err := WB.File.ReadAt(bufferBit , WB.lenFields + (byteLine * WB.SizeLine) + size[0])
				if err != nil{

					bufferBit = []byte{0}
					err = nil
				}
			

				//log.Printf("Binary Before: %b", bufferBit ) 

				switch string(*WB.Buffer) {

				case "on":
					bufferBit = writeBit(bitLine ,true , bufferBit )
				case "off":
					bufferBit = writeBit(bitLine ,false , bufferBit )
				}

			//	log.Printf("Binary After: %b", bufferBit ) 

				_ , err = WB.File.WriteAt(bufferBit, WB.lenFields + (byteLine * WB.SizeLine) + size[0] )	
				if err != nil{

					log.Println("Error de buffer en Wspace")

				}
				return &line
			}
		}
		return nil
	}

	//Buffer de bytes
	if CheckFileTypeBuffer(WB.typeBuff , BuffChan ){
	
		for CHAN := range WB.Channel {

		
			if WB.IndexSizeFields != nil && CHAN.WRangues != nil {

				size, found := WB.IndexSizeFields[CHAN.ColName]
				if found {
					
					WB.WriteIndexSizeField(CHAN.ColName,size,*CHAN.WRangues, &CHAN.Buffer)
					continue
				}
				
			}

			if WB.IndexSizeColumns != nil  && CHAN.WLines != nil { 

				size , found := WB.IndexSizeColumns[CHAN.ColName]
				if found {
				
					var byteLine int64 =  CHAN.Line / 8
					var bitLine  int64  =  CHAN.Line % 8 
					bufferBit := make([]byte , 1 )
					WB.Lock()
					
					_ , err := WB.File.ReadAt(bufferBit , WB.lenFields + (byteLine * WB.SizeLine) + size[0])
					if err != nil{
	
						bufferBit = []byte{0}
						err = nil
					}

					switch string(CHAN.Buffer) {

						case "on":
							bufferBit = writeBit(bitLine ,true , bufferBit )
						case "off":
							bufferBit = writeBit(bitLine ,false , bufferBit )
					}

					_, err = WB.File.WriteAt(bufferBit , WB.lenFields + (WB.SizeLine * byteLine) + size[0])
					if err != nil{

						log.Println("Error de buffer en Wspace")

					}

					WB.Unlock()

					continue
				}
			}
		}

	}

	return nil
}
