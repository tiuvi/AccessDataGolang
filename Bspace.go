package bd


import (
	"log"
	"sync/atomic"
//	"net/http"
	)


type FileTypeBuffer int64
const(
	BuffMap FileTypeBuffer = 1 << iota
	BuffChan
	BuffBytes
)

type RChanBuf struct{
	Line int64
	ColName string
	Buffer 	[]byte
}

type RBuffer struct {

	StartLine int64
	EndLine   int64
	SizeLine  int64
	typeBuff FileTypeBuffer

	Buffer []byte
	BufferMap    map[string][][]byte
	Channel chan RChanBuf
}


//Inicia el numero de columnas en una tabla
func (sp *Space) BRspace (typeBuff FileTypeBuffer,startLine int64,endLine int64,data ...string )(buf *RBuffer){
	 

	spaceFile := sp.OSpace()

	if  *spaceFile.SizeFileLine < endLine {
		
		log.Fatalln("Error de buffer archivo: ",spaceFile.Url , " Linea final: ",endLine  , " Numero de lineas del archivo: ",spaceFile.SizeFileLine)
		
	}

	if CheckFileTipeByte(sp.FileTipeByte, OneColumn) {
		
		if len(data) > 1 {
			log.Fatalln("Error el archivo solo tiene una columna.",spaceFile.Url   )
		}
	}

	if CheckFileTipeByte(sp.FileTipeByte, MultiColumn) {
		
		if len(data) > len(sp.IndexSizeColumns) {
			log.Fatalln("Error el archivo solo tiene: ",len(sp.IndexSizeColumns),spaceFile.Url   )
		}
	}

	buf = &RBuffer{
		StartLine: startLine,
		EndLine:   endLine + 1,
		SizeLine: sp.SizeLine,
		typeBuff:typeBuff,
	}


	//Buffer de bytes
	if CheckFileTypeBuffer(typeBuff, BuffBytes ){

		buf.Buffer = make([]byte ,sp.SizeLine)
		return
	}
	
	if CheckFileTypeBuffer(typeBuff, BuffChan){

		buf.BufferMap = make(map[string][][]byte)

		for _, value := range data {
		
			buf.BufferMap[value] = nil
		}
	
		buf.Buffer = make([]byte ,sp.SizeLine)
		buf.Channel =  make(chan RChanBuf,1)

		return
	}
	
	if CheckFileTypeBuffer(typeBuff, BuffMap ){


		buf.BufferMap = make(map[string][][]byte)

		for _, value := range data {
	 
			buf.BufferMap[value] = make([][]byte ,0)
		}
	 
		buf.BufferMap["buffer"]    = make([][]byte ,1)
		buf.BufferMap["buffer"][0] = make([]byte ,sp.SizeLine * (buf.EndLine - buf.StartLine ))
		return
	}


	log.Fatalln("Buffer flags definidas incorrectamente: ",typeBuff)
   	return
}

func (buf *RBuffer) NewChanBuffer (){

	buf.Buffer = make([]byte ,buf.SizeLine)
}


//Migrar searchbit al buffer read normal
func (sp *Space) NewSearchBitSpace (line int64, data ...string )(buf *RBuffer){

	buf = &RBuffer{
		StartLine: line,
		BufferMap: make(map[string][][]byte),
	}

	buf.BufferMap["buffer"]    = make([][]byte , 1)
	buf.BufferMap["buffer"][0] = make([]byte   , 1)

	return
}

func CheckFileTypeBuffer(base FileTypeBuffer, compare FileTypeBuffer)(bool){

	if (base & compare) != 0 {

		return true

	}
	return false
}

/*
//Test Funcion
func ReadChanBuffer(resp http.ResponseWriter, req *http.Request, buffer Buffer){


	for chanBufer := range buffer.Channel {

		resp.Write([]byte("Range chanBuffer: "+ string(chanBufer.Buffer) + "<br><br>" ) )

	}
}
*/


type WChanBuf struct{
	Line int64
	ColName string
	Buffer 	[]byte
}

type WBuffer struct {

	Line int64
	ColumnName string
	SizeLine  int64
	typeBuff FileTypeBuffer

	Buffer *[]byte
	BufferMap map[string][]byte
	Channel chan WChanBuf
}



func (sp *Space) BWspaceBuff(line int64 ,columnName string, buff []byte)(*WBuffer){

	 _ , found := sp.IndexSizeColumns[columnName]
	if !found {

		log.Fatalln("La columna: " + columnName + " no existe en ese archivo",sp.url)
		return nil
	}
	
	return &WBuffer{
		typeBuff: BuffBytes,
		Line: line,
		ColumnName: columnName,
		Buffer: &buff,
	}

}


func (sp *Space) BWspaceBuffMap(line int64 , buff map[string][]byte)(*WBuffer){

	for val := range buff {
		_ , found := sp.IndexSizeColumns[val]
		if !found {
	 
			log.Fatalln("La columna: " + val + " no existe en ese archivo",sp.url)
			return nil
		}
	}

   return &WBuffer{
		typeBuff: BuffMap,
	    Line: line,
	    BufferMap: buff,
   }

}

func  BWChanBuf()(*WBuffer){

	return &WBuffer{
		typeBuff: BuffChan,
		Channel: make(chan WChanBuf, 0),
	}
}


func (sp *Space)BWspaceSendchan(line int64, columnName string , buf *[]byte, WBuffer *WBuffer)int64{

	sF := sp.OSpace()

	if line == -1 {
	
		line = atomic.AddInt64(sF.SizeFileLine, 1)

	}

	if line > *sF.SizeFileLine {

		atomic.AddInt64(sF.SizeFileLine, line - *sF.SizeFileLine )
		
	}

	WBuffer.Buffer = buf
	WBuffer.Channel <- WChanBuf{line,columnName, *WBuffer.Buffer }

	return line
}

func (WBuffer *WBuffer) BWspaceClosechan(){
	
	close(WBuffer.Channel)
}

