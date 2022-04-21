package bd

import (
	"bytes"
	"log"
)



//writeBit: modifica un bit dentro de un byte activandolo o
//desactivandolo de izquierda a derecha.
// 1 2 3 4 5 6 7 8 9
// 0 0 0 0 0 0 0 0 0 
//#bd/core.go
func writeBit(id int64,turn bool ,bufferBit []byte)([]byte){

	 bitId := []uint8{128,64,32,16,8,4,2,1}

	var bufferuint uint8 = bufferBit[0]

	switch turn {
	
		case true:
			if (bufferuint & bitId[id]) == 0 {
				
				bufferuint = bufferuint | bitId[id]

			} 

		case false:
		
			if (bufferuint & bitId[id]) != 0 {

				bufferuint = bufferuint ^ bitId[id]
			
			} 

	}
	
	bufferBit[0] = bufferuint

	return bufferBit
	
}

//readBit: lee un unico bit de izquierda a derecha y te dice si 
//esta activado o no.
// 1 2 3 4 5 6 7 8 9
// 0 0 0 0 0 0 0 0 0 
//#bd/core.go
func readBit(id int64,bufferBit []byte)(bool){

	bitId := []uint8{128,64,32,16,8,4,2,1}

   var bufferuint uint8 = bufferBit[0]

	if (bufferuint & bitId[id]) != 0 {

		return true
	
	} 

	//if (bufferuint & bitId[id]) == 0 {}
	return false

}





//CheckBit: compara dos int64 y te dice si ese bit 
//esta activado a nivel de bytes.
//#bd/core.go
func CheckBit(base int64, compare int64)(bool){


	if (base & compare) != 0 {

		return true

	}
	return false
}




//hookerPreFormatPointer: Da preformato a la entrada de bytes en los 
//archivos tanto en columnas como fields.
//#bd/core.go
func (sp *spaceFile) hookerPreFormatPointer(bufByte *[]byte,colName string){

	//Preformat por columnas
	function, exist := sp.Hooker[Preformat + colName]
	if exist{

		 function(bufByte)

	} else {

		//Preformat global
		function, exist = sp.Hooker[Preformat]
		if exist {

			function(bufByte)

		}
	}
}


//hookerPostFormatPointer: Da postformato a la salida de bytes en los 
//archivos tanto en columnas como fields.
//#bd/core.go
func (sp *spaceFile) hookerPostFormatPointer(bufByte *[]byte,colName string){

	//Postformat por columnas
	function, exist := sp.Hooker[Postformat + colName]
	if exist{

		function(bufByte)

	} else {

		//Postformat global
		function, exist = sp.Hooker[Postformat]
		if exist {

			function(bufByte)
		
		}
	}
}


//spacePaddingPointer: Genera espacios tanto para fields como para columnas.
//#bd/core.go
func (sP *Space)spacePaddingPointer(buf *[]byte , colName string){

	//Contamos el array de bytes
	var text_count = int64(len(*buf))


	var sizeColumn int64 = 0

	if sP.IndexSizeColumns != nil {

		_, found := sP.IndexSizeColumns[colName]
		if found {

			sizeColumn = sP.IndexSizeColumns[colName][1] - sP.IndexSizeColumns[colName][0]

		}
		if !found {
		
			if sP.IndexSizeFields != nil {

				_, found := sP.IndexSizeFields[colName]
				if found {

					sizeColumn = sP.IndexSizeFields[colName][1] - sP.IndexSizeFields[colName][0]

				}
				if !found {
		
					log.Fatalln("Archivo: core.go ; Funcion: spacePaddingPointer ; No hubo coincidencias en fields o columnas para añadir o quitar padding." + sP.Dir)
				
				}
			}
		}
	}

	if sP.IndexSizeColumns == nil {

		if sP.IndexSizeFields != nil {

			_, found := sP.IndexSizeFields[colName]
			if found {

				sizeColumn = sP.IndexSizeFields[colName][1] - sP.IndexSizeFields[colName][0]

			}
			if !found {
	
				log.Fatalln("Archivo: core.go ; Funcion: spacePaddingPointer ; No hubo coincidencias en fields o columnas para añadir o quitar padding." + sP.Dir)
			
			}
		}
	}

	if sP.IndexSizeColumns == nil && sP.IndexSizeFields == nil {

		log.Fatalln("Archivo: core.go ; Funcion: spacePaddingPointer ; El archivo no tiene columnas o campos." + sP.Dir)

	}

	if sizeColumn == 0 {

		log.Fatalln("Archivo: core.go ; Funcion: spacePaddingPointer ; Columnas o campos invalidos para añadir o quitar padding." + sP.Dir)

	}

	if text_count < sizeColumn {

		whitespace := bytes.Repeat( []byte(" ") , int(sizeColumn - text_count)) 
					
		*buf = append(*buf ,  whitespace... )
	}

	if text_count > sizeColumn {

		*buf = (*buf)[:sizeColumn]
	}

}

//spacePaddingPointer: Genera espacios tanto para fields como para columnas.
//#bd/core.go
func (sP *Space)spaceTrimPointer(buf *[]byte){

		//Limpiamos espacios en blanco
		for len(*buf) > 0 && (*buf)[len(*buf)-1] == 32 {

			*buf = (*buf)[:len(*buf)-1]

		}
		
}



