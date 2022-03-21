package bd

import (
	"bytes"
	"log"

)

func (sp *Space) WriteColumnSpace(line int64, column map[string][]byte){

	//Actualizamos el numero de lineas
	sp.complete_file_lines(line)

	//Actualizamos el tamaño del archivo
	sp.update_size_file(line)

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
		
		//if text_count < sp.Size_line {
		if text_count < sp.IndexSizeColumns[val][1] {
			//whitespace := bytes.Repeat( []byte{32} , int(sp.Size_line - text_count)) 
			//column[val] = append(column[val] ,  whitespace... )

			whitespace := bytes.Repeat( []byte(" ") , int(sp.IndexSizeColumns[val][1] - text_count)) 
						
			column[val] = append(column[val] ,  whitespace... )
		}

		//Segundo caso el string es mayor que el tamaño de la linea
		//	if text_count > sp.Size_line {
		if text_count > sp.IndexSizeColumns[val][1] {

			//column[val] = column[val][:sp.Size_line]
			column[val] = column[val][:sp.IndexSizeColumns[val][1]]
		}
		
	



		sp.File.WriteAt(column[val], sp.Size_line * line + sp.IndexSizeColumns[val][0])
		
		
		//🔥🔥🔥Actualizamos ram
		if (sp.FileNativeType & RamSearch) != 0 && sp.IndexSizeColumns[val][0] == 0  {

			sp.updateRamMap(column[val], line)

		}
		
		if (sp.FileNativeType & RamIndex) != 0 && sp.IndexSizeColumns[val][0] == 0 {

			sp.updateRamIndex(column[val], line)

		}

	}

}


func (sp *Space) WriteListBitSpace(line int64, column map[string][]byte){
	
	var byteLine int64 =  line / 8

	sp.complete_file_lines_bit(byteLine)



		bufferBit := make([]byte , 1 )
		_ , sp.err = sp.File.ReadAt(bufferBit , byteLine)
		
		var bitLine int64 =  line - ((line / 8) * 8)

		log.Println("Bitline: ", bitLine)
		//Antes
		log.Printf("Antes: %08b", bufferBit)

		for val, ind := range sp.IndexSizeColumns {

			if ind[0] == 0 {

				switch string(column[val]) {
	
				case "on":
					sp.writeBit(bitLine ,true , bufferBit )
				case "off":
					sp.writeBit(bitLine ,false , bufferBit )
				}
				
				break
			}
		}

		//Despues
		log.Printf("Despues: %08b", bufferBit)

		sp.File.WriteAt(bufferBit, byteLine )

		return
	}





