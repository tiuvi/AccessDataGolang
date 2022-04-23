package bd

import (
	"sync/atomic"
	"log"
)




type WLines struct{
	Line  int64
}

type WRangues struct {
	Rangue      int64 
	RangeBytes  int64
	TotalRangue int64 
}

//El canal de escritura linea , nombre de columna y el buffer.
type WChanBuf struct{
	*WLines
	*WRangues
	ColName string
	Buffer 	[]byte
}

//Buffer de escritura con tres tipos de buffer
//Tipo buffer unicamente puede escribir en una columna o un field.
//Tipo mapaBuffer puede escribir simultaneametne en columnas y fields.
//Abre un canal que puede actualizar tantno columnas como fields.
type WBuffer struct {
	*spaceFile
	ColumnName string
	typeBuff FileTypeBuffer
	PreFormat bool
	*WLines
	*WRangues

	Buffer *[]byte
	BufferMap map[string][]byte
	Channel chan WChanBuf
}

//Activa o desactiva los errores y chequeos adicionales.
func (WB *WBuffer)CheckWRspace(active bool){

	WB.check = active
}

//Crea un nuevo espacio de buffer, que es una unica referencia.
func (sF *spaceFile)NewWBspace(typeBuff FileTypeBuffer)*WBuffer{

	WB := &WBuffer{
			spaceFile: sF,
			PreFormat: true,
			typeBuff:typeBuff,
	}

	//Buffer de bytes
	if CheckFileTypeBuffer(WB.typeBuff, BuffBytes ){

		WB.Buffer = new([]byte)
		return WB
	}

	//Buffer de bytes
	if CheckFileTypeBuffer(WB.typeBuff, BuffMap ){

		WB.BufferMap = make(map[string][]byte)
		return WB
	}

	//Buffer de bytes
	if CheckFileTypeBuffer(WB.typeBuff, BuffChan ){

		WB.Channel = make(chan WChanBuf, 0)
		return WB
	}

	return nil
}

//Activa o desactiva el preformateado de los datos.
func (WB *WBuffer)PreFormatBRspace(active bool){

	WB.PreFormat = active
}

//Escribe una nueva linea en el archivo.
func (WB *WBuffer)NewLineWBspace(){

	WB.WLines = new(WLines)
	WB.Line = -1
}

//Escribe en un numero determinado de linea.
func (WB *WBuffer)UpdateLineWBspace(line int64){

	if WB.check &&  line < 0 {
		
		log.Fatalln("Error no se puede enviar una linea inferior a 0",WB.Url)
		
	}
	WB.WLines = new(WLines)
	WB.Line = line
}

//Escribe en la siguiente linea
func (WB *WBuffer)NextLineWBspace(){

	WB.Line += 1
}

//Añade el formato de los rangos
func (WB *WBuffer)NewRangeWBspace(RangeBytes int64,TotalRangue int64 ){

	WB.WRangues = new(WRangues)
	WB.Rangue   =  0       
	WB.RangeBytes  = RangeBytes
	WB.TotalRangue = TotalRangue
	
}

//Escribe en un rango
func (WB *WBuffer)RangeWBspace(Range int64){

	WB.Rangue   =  Range
	
}

//Escribe en el siguiente rango
func (WB *WBuffer)NextRangeWBspace(){

	WB.Rangue   +=  1
	
}

//CAlcula el tamaño de un campo
func (WB *WBuffer)CalcSizeFieldBWspace(field string)*int64{

	if WB.IndexSizeFields != nil {

		size, found := WB.IndexSizeFields[field]
		if found {

			sizeTotal   := size[1] - size[0]
		
			return &sizeTotal
		}
	}
	return nil
}

//CAlcula el tamaño de una columna
func (WB *WBuffer)CalcSizeColumnBWspace(field string)*int64{

	if WB.IndexSizeColumns != nil {

		size, found := WB.IndexSizeColumns[field]
		if found {

			sizeTotal   := size[1] - size[0]
		
			return &sizeTotal
		}
	}
	return nil
}

//  if WB.IndexSizeFields != nil && WB.WRangues != nil {

//  if WB.IndexSizeColumns != nil  && WB.WLines != nil {  

func(WB *WBuffer)SendBWspace(columnName string, bufferBytes *[]byte)*int64{

		if WB.check  {
			
			WB.checkColFil(columnName, "Archivo: BWspace.go ; Funcion: BWspaceBuff")
		}

		//Buffer de bytes
		if CheckFileTypeBuffer(WB.typeBuff, BuffBytes ){

			WB.ColumnName = columnName
			WB.Buffer     = bufferBytes
			return nil
		}

		//Buffer de bytes
		if CheckFileTypeBuffer(WB.typeBuff, BuffMap ){

			if WB.BufferMap == nil {

				WB.BufferMap = make(map[string][]byte)

			}

			WB.BufferMap[columnName] = *bufferBytes
			return nil
		}

		//Buffer de bytes
		if CheckFileTypeBuffer(WB.typeBuff, BuffChan ){

			if WB.IndexSizeFields != nil && WB.WRangues != nil {

				WB.Channel <- WChanBuf{nil, WB.WRangues , columnName, *bufferBytes }
				return nil
			}

			if WB.IndexSizeColumns != nil  && WB.WLines != nil {  
			
				if WB.Line == -1 {
		
					WB.Line = atomic.AddInt64(WB.SizeFileLine, 1)
		
				}
		
				if WB.Line > *WB.SizeFileLine {
		
					atomic.AddInt64(WB.SizeFileLine, WB.Line - *WB.SizeFileLine )
						
				}

				WB.Channel <- WChanBuf{WB.WLines, nil , columnName, *bufferBytes }
				return &WB.Line
			}

		}

		return nil
}


//Cierra un canal.
func (WBuffer *WBuffer)BWspaceClosechan(){
	
	close(WBuffer.Channel)
}