//WriteIndexSizeField: Escribe los fields en los archivos.
//#bd/core.go
func (sF *spaceFile) WriteIndexSizeField(columnName string,buff *[]byte ){

	if sF.Hooker != nil {

		sF.hookerPreFormatPointer(buff, columnName)

	}

	//Contamos el array de bytes
	var text_count = int64(len(*buff))
	
	//Obtenemos el campo
	sizeField , _ := sF.IndexSizeFields[columnName]

	//Obtenermos el valor de las columnas
	sizeColumn    := sizeField[1] - sizeField[0]
	//Primer caso el texto es menor que el tamaño de la linea
	//En este caso añadimos un padding de espacios al final
	

	if text_count < sizeColumn {

		whitespace := bytes.Repeat( []byte(" ") , int(sizeColumn - text_count)) 
					
		*buff = append(*buff ,  whitespace... )
	}

	if text_count > sizeColumn {

		*buff = (*buff)[:sizeColumn]
	}
	
	
	sF.File.WriteAt(*buff,   + sizeField[0])


}


//columnSpacePadding: Da spacios actualmente a la funcion de columns
//#bd/core.go
func (sp *spaceFile)columnSpacePadding(colName string, buf *[]byte){

	//Contamos el array de bytes
	var text_count = int64(len(*buf))

	//Primer caso el texto es menor que el tamaño de la linea
	//En este caso añadimos un padding de espacios al final
	sizeColumn := sp.IndexSizeColumns[colName][1] - sp.IndexSizeColumns[colName][0]

	if text_count < sizeColumn {

		whitespace := bytes.Repeat( []byte(" ") , int(sizeColumn - text_count)) 
					
		*buf = append(*buf ,  whitespace... )
	}

	if text_count > sizeColumn {

		*buf = (*buf)[:sizeColumn]
	}
}

//check: revisa las columnas y los fields haber si existen como columna si no existe da error fatal.
//#core.go/check
func (sP *Space) checkColFil(name string, err string){

	mensaje := err + " ; La columna o field: " + name + " no existe en el archivo ; " + sP.Dir

	if sP.IndexSizeColumns != nil {

		_, found := sP.IndexSizeColumns[name]
		if !found {
		
			if sP.IndexSizeFields != nil {

				_, found := sP.IndexSizeFields[name]
				if !found {
		
					log.Fatalln(mensaje)
				
				}
				return
			}
		}
	}

	if sP.IndexSizeColumns == nil {

		if sP.IndexSizeFields != nil {

			_, found := sP.IndexSizeFields[name]
			if !found {
	
				log.Fatalln(mensaje)
			
			}
			return
		}
	}

	if sP.IndexSizeColumns == nil && sP.IndexSizeFields == nil {

		log.Fatalln("El archivo no tiene columnas o campos que sincronizar." + sP.Dir)

	}

}




//WriteIndexSizeField: Escribe los fields en los archivos.
//#bd/core.go
func (buf *RBuffer)readIndexSizeFieldPointer(colName string,  size [2]int64 ){

	if CheckFileTypeBuffer(buf.typeBuff , BuffChan ) {

		buf.FieldBuffer = new([]byte)

	}


	sizeTotal := size[1] - size[0]

	if buf.RangeBytes < sizeTotal && buf.RangeBytes > 0 {
		
		buf.TotalRangue = sizeTotal / buf.RangeBytes
		restoRangue := sizeTotal % buf.RangeBytes


		if buf.Rangue < buf.TotalRangue {

			if len(*buf.FieldBuffer) != int(buf.RangeBytes) {

				*buf.FieldBuffer = make([]byte, buf.RangeBytes)

			}
			

			_ , err := buf.File.ReadAt(*buf.FieldBuffer , size[0] + (buf.RangeBytes * buf.Rangue) )
			if err != nil {
				log.Println(err)
				
			}

			//Limpiamos el buffer de espacios
			if buf.PostFormat == true {

				buf.spaceTrimPointer(buf.FieldBuffer)

			}
			
			//Activamos PostFormat si existe
			if buf.Hooker != nil && buf.PostFormat == true {
			
				buf.hookerPostFormatPointer(buf.FieldBuffer ,colName)

			}

			if restoRangue != 0 {
			
				buf.TotalRangue += 1
			}


		}

		
		if  buf.Rangue == buf.TotalRangue && restoRangue != 0 {


			*buf.FieldBuffer = make([]byte, restoRangue)

			_ , err := buf.File.ReadAt(*buf.FieldBuffer , size[0] + (buf.RangeBytes * buf.Rangue) )
			if err != nil {

				log.Println(err)
				
			}

			//Limpiamos el buffer de espacios
			if buf.PostFormat == true {

				buf.spaceTrimPointer(buf.FieldBuffer)

			}
			
			//Activamos PostFormat si existe
			if buf.Hooker != nil && buf.PostFormat == true {
			
				buf.hookerPostFormatPointer(buf.FieldBuffer ,colName)

			}

			buf.TotalRangue += 1
			
		}	
	}


	if buf.RangeBytes >= sizeTotal || buf.RangeBytes <= 0 {

		*buf.FieldBuffer = make([]byte, size[1] - size[0])

		_ , err := buf.File.ReadAt(*buf.FieldBuffer , size[0])
		if err != nil {
			log.Println(err)
			
		}

		//Limpiamos el buffer de espacios
		if buf.PostFormat == true {

			buf.spaceTrimPointer(buf.FieldBuffer)

		}
		
		//Activamos PostFormat si existe
		if buf.Hooker != nil && buf.PostFormat == true {
		
			buf.hookerPostFormatPointer(buf.FieldBuffer ,colName)

		}
	
	}


}