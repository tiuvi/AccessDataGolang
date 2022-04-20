package bd

import (
	"log"
)


func (buf *RBuffer) readByteSpace(){

	var err error
	
	
	if CheckFileTypeBuffer(buf.typeBuff , BuffBytes ){

		if buf.IndexSizeFields != nil {

			size, found := buf.IndexSizeFields[(*buf.ColName)[0]]
			if found {

				buf.readIndexSizeFieldPointer((*buf.ColName)[0], size)
				return
			}
		}

		if buf.IndexSizeColumns != nil {

			size, found := buf.IndexSizeColumns[(*buf.ColName)[0]]
			if found {

				_ , err = buf.File.ReadAt(*buf.Buffer , buf.lenFields + (buf.StartLine * buf.SizeLine) + size[0])
				if err != nil {
					log.Println(err)
					
				}
				
				//Limpiamos el buffer de espacios
				if buf.PostFormat == true {

					buf.spaceTrimPointer(buf.Buffer)

				}


				//Activamos PostFormat si existe
				if buf.Hooker != nil && buf.PostFormat == true {
					
					buf.hookerPostFormatPointer(buf.Buffer ,(*buf.ColName)[0])

				}
			}
			return
		}


		//Cerramos el script y pasamos datos por referencia.
		return
	}
	
	
	if CheckFileTypeBuffer(buf.typeBuff , BuffChan ){


		if buf.IndexSizeFields != nil {

			for _ , colName := range (*buf.ColName) {
			
				size, found := buf.IndexSizeFields[colName]
				if found {
	
			
					buf.Rangue      = 0
					buf.TotalRangue = 1

					for buf.Rangue < buf.TotalRangue {

						buf.readIndexSizeFieldPointer(colName, size)

						buf.Channel <- RChanBuf{buf.Rangue ,colName,*buf.FieldBuffer}

						buf.Rangue++
					}
				}
			}
		}

		if buf.IndexSizeColumns != nil {

			for buf.StartLine < buf.EndLine {

				//El buffer por referencia crea errores en los canales
				buf.NewChanBuffer()
				//Leemos una sola linea
				_ , err = buf.File.ReadAt(*buf.Buffer , buf.lenFields + buf.StartLine * buf.SizeLine)
				if err != nil {

					log.Println(err)

				}

				for _ , colName := range (*buf.ColName) {

					size, found := buf.IndexSizeColumns[colName]
					if found {

						bufferChan := make([]byte, (size[1] - size[0])  )
						bufferChan  = (*buf.Buffer)[size[0]:size[1]]
			
						//Limpiamos el buffer de espacios
						if buf.PostFormat == true {
		
							buf.spaceTrimPointer(&bufferChan)

						}
						
						//Activamos PostFormat si existe
						if buf.Hooker != nil && buf.PostFormat == true {
						
							buf.hookerPostFormatPointer(&bufferChan ,colName)

						}
						
						
						buf.Channel <- RChanBuf{buf.StartLine,colName,bufferChan}
					}

				}


				buf.StartLine += 1
			}
		}


		//Cerramos el canal.
		close(buf.Channel)

		return
	}
	
	if CheckFileTypeBuffer(buf.typeBuff , BuffMap ){


		if buf.IndexSizeFields != nil {
			for _ , colName := range (*buf.ColName) {
			
				size, found := buf.IndexSizeFields[colName]
				if found {
				
					buf.Rangue      = 0
					buf.TotalRangue = 1

					for buf.Rangue < buf.TotalRangue {

						buf.readIndexSizeFieldPointer(colName, size)
					
						buf.BufferMap[colName]       = append(buf.BufferMap[colName],  *buf.FieldBuffer )

						buf.Rangue++

					}
					
				}
			}
		}
		
		if buf.IndexSizeColumns != nil {

			bufMapFile := &buf.BufferMap["buffer"][0]
			_ , err = buf.File.ReadAt(*bufMapFile , buf.lenFields + buf.StartLine * buf.SizeLine )
			if err != nil {
				log.Println("Error: ",err)
			
			}

			for buf.StartLine < buf.EndLine {
				
				buf.StartLine += 1

				for _ , colName := range (*buf.ColName) {


					if colName == "buffer" {

						continue

					}

					size, found := buf.IndexSizeColumns[colName]
					if found {

						if _, found :=  buf.BufferMap[colName]; !found {

							buf.BufferMap[colName] = make([][]byte ,0)

						}
						

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

	if CheckFileTypeBuffer(buf.typeBuff , BuffBytes ){

		if buf.IndexSizeFields != nil {

			size, found := buf.IndexSizeFields[(*buf.ColName)[0]]
			if found {

				buf.readIndexSizeFieldPointer((*buf.ColName)[0], size)
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

			*buf.Buffer = []byte("off")
			return
		}

		var bitLine int64 =  buf.StartLine % 8 

		turn := readBit(bitLine,bufferBit)

		if turn {

			*buf.Buffer = []byte("on")

		}else{

			*buf.Buffer = []byte("off")

		}
		return
	}
	
	if CheckFileTypeBuffer(buf.typeBuff , BuffChan ){


		if buf.IndexSizeFields != nil {

			for _ , colName := range (*buf.ColName) {
			
				size, found := buf.IndexSizeFields[colName]
				if found {
	
					buf.Rangue      = 0
					buf.TotalRangue = 1

					for buf.Rangue < buf.TotalRangue {

						buf.readIndexSizeFieldPointer(colName, size)

						buf.Channel <- RChanBuf{buf.Rangue ,colName,*buf.FieldBuffer}

						buf.Rangue++
					}
				}
			}
		}


		for _ , colName := range (*buf.ColName) {

			_, found := buf.IndexSizeColumns[colName]
			if found {

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
			
				}else{

					buf.Channel <- RChanBuf{0 , "ListBit", []byte("off")}
					
				}
			}
		}

		close(buf.Channel)
		return
	}
	
	if CheckFileTypeBuffer(buf.typeBuff , BuffMap ){

	
		if buf.IndexSizeFields != nil {
			for _ , colName := range (*buf.ColName) {
			
				size, found := buf.IndexSizeFields[colName]
				if found {
				
					buf.Rangue      = 0
					buf.TotalRangue = 1

					for buf.Rangue < buf.TotalRangue {

						buf.readIndexSizeFieldPointer(colName, size)
					
						buf.BufferMap[colName]       = append(buf.BufferMap[colName],  *buf.FieldBuffer )

						buf.Rangue++

					}
					
				}
			}
		}

		for _ , colName := range (*buf.ColName) {


			if colName == "buffer" {

				continue

			}

			_, found := buf.IndexSizeColumns[colName]
			if found {

				if _, found :=  buf.BufferMap[colName]; !found {

					buf.BufferMap[colName] = make([][]byte ,0)

				}
	
				var byteLine int64 =  buf.StartLine / 8

				bufferBit := make([]byte , 1 )

				_ , err := buf.File.ReadAt(bufferBit , buf.lenFields + byteLine)
				if err != nil {

					buf.BufferMap[colName] = append(buf.BufferMap[colName] , []byte("off"))
					return
				}

				var bitLine int64 =  buf.StartLine % 8 

				turn := readBit(bitLine,bufferBit)

				if turn {

					buf.BufferMap[colName] = append(buf.BufferMap[colName] , []byte("on"))

				}else{

					buf.BufferMap[colName] = append(buf.BufferMap[colName] , []byte("off"))

				}
			}
		}
		
		return
	}

	

	log.Fatalln("Error Grave, Uspace.go ; Rspace.go ; Funcion: ListBitSpace ;" +
	"No Hubo coincidencias")
	return
}

