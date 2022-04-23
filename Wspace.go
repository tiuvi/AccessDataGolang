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

		
				if WB.Line == -1 {
			
					WB.Line = atomic.AddInt64(WB.SizeFileLine, 1)
			
				}
			
				if WB.Line > *WB.SizeFileLine {
			
					atomic.AddInt64(WB.SizeFileLine, WB.Line - *WB.SizeFileLine )
					
				}

				if WB.Hooker != nil {

					WB.hookerPreFormatPointer(WB.Buffer,WB.ColumnName)

				}
				
				WB.spacePaddingPointer(WB.Buffer , size)

				_ , err := WB.File.WriteAt(*WB.Buffer , WB.lenFields + (WB.SizeLine * WB.Line) + size[0])
				if err != nil{

					log.Println("Error de buffer en Wspace")

				}
				return &WB.Line
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

			//Añade una linea.
			if WB.Line == -1 {
		
				WB.Line = atomic.AddInt64(WB.SizeFileLine, 1)
		
			}
			//Añade lineas hasta llegar a la linea indicada.
			if WB.Line > *WB.SizeFileLine {
		
				atomic.AddInt64(WB.SizeFileLine, WB.Line - *WB.SizeFileLine )
				
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
					
					_ , err := WB.File.WriteAt(bufBytes, WB.lenFields + (WB.SizeLine * WB.Line) + size[0])
					if err != nil{

						log.Println("Error de buffer en Wspace")
	
					}
					return &WB.Line
				}
			}
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

				size, found := WB.IndexSizeFields[CHAN.ColName]
				if found {
				
					if WB.Hooker != nil {
						
						WB.hookerPreFormatPointer(&CHAN.Buffer, CHAN.ColName)
					
					}
			
					WB.spacePaddingPointer(&CHAN.Buffer , size)
						
				
					_, err := WB.File.WriteAt(CHAN.Buffer, WB.lenFields + WB.SizeLine * CHAN.Line + size[0])
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

		WB.Lock()
		defer WB.Unlock()

		for colName, buffer := range WB.BufferMap {

			if WB.IndexSizeFields != nil && WB.WRangues != nil {

				size , found := WB.IndexSizeFields[colName]
				if found {
				
					WB.WriteIndexSizeField(WB.ColumnName,size,*WB.WRangues,WB.Buffer )
					continue
				}
				
			}
	
			if WB.IndexSizeColumns != nil  && WB.WLines != nil { 

				size , found := WB.IndexSizeColumns[colName]
				if found {

					if WB.Line == -1 {
				
						WB.Line = atomic.AddInt64(WB.SizeFileLine, 1)
				
					}
				
					if WB.Line > *WB.SizeFileLine {
				
						atomic.AddInt64(WB.SizeFileLine, WB.Line - *WB.SizeFileLine )
						
					}


					
					var byteLine int64 =  WB.Line / 8

						
					bufferBit := make([]byte , 1 )
					

					_ , err := WB.File.ReadAt(bufferBit , WB.lenFields + byteLine + size[0])
					if err != nil{

						bufferBit = []byte{0}
						err = nil
					}


					var bitLine int64 =  WB.Line - ((WB.Line / 8) * 8)


					switch string(buffer) {

					case "on":
						writeBit(bitLine ,true , bufferBit )
					case "off":
						writeBit(bitLine ,false , bufferBit )
					}
					
					_ , err = WB.File.WriteAt(bufferBit, WB.lenFields + byteLine + size[0] )	
					if err != nil{

						log.Println("Error de buffer en Wspace")

					}
					return &WB.Line
				}
			}
		}
		return nil
	}



	//Buffer de bytes
	if CheckFileTypeBuffer(WB.typeBuff , BuffBytes ){

		if WB.IndexSizeFields != nil && WB.WRangues != nil {

			size, found := WB.IndexSizeFields[WB.ColumnName]
			if found {
				_ = size

			}
		}

		if WB.IndexSizeColumns != nil  && WB.WLines != nil { 

			size , found := WB.IndexSizeColumns[WB.ColumnName]
			if found {
				_ = size

			}
		}

	}

	//Buffer de bytes
	if CheckFileTypeBuffer(WB.typeBuff , BuffChan ){
	
		for CHAN := range WB.Channel {

			if WB.IndexSizeFields != nil && CHAN.WRangues != nil {

				size, found := WB.IndexSizeFields[CHAN.ColName]
				if found {
					_ = size

				}
			}

			if WB.IndexSizeColumns != nil  && CHAN.WLines != nil { 

				size , found := WB.IndexSizeColumns[CHAN.ColName]
				if found {
					_ = size

				}
			}
		}

	}

	return nil
}
