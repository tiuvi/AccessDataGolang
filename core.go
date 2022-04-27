package bd

import (
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
func (sp *spaceFile)hookerPreFormatPointer(bufByte *[]byte,colName string){

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
func (sp *spaceFile)hookerPostFormatPointer(bufByte *[]byte,colName string){

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
func (sP *Space)spacePaddingPointer(buf *[]byte , size [2]int64){

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
func (sP *Space)spaceTrimPointer(buf *[]byte){

		//Limpiamos nulos
		for len(*buf) > 0 && (*buf)[len(*buf)-1] == 0 {

			*buf = (*buf)[:len(*buf)-1]

		}
		
}



//WriteIndexSizeField: Escribe los fields en los archivos.
//#bd/core.go
func (WB *WBuffer)WriteIndexSizeField(colName string, size [2]int64,rangues WRangues,fieldBuffer *[]byte){

	rangueTotal := size[1] - size[0]
	
	
	if WB.Hooker != nil && WB.PreFormat {

		WB.hookerPreFormatPointer(fieldBuffer, colName)

	}


	if  rangues.RangeBytes < rangueTotal && rangues.RangeBytes > 0 {
	
		rangueFinal := (rangues.Rangue * rangues.RangeBytes) + rangues.RangeBytes

		if rangueFinal <= rangueTotal {

			WB.spacePaddingPointer(fieldBuffer , [2]int64{0 , rangues.RangeBytes})
	
		}
	
		if rangueFinal > rangueTotal {
	
			restRangue := rangues.RangeBytes - (rangueFinal - rangueTotal)
			if restRangue > 0 {

				WB.spacePaddingPointer(fieldBuffer , [2]int64{0 , restRangue })
			} 
			if restRangue <= 0 {
				return
			}
		}

	}

	if rangues.RangeBytes >= rangueTotal || rangues.RangeBytes <= 0 {

		WB.spacePaddingPointer(fieldBuffer , size)

	}
	
	WB.File.WriteAt(*fieldBuffer,   + size[0] + (rangues.Rangue * rangues.RangeBytes))


}



//WriteIndexSizeField: Escribe los fields en los archivos.
//#bd/core.go
func (buf *RBuffer)readIndexSizeFieldPointer(colName string,  size [2]int64 ){

	if CheckFileTypeBuffer(buf.typeBuff , BuffChan ) {

		buf.FieldBuffer = new([]byte)

	}


	sizeTotal := size[1] - size[0]

	if  buf.RangeBytes < sizeTotal && buf.RangeBytes > 0 {
		
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