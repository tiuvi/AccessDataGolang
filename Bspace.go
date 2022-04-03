package bd


import (
	"log"
//	"net/http"
	)


type FileTypeBuffer int64
const(
	BuffMap FileTypeBuffer = 1 << iota
	BuffChan
	BuffBytes
	BuffDir
)

type ChanBuf struct{
	Line int64
	ColName string
	Buffer 	[]byte
}

type Buffer struct {

	StartLine int64
	EndLine   int64
	SizeLine  int64
	typeBuff FileTypeBuffer

	BufferMap    map[string][][]byte

	Buffer []byte

	Channel chan ChanBuf
}


//Inicia el numero de columnas en una tabla
func (sp *Space) Bspace (typeBuff FileTypeBuffer,startLine int64,endLine int64,data ...string )(buf *Buffer){
	 

	spaceFile := sp.OSpace()

	if  *spaceFile.SizeFileLine < endLine {
		
		log.Fatalln("Error de buffer archivo: ",spaceFile.Url , " Linea final: ",endLine  , " Numero de lineas del archivo: ",spaceFile.SizeFileLine)
		
	}

	if CheckBit(int64(sp.FileTipeByte), int64(OneColumn)) {
		
		if len(data) > 1 {
			log.Fatalln("Error el archivo solo tiene una columna.",spaceFile.Url   )
		}
	}

	if CheckBit(int64(sp.FileTipeByte), int64(MultiColumn)) {
		
		if len(data) > len(sp.IndexSizeColumns) {
			log.Fatalln("Error el archivo solo tiene: ",len(sp.IndexSizeColumns),spaceFile.Url   )
		}
	}

	buf = &Buffer{
		StartLine: startLine,
		EndLine:   endLine + 1,
		typeBuff:typeBuff,
	}


	//tipos de buffer para directorios
	if CheckBit(int64(typeBuff),int64(BuffDir)) {

		if CheckBit(int64(sp.FileNativeType),int64(Directory)) {		
	
			buf.BufferMap = make(map[string][][]byte)
			buf.BufferMap["buffer"]    = make([][]byte ,1)
			return
		}	
	}

	//Buffer de bytes
	if CheckBit(int64(typeBuff), int64(BuffBytes) ){

		buf.Buffer = make([]byte ,sp.SizeLine)
		return
	}
	
	if CheckBit(int64(typeBuff), int64(BuffChan) ){

		buf.BufferMap = make(map[string][][]byte)

		for _, value := range data {
		
			buf.BufferMap[value] = make([][]byte ,0)
		}
	
		buf.SizeLine = sp.SizeLine
		buf.Buffer = make([]byte ,sp.SizeLine)
		buf.Channel =  make(chan ChanBuf,1)

		return
	}
	
	if CheckBit(int64(typeBuff), int64(BuffMap) ){


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

func (buf *Buffer) NewChanBuffer (){

	buf.Buffer = make([]byte ,buf.SizeLine)
}



func (sp *Space) NewSearchBitSpace (line int64, data ...string )(buf *Buffer){

	buf = &Buffer{
		StartLine: line,
		BufferMap: make(map[string][][]byte),
	}

	buf.BufferMap["buffer"]    = make([][]byte , 1)
	buf.BufferMap["buffer"][0] = make([]byte   , 1)

	return
}

/*
//Test Funcion
func ReadChanBuffer(resp http.ResponseWriter, req *http.Request, buffer Buffer){


	for chanBufer := range buffer.Channel {

		resp.Write([]byte("Range chanBuffer: "+ string(chanBufer.Buffer) + "<br><br>" ) )

	}
}
*/