package bd

import (
	"sync/atomic"
	"log"
)






//Crea un nuevo espacio de buffer, que es una unica referencia.
func (sF *spaceFile)NewWBspace(typeBuff fileTypeBuffer)*WBuffer{

	WB := &WBuffer{
			spaceFile: sF,
			preFormat: true,
			typeBuff:typeBuff,
	}

	//Buffer de bytes
	if checkFileTypeBuffer(WB.typeBuff, buffBytes ){

		WB.buffer = new([]byte)
		return WB
	}

	//Buffer de bytes
	if checkFileTypeBuffer(WB.typeBuff, buffMap ){

		WB.bufferMap = make(map[string][]byte)
		return WB
	}

	//Buffer de bytes
	if checkFileTypeBuffer(WB.typeBuff, buffChan ){

		WB.channel = make(chan wChanBuf, 0)
		return WB
	}

	return nil
}

//Activa o desactiva el preformateado de los datos.
func (WB *WBuffer)PreFormatBRspace(active bool){

	WB.preFormat = active
}

//Escribe una nueva linea en el archivo.
func (WB *WBuffer)NewLineWBspace(){

	WB.wLines = &wLines{
		line: -1,
	}
}

//Escribe en un numero determinado de linea.
func (WB *WBuffer)UpdateLineWBspace(line int64){

	if WB.spaceErrors != nil &&  line < 0 {
		
		log.Fatalln("Error no se puede enviar una linea inferior a 0",WB.url)
		
	}

	WB.wLines = &wLines{
		line: line,
	}
}

//Escribe en la siguiente linea
func (WB *WBuffer)NextLineWBspace(){

	*WB.wLines = wLines{
		line: WB.line + 1,
	}

}


//Añade el formato de los rangos
func (WB *WBuffer)NewRangeWBspace(RangeBytes int64 ){

	
	WB.wRangues = &wRangues{
		rangue:      0,       
		rangeBytes:  RangeBytes,
	}	
}

func (WB *WBuffer)NewNoRangeWBspace(){

	WB.wRangues = &wRangues{
		rangue:      0,       
		rangeBytes:  0,
	}	
}

//Escribe en un rango
func (WB *WBuffer)RangeWBspace(Range int64){

	WB.wRangues = &wRangues{
		rangue:      Range,       
		rangeBytes:  WB.rangeBytes,
	}

	//WB.Rangue   =  Range
	
}

//Escribe en el siguiente rango
func (WB *WBuffer)NextRangeWBspace(){

	WB.wRangues = &wRangues{
		rangue:      WB.rangue + 1,       
		rangeBytes:  WB.rangeBytes,
	}
	//WB.Rangue   +=  1
	
}

//Funcion de soporte para calcular los rangos
func calcRanges(lenBuffer int64,RangeBytes int64)*int64{

	if lenBuffer < RangeBytes {
		return nil
	}

	if RangeBytes <= 0 {
		return nil
	}

	TotalRangue := lenBuffer / RangeBytes
	restoRangue := lenBuffer % RangeBytes
	if restoRangue != 0 {

		TotalRangue += 1
	}
	return &TotalRangue
}

//CAlcula el tamaño de un campo
func (WB *WBuffer)CalcSizeFieldBWspace(field string)*int64{

	if WB.indexSizeFields != nil {

		size, found := WB.indexSizeFields[field]
		if found {

			sizeTotal   := size[1] - size[0]
		
			return &sizeTotal
		}
	}
	return nil
}

//CAlcula el tamaño de una columna
func (WB *WBuffer)CalcSizeColumnBWspace(field string)*int64{

	if WB.indexSizeColumns != nil {

		size, found := WB.indexSizeColumns[field]
		if found {

			sizeTotal   := size[1] - size[0]
		
			return &sizeTotal
		}
	}
	return nil
}


func(WB *WBuffer)SendBWspace(columnName string, bufferBytes *[]byte)*int64{

		if WB.spaceErrors != nil  {
			
			WB.checkColFil(columnName, "Archivo: BWspace.go ; Funcion: BWspaceBuff")
		}

		//Buffer de bytes
		if checkFileTypeBuffer(WB.typeBuff, buffBytes ){

			WB.ColumnName = columnName
			WB.buffer     = bufferBytes
			return nil
		}

		//Buffer de bytes
		if checkFileTypeBuffer(WB.typeBuff, buffMap ){

			if WB.bufferMap == nil {

				WB.bufferMap = make(map[string][]byte)

			}

			WB.bufferMap[columnName] = *bufferBytes
			return nil
		}

		//Buffer de bytes
		if checkFileTypeBuffer(WB.typeBuff, buffChan ){
		
			if WB.indexSizeFields != nil && WB.wRangues != nil  {
			
					WB.channel <- wChanBuf{nil, WB.wRangues , columnName, *bufferBytes }
					return nil
			
			}

			if WB.indexSizeColumns != nil  && WB.wLines != nil {  
			
			
				newLine := &wLines{
					line:WB.line,
				}
				if WB.line == -1 {
					
					 newLine.line = atomic.AddInt64(WB.sizeFileLine, 1)
		
				}
		
				if WB.line > *WB.sizeFileLine {
		
					atomic.AddInt64(WB.sizeFileLine, newLine.line - *WB.sizeFileLine )
						
				}


				WB.channel <- wChanBuf{newLine , nil , columnName, *bufferBytes }
				return &WB.line
			}

		}

		return nil
}


//Cierra un canal.
func (WBuffer *WBuffer)BWspaceClosechan(){
	
	close(WBuffer.channel)
}

