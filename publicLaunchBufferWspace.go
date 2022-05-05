package dac

import ("sync/atomic")





//Crea un nuevo buffer de bytes de escritura
func (sF *spaceFile) NewWriterBytes()*WBuffer {

	return &WBuffer{
		spaceFile: sF,
		preFormat: true,
		typeBuff:buffBytes,
		buffer: new([]byte),
	}

}

//Crea un nuevo buffer de bytes de escritura
func (sF *spaceFile) NewWriterMapBytes()*WBuffer {

	return &WBuffer{
		spaceFile: sF,
		preFormat: true,
		typeBuff:buffMap,
		bufferMap: make(map[string][]byte),
	}

}

//Crea un nuevo canal de escritura
func (sF *spaceFile) NewWriterChan()*WBuffer {

	return &WBuffer{
		spaceFile: sF,
		preFormat: true,
		typeBuff:buffChan,
		channel: make(chan wChanBuf, 0),
	}

}


//Activa o desactiva el preformateado de los datos.
func (WB *WBuffer)OffPreFormat(){

	WB.ECSFD(WB.preFormat == false, "Ya se desactivo el preformateo.")

	WB.preFormat = false
}

//Escribe una nueva linea en el archivo.
func (WB *WBuffer)NewLineWBspace(){

	WB.ECSFD(WB.wLines != nil, "Llamada dos veces a la misma funcion, Accion incompatible en paralelo, crea un nuevo buffer.")

	WB.wLines = &wLines{
		line: -1,
	}
}

//Escribe en un numero determinado de linea.
func (WB *WBuffer)UpdateLineWBspace(line int64){

	WB.ECSFD(WB.wLines != nil, "Llamada dos veces a la misma funcion, Accion incompatible en paralelo, crea un nuevo buffer.")
	WB.ECSFD(line < 0, "El numero de linea no puede ser inferior a cero.")

	WB.wLines = &wLines{
		line: line,
	}
}


//Escribe en un rango
func (WB *WBuffer)RangeWBspace(Range int64,RangeBytes int64){

	WB.ECSFD(WB.wRangues != nil, "Llamada dos veces a la misma funcion, Accion incompatible en paralelo, crea un nuevo buffer.")
	WB.ECSFD(Range < 0, "El numero de rango no puede ser inferior a cero.")
	WB.ECSFD(RangeBytes < 0, "El numero de rango de bytes no puede ser inferior a cero.")
	WB.ECSFD(Range > RangeBytes, "El numero de rango no puede ser mayor que el rango de bytes")

	WB.wRangues = &wRangues{
		rangue:      Range,       
		rangeBytes:  WB.rangeBytes,
	}
}

func (WB *WBuffer)NewNoRangeWBspace(){

	WB.ECSFD(WB.wRangues != nil, "Llamada dos veces a la misma funcion, Accion incompatible en paralelo, crea un nuevo buffer.")

	WB.wRangues = &wRangues{
		rangue:      0,       
		rangeBytes:  0,
	}	
}





func(WB *WBuffer)SendBWspace(columnName string, bufferBytes *[]byte)*int64{

	if EDAC && WB.ECSFD(WB.IsNotColFil(columnName) , "Se ha iniciado un buffer de escritura con una columna o field que no existe."){}
	WB.ECSFD(WB.wRangues != nil && WB.wLines != nil,"Usar buffer de lineas con buffer de rangos es propenso a errores.")
	WB.ECSFD(WB.wRangues == nil && WB.wLines == nil,"Iniciaste un buffer vacio.")
	

		//Buffer de bytes
		if checkFileTypeBuffer(WB.typeBuff, buffBytes ){

			WB.columnName = columnName
			WB.buffer     = bufferBytes
			return nil
		}

		//Buffer de bytes
		if checkFileTypeBuffer(WB.typeBuff, buffMap ){

	
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

