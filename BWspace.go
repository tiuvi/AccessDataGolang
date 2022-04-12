package bd

import (
	"log"
	"sync/atomic"
	)

type WChanBuf struct{
	Line int64
	ColName string
	Buffer 	[]byte
}

type WBuffer struct {
	*spaceFile
	Line int64
	ColumnName string
	SizeLine  int64
	typeBuff FileTypeBuffer

	Buffer *[]byte
	BufferMap map[string][]byte
	Channel chan WChanBuf
}



func (sp *spaceFile) BWspaceBuff(line int64 ,columnName string, buff []byte)(*WBuffer){




	if sp.IndexSizeColumns != nil {

		_, found := sp.IndexSizeColumns[columnName]
		if !found {
			
			if sp.IndexSizeFields != nil {
	
				_, found := sp.IndexSizeFields[columnName]
				if !found {
		
					log.Fatalln("Bspace.go - 222; funcion: BWspaceBuff" +
					"La columna: " + columnName + " no existe en ese archivo",sp.Url)
					return nil
				}
			}
		}
	}

	return &WBuffer{
		spaceFile: sp,
		typeBuff: BuffBytes,
		Line: line,
		ColumnName: columnName,
		Buffer: &buff,
	}

}


func (sp *spaceFile) BWspaceBuffMap(line int64 , buff map[string][]byte)(*WBuffer){

	for columnName := range buff {
	
		if sp.IndexSizeColumns != nil {

			_, found := sp.IndexSizeColumns[columnName]
			if !found {
				
				if sp.IndexSizeFields != nil {
		
					_, found := sp.IndexSizeFields[columnName]
					if !found {
			
						log.Fatalln("Bspace.go - 222; funcion: BWspaceBuff" +
						"La columna: " + columnName + " no existe en ese archivo",sp.Url)
						return nil
					}
				}
			}
		}
	}

   return &WBuffer{
		spaceFile: sp,
		typeBuff: BuffMap,
	    Line: line,
	    BufferMap: buff,
   }

}

func (sp *spaceFile)BWChanBuf()(*WBuffer){

	return &WBuffer{
		spaceFile: sp,
		typeBuff: BuffChan,
		Channel: make(chan WChanBuf, 0),
	}
}


func (WBuffer *WBuffer)BWspaceSendchan(line int64, columnName string , buf *[]byte)int64{



	if WBuffer.IndexSizeColumns != nil {

		_, found := WBuffer.IndexSizeColumns[columnName]
		if !found {
			
			if WBuffer.IndexSizeFields != nil {
	
				_, found := WBuffer.IndexSizeFields[columnName]
				if !found {
		
					log.Fatalln("Bspace.go - 222; funcion: BWspaceBuff" +
					"La columna: " + columnName + " no existe en ese archivo",WBuffer.Url)
					return -2
				}
			}
		}
	}

	//Actualización de campos sin lineas.
	if line != -2 {

		if line == -1 {
		
			line = atomic.AddInt64(WBuffer.SizeFileLine, 1)

		}

		if line > *WBuffer.SizeFileLine {

			atomic.AddInt64(WBuffer.SizeFileLine, line - *WBuffer.SizeFileLine )
			
		}
	}

	WBuffer.Buffer = buf
	WBuffer.Channel <- WChanBuf{line,columnName, *WBuffer.Buffer }

	return line
}

func (WBuffer *WBuffer) BWspaceClosechan(){
	
	close(WBuffer.Channel)
}

