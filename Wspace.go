package bd

import (
	"bytes"
	"log"
	"sync/atomic"
)

func (buf *WBuffer) writeByteSpace()int64{

	if CheckFileTypeBuffer(buf.typeBuff , BuffMap ){

		if buf.IndexSizeFields != nil {

			for columnName, bufMap := range buf.BufferMap {
		
				_, found := buf.IndexSizeFields[columnName]
				if found {
				
					buf.WriteIndexSizeField(columnName, &bufMap)
					
				}
			}
		}

		//Actualización de campos sin lineas.
		if buf.Line == -2 {
	
			return -2
	
		}

		//Añade una linea.
		if buf.Line == -1 {
	
			buf.Line = atomic.AddInt64(buf.SizeFileLine, 1)
	
		}
		//Añade lineas hasta llegar a la linea indicada.
		if buf.Line > *buf.SizeFileLine {
	
			atomic.AddInt64(buf.SizeFileLine, buf.Line - *buf.SizeFileLine )
			
		}
		
		//ind -> index val -> valor
		for val := range buf.IndexSizeColumns {
				
			//value-> valor found -> Encontrado en el mapa
			bufBytes, found := buf.BufferMap[val]
			//Si no encontramos la columna seguimos con el ciclo for
			if !found {
				
				continue
	
			}
			if buf.Hooker != nil {
	
				buf.hookerPreFormatPointer(&bufBytes,val)
			}
			
		
			//Contamos el array de bytes
			var text_count = int64(len(bufBytes))
	
			//Primer caso el texto es menor que el tamaño de la linea
			//En este caso añadimos un padding de espacios al final
			sizeColumn := buf.IndexSizeColumns[val][1] - buf.IndexSizeColumns[val][0]
	
			if text_count < sizeColumn {
	
				whitespace := bytes.Repeat( []byte(" ") , int(sizeColumn - text_count)) 
							
				bufBytes = append(bufBytes ,  whitespace... )
			}
	
			if text_count > sizeColumn {
	
				bufBytes = bufBytes[:sizeColumn]
			}
			
			
			buf.File.WriteAt(bufBytes, buf.lenFields + (buf.SizeLine * buf.Line) + buf.IndexSizeColumns[val][0])
	
		
		
			
	
		}
	}

	//Buffer de bytes
	if CheckFileTypeBuffer(buf.typeBuff , BuffBytes ){


		if buf.IndexSizeFields != nil {

			_, found := buf.IndexSizeFields[buf.ColumnName]
			if found {
				buf.WriteIndexSizeField(buf.ColumnName,buf.Buffer)
				return -2
			}

		}
		
		//Actualización de campos sin lineas.
		if buf.Line == -2 {

			return -2
	
		}

		if buf.Line == -1 {
	
			buf.Line = atomic.AddInt64(buf.SizeFileLine, 1)
	
		}
	
		if buf.Line > *buf.SizeFileLine {
	
			atomic.AddInt64(buf.SizeFileLine, buf.Line - *buf.SizeFileLine )
			
		}

		if buf.Hooker != nil {

			buf.hookerPreFormatPointer(buf.Buffer,buf.ColumnName)

		}

		//Contamos el array de bytes
		var text_count = int64(len(*buf.Buffer))
		val := buf.ColumnName
		//Primer caso el texto es menor que el tamaño de la linea
		//En este caso añadimos un padding de espacios al final
		sizeColumn := buf.IndexSizeColumns[val][1] - buf.IndexSizeColumns[val][0]

		if text_count < sizeColumn {

			whitespace := bytes.Repeat( []byte(" ") , int(sizeColumn - text_count)) 
						
			*buf.Buffer = append(*buf.Buffer ,  whitespace... )
		}

		if text_count > sizeColumn {

			*buf.Buffer = (*buf.Buffer)[:sizeColumn]
		}
		
		
		buf.File.WriteAt(*buf.Buffer, buf.lenFields + (buf.SizeLine * buf.Line) + buf.IndexSizeColumns[val][0])
	
		

	}


	if CheckFileTypeBuffer(buf.typeBuff , BuffChan ){

		for Ch := range buf.Channel {

			colName := Ch.ColName
			line    := Ch.Line
			bufferChan  := &Ch.Buffer

			log.Println("writefieldschan")
			
			if buf.IndexSizeFields != nil {

				_, found := buf.IndexSizeFields[colName]
				if found {
				
					buf.WriteIndexSizeField(colName, bufferChan)
					continue
				}
				
			}

			if line == -2 {

				continue

			}

			if buf.Hooker != nil {
				
				buf.hookerPreFormatPointer(bufferChan, colName)
			
			}
	

			buf.columnSpacePadding(colName , bufferChan)

		
			buf.File.WriteAt(*bufferChan, buf.lenFields + buf.SizeLine * line + buf.IndexSizeColumns[colName][0])
			
	
		}
	}

	return buf.Line
}


func (buf *WBuffer) writeBitSpace()int64{
	
	if CheckFileTypeBuffer(buf.typeBuff , BuffMap ){

		for colName, buffer := range buf.BufferMap {

			if buf.IndexSizeFields != nil {

				_, found := buf.IndexSizeFields[colName]
				if found {
				
					buf.WriteIndexSizeField(colName, &buffer)
					continue
				}
				
			}
	
			_ , found := buf.IndexSizeColumns[colName]
			if found {

				if buf.Line == -1 {
			
					buf.Line = atomic.AddInt64(buf.SizeFileLine, 1)
			
				}
			
				if buf.Line > *buf.SizeFileLine {
			
					atomic.AddInt64(buf.SizeFileLine, buf.Line - *buf.SizeFileLine )
					
				}

				buf.Lock()

				defer buf.Unlock()
				
				var byteLine int64 =  buf.Line / 8

					
				bufferBit := make([]byte , 1 )
				

				_ , err := buf.File.ReadAt(bufferBit , buf.lenFields + byteLine)
				if err != nil{

					bufferBit = []byte{0}

				}


				var bitLine int64 =  buf.Line - ((buf.Line / 8) * 8)


				switch string(buffer) {

				case "on":
					writeBit(bitLine ,true , bufferBit )
				case "off":
					writeBit(bitLine ,false , bufferBit )
				}
				
				buf.File.WriteAt(bufferBit, buf.lenFields + byteLine )	
				
			}
		}

	}

	return buf.Line
}

