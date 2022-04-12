package bd


import (
	"log"

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
	*spaceFile
	StartLine int64
	EndLine   int64
	SizeLine  int64
	ColName   string
	typeBuff  FileTypeBuffer

	Buffer []byte
	BufferMap    map[string][][]byte
	Channel chan RChanBuf
}


//Inicia el numero de columnas en una tabla
func (sp *spaceFile) BRspace (typeBuff FileTypeBuffer,startLine int64,endLine int64,data ...string )(buf *RBuffer){
	 

	var lenData int64 = int64(len(data))

	if  *sp.SizeFileLine < endLine {
		
		log.Fatalln("Error de buffer archivo: ",sp.Url , " Linea final: ",endLine  , " Numero de lineas del archivo: ",sp.SizeFileLine)
		
	}


	if lenData > sp.lenColumns + sp.lenFields {

		log.Fatalln("Error el archivo solo tiene: ",sp.lenColumns + sp.lenFields,"columnas y campos",sp.Url   )
	
	}

	if lenData == 0 {

		log.Fatalln("No se puede enviar un buffer vacio en:",sp.Url   )
	
	}

	buf = &RBuffer{
		spaceFile: sp,
		StartLine: startLine,
		EndLine:   endLine + 1,
		SizeLine: sp.SizeLine,
		typeBuff:typeBuff,
	}


	//Buffer de bytes
	if CheckFileTypeBuffer(typeBuff, BuffBytes ){

		if lenData > 1 {

			log.Fatalln("El Buffer de Bytes solo es compatible con un unico campo.",sp.Url   )
		
		}

		for _, colname := range data {
			
			buf.ColName = colname

			if sp.IndexSizeFields != nil {
				size, found := sp.IndexSizeFields[buf.ColName]
				if found {

					buf.Buffer = make([]byte ,size[1] - size[0])
					return
				}
			}

			if sp.IndexSizeColumns != nil {
				size, found := sp.IndexSizeColumns[buf.ColName]
				if found {
					
					buf.Buffer = make([]byte ,size[1] - size[0])
					return
				}
			}
		}

		log.Fatalln("Buffer de Bytes sin coincidencias de campos o columnas.",sp.Url   )
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

//Funcion de soporte para la lectura y posterior envio de datos en un canal.
//Si no se crea un buffer nuevo antes de cada envio el buffer falla al intentar
//pasarlo por referencia
func (buf *RBuffer) NewChanBuffer (){

	buf.Buffer = make([]byte ,buf.SizeLine)
}


//Migrar searchbit al buffer read normal
func (sp *Space) NewSearchBitSpace(line int64, data ...string )(buf *RBuffer){

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

