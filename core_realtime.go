package bd

import (
	"bytes"
	"log"
)

func (obj *Space ) complete_file_lines(line int64){

	if line  * obj.Size_line > obj.Size_file {
		
		save_table := int(  ( (obj.Size_line * line ) - obj.Size_file) )
		whitespace := bytes.Repeat( []byte(" ") , save_table ) 
	
		obj.File.WriteAt(whitespace , obj.Size_file)				
		
		obj.Size_file = obj.Size_line * line
	
	}
}

func (obj *Space ) complete_file_lines_bit(line int64){

	log.Println("Complete file lines bit: ",line, obj.Size_file)
	if line > obj.Size_file {
	
		log.Println("run..." , []byte{00})

		n, err := obj.File.WriteAt( []byte{0} , line   )				
		log.Println(n)
		log.Println(err)
		
		obj.Size_file = line
		log.Println("Complete file lines bit: ",line, obj.Size_file)

	}
}

func (obj *Space ) writeBit(id int64,turn bool ,bufferBit []byte)([]byte){

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

func (obj *Space ) readBit(id int64,bufferBit []byte)(bool){

	bitId := []uint8{128,64,32,16,8,4,2,1}

   var bufferuint uint8 = bufferBit[0]

	if (bufferuint & bitId[id]) != 0 {

		return true
	
	} 

	//if (bufferuint & bitId[id]) == 0 {}
	return false

}

func (obj *Space ) update_size_file(line int64){

	if line  * obj.Size_line == obj.Size_file {
		
		obj.Size_file = obj.Size_line * (line + 1)	

	}
}

func (obj *Space ) updateRamMap(str_write []byte, line int64){

	str_write = bytes.Trim(str_write, " ")
	obj.Lock()
	//Transformacion a string
	obj.Search[string(str_write)] = line
	obj.Unlock()

}

func (obj *Space ) updateRamIndex(str_write []byte, line int64){

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