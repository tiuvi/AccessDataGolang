package bd

import (
	"log"
	"regexp"
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
func (sp *spaceFile)hookerPreFormatPointer(bufByte *[]byte,colName string){

	//Preformat por columnas
	function, exist := sp.hooker[preformat + colName]
	if exist{

		 function(bufByte)

	} else {

		//Preformat global
		function, exist = sp.hooker[preformat]
		if exist {

			function(bufByte)

		}
	}
}


//hookerPostFormatPointer: Da postformato a la salida de bytes en los 
//archivos tanto en columnas como fields.
//#bd/core.go
func (sp *spaceFile)hookerPostFormatPointer(bufByte *[]byte,colName string){

	//Postformat por columnas
	function, exist := sp.hooker[postformat + colName]
	if exist{

		function(bufByte)

	} else {

		//Postformat global
		function, exist = sp.hooker[postformat]
		if exist {

			function(bufByte)
		
		}
	}
}


//spacePaddingPointer: Genera espacios tanto para fields como para columnas.
//#bd/core.go
func (sP *space)spacePaddingPointer(buf *[]byte , size [2]int64){

	//Contamos el array de bytes
	textCount  := int64(len(*buf))
	sizeColumn := size[1] - size[0]

	if textCount < sizeColumn {

		newBuf := make([]byte, sizeColumn)
		copy(newBuf, *buf)
		*buf = newBuf
	}

	if textCount > sizeColumn {

		*buf = (*buf)[:sizeColumn]
	}

}

//spacePaddingPointer: Genera espacios tanto para fields como para columnas.
//#bd/core.go
func (sP *space)spaceTrimPointer(buf *[]byte){

		//Limpiamos nulos
		for len(*buf) > 0 && (*buf)[len(*buf)-1] == 0 {

			*buf = (*buf)[:len(*buf)-1]

		}
		
}



//WriteIndexSizeField: Escribe los fields en los archivos.
//#bd/core.go
func (WB *WBuffer)WriteIndexSizeField(colName string, size [2]int64,rangues wRangues,fieldBuffer *[]byte){

	rangueTotal := size[1] - size[0]
	
	
	if WB.hooker != nil && WB.preFormat {

		WB.hookerPreFormatPointer(fieldBuffer, colName)

	}


	if  rangues.rangeBytes < rangueTotal && rangues.rangeBytes > 0 {
	
		rangueFinal := (rangues.rangue * rangues.rangeBytes) + rangues.rangeBytes

		if rangueFinal <= rangueTotal {

			WB.spacePaddingPointer(fieldBuffer , [2]int64{0 , rangues.rangeBytes})
	
		}
	
		if rangueFinal > rangueTotal {
	
			restRangue := rangues.rangeBytes - (rangueFinal - rangueTotal)
			if restRangue > 0 {

				WB.spacePaddingPointer(fieldBuffer , [2]int64{0 , restRangue })
			} 
			if restRangue <= 0 {
				return
			}
		}

	}

	if rangues.rangeBytes >= rangueTotal || rangues.rangeBytes <= 0 {

		WB.spacePaddingPointer(fieldBuffer , size)

	}
	
	WB.file.WriteAt(*fieldBuffer,   + size[0] + (rangues.rangue * rangues.rangeBytes))


}



//WriteIndexSizeField: Escribe los fields en los archivos.
//#bd/core.go
func (buf *RBuffer)readIndexSizeFieldPointer(colName string,  size [2]int64 ){

	if checkFileTypeBuffer(buf.typeBuff , buffChan ) {

		buf.FieldBuffer = new([]byte)

	}


	sizeTotal := size[1] - size[0]

	if  buf.rangeBytes < sizeTotal && buf.rangeBytes > 0 {
		
		buf.totalRangue = sizeTotal / buf.rangeBytes
		restoRangue := sizeTotal % buf.rangeBytes


		if buf.rangue < buf.totalRangue {

			if len(*buf.FieldBuffer) != int(buf.rangeBytes) {

				*buf.FieldBuffer = make([]byte, buf.rangeBytes)

			}
			

			_ , err := buf.file.ReadAt(*buf.FieldBuffer , size[0] + (buf.rangeBytes * buf.rangue) )
			if err != nil {
				log.Println(err)
				
			}

			//Limpiamos el buffer de espacios
			if buf.postFormat == true {

				buf.spaceTrimPointer(buf.FieldBuffer)

			}
			
			//Activamos PostFormat si existe
			if buf.hooker != nil && buf.postFormat == true {
			
				buf.hookerPostFormatPointer(buf.FieldBuffer ,colName)

			}

			if restoRangue != 0 {
			
				buf.totalRangue += 1
			}


		}

		
		if  buf.rangue == buf.totalRangue && restoRangue != 0 {


			*buf.FieldBuffer = make([]byte, restoRangue)

			_ , err := buf.file.ReadAt(*buf.FieldBuffer , size[0] + (buf.rangeBytes * buf.rangue) )
			if err != nil {

				log.Println(err)
				
			}

			//Limpiamos el buffer de espacios
			if buf.postFormat == true {

				buf.spaceTrimPointer(buf.FieldBuffer)

			}
			
			//Activamos PostFormat si existe
			if buf.hooker != nil && buf.postFormat == true {
			
				buf.hookerPostFormatPointer(buf.FieldBuffer ,colName)

			}

			buf.totalRangue += 1
			
		}	
	}


	if buf.rangeBytes >= sizeTotal || buf.rangeBytes <= 0 {

		*buf.FieldBuffer = make([]byte, size[1] - size[0])

		_ , err := buf.file.ReadAt(*buf.FieldBuffer , size[0])
		if err != nil {
			log.Println(err)
			
		}

		//Limpiamos el buffer de espacios
		if buf.postFormat == true {

			buf.spaceTrimPointer(buf.FieldBuffer)

		}
		
		//Activamos PostFormat si existe
		if buf.hooker != nil && buf.postFormat == true {
		
			buf.hookerPostFormatPointer(buf.FieldBuffer ,colName)

		}
	
	}


}



func checkFileNativeType(base fileNativeType, compare fileNativeType)(bool){

	if (base & compare) != 0 {

		return true

	}
	return false
}

func checkFileCoding(base fileCoding, compare fileCoding)(bool){

	if (base & compare) != 0 {

		return true

	}
	return false
}

//Revisa el tipo de buffer y devuelve true o false dependiendo de si hay coincidencia.
func checkFileTypeBuffer(base fileTypeBuffer, compare fileTypeBuffer) bool {

	if (base & compare) != 0 {

		return true

	}

	return false

}

func regexPathGlobal(path string) string {

	//Filtrado solo se permite letras mayusculas, letras minusculas, numeros y barras.
	return regexp.MustCompile(`[^a-zA-Z0-9/]`).ReplaceAllString(path, "")

}