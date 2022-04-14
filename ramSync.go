package bd

import (
//	"log"
)



type SpaceSGMStr struct  {
	*spaceFile
	col string
	Map map[string]int64
}

//SGMapStringLineInit: Sincroniza un mapa superglobal con un archivo.
//#bd/ramSync/SGMstringLine
func (sF *spaceFile) SGMapStringLineInit(colName string)*SpaceSGMStr {


	sF.check(colName , "Archivo: ramSync.go ; Funcion: SGMstringLine")

	SGMS := &SpaceSGMStr{
		spaceFile: sF,
		col: colName,
		Map: make(map[string]int64),
	}
	
	sF.Lock()
	defer sF.Unlock()

	mapColumn := SGMS.BRspace( BuffMap, 0, *sF.SizeFileLine  , colName)
	mapColumn.Rspace()

	var x int64
	for x = 0 ; x <= *sF.SizeFileLine; x++{


		SGMS.Map[ string( mapColumn.BufferMap[colName][x] ) ] = x
		
	}

	return SGMS
}

func (SGMS *SpaceSGMStr) SGMapStringLineRead(bufferBytes *[]byte)bool {
	
	
	if SGMS.Hooker != nil {

		SGMS.hookerPreFormatPointer(bufferBytes , SGMS.col)

	}

	SGMS.spaceTrimPointer(bufferBytes)

	valueStr := string(*bufferBytes)

	SGMS.RLock()
	defer SGMS.RUnlock()

	_ , found := SGMS.Map[valueStr]
	if found {

		return true

	}

	return false
}



//SGMapStringLineUpd: Actualiza un mapa superglobal
//#bd/ramSync/SGMstringLineUpd
func (SGMS *SpaceSGMStr)SGMapStringLineUpd(line int64 , bufferBytes *[]byte)(int64, bool){


	if SGMS.Hooker != nil {

		SGMS.hookerPreFormatPointer(bufferBytes , SGMS.col)

	}

    WBuf := SGMS.BWspaceBuff(line, SGMS.col,*bufferBytes)
	strBuffer := string(*bufferBytes)

	SGMS.Lock()
	defer SGMS.Unlock()

	if line == -1 {

		_ , found := SGMS.Map[strBuffer]
		if !found {
	
			line := WBuf.Wspace()
	
			WBuf.spaceTrimPointer(WBuf.Buffer)
	
			SGMS.Map[strBuffer] = line
		
			return line, true
	
		}

	}

	if line > -1 {

		delete(SGMS.Map , strBuffer)

		line := WBuf.Wspace()
	
		WBuf.spaceTrimPointer(WBuf.Buffer)

		SGMS.Map[strBuffer] = line
	
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