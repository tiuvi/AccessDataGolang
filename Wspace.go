package bd

import (
	"bytes"
	"log"

	"os"
	"sync/atomic"
	//	"time"
)

func (sp *spaceFile) WriteColumnSpace(buf *WBuffer)int64{

	if CheckFileTypeBuffer(buf.typeBuff , BuffMap ){

		if sp.IndexSizeFields != nil {

			for columnName, bufMap := range buf.BufferMap {
		
				_, found := sp.IndexSizeFields[columnName]
				if found {
				
					sp.WriteIndexSizeField(columnName, &bufMap)
					
				}
			}
		}

		//Actualización de campos sin lineas.
		if buf.Line == -2 {
	
			return -2
	
		}

		//Añade una linea.
		if buf.Line == -1 {
	
			buf.Line = atomic.AddInt64(sp.SizeFileLine, 1)
	
		}
		//Añade lineas hasta llegar a la linea indicada.
		if buf.Line > *sp.SizeFileLine {
	
			atomic.AddInt64(sp.SizeFileLine, buf.Line - *sp.SizeFileLine )
			
		}
		
		//ind -> index val -> valor
		for val := range sp.IndexSizeColumns {
				
			//value-> valor found -> Encontrado en el mapa
			bufBytes, found := buf.BufferMap[val]
			//Si no encontramos la columna seguimos con el ciclo for
			if !found {
				
				continue
	
			}
			if sp.Hooker != nil {
	
				sp.hookerPreFormatPointer(&bufBytes,val)
			}
			
		
			//Contamos el array de bytes
			var text_count = int64(len(bufBytes))
	
			//Primer caso el texto es menor que el tamaño de la linea
			//En este caso añadimos un padding de espacios al final
			sizeColumn := sp.IndexSizeColumns[val][1] - sp.IndexSizeColumns[val][0]
	
			if text_count < sizeColumn {
	
				whitespace := bytes.Repeat( []byte(" ") , int(sizeColumn - text_count)) 
							
				bufBytes = append(bufBytes ,  whitespace... )
			}
	
			if text_count > sizeColumn {
	
				bufBytes = bufBytes[:sizeColumn]
			}
			
			
			sp.File.WriteAt(bufBytes, sp.SizeLine * buf.Line + sp.IndexSizeColumns[val][0])
	
		
			
			//🔥🔥🔥Actualizamos ram
			if (sp.FileNativeType & RamSearch) != 0 && sp.IndexSizeColumns[val][0] == 0  {
	
				sp.updateRamMap(bufBytes, buf.Line)
	
			}
	
			
			if (sp.FileNativeType & RamIndex) != 0 && sp.IndexSizeColumns[val][0] == 0 {
	
				sp.updateRamIndex(bufBytes, buf.Line)
	
			}
			
	
		}
	}

	//Buffer de bytes
	if CheckFileTypeBuffer(buf.typeBuff , BuffBytes ){


		if sp.IndexSizeFields != nil {

			_, found := sp.IndexSizeFields[buf.ColumnName]
			if found {
				sp.WriteIndexSizeField(buf.ColumnName,buf.Buffer)
				return -2
			}

		}
		
		//Actualización de campos sin lineas.
		if buf.Line == -2 {

			return -2
	
		}

		if buf.Line == -1 {
	
			buf.Line = atomic.AddInt64(sp.SizeFileLine, 1)
	
		}
	
		if buf.Line > *sp.SizeFileLine {
	
			atomic.AddInt64(sp.SizeFileLine, buf.Line - *sp.SizeFileLine )
			
		}

		if sp.Hooker != nil {

			sp.hookerPreFormatPointer(buf.Buffer,buf.ColumnName)

		}

		//Contamos el array de bytes
		var text_count = int64(len(*buf.Buffer))
		val := buf.ColumnName
		//Primer caso el texto es menor que el tamaño de la linea
		//En este caso añadimos un padding de espacios al final
		sizeColumn := sp.IndexSizeColumns[val][1] - sp.IndexSizeColumns[val][0]

		if text_count < sizeColumn {

			whitespace := bytes.Repeat( []byte(" ") , int(sizeColumn - text_count)) 
						
			*buf.Buffer = append(*buf.Buffer ,  whitespace... )
		}

		if text_count > sizeColumn {

			*buf.Buffer = (*buf.Buffer)[:sizeColumn]
		}
		
		
		sp.File.WriteAt(*buf.Buffer, sp.lenFields + (sp.SizeLine * buf.Line) + sp.IndexSizeColumns[val][0])
	
		
			
			//🔥🔥🔥Actualizamos ram
			if (sp.FileNativeType & RamSearch) != 0 && sp.IndexSizeColumns[val][0] == 0  {
	
				sp.updateRamMap(*buf.Buffer, buf.Line)
	
			}
	
			
			if (sp.FileNativeType & RamIndex) != 0 && sp.IndexSizeColumns[val][0] == 0 {
	
				sp.updateRamIndex(*buf.Buffer, buf.Line)
	
			}

	}


	if CheckFileTypeBuffer(buf.typeBuff , BuffChan ){

		for Ch := range buf.Channel {

			colName := Ch.ColName
			line    := Ch.Line
			buf     := &Ch.Buffer

			log.Println("writefieldschan")
			
			if sp.IndexSizeFields != nil {

				_, found := sp.IndexSizeFields[colName]
				if found {
				
					sp.WriteIndexSizeField(colName, buf)
					continue
				}
				
			}

			if line == -2 {

				continue

			}

			if sp.Hooker != nil {
				
				sp.hookerPreFormatPointer(buf, colName)
			
			}
	

			sp.columnSpacePadding(colName , buf)

		
			sp.File.WriteAt(*buf, sp.lenFields + sp.SizeLine * line + sp.IndexSizeColumns[colName][0])
			
		
			
			//🔥🔥🔥Actualizamos ram
			if (sp.FileNativeType & RamSearch) != 0 && sp.IndexSizeColumns[colName][0] == 0  {
	
				sp.updateRamMap(*buf, line)
	
			}
	
			
			if (sp.FileNativeType & RamIndex) != 0 && sp.IndexSizeColumns[colName][0] == 0 {
	
				sp.updateRamIndex(*buf, line)
	
			}
	
		}
	}

	return buf.Line
}


func (sp *spaceFile) WriteListBitSpace(buf *WBuffer)int64{
	

	sp.Lock()
	
	sp.File , _ = os.OpenFile(sp.Url , os.O_RDWR | os.O_CREATE, 0666)
	
	defer sp.File.Close()
	defer sp.Unlock()
	
	var byteLine int64 =  buf.Line / 8

		
	bufferBit := make([]byte , 1 )
	

	_ , err := sp.File.ReadAt(bufferBit , byteLine)
	if err != nil{

		bufferBit = []byte{0}

	}


	var bitLine int64 =  buf.Line - ((buf.Line / 8) * 8)



	for val, ind := range sp.IndexSizeColumns {

		if ind[0] == 0 {

			switch string(buf.BufferMap[val]) {

			case "on":
				writeBit(bitLine ,true , bufferBit )
			case "off":
				writeBit(bitLine ,false , bufferBit )
			}
			
			break
		}
	}

	sp.File.WriteAt(bufferBit, byteLine )	
	
	return buf.Line
}




/*
func (sp *spaceFile) WriteEmptyDirSpace(buf *WBuffer){

	var err error
	var value []byte
	var found bool

	_ , found = column["newBuffer"]
	if found {

		//sp.File, err = os.OpenFile(sp.Name + strconv.FormatInt(line,10) + sp.Extension , os.O_RDWR | os.O_CREATE | os.O_TRUNC, 0666)
		sp.File, err = os.OpenFile(sp.Url , os.O_RDWR | os.O_CREATE | os.O_TRUNC, 0666)
		if err != nil {

			log.Print(err)
		
		}
	}

	value , found = column["appendBuffer"]
	if found {

		//sp.File, err = os.OpenFile(sp.Name + strconv.FormatInt(line,10) + sp.Extension , os.O_RDWR | os.O_APPEND, 0666)
		sp.File, err = os.OpenFile(sp.Url , os.O_RDWR | os.O_APPEND, 0666)
		if err != nil {

			log.Print(err)
		
		}
		if _, err := sp.File.Write(value); 
		err != nil {

			log.Print(err)
	
		}
	}

	_ , found = column["endBuffer"]
	if found {

		defer sp.File.Close()

	}



}
*/