package dac

import (
	"fmt"

	"time"
)


func (RB *RBuffer) readByteSpace(){


	if EDAC &&
	RB.logTimeReadFile && !RB.isErrorFile {
		defer RB.NewLogDeferTimeMemorySF("readBytes", time.Now())
	}

	if checkFileTypeBuffer(RB.typeBuff , buffBytes ){

		if RB.indexSizeFields != nil && RB.rRangues != nil {

			size, found := RB.indexSizeFields[(*RB.colName)[0]]
			if found {

				RB.readIndexSizeFieldPointer((*RB.colName)[0], size)
		
			}
		}

		if RB.indexSizeColumns != nil  && RB.rLines != nil {

			size, found := RB.indexSizeColumns[(*RB.colName)[0]]
			if found {

				_ , err := RB.file.ReadAt(*RB.Buffer , RB.lenFields + (RB.startLine * RB.sizeLine) + size[0])
				if err != nil && EDAC && 
				RB.ECSFD( true,"Error al leer en el archivo \n\r" + fmt.Sprintln(err)){}
				
				
				//Limpiamos el buffer de espacios
				if RB.postFormat == true {

					SpaceTrimPointer(RB.Buffer)

				}


				//Activamos PostFormat si existe
				if RB.hooker != nil && RB.postFormat == true {
					
					RB.hookerPostFormatPointer(RB.Buffer ,(*RB.colName)[0])

				}
			}	
		}
		return
	}
	
	
	if checkFileTypeBuffer(RB.typeBuff , buffChan ){


		if RB.indexSizeFields != nil && RB.rRangues != nil {
			
			for _ , colName := range (*RB.colName) {
			
				size, found := RB.indexSizeFields[colName]
				if found {
	
					RB.rangue      = 0
					RB.totalRangue = 1

					for RB.rangue < RB.totalRangue {
					
						RB.readIndexSizeFieldPointer(colName, size)
					
						RB.Channel <- RChanBuf{RB.rangue ,colName,*RB.FieldBuffer}

						RB.rangue++
					}
				}

			}
		}

		

		if RB.indexSizeColumns != nil  && RB.rLines != nil {

			for RB.startLine < RB.endLine {

				//El buffer por referencia crea errores en los canales
				RB.Buffer = new([]byte)
				*RB.Buffer = make([]byte ,RB.sizeLine)
				//Leemos una sola linea
				_ , err := RB.file.ReadAt(*RB.Buffer , RB.lenFields + RB.startLine * RB.sizeLine)
				if err != nil && EDAC && 
				RB.ECSFD( true,"Error al leer en el archivo \n\r" + fmt.Sprintln(err)){}
				

				for _ , colName := range (*RB.colName) {

					size, found := RB.indexSizeColumns[colName]
					if found {

						bufferChan := make([]byte, (size[1] - size[0])  )
						bufferChan  = (*RB.Buffer)[size[0]:size[1]]
			
						//Limpiamos el buffer de espacios
						if RB.postFormat == true {
		
							SpaceTrimPointer(&bufferChan)

						}
						
						//Activamos PostFormat si existe
						if RB.hooker != nil && RB.postFormat == true {
						
							RB.hookerPostFormatPointer(&bufferChan ,colName)

						}
						
						
						RB.Channel <- RChanBuf{RB.startLine,colName,bufferChan}
					}

				}

				RB.startLine += 1
			}
		}

		
		//Cerramos el canal.
		close(RB.Channel)

		return
	}
	
	if checkFileTypeBuffer(RB.typeBuff , buffMap ){

		if RB.indexSizeFields != nil && RB.rRangues != nil {

			for _ , colName := range (*RB.colName) {
			
				size, found := RB.indexSizeFields[colName]
				if found {
				
					RB.rangue      = 0
					RB.totalRangue = 1

					//Obtimizacion especial para mapas.
					sizeTotal := size[1] - size[0]
					if  RB.rangeBytes < sizeTotal && RB.rangeBytes > 0 {

						TotalRangue := sizeTotal / RB.rangeBytes
						restoRangue := sizeTotal % RB.rangeBytes
						if restoRangue != 0 {

							TotalRangue += 1
						}
						RB.BufferMap[colName] = make([][]byte ,TotalRangue)

					} else {

						RB.BufferMap[colName] = make([][]byte , 1)
					}

					for RB.rangue < RB.totalRangue {

						RB.readIndexSizeFieldPointer(colName, size)
		
						RB.BufferMap[colName][RB.rangue] = append( RB.BufferMap[colName][RB.rangue],  *RB.FieldBuffer... )
						//RB.BufferMap[colName][RB.rangue]   =  *RB.FieldBuffer
				
						RB.rangue++

					}
					
				}
			}
		}
		
		if RB.indexSizeColumns != nil  && RB.rLines != nil {


			bufMapFile := &RB.BufferMap["buffer"][0]
			_ , err := RB.file.ReadAt(*bufMapFile , RB.lenFields + RB.startLine * RB.sizeLine )
			if err != nil && EDAC && 
			RB.ECSFD( true,"Error al leer en el archivo \n\r" + fmt.Sprintln(err)){}
				

			var mapCount int64

			for RB.startLine < RB.endLine {
				

				for _ , colName := range (*RB.colName) {


					if colName == "buffer" {

						continue

					}

					size, found := RB.indexSizeColumns[colName]
					if found {

						if _, found :=  RB.BufferMap[colName]; !found {
						
							RB.BufferMap[colName] = make([][]byte ,(RB.endLine - RB.startLine))

						}
						

						//Limpiamos el buffer de espacios
						//RB.BufferMap[colName][RB.startLine] = append(RB.BufferMap[colName][RB.startLine],  (*bufMapFile)[size[0]:size[1]]... )
						RB.BufferMap[colName][mapCount] =  (*bufMapFile)[size[0]:size[1]]
						
						
						
					
						//Limpiamos el buffer de espacios
						if RB.postFormat == true {

							SpaceTrimPointer(&RB.BufferMap[colName][len(RB.BufferMap[colName])-1])

						}
					
						//Activamos PostFormat si existe
						if RB.hooker !=nil && RB.postFormat == true {

							RB.hookerPostFormatPointer(&RB.BufferMap[colName][len(RB.BufferMap[colName])-1], colName)
							
						}

					}
				}

				mapCount++
				RB.startLine++

				//End bucle
				*bufMapFile = (*bufMapFile)[ RB.sizeLine:]
			}
	
			delete(RB.BufferMap,"buffer")
		}

		return
	}


	if EDAC && 
	RB.ECSFD( true,"Error fatal sin coincidencias filetypebuffer"){}
	
	return
}



func (RB *RBuffer) readBitSpace() {

	if EDAC &&
	RB.logTimeReadFile && !RB.isErrorFile {
		defer RB.NewLogDeferTimeMemory("readBits", time.Now())
	}

	if checkFileTypeBuffer(RB.typeBuff , buffBytes ){

		if RB.indexSizeFields != nil && RB.rRangues != nil {

			size, found := RB.indexSizeFields[(*RB.colName)[0]]
			if found {

				RB.readIndexSizeFieldPointer((*RB.colName)[0], size)
			}
		}

		if RB.indexSizeColumns != nil  && RB.rLines != nil {

			size , found := RB.indexSizeColumns[(*RB.colName)[0]]
			if found {

				var byteLine int64  =  RB.startLine / 8
				var bitLine  int64  =  RB.startLine % 8 

				//bufferBit := make([]byte , 1 )

				*RB.Buffer = (*RB.Buffer)[:1]

				_ , err := RB.file.ReadAt(*RB.Buffer ,RB.lenFields + (byteLine * RB.sizeLine) + size[0])
				if err != nil && EDAC && 
				RB.ECSFD( true,"Error al leer en el archivo de bit \n\r" + fmt.Sprintln(err)){}
				

				turn := readBit(bitLine,*RB.Buffer)

				if turn {

					*RB.Buffer = []byte("on")

				}else{

					*RB.Buffer = []byte("off")

				}
			}
		}
		return
	}
	
	if checkFileTypeBuffer(RB.typeBuff , buffChan ){

		if RB.indexSizeFields != nil && RB.rRangues != nil {

			for _ , colName := range (*RB.colName) {
			
				size, found := RB.indexSizeFields[colName]
				if found {
	
					RB.rangue      = 0
					RB.totalRangue = 1

					for RB.rangue < RB.totalRangue {

						RB.readIndexSizeFieldPointer(colName, size)

						RB.Channel <- RChanBuf{RB.rangue ,colName,*RB.FieldBuffer}

						RB.rangue++
					}
				}
			}
		}

		if RB.indexSizeColumns != nil  && RB.rLines != nil {

			for RB.startLine < RB.endLine {
					
				for _ , colName := range (*RB.colName) {

					size , found := RB.indexSizeColumns[colName]
					if found {

						var byteLine int64 =  RB.startLine / 8
						var bitLine  int64 =  RB.startLine % 8 


						_ , err := RB.file.ReadAt(*RB.Buffer , RB.lenFields + (byteLine * RB.sizeLine) + size[0])
						if err != nil && EDAC && 
						RB.ECSFD( true,"Error al leer en el archivo de bits \n\r" + fmt.Sprintln(err)){}
				

						turn := readBit(bitLine,*RB.Buffer)

						if turn {

							RB.Channel <- RChanBuf{RB.startLine , colName, []byte("on")}
					
						}else{

							RB.Channel <- RChanBuf{RB.startLine , colName, []byte("off")}
							
						}
					}
				}
				RB.startLine++
			}

			close(RB.Channel)
		}
		return
	}
	
	if checkFileTypeBuffer(RB.typeBuff , buffMap ){

		if RB.indexSizeFields != nil && RB.rRangues != nil {

			for _ , colName := range (*RB.colName) {
			
				size, found := RB.indexSizeFields[colName]
				if found {
				
					RB.rangue      = 0
					RB.totalRangue = 1

					sizeTotal := size[1] - size[0]
					if  RB.rangeBytes < sizeTotal && RB.rangeBytes > 0 {

						TotalRangue := sizeTotal / RB.rangeBytes
						restoRangue := sizeTotal % RB.rangeBytes
						if restoRangue != 0 {

							TotalRangue += 1
						}

						RB.BufferMap[colName] = make([][]byte ,TotalRangue)

					} else {

						RB.BufferMap[colName] = make([][]byte ,1)

					}

					for RB.rangue < RB.totalRangue {

						RB.readIndexSizeFieldPointer(colName, size)
		
						RB.BufferMap[colName][RB.rangue] = append( RB.BufferMap[colName][RB.rangue],  *RB.FieldBuffer... )
				
						RB.rangue++

					}
					
				}
			}
		}


		if RB.indexSizeColumns != nil  && RB.rLines != nil {

			var mapCount int64
			bufferBit := &RB.BufferMap["buffer"][0]

			for RB.startLine < RB.endLine {
					
				var byteLine int64 =  RB.startLine / 8
				var bitLine  int64 =  RB.startLine % 8 
				

				for _ , colName := range (*RB.colName) {


					if colName == "buffer" {

						continue

					}

					size , found := RB.indexSizeColumns[colName]
					if found {

						if _, found :=  RB.BufferMap[colName]; !found {

							RB.BufferMap[colName] = make([][]byte , (RB.endLine - RB.startLine))

						}
							
						_ , err := RB.file.ReadAt(*bufferBit , RB.lenFields + (byteLine * RB.sizeLine) + size[0])
						if err != nil && EDAC && 
						RB.ECSFD( true,"Error al leer en el archivo de bits \n\r" + fmt.Sprintln(err)){}
				

						turn := readBit(bitLine,*bufferBit)
						
						if turn {

							//RB.BufferMap[colName][mapCount] = append(RB.BufferMap[colName][mapCount] , []byte("on")...)
							RB.BufferMap[colName][mapCount] = []byte("on")
						}else{

							//RB.BufferMap[colName][mapCount] = append(RB.BufferMap[colName][mapCount] , []byte("off")...)
							RB.BufferMap[colName][mapCount] = []byte("off")
						}
					}
				}

				mapCount++
				RB.startLine++
			}
		}
		return
	}
	

	if EDAC && 
	RB.ECSFD( true,"Error fatal sin coincidencias filetypebuffer"){}
	return
}

