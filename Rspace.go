package bd

import (
	"log"
)


func (buf *RBuffer) readByteSpace(){

	var err error
	
	
	if checkFileTypeBuffer(buf.typeBuff , buffBytes ){

		if buf.indexSizeFields != nil && buf.rRangues != nil {

			size, found := buf.indexSizeFields[(*buf.colName)[0]]
			if found {

				buf.readIndexSizeFieldPointer((*buf.colName)[0], size)
		
			}
		}

		if buf.indexSizeColumns != nil  && buf.rLines != nil {

			size, found := buf.indexSizeColumns[(*buf.colName)[0]]
			if found {

				_ , err = buf.file.ReadAt(*buf.Buffer , buf.lenFields + (buf.startLine * buf.sizeLine) + size[0])
				if err != nil {
					log.Println(err)
					
				}
				
				//Limpiamos el buffer de espacios
				if buf.postFormat == true {

					buf.spaceTrimPointer(buf.Buffer)

				}


				//Activamos PostFormat si existe
				if buf.hooker != nil && buf.postFormat == true {
					
					buf.hookerPostFormatPointer(buf.Buffer ,(*buf.colName)[0])

				}
			}	
		}
		return
	}
	
	
	if checkFileTypeBuffer(buf.typeBuff , buffChan ){


		if buf.indexSizeFields != nil && buf.rRangues != nil {

			for _ , colName := range (*buf.colName) {
			
				size, found := buf.indexSizeFields[colName]
				if found {
	
					buf.rangue      = 0
					buf.totalRangue = 1

					for buf.rangue < buf.totalRangue {

						buf.readIndexSizeFieldPointer(colName, size)

						buf.Channel <- RChanBuf{buf.rangue ,colName,*buf.FieldBuffer}

						buf.rangue++
					}
				}
			}
		}

		if buf.indexSizeColumns != nil  && buf.rLines != nil {

			for buf.startLine < buf.endLine {

				//El buffer por referencia crea errores en los canales
				buf.Buffer = new([]byte)
				*buf.Buffer = make([]byte ,buf.sizeLine)
				//Leemos una sola linea
				_ , err = buf.file.ReadAt(*buf.Buffer , buf.lenFields + buf.startLine * buf.sizeLine)
				if err != nil {

					log.Println(err)

				}

				for _ , colName := range (*buf.colName) {

					size, found := buf.indexSizeColumns[colName]
					if found {

						bufferChan := make([]byte, (size[1] - size[0])  )
						bufferChan  = (*buf.Buffer)[size[0]:size[1]]
			
						//Limpiamos el buffer de espacios
						if buf.postFormat == true {
		
							buf.spaceTrimPointer(&bufferChan)

						}
						
						//Activamos PostFormat si existe
						if buf.hooker != nil && buf.postFormat == true {
						
							buf.hookerPostFormatPointer(&bufferChan ,colName)

						}
						
						
						buf.Channel <- RChanBuf{buf.startLine,colName,bufferChan}
					}

				}

				buf.startLine += 1
			}
		}

		//Cerramos el canal.
		close(buf.Channel)

		return
	}
	
	if checkFileTypeBuffer(buf.typeBuff , buffMap ){

		if buf.indexSizeFields != nil && buf.rRangues != nil {

			for _ , colName := range (*buf.colName) {
			
				size, found := buf.indexSizeFields[colName]
				if found {
				
					buf.rangue      = 0
					buf.totalRangue = 1

					//Obtimizacion especial para mapas.
					sizeTotal := size[1] - size[0]
					if  buf.rangeBytes < sizeTotal && buf.rangeBytes > 0 {

						TotalRangue := sizeTotal / buf.rangeBytes
						restoRangue := sizeTotal % buf.rangeBytes
						if restoRangue != 0 {

							TotalRangue += 1
						}
						buf.BufferMap[colName] = make([][]byte ,TotalRangue)

					} else {

						buf.BufferMap[colName] = make([][]byte , 1)
					}

					for buf.rangue < buf.totalRangue {

						buf.readIndexSizeFieldPointer(colName, size)
		
						buf.BufferMap[colName][buf.rangue] = append( buf.BufferMap[colName][buf.rangue],  *buf.FieldBuffer... )
						//buf.BufferMap[colName][buf.rangue]   =  *buf.FieldBuffer
				
						buf.rangue++

					}
					
				}
			}
		}
		
		if buf.indexSizeColumns != nil  && buf.rLines != nil {


			bufMapFile := &buf.BufferMap["buffer"][0]
			_ , err = buf.file.ReadAt(*bufMapFile , buf.lenFields + buf.startLine * buf.sizeLine )
			if err != nil {
				log.Println("Error: ",err)
			
			}

			var mapCount int64

			for buf.startLine < buf.endLine {
				

				for _ , colName := range (*buf.colName) {


					if colName == "buffer" {

						continue

					}

					size, found := buf.indexSizeColumns[colName]
					if found {

						if _, found :=  buf.BufferMap[colName]; !found {
						
							buf.BufferMap[colName] = make([][]byte ,(buf.endLine - buf.startLine))

						}
						

						//Limpiamos el buffer de espacios
						//buf.BufferMap[colName][buf.startLine] = append(buf.BufferMap[colName][buf.startLine],  (*bufMapFile)[size[0]:size[1]]... )
						buf.BufferMap[colName][mapCount] =  (*bufMapFile)[size[0]:size[1]]
						
						
						
					
						//Limpiamos el buffer de espacios
						if buf.postFormat == true {

							buf.spaceTrimPointer(&buf.BufferMap[colName][len(buf.BufferMap[colName])-1])

						}
					
						//Activamos PostFormat si existe
						if buf.hooker !=nil && buf.postFormat == true {

							buf.hookerPostFormatPointer(&buf.BufferMap[colName][len(buf.BufferMap[colName])-1], colName)
							
						}

					}
				}

				mapCount++
				buf.startLine++

				//End bucle
				*bufMapFile = (*bufMapFile)[ buf.sizeLine:]
			}
	
			delete(buf.BufferMap,"buffer")
		}

		return
	}


	log.Fatalln("Error Grave, Uspace.go ; Rspace.go ; Funcion: MultiColumnSpace ;" +
	"No Hubo coincidencias")
	return
}



func (buf *RBuffer) readBitSpace() {

	if checkFileTypeBuffer(buf.typeBuff , buffBytes ){

		if buf.indexSizeFields != nil && buf.rRangues != nil {

			size, found := buf.indexSizeFields[(*buf.colName)[0]]
			if found {

				buf.readIndexSizeFieldPointer((*buf.colName)[0], size)
			}
		}

		if buf.indexSizeColumns != nil  && buf.rLines != nil {

			size , found := buf.indexSizeColumns[(*buf.colName)[0]]
			if found {

				var byteLine int64  =  buf.startLine / 8
				var bitLine  int64  =  buf.startLine % 8 

				//bufferBit := make([]byte , 1 )

				*buf.Buffer = (*buf.Buffer)[:1]

				_ , err := buf.file.ReadAt(*buf.Buffer ,buf.lenFields + (byteLine * buf.sizeLine) + size[0])
				if err !=nil {

					*buf.Buffer = []byte("off")
					return
				}

				turn := readBit(bitLine,*buf.Buffer)

				if turn {

					*buf.Buffer = []byte("on")

				}else{

					*buf.Buffer = []byte("off")

				}
			}
		}
		return
	}
	
	if checkFileTypeBuffer(buf.typeBuff , buffChan ){

		if buf.indexSizeFields != nil && buf.rRangues != nil {

			for _ , colName := range (*buf.colName) {
			
				size, found := buf.indexSizeFields[colName]
				if found {
	
					buf.rangue      = 0
					buf.totalRangue = 1

					for buf.rangue < buf.totalRangue {

						buf.readIndexSizeFieldPointer(colName, size)

						buf.Channel <- RChanBuf{buf.rangue ,colName,*buf.FieldBuffer}

						buf.rangue++
					}
				}
			}
		}

		if buf.indexSizeColumns != nil  && buf.rLines != nil {

			for buf.startLine < buf.endLine {
					
				for _ , colName := range (*buf.colName) {

					size , found := buf.indexSizeColumns[colName]
					if found {

						var byteLine int64 =  buf.startLine / 8
						var bitLine  int64 =  buf.startLine % 8 


						_ , err := buf.file.ReadAt(*buf.Buffer , buf.lenFields + (byteLine * buf.sizeLine) + size[0])
						if err !=nil {

							buf.Channel <- RChanBuf{buf.startLine , colName, []byte("off")}
					
						}

						turn := readBit(bitLine,*buf.Buffer)

						if turn {

							buf.Channel <- RChanBuf{buf.startLine , colName, []byte("on")}
					
						}else{

							buf.Channel <- RChanBuf{buf.startLine , colName, []byte("off")}
							
						}
					}
				}
				buf.startLine++
			}

			close(buf.Channel)
		}
		return
	}
	
	if checkFileTypeBuffer(buf.typeBuff , buffMap ){

		if buf.indexSizeFields != nil && buf.rRangues != nil {

			for _ , colName := range (*buf.colName) {
			
				size, found := buf.indexSizeFields[colName]
				if found {
				
					buf.rangue      = 0
					buf.totalRangue = 1

					sizeTotal := size[1] - size[0]
					if  buf.rangeBytes < sizeTotal && buf.rangeBytes > 0 {

						TotalRangue := sizeTotal / buf.rangeBytes
						restoRangue := sizeTotal % buf.rangeBytes
						if restoRangue != 0 {

							TotalRangue += 1
						}

						buf.BufferMap[colName] = make([][]byte ,TotalRangue)

					} else {

						buf.BufferMap[colName] = make([][]byte ,1)

					}

					for buf.rangue < buf.totalRangue {

						buf.readIndexSizeFieldPointer(colName, size)
		
						buf.BufferMap[colName][buf.rangue] = append( buf.BufferMap[colName][buf.rangue],  *buf.FieldBuffer... )
				
						buf.rangue++

					}
					
				}
			}
		}


		if buf.indexSizeColumns != nil  && buf.rLines != nil {

			var mapCount int64
			bufferBit := &buf.BufferMap["buffer"][0]

			for buf.startLine < buf.endLine {
					
				var byteLine int64 =  buf.startLine / 8
				var bitLine  int64 =  buf.startLine % 8 
				

				for _ , colName := range (*buf.colName) {


					if colName == "buffer" {

						continue

					}

					size , found := buf.indexSizeColumns[colName]
					if found {

						if _, found :=  buf.BufferMap[colName]; !found {

							buf.BufferMap[colName] = make([][]byte , (buf.endLine - buf.startLine))

						}
							
						_ , err := buf.file.ReadAt(*bufferBit , buf.lenFields + (byteLine * buf.sizeLine) + size[0])
						if err != nil {
							
							buf.BufferMap[colName] = append(buf.BufferMap[colName] , []byte("off"))
						
						}

						turn := readBit(bitLine,*bufferBit)
						
						if turn {

							//buf.BufferMap[colName][mapCount] = append(buf.BufferMap[colName][mapCount] , []byte("on")...)
							buf.BufferMap[colName][mapCount] = []byte("on")
						}else{

							//buf.BufferMap[colName][mapCount] = append(buf.BufferMap[colName][mapCount] , []byte("off")...)
							buf.BufferMap[colName][mapCount] = []byte("off")
						}
					}
				}

				mapCount++
				buf.startLine++
			}
		}
		return
	}
	

	log.Fatalln("Error Grave, Uspace.go ; Rspace.go ; Funcion: ListBitSpace ;" +
	"No Hubo coincidencias")
	return
}

