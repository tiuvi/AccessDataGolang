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

//Revisa el tipo de buffer y devuelve true o false dependiendo de si hay coincidencia.
func CheckFileTypeBuffer(base FileTypeBuffer, compare FileTypeBuffer)(bool){

	if (base & compare) != 0 {

		return true

	}

	return false

}

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
type RLines struct{
	StartLine  int64
	EndLine    int64
}

type RRangues struct {
	Rangue      int64 
	TotalRangue int64 
	RangeBytes  int64
}

type RBuffer struct {
	
	*spaceFile
	typeBuff   FileTypeBuffer
	PostFormat bool
	ColName    *[]string

	*RLines
	*RRangues

	FieldBuffer *[]byte
	Buffer      *[]byte
	BufferMap    map[string][][]byte
	Channel      chan RChanBuf
}

//Crea un nuevo espacio de buffer, que es una unica referencia.
func (sF *spaceFile)NewBRspace(typeBuff FileTypeBuffer)(buf *RBuffer){

	return &RBuffer{
			spaceFile: sF,
			typeBuff:typeBuff,
	}
}

//Activa o desactiva los errores y chequeos adicionales.
func (RB *RBuffer)CheckBRspace(active bool){

	RB.check = active
}

//Activa o desactiva el postformateado de los datos.
func (RB *RBuffer)PostFormatBRspace(active bool){

	RB.PostFormat = active
}

//Lee una sola linea de buffer.(Columnas)
func (RB *RBuffer)OneLineBRspace(line int64){

	if RB.check && *RB.SizeFileLine < line {
		
		log.Fatalln("Error de buffer archivo: ",RB.Url , " Linea final: ",line  , " Numero de lineas del archivo: ", *RB.SizeFileLine)
		
	}
	RB.RLines = new(RLines)
	RB.StartLine = line
	RB.EndLine   = line + 1
}

//Lee multiples lineas de buffer.(Columnas)
func (RB *RBuffer)MultiLineBRspace(startLine int64,endLine int64){

	if RB.check && *RB.SizeFileLine < endLine {
		
		log.Fatalln("Error de buffer archivo: ",RB.Url , " Linea final: ",endLine  , " Numero de lineas del archivo: ", *RB.SizeFileLine)
		
	}
	RB.RLines = new(RLines)
	RB.StartLine = startLine
	RB.EndLine   = endLine + 1
}

//Lee todas las lineas de un archivo(Columnas) 
func (RB *RBuffer)AllLinesBRspace(){

	RB.RLines = new(RLines)
	RB.StartLine = 0
	RB.EndLine   = *RB.SizeFileLine + 1
}

//Lee un field sin rangos.(Campos)
func (RB *RBuffer)NoRangeFieldsBRspace(){
	
	RB.RRangues = new(RRangues)
	RB.RangeBytes  = 0
	RB.TotalRangue = 1
	RB.Rangue      = 0
}

//Lee un fields por rangos, de manera dinamica.(Campos)
func (RB *RBuffer)RangeFieldsBRspace(RangeBytes int64){
	
	RB.RRangues = new(RRangues)
	RB.RangeBytes  = RangeBytes
	RB.TotalRangue = 1
	RB.Rangue      = 0
}

//CAlcula el numero de rangos de un field.
func (RB *RBuffer)CalcRangeFieldBRspace(field string, RangeBytes int64)*int64{
	
	if RB.IndexSizeFields != nil {

		size, found := RB.IndexSizeFields[field]
		if found {

			sizeTotal   := size[1] - size[0]
			if  RangeBytes < sizeTotal && RangeBytes > 0 {

				TotalRangue := sizeTotal / RB.RangeBytes
				restoRangue := sizeTotal % RB.RangeBytes
				if restoRangue != 0 {

					TotalRangue += 1
				}
				return &TotalRangue
			}
		}
	}
	return nil
}

//Verifica si el string, es un campo en ese espacio.
func (RB *RBuffer)isField(field string)*bool{

	if RB.IndexSizeFields != nil {

		_, found := RB.IndexSizeFields[field]
		if found {

			exist := true
			return &exist	
		}

		exist := false
		return &exist	
	}

	return nil
}
	
//Verifica si el string, es una columna en ese espacio.
func (RB *RBuffer)isColumn(column string)*bool{

	if RB.IndexSizeColumns != nil {

		_, found := RB.IndexSizeColumns[column]
		if found {

			exist := true
			return &exist	
		}

		exist := false
		return &exist	
	}

	return nil
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

	
		if RB.IndexSizeFields != nil && RB.RRangues != nil {

			_, found := RB.IndexSizeFields[(*RB.ColName)[0]]
			if found {

				RB.FieldBuffer = new([]byte)	
			}
		}

		if RB.IndexSizeColumns != nil  && RB.RLines != nil {  

			size, found := RB.IndexSizeColumns[(*RB.ColName)[0]]
			if found {

				RB.Buffer = new([]byte)
				*RB.Buffer = make([]byte ,size[1] - size[0])
			}
		}
		return
	}
	

	if CheckFileTypeBuffer(RB.typeBuff, BuffChan){

		if RB.IndexSizeFields != nil && RB.RRangues != nil {

			RB.FieldBuffer = new([]byte)
		}

		if RB.IndexSizeColumns != nil  && RB.RLines != nil {  

			RB.Buffer      = new([]byte)
			*RB.Buffer     = make([]byte ,RB.SizeLine)
			RB.Channel     =  make(chan RChanBuf,1)
		}
		return
	}
	

	if CheckFileTypeBuffer(RB.typeBuff, BuffMap ){

		RB.BufferMap    = make(map[string][][]byte)

		if RB.IndexSizeFields != nil && RB.RRangues != nil {

			RB.FieldBuffer  = new([]byte)
		}

		if RB.IndexSizeColumns != nil  && RB.RLines != nil {  

			RB.BufferMap["buffer"]    = make([][]byte ,1)
			RB.BufferMap["buffer"][0] = make([]byte ,RB.SizeLine * (RB.EndLine - RB.StartLine ))
		}
		return
	}
}







