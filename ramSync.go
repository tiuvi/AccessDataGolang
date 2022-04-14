package bd

import (
	"log"
)

//SGMapStringLineInit: Sincroniza un mapa superglobal con un archivo.
//#bd/ramSync/SGMstringLine
func (sF *spaceFile) SGMapStringLineInit(name string,mapStr map[string]int64)map[string]int64{

	if mapStr == nil {

		mapStr = make(map[string]int64)
	
	}

	sF.check(name , "Archivo: ramSync.go ; Funcion: SGMstringLine")

	mapColumn := sF.BRspace( BuffMap, 0, *sF.SizeFileLine  , name)
	mapColumn.Rspace()

	var x int64
	for x = 0 ; x <= *sF.SizeFileLine; x++{


		mapStr[ string( mapColumn.BufferMap[name][x] ) ] = x
		
	}

	return mapStr
}

//SGMapStringLineUpd: Actualiza un mapa superglobal
//#bd/ramSync/SGMstringLineUpd
func (WBuf *WBuffer)SGMapStringLineUpd(mapStr map[string]int64)(int64, bool){


	if WBuf.typeBuff != BuffBytes {

		log.Fatalln("Sincronizacion de ram compatible unicamente con buffer de bytes.")

	}

	WBuf.Lock()
	defer WBuf.Unlock()
	log.Println(mapStr)
	_ , found := mapStr[string(*WBuf.Buffer)]
	if !found {

		line := WBuf.Wspace()

		if WBuf.Hooker != nil {

			WBuf.hookerPostFormatPointer(WBuf.Buffer, WBuf.ColumnName)
	
		}

		WBuf.spaceTrimPointer(WBuf.Buffer)

		mapStr[string(*WBuf.Buffer)] = line
	
		return line, true

	}

	return -1 ,false
}

/*
func (obj *spaceFile ) Rmapspace(str_write string)(value int64, found bool){

	if int64(len(str_write)) > obj.SizeLine{

		str_write = str_write[:obj.SizeLine]

	}
	
	if obj.Hooker != nil {

		bufferByte := []byte(str_write)
		obj.hookerPreFormatPointer(&bufferByte, Preformat)
		str_write = string(bufferByte)
	}


	obj.RLock()
	value, found = obj.Search[str_write]
	obj.RUnlock()
	return
}
*/







/*

func (sF *spaceFile ) ospaceCompilationFileRamIndex() {

	sF.FileNativeType |= RamIndex

	var field string

	for val, ind := range sF.IndexSizeColumns {

		if ind[0] == 0 {

			field = val
			break
		}
	}

	mapColumn := sF.BRspace(BuffMap,0, *sF.SizeFileLine, field)
	mapColumn.Rspace()

	sF.Index = make([]string ,0)
	
	var x int64
	for x = 0 ; x <= *sF.SizeFileLine; x++{
		
		sF.Index = append(sF.Index, string( mapColumn.BufferMap[field][x] ))

	}

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


func (obj *spaceFile ) Rindexspace(line int64)(value string, found bool){

	if  int64(len(obj.Index) ) > line {

		value, found = obj.Index[line] , true

	} else {

		value, found = "", false
	}

	return
}

*/