package bd

import (
	"bytes"

	//"log"
)






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

func readBit(id int64,bufferBit []byte)(bool){

	bitId := []uint8{128,64,32,16,8,4,2,1}

   var bufferuint uint8 = bufferBit[0]

	if (bufferuint & bitId[id]) != 0 {

		return true
	
	} 

	//if (bufferuint & bitId[id]) == 0 {}
	return false

}


//#bd/corerealtime
func (obj *spaceFile ) updateRamMap(str_write []byte, line int64){

	str_write = bytes.Trim(str_write, " ")
	obj.Lock()
	//Transformacion a string
	obj.Search[string(str_write)] = line
	obj.Unlock()

}


//#bd/corerealtime
func (obj *spaceFile ) updateRamIndex(str_write []byte, line int64){

	str := string(bytes.Trim(str_write, " "))

	lenIndex := int64(len(obj.Index))
	lineAdd  := line + 1
	if  lenIndex <= lineAdd {
		
		for lenIndex <= lineAdd {

			obj.Index = append(obj.Index, "")
			lenIndex = lenIndex + 1
		}
		
	} 

	obj.Index[line] = str

}


func CheckBit(base int64, compare int64)(bool){


	if (base & compare) != 0 {

		return true

	}
	return false
}





/*
*
* hooker preformat
*/

func (sp *spaceFile) hookerPreFormatMap(buf *WBuffer,val string){

	//Preformat por columnas
	function, exist := sp.Hooker[Preformat + val]
	if exist{

		buf.BufferMap[val] = function(buf.BufferMap[val])

	} else {

		//Preformat global
		function, exist = sp.Hooker[Preformat]
		if exist {

			buf.BufferMap[val] = function(buf.BufferMap[val])

		}
	}
}

func (sp *spaceFile) hookerPreFormatBuf(buf *WBuffer){

	//Preformat por columnas
	function, exist := sp.Hooker[Preformat + buf.ColumnName]
	if exist{

		*buf.Buffer = function(*buf.Buffer)

	} else {

		//Preformat global
		function, exist = sp.Hooker[Preformat]
		if exist {

			*buf.Buffer = function(*buf.Buffer)

		}
	}
}


func (sp *spaceFile) hookerPreFormatChan(name string, buf *[]byte){

	//Preformat por columnas
	function, exist := sp.Hooker[Preformat + name]
	if exist{

		*buf = function(*buf)

	} else {

		//Preformat global
		function, exist = sp.Hooker[Preformat]
		if exist {

			*buf = function(*buf)

		}
	}
}

/*
* Hooker post format
*
*
*/

func (sp *spaceFile) hookerPostFormatMap(buf *RBuffer,val string){

	//Postformat por columnas
	function, exist := sp.Hooker[Postformat + val]
	if exist{

		buf.BufferMap[val][len(buf.BufferMap[val])-1] = function(buf.BufferMap[val][len(buf.BufferMap[val])-1])

	} else {

		//Postformat global
		function, exist = sp.Hooker[Postformat]
		if exist {

			buf.BufferMap[val][len(buf.BufferMap[val])-1] = function(buf.BufferMap[val][len(buf.BufferMap[val])-1])

		}
	}

}


func (sp *spaceFile) hookerPostFormatBuff(buf *RBuffer,val string){

	//Postformat por columnas
	function, exist := sp.Hooker[Postformat + val]
	if exist{

		buf.Buffer = function(buf.Buffer)

	} else {

		//Postformat global
		function, exist = sp.Hooker[Postformat]
		if exist {

			buf.Buffer = function(buf.Buffer)
		
		}
	}
}

func (sp *spaceFile) hookerPostFormatBuffMultiColumn(buf *[]byte,val string){

	//Postformat por columnas
	function, exist := sp.Hooker[Postformat + val]
	if exist{

		testBuf := *buf
		testBuf  = function(testBuf)
		*buf      = testBuf

	} else {

		//Postformat global
		function, exist = sp.Hooker[Postformat]
		if exist {

			testBuf := *buf
			testBuf  = function(testBuf)
			*buf      = testBuf
		}
	}
	
}



/*
* Paddding column
*
*
*
**/
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