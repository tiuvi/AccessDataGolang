package bd

import (
	"log"
)


func (buf *RBuffer) readByteSpace(){

	var err error
	
	
	if CheckFileTypeBuffer(buf.typeBuff , BuffBytes ){

		if buf.IndexSizeFields != nil && buf.RRangues != nil {

			size, found := buf.IndexSizeFields[(*buf.ColName)[0]]
			if found {

				buf.readIndexSizeFieldPointer((*buf.ColName)[0], size)
		
			}
		}

		if buf.IndexSizeColumns != nil  && buf.RLines != nil {

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
		}
		return
	}
	
	
	if CheckFileTypeBuffer(buf.typeBuff , BuffChan ){


		if buf.IndexSizeFields != nil && buf.RRangues != nil {

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

		if buf.IndexSizeColumns != nil  && buf.RLines != nil {

			for buf.StartLine < buf.EndLine {

				//El buffer por referencia crea errores en los canales
				buf.Buffer = new([]byte)
				*buf.Buffer = make([]byte ,buf.SizeLine)
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

		if buf.IndexSizeFields != nil && buf.RRangues != nil {

			for _ , colName := range (*buf.ColName) {
			
				size, found := buf.IndexSizeFields[colName]
				if found {
				
					buf.Rangue      = 0
					buf.TotalRangue = 1

					//Obtimizacion especial para mapas.
					sizeTotal := size[1] - size[0]
					if  buf.RangeBytes < sizeTotal && buf.RangeBytes > 0 {

						TotalRangue := sizeTotal / buf.RangeBytes
						restoRangue := sizeTotal % buf.RangeBytes
						if restoRangue != 0 {

							TotalRangue += 1
						}
						buf.BufferMap[colName] = make([][]byte ,TotalRangue)

					} else {

						buf.BufferMap[colName] = make([][]byte , 1)
					}

					for buf.Rangue < buf.TotalRangue {

						buf.readIndexSizeFieldPointer(colName, size)
		
						buf.BufferMap[colName][buf.Rangue] = append( buf.BufferMap[colName][buf.Rangue],  *buf.FieldBuffer... )
						//buf.BufferMap[colName][buf.Rangue]   =  *buf.FieldBuffer
				
						buf.Rangue++

					}
					
				}
			}
		}
		
		if buf.IndexSizeColumns != nil  && buf.RLines != nil {


			bufMapFile := &buf.BufferMap["buffer"][0]
			_ , err = buf.File.ReadAt(*bufMapFile , buf.lenFields + buf.StartLine * buf.SizeLine )
			if err != nil {
				log.Println("Error: ",err)
			
			}

			var mapCount int64

			for buf.StartLine < buf.EndLine {
				

				for _ , colName := range (*buf.ColName) {


					if colName == "buffer" {

						continue

					}

					size, found := buf.IndexSizeColumns[colName]
					if found {

						if _, found :=  buf.BufferMap[colName]; !found {
						
							buf.BufferMap[colName] = make([][]byte ,(buf.EndLine - buf.StartLine))

						}
						

						//Limpiamos el buffer de espacios
						//buf.BufferMap[colName][buf.StartLine] = append(buf.BufferMap[colName][buf.StartLine],  (*bufMapFile)[size[0]:size[1]]... )
						buf.BufferMap[colName][mapCount] =  (*bufMapFile)[size[0]:size[1]]
						
						
						
					
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

				mapCount++
				buf.StartLine++

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

		if buf.IndexSizeFields != nil && buf.RRangues != nil {

			size, found := buf.IndexSizeFields[(*buf.ColName)[0]]
			if found {

				buf.readIndexSizeFieldPointer((*buf.ColName)[0], size)
			}
		}

		if buf.IndexSizeColumns != nil  && buf.RLines != nil {

			size , found := buf.IndexSizeColumns[(*buf.ColName)[0]]
			if found {

				var byteLine int64  =  buf.StartLine / 8
				var bitLine  int64  =  buf.StartLine % 8 

				//bufferBit := make([]byte , 1 )

				*buf.Buffer = (*buf.Buffer)[:1]

				_ , err := buf.File.ReadAt(*buf.Buffer ,buf.lenFields + byteLine + size[0])
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
	
	if CheckFileTypeBuffer(buf.typeBuff , BuffChan ){

		if buf.IndexSizeFields != nil && buf.RRangues != nil {

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

		if buf.IndexSizeColumns != nil  && buf.RLines != nil {

			for buf.StartLine < buf.EndLine {
					
				for _ , colName := range (*buf.ColName) {

					size , found := buf.IndexSizeColumns[colName]
					if found {

						var byteLine int64 =  buf.StartLine / 8
						var bitLine  int64 =  buf.StartLine % 8 


						_ , err := buf.File.ReadAt(*buf.Buffer , buf.lenFields + byteLine + size[0])
						if err !=nil {

							buf.Channel <- RChanBuf{buf.StartLine , colName, []byte("off")}
							close(buf.Channel)
							return
						}

						turn := readBit(bitLine,*buf.Buffer)

						if turn {

							buf.Channel <- RChanBuf{buf.StartLine , colName, []byte("on")}
					
						}else{

							buf.Channel <- RChanBuf{buf.StartLine , colName, []byte("off")}
							
						}
					}
				}
				buf.StartLine++
			}

			close(buf.Channel)
		}
		return
	}
	
	if CheckFileTypeBuffer(buf.typeBuff , BuffMap ){

		if buf.IndexSizeFields != nil && buf.RRangues != nil {

			for _ , colName := range (*buf.ColName) {
			
				size, found := buf.IndexSizeFields[colName]
				if found {
				
					buf.Rangue      = 0
					buf.TotalRangue = 1

					sizeTotal := size[1] - size[0]
					if  buf.RangeBytes < sizeTotal && buf.RangeBytes > 0 {

						TotalRangue := sizeTotal / buf.RangeBytes
						restoRangue := sizeTotal % buf.RangeBytes
						if restoRangue != 0 {

							TotalRangue += 1
						}

						buf.BufferMap[colName] = make([][]byte ,TotalRangue)

					} else {

						buf.BufferMap[colName] = make([][]byte ,1)

					}

					for buf.Rangue < buf.TotalRangue {

						buf.readIndexSizeFieldPointer(colName, size)
		
						buf.BufferMap[colName][buf.Rangue] = append( buf.BufferMap[colName][buf.Rangue],  *buf.FieldBuffer... )
				
						buf.Rangue++

					}
					
				}
			}
		}


		if buf.IndexSizeColumns != nil  && buf.RLines != nil {

			var mapCount int64
			bufferBit := &buf.BufferMap["buffer"][0]

			for buf.StartLine < buf.EndLine {
					
				var byteLine int64 =  buf.StartLine / 8
				var bitLine  int64 =  buf.StartLine % 8 
				

				for _ , colName := range (*buf.ColName) {


					if colName == "buffer" {

						continue

					}

					size , found := buf.IndexSizeColumns[colName]
					if found {

						if _, found :=  buf.BufferMap[colName]; !found {

							buf.BufferMap[colName] = make([][]byte , (buf.EndLine - buf.StartLine))

						}
							
						_ , err := buf.File.ReadAt(*bufferBit , buf.lenFields + byteLine + size[0])
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
				buf.StartLine++
			}
		}
		return
	}
	

	log.Fatalln("Error Grave, Uspace.go ; Rspace.go ; Funcion: ListBitSpace ;" +
	"No Hubo coincidencias")
	return
}

