package bd


import (
	"log"

)

//Tipos de buffer para las funciones de lectura y escritura.
type FileTypeBuffer int64
const(
	BuffMap FileTypeBuffer = 1 << iota
	BuffChan
	BuffBytes
)

//canal del buffer de lectura.
type RChanBuf struct{
	Line int64
	ColName string
	Buffer 	[]byte
}

//Buffer de lectura compatible con multiples lineas y fieldas.
//El buffer solo puede leer una columna o fields tampoco es compatible con multples lineas.
//Buffer map es compatible con multiples lineas, fields y columnas
//El canal de buffer es compatible con mutiples lineas fields y columnas.
type RBuffer struct {
	*spaceFile
	StartLine int64
	EndLine   int64

	ColName   string
	typeBuff  FileTypeBuffer
	PostFormat bool

	Rangue int64 
	TotalRangue int64 
	RangeBytes int64

	Buffer *[]byte
	BufferMap    map[string][][]byte
	Channel chan RChanBuf
}


//BRspace: crea un buffer de lectura, se puede elegir si añadir el postformat de los campos.
//Las lineas estan precalculadas inicio 0 - fin - 0 equivale a la linea 0, 0 - 1 equivale a la linea 0 y 1.
//data son los fields y las columnas que se desean.
func (sp *spaceFile) BRspace(typeBuff FileTypeBuffer,PostFormat bool, startLine int64,endLine int64,  data ...string )(buf *RBuffer){
	 

	var lenData int64 = int64(len(data))

	if  *sp.SizeFileLine < endLine {
		
		log.Fatalln("Error de buffer archivo: ",sp.Url , " Linea final: ",endLine  , " Numero de lineas del archivo: ", *sp.SizeFileLine)
		
	}


	if lenData > sp.lenColumns + sp.lenFields {

		log.Fatalln("Error el archivo solo tiene: ",sp.lenColumns + sp.lenFields,"columnas y campos",sp.Url   )
	
	}

	if lenData == 0 {

		log.Fatalln("No se puede enviar un buffer vacio en:",sp.Url   )
	
	}

	buf = &RBuffer{
		spaceFile: sp,
		typeBuff:typeBuff,
		PostFormat:PostFormat,
		StartLine: startLine,
		EndLine:   endLine + 1,
		
		TotalRangue: 1,
		Rangue: 0,
	}


	//Buffer de bytes
	if CheckFileTypeBuffer(typeBuff, BuffBytes ){

		if lenData > 1 {

			log.Fatalln("El Buffer de Bytes solo es compatible con un unico campo.",sp.Url   )
		
		}

		for _, colname := range data {

			sp.check(colname, "Archivo: BRspace.go ; Funcion: BRspace")

			buf.ColName = colname

			if sp.IndexSizeFields != nil {
				size, found := sp.IndexSizeFields[buf.ColName]
				if found {
					 buf.Buffer = new([]byte)
					*buf.Buffer = make([]byte ,size[1] - size[0])
					return
				}
			}

			if sp.IndexSizeColumns != nil {
				size, found := sp.IndexSizeColumns[buf.ColName]
				if found {
					 buf.Buffer = new([]byte)
					*buf.Buffer = make([]byte ,size[1] - size[0])
					return
				}
			}
		}

		log.Fatalln("Buffer de Bytes sin coincidencias de campos o columnas.",sp.Url   )
		return
	}
	
	if CheckFileTypeBuffer(typeBuff, BuffChan){

		buf.BufferMap = make(map[string][][]byte)

		for _, colname := range data {

			sp.check(colname, "Archivo: BRspace.go ; Funcion: BRspace")

			buf.BufferMap[colname] = nil
		}
		buf.Buffer = new([]byte)
		*buf.Buffer = make([]byte ,sp.SizeLine)
		buf.Channel =  make(chan RChanBuf,1)

		return
	}
	
	if CheckFileTypeBuffer(typeBuff, BuffMap ){


		buf.BufferMap = make(map[string][][]byte)

		for _, colname := range data {

			sp.check(colname, "Archivo: BRspace.go ; Funcion: BRspace")

			buf.BufferMap[colname] = make([][]byte ,0)

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
	 buf.Buffer = new([]byte)
	*buf.Buffer = make([]byte ,buf.SizeLine)
}

/*
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
*/

func CheckFileTypeBuffer(base FileTypeBuffer, compare FileTypeBuffer)(bool){

	if (base & compare) != 0 {

		return true

	}
	return false
}



