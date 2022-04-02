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

func (sp *spaceFile) hookerPostFormatMap(buf *Buffer,val string){

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

func (sp *spaceFile) hookerPostFormatBuff(buf *Buffer){

		//Postformat global
		function, exist := sp.Hooker[Postformat]
		if exist {

			buf.Buffer = function(buf.Buffer)

		}
}