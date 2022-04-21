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
	typeBuff   FileTypeBuffer
	PostFormat bool

	StartLine  int64
	EndLine    int64
	ColName    *[]string

	Rangue      int64 
	TotalRangue int64 
	RangeBytes  int64

	FieldBuffer *[]byte
	Buffer      *[]byte
	BufferMap    map[string][][]byte
	Channel      chan RChanBuf
}


func (sF *spaceFile)NewBRspace(typeBuff FileTypeBuffer)(buf *RBuffer){

	return &RBuffer{
			spaceFile: sF,
			typeBuff:typeBuff,
	}
}

func (RB *RBuffer)CheckBRspace(active bool){

	RB.check = active
}

func (RB *RBuffer)PostFormatBRspace(active bool){

	RB.PostFormat = active
}

func (RB *RBuffer)OneLineBRspace(line int64){

	if RB.check && *RB.SizeFileLine < line {
		
		log.Fatalln("Error de buffer archivo: ",RB.Url , " Linea final: ",line  , " Numero de lineas del archivo: ", *RB.SizeFileLine)
		
	}

	RB.StartLine = line
	RB.EndLine   = line + 1
}

func (RB *RBuffer)MultiLineBRspace(startLine int64,endLine int64){

	if RB.check && *RB.SizeFileLine < endLine {
		
		log.Fatalln("Error de buffer archivo: ",RB.Url , " Linea final: ",endLine  , " Numero de lineas del archivo: ", *RB.SizeFileLine)
		
	}

	RB.StartLine = startLine
	RB.EndLine   = endLine + 1
}

func (RB *RBuffer)AllLinesBRspace(){

	RB.StartLine = 0
	RB.EndLine   = *RB.SizeFileLine + 1
}

func (RB *RBuffer)RangeBytesBRspace(RangeBytes int64){

	RB.RangeBytes  = RangeBytes
	RB.TotalRangue = 1
	RB.Rangue      = 0
}

//BRspace: crea un buffer de lectura, se puede elegir si añadir el postformat de los campos.
//Las lineas estan precalculadas inicio 0 - fin - 0 equivale a la linea 0, 0 - 1 equivale a la linea 0 y 1.
//data son los fields y las columnas que se desean.
func (RB *RBuffer)BRspace(data ...string ){
	 
	if RB.check {

		if  int64(len(data)) > RB.lenColumns + RB.lenFields {

			log.Fatalln("Error el archivo solo tiene: ",RB.lenColumns + RB.lenFields,"columnas y campos",RB.Url   )
		
		}
	
		if  len(data) == 0 {
	
			log.Fatalln("No se puede enviar un buffer vacio en:",RB.Url   )
		
		}

		if CheckFileTypeBuffer(RB.typeBuff, BuffBytes ){

			if RB.check && len(data) > 1 {
	
				log.Fatalln("El Buffer de Bytes solo es compatible con un unico campo.",RB.Url   )
			
			}
		}

		for _, colname := range data {

			if RB.check {

				RB.checkColFil(colname, "Archivo: BRspace.go ; Funcion: BRspace")
			}
		}

	}

	RB.ColName = &data

	//Buffer de bytes
	if CheckFileTypeBuffer(RB.typeBuff, BuffBytes ){

	
		if RB.IndexSizeFields != nil {
			_, found := RB.IndexSizeFields[(*RB.ColName)[0]]
			if found {
				RB.FieldBuffer = new([]byte)
				return
			}
		}

		if RB.IndexSizeColumns != nil {
			size, found := RB.IndexSizeColumns[(*RB.ColName)[0]]
			if found {
				RB.Buffer = new([]byte)
				*RB.Buffer = make([]byte ,size[1] - size[0])
				return
			}
		}
		
		return
	}
	

	if CheckFileTypeBuffer(RB.typeBuff, BuffChan){

		RB.FieldBuffer = new([]byte)
		RB.Buffer      = new([]byte)
		*RB.Buffer     = make([]byte ,RB.SizeLine)
		RB.Channel     =  make(chan RChanBuf,1)

		return
	}
	

	if CheckFileTypeBuffer(RB.typeBuff, BuffMap ){

		RB.FieldBuffer  = new([]byte)

		RB.BufferMap    = make(map[string][][]byte)
		RB.BufferMap["buffer"]    = make([][]byte ,1)
		RB.BufferMap["buffer"][0] = make([]byte ,RB.SizeLine * (RB.EndLine - RB.StartLine ))

		return
	}

}

//Funcion de soporte para la lectura y posterior envio de datos en un canal.
//Si no se crea un buffer nuevo antes de cada envio el buffer falla al intentar
//pasarlo por referencia
func (buf *RBuffer) NewChanBuffer (){

	 buf.Buffer = new([]byte)
	*buf.Buffer = make([]byte ,buf.SizeLine)

}

//Revisa el tipo de buffer y devuelve true o false dependiendo de si hay coincidencia.
func CheckFileTypeBuffer(base FileTypeBuffer, compare FileTypeBuffer)(bool){

	if (base & compare) != 0 {

		return true

	}

	return false

}



