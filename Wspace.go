package bd

import (
	"bytes"
	"log"

	"os"
	"sync/atomic"
//	"time"
)

func (sp *spaceFile) WriteColumnSpace(line int64, column map[string][]byte)int64{


	if line == -1 {
		log.Println("SizeFileLine",*sp.SizeFileLine)
		line = atomic.AddInt64(sp.SizeFileLine, 1)

	}

	if line > *sp.SizeFileLine {

		atomic.AddInt64(sp.SizeFileLine, line - *sp.SizeFileLine )
		
	}

	//ind -> index val -> valor
	for val := range sp.IndexSizeColumns {
			
		//value-> valor found -> Encontrado en el mapa
		_, found := column[val]
		//Si no encontramos la columna seguimos con el ciclo for
		if !found {
			
			continue

		}

		//Preformat por columnas
		function, exist := sp.Hooker[Preformat + val]
		if exist{

			column[val] = function(column[val])

		} else {

			//Preformat global
			function, exist = sp.Hooker[Preformat]
			if exist {

				column[val] = function(column[val])

			}
		}
	
		//Contamos el array de bytes
		var text_count = int64(len(column[val]))

		//Primer caso el texto es menor que el tamaño de la linea
		//En este caso añadimos un padding de espacios al final
		sizeColumn := sp.IndexSizeColumns[val][1] - sp.IndexSizeColumns[val][0]

		if text_count < sizeColumn {

			whitespace := bytes.Repeat( []byte(" ") , int(sizeColumn - text_count)) 
						
			column[val] = append(column[val] ,  whitespace... )
		}

		if text_count > sizeColumn {

			column[val] = column[val][:sizeColumn]
		}
		
		
		sp.File.WriteAt(column[val], sp.SizeLine * line + sp.IndexSizeColumns[val][0])

	
		
		//🔥🔥🔥Actualizamos ram
		if (sp.Space.FileNativeType & RamSearch) != 0 && sp.IndexSizeColumns[val][0] == 0  {

			sp.updateRamMap(column[val], line)

		}

		/*
		if (sp.FileNativeType & RamIndex) != 0 && sp.IndexSizeColumns[val][0] == 0 {

			sp.updateRamIndex(column[val], line)

		}
		*/

}
//	log.Println("tiempo de escritura simultanea: ", time.Since(test).Nanoseconds())
    

	return line
}


func (sp *spaceFile) WriteListBitSpace(line int64, column map[string][]byte){
	

	sp.Lock()
	
	sp.File , _ = os.OpenFile(sp.Url , os.O_RDWR | os.O_CREATE, 0666)
	
	defer sp.File.Close()
	defer sp.Unlock()
	
	var byteLine int64 =  line / 8

		
	bufferBit := make([]byte , 1 )
	

	_ , err := sp.File.ReadAt(bufferBit , byteLine)
	if err != nil{

		bufferBit = []byte{0}

	}


	var bitLine int64 =  line - ((line / 8) * 8)



	for val, ind := range sp.IndexSizeColumns {

		if ind[0] == 0 {

			switch string(column[val]) {

			case "on":
				writeBit(bitLine ,true , bufferBit )
			case "off":
				writeBit(bitLine ,false , bufferBit )
			}
			
			break
		}
	}

	sp.File.WriteAt(bufferBit, byteLine )	
	

	return
}





func (sp *spaceFile) WriteEmptyDirSpace(line int64, column map[string][]byte){

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