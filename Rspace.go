package bd

import (
	"log"
)



func (buf *RBuffer) readByteSpace(){

	var err error
	startLine := buf.StartLine
	endLine   := buf.EndLine


		if CheckBit(int64(buf.typeBuff), int64(BuffBytes) ){

			if buf.IndexSizeFields != nil {

				size, found := buf.IndexSizeFields[buf.ColName]
				if found {

					_ , err = buf.File.ReadAt(buf.Buffer , size[0])
					if err != nil {
						log.Println(err)
						
					}

				}
			}

			if buf.IndexSizeColumns != nil {

				size, found := buf.IndexSizeColumns[buf.ColName]
				if found {
				
					_ , err = buf.File.ReadAt(buf.Buffer , buf.lenFields + (startLine * buf.SizeLine) + size[0])
					if err != nil {
						log.Println(err)
						
					}
				}
			}
	
			//Limpiamos el buffer de espacios
			if buf.PostFormat == true {

				buf.spaceTrimPointer(&buf.Buffer)

			}


			//Activamos PostFormat si existe
			if buf.Hooker != nil && buf.PostFormat == true {
				
				buf.hookerPostFormatPointer(&buf.Buffer ,buf.ColName)

			}

			//Cerramos el script y pasamos datos por referencia.
			return
		}
	

	if CheckBit(int64(buf.typeBuff), int64(BuffChan) ){


		if buf.IndexSizeFields != nil {

			for colName := range buf.BufferMap {
			
				size, found := buf.IndexSizeFields[colName]
				if found {

					sizeTotal := size[1] - size[0]


					if buf.rangeBytes < sizeTotal && buf.rangeBytes > 0 {
						
						totalRangue := sizeTotal / buf.rangeBytes
						restoRangue := sizeTotal % buf.rangeBytes

						var x int64
						for x = 0 ; x < totalRangue; x++ {

							fieldBuffer := make([]byte, buf.rangeBytes)
							_ , err = buf.File.ReadAt(fieldBuffer , size[0] + (buf.rangeBytes * x) )
							if err != nil {
								log.Println(err)
								
							}

							//Limpiamos el buffer de espacios
							if buf.PostFormat == true {

								buf.spaceTrimPointer(&fieldBuffer)

							}
							
							//Activamos PostFormat si existe
							if buf.Hooker != nil && buf.PostFormat == true {
							
								buf.hookerPostFormatPointer(&fieldBuffer ,colName)

							}

							buf.Channel <- RChanBuf{x ,colName,fieldBuffer}
						}
			
						
						if restoRangue != 0 {

			
							fieldBuffer := make([]byte, restoRangue)

							_ , err = buf.File.ReadAt(fieldBuffer , size[0] + (buf.rangeBytes * x) )
							if err != nil {

								log.Println(err)
								
							}

							//Limpiamos el buffer de espacios
							if buf.PostFormat == true {

								buf.spaceTrimPointer(&fieldBuffer)

							}
							
							//Activamos PostFormat si existe
							if buf.Hooker != nil && buf.PostFormat == true {
							
								buf.hookerPostFormatPointer(&fieldBuffer ,colName)

							}

							buf.Channel <- RChanBuf{x ,colName,fieldBuffer}
						}	
					}


					if buf.rangeBytes >= sizeTotal || buf.rangeBytes <= 0 {

						fieldBuffer := make([]byte, size[1] - size[0])
						_ , err = buf.File.ReadAt(fieldBuffer , size[0])
						if err != nil {
							log.Println(err)
							
						}
	
						//Limpiamos el buffer de espacios
						if buf.PostFormat == true {

							buf.spaceTrimPointer(&fieldBuffer)

						}
						
						//Activamos PostFormat si existe
						if buf.Hooker != nil && buf.PostFormat == true {
						
							buf.hookerPostFormatPointer(&fieldBuffer ,colName)

						}

						buf.Channel <- RChanBuf{0 ,colName,fieldBuffer}
	
					}
				}
			}
		}

		if buf.IndexSizeColumns != nil {

			for startLine < endLine {

				//El buffer por referencia crea errores en los canales
				buf.NewChanBuffer()
				//Leemos una sola linea
				_ , err = buf.File.ReadAt(buf.Buffer , buf.lenFields + startLine * buf.SizeLine)
				if err != nil {

					log.Println(err)

				}

				for colName := range buf.BufferMap {

					size, found := buf.IndexSizeColumns[colName]
					if found {

						bufferChan := make([]byte, (size[1] - size[0])  )
						bufferChan  = buf.Buffer[size[0]:size[1]]
			
						//Limpiamos el buffer de espacios
						if buf.PostFormat == true {
		
							buf.spaceTrimPointer(&bufferChan)

						}
						
						//Activamos PostFormat si existe
						if buf.Hooker != nil && buf.PostFormat == true {
						
							buf.hookerPostFormatPointer(&bufferChan ,colName)

						}
						
						
						buf.Channel <- RChanBuf{startLine,colName,bufferChan}
					}

				}


				startLine += 1
			}
		}


		//Cerramos el canal.
		close(buf.Channel)

		return
	}
	
	if CheckBit(int64(buf.typeBuff), int64(BuffMap) ){

	
		if buf.IndexSizeFields != nil {
			for colName := range buf.BufferMap {
			
				size, found := buf.IndexSizeFields[colName]
				if found {
					_ = size

				}
			}
		}
		
		if buf.IndexSizeColumns != nil {

			bufMapFile := &buf.BufferMap["buffer"][0]
			_ , err = buf.File.ReadAt(*bufMapFile , buf.lenFields + startLine * buf.SizeLine )
			if err != nil {
				log.Println("Error: ",err)
				//return
			}

			for startLine < endLine {
				
				startLine += 1

				for colName := range buf.BufferMap {


					if colName == "buffer" {

						continue

					}

					size, found := buf.IndexSizeColumns[colName]
					if found {

						//buf.BufferMap[val] = append(buf.BufferMap[val], bytes.Trim((*bufMapFile)[  buf.IndexSizeColumns[val][0] : buf.IndexSizeColumns[val][1]  ], " "))
					
						//Limpiamos el buffer de espacios
						buf.BufferMap[colName] = append(buf.BufferMap[colName],  (*bufMapFile)[size[0]:size[1]] )
						
						
						
					
						//Limpiamos el buffer de espacios
						if buf.PostFormat == true {

							buf.spaceTrimPointer(&buf.BufferMap[colName][len(buf.BufferMap[colName])-1])

						}
					
						//Activamos PostFormat si existe
						if buf.Hooker !=nil && buf.PostFormat == true {

							buf.hookerPostFormatPointer(&buf.BufferMap[colName][len(buf.BufferMap[colName])-1], colName)
							
						}

					}
				}
			
				//End bucle
				*bufMapFile = (*bufMapFile)[ buf.SizeLine:]
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

	if CheckBit(int64(buf.typeBuff), int64(BuffBytes) ){

		if buf.IndexSizeFields != nil {

			size, found := buf.IndexSizeFields[buf.ColName]
			if found {
			
				_ , err := buf.File.ReadAt(buf.Buffer , size[0])
				if err != nil {
					log.Println(err)
					return
				}
				return
			}
		}

		if buf.SizeLine == -2 || buf.EndLine == -2 {

			return

		}

		var byteLine int64 =  buf.StartLine / 8
	
		bufferBit := make([]byte , 1 )

		_ , err := buf.File.ReadAt(bufferBit ,buf.lenFields + byteLine)
		if err !=nil {

			buf.Buffer = []byte("off")
			return
		}

		var bitLine int64 =  buf.StartLine % 8 

		turn := readBit(bitLine,bufferBit)

		if turn {

			buf.Buffer = []byte("on")

		}else{

			buf.Buffer = []byte("off")

		}
		return
	}
	
	if CheckBit(int64(buf.typeBuff), int64(BuffChan) ){
	
		var byteLine int64 =  buf.StartLine / 8
	
		bufferBit := make([]byte , 1 )

		_ , err := buf.File.ReadAt(bufferBit , buf.lenFields + byteLine)
		if err !=nil {

			buf.Channel <- RChanBuf{0 , "ListBit", []byte("off")}
			close(buf.Channel)
			return
		}

		var bitLine int64 =  buf.StartLine % 8 

		turn := readBit(bitLine,bufferBit)

		if turn {

			buf.Channel <- RChanBuf{1 , "ListBit", []byte("on")}
			close(buf.Channel)

		}else{

			buf.Channel <- RChanBuf{0 , "ListBit", []byte("off")}
			close(buf.Channel)
		}
		return
	}
	
	if CheckBit(int64(buf.typeBuff), int64(BuffMap) ){
	
		for val, ind := range buf.IndexSizeColumns {

			if ind[0] == 0 {
	
				var byteLine int64 =  buf.StartLine / 8
	
				bufferBit := make([]byte , 1 )
	
				_ , err := buf.File.ReadAt(bufferBit , buf.lenFields + byteLine)
				if err !=nil {
	
					buf.BufferMap[val] = append(buf.BufferMap[val] , []byte("off"))
					return
				}
	
				var bitLine int64 =  buf.StartLine % 8 
	
				turn := readBit(bitLine,bufferBit)
	
				if turn {
	
					buf.BufferMap[val] = append(buf.BufferMap[val] , []byte("on"))
	
				}else{
	
					buf.BufferMap[val] = append(buf.BufferMap[val] , []byte("off"))
	
				}
			}
			return
		}
	}

	

	log.Fatalln("Error Grave, Uspace.go ; Rspace.go ; Funcion: ListBitSpace ;" +
	"No Hubo coincidencias")
	return
}

