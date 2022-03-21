package bd


import (
	"log"
	)


type Buffer struct {
	StartLine int64
	EndLine   int64
	Buffer    map[string][][]byte
}


//Inicia el numero de columnas en una tabla
func (sp *Space) NewSearchSpace (startLine int64,endLine int64,data ...string )(buf *Buffer){
	 
	if  sp.Size_file / sp.Size_line <= endLine {
		
		log.Fatalln("NewSearchSpace: error endLine",sp.Url ,endLine ,sp.Size_file / sp.Size_line)
		
	}

	buf = &Buffer{
	   StartLine: startLine,
	   EndLine:   endLine + 1,
	   Buffer: make(map[string][][]byte),
   }
   
   for _, value := range data {

	   buf.Buffer[value] = make([][]byte ,0)
   }

   buf.Buffer["buffer"]    = make([][]byte ,1)
   buf.Buffer["buffer"][0] = make([]byte ,sp.Size_line * (buf.EndLine - buf.StartLine ))

   return
}

func (sp *Space) NewSearchBitSpace (line int64, data ...string )(buf *Buffer){

	buf = &Buffer{
		StartLine: line,
		Buffer: make(map[string][][]byte),
	}

	buf.Buffer["buffer"]    = make([][]byte , 1)
	buf.Buffer["buffer"][0] = make([]byte   , 1)

	return
}