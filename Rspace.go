package bd

import (
	"bytes"
	"log"
	"os"
)

/*
if CheckBit(int64(buf.typeBuff), int64(BuffBytes) ){

	return
}

if CheckBit(int64(buf.typeBuff), int64(BuffChan) ){

	return
}

if CheckBit(int64(buf.typeBuff), int64(BuffMap) ){
}
*/


func (sp *spaceFile) OneColumnSpace(buf *Buffer){
	
	var err error
	startLine := buf.StartLine
	endLine   := buf.EndLine


	//Primer caso cuando solo hay que leer una linea.
	if (endLine - startLine ) == 1 {


		if CheckBit(int64(buf.typeBuff), int64(BuffBytes) ){

			var colName string
			for val := range buf.BufferMap{
		
				colName = val		
			
			}

			//Leemos una sola linea
			_ , err = sp.File.ReadAt(buf.Buffer , startLine * sp.SizeLine )
			if err != nil {
				log.Println(err)
				return
			}

			//Limpiamos espacios en blanco
			buf.Buffer = bytes.Trim(buf.Buffer , " ")

			//Activamos PostFormat si existe
			if sp.Hooker != nil {
				
				sp.hookerPostFormatBuff(buf,colName)

			}

			//Cerramos el script y pasamos datos por referencia.
			return
		}


		if CheckBit(int64(buf.typeBuff), int64(BuffChan) ){

			var colName string
			for val := range buf.BufferMap{
		
				colName = val		
			
			}

			//Leemos una sola linea
			_ , err = sp.File.ReadAt(buf.Buffer , startLine * sp.SizeLine )
			if err != nil {
				log.Println(err)

			}

			//Limpiamos el buffer de espacios
			buf.Buffer = bytes.Trim(buf.Buffer , " ")

			//Activamos PostFormat si existe
			if sp.Hooker != nil {
	
				for val , ind := range sp.IndexSizeColumns{
					
					if ind[0] == 0 {

						sp.hookerPostFormatBuff(buf,val)

					}
				}
			}

			//Pasamos el buffer por el canal
			buf.Channel <- ChanBuf{startLine,colName,buf.Buffer}

			//Cerramos el canal.
			close(buf.Channel)
			return
		
		}

		if CheckBit(int64(buf.typeBuff), int64(BuffMap) ){

			//Leemos una sola linea en el buffer del mapa
			//La razon de esto es que es mas eficiente que usar el otro buffer.
			_ , err = sp.File.ReadAt(buf.BufferMap["buffer"][0] , startLine * sp.SizeLine )
			if err != nil {
				log.Println(err)
				return
			}

			//Reccorremos los valores que nos han pedido
			for val := range buf.BufferMap {

				//Si el valor es buffer, nos lo saltamos
				if val == "buffer" {
					continue
				}

				//Anexamos el valor al mapa
				buf.BufferMap[val] = append(buf.BufferMap[val], bytes.Trim(buf.BufferMap["buffer"][0] , " "))
			
				//Activamos PostFormat si existe
				if sp.Hooker !=nil {
					sp.hookerPostFormatMap(buf,val)
				}
				 
			}

			//No lo borro para que se pueda reutilizar en multiples lecturas.
			//delete(buf.BufferMap,"buffer")
			return
		}
	}




	
	//Segundo caso cuando hay que leer mas de una linea
	if (endLine - startLine) > 1 {

		if CheckBit(int64(buf.typeBuff), int64(BuffBytes) ){

			log.Fatalln("No se puede leer multiples lineas en un buffer de bytes.")
			return
		}


		if CheckBit(int64(buf.typeBuff), int64(BuffChan) ){

			var colName string
			for val := range buf.BufferMap{
		
				colName = val		
			
			}

			for startLine < endLine {

				//El buffer por referencia crea errores en los canales
				buf.NewChanBuffer()
				//Leemos una sola linea
				_ , err = sp.File.ReadAt(buf.Buffer , startLine * sp.SizeLine )
				if err != nil {
					log.Println(err)

				}

				//Limpiamos el buffer de espacios
				buf.Buffer = bytes.Trim(buf.Buffer , " ")

				
				//Activamos PostFormat si existe
				if sp.Hooker != nil {

					for val , ind := range sp.IndexSizeColumns{
						
						if ind[0] == 0 {

							sp.hookerPostFormatBuff(buf,val)

						}
					}
				}
				
				buf.Channel <- ChanBuf{startLine,colName,buf.Buffer}

				startLine += 1
				
			}
			//Cerramos el canal.
			close(buf.Channel)
			return
		}
		


		if CheckBit(int64(buf.typeBuff), int64(BuffMap) ){
		
			if buf.BufferMap["buffer"] == nil {
				log.Fatalln("No se puede reciclar un buffer de mapa multilinea y multicolumna")
				return
			}

			_ , err = sp.File.ReadAt(buf.BufferMap["buffer"][0] , startLine * sp.SizeLine )
			if err != nil {
				log.Println(err)
				return
			}

			for startLine < endLine {
		
				startLine += 1
		
				for val := range buf.BufferMap {
	
					if val == "buffer" {
							continue
						}
		
					buf.BufferMap[val] = append(buf.BufferMap[val], bytes.Trim(buf.BufferMap["buffer"][0][:sp.SizeLine] , " "))
	
					//Activamos PostFormat si existe
					if sp.Hooker !=nil {
						sp.hookerPostFormatMap(buf,val)
					}
	
					
				}
	
				//Cuando terminamos el bucle hemos terminado de hacer range a
				//todos los campos de esa linea entonces borramos la linea del buffer.
				buf.BufferMap["buffer"][0] = buf.BufferMap["buffer"][0][sp.SizeLine:]
			
			}
			//Borramos el buffer
			delete(buf.BufferMap,"buffer")
			return
		}

		
	}



	log.Fatalln("Error Grave, Uspace.go ; Rspace.go ; Funcion: OneColumnSpace ;" +
	"No Hubo coincidencias")
	return
	
}



func (sp *spaceFile) MultiColumnSpace(buf *Buffer){

	var err error
	startLine := buf.StartLine
	endLine   := buf.EndLine


	if CheckBit(int64(buf.typeBuff), int64(BuffBytes) ){
		log.Fatalln("Error Grave, Uspace.go ; Rspace.go ; Funcion: MultiColumnSpace ;" +
		"Buffer no soporta multicolumnas")
		return
	}
	

	if CheckBit(int64(buf.typeBuff), int64(BuffChan) ){

		for startLine < endLine {

			//El buffer por referencia crea errores en los canales
			buf.NewChanBuffer()
			//Leemos una sola linea
			_ , err = sp.File.ReadAt(buf.Buffer , startLine * sp.SizeLine)
			if err != nil {

				log.Println(err)

			}

			for val := range buf.BufferMap {
			
				value := make([]byte,0)
				value = append(buf.Buffer[sp.IndexSizeColumns[val][0]:sp.IndexSizeColumns[val][1]])
	
				//Limpiamos el buffer de espacios
				value = bytes.Trim(value , " ")

				
				//Activamos PostFormat si existe
				if sp.Hooker != nil {

					sp.hookerPostFormatBuffMultiColumn(&value,val)

				}
				
				log.Println("Multicolumn RAmchan: ",startLine,val,value)
				buf.Channel <- ChanBuf{startLine,val,value}

				
			}


			startLine += 1
		}
	
		//Cerramos el canal.
		close(buf.Channel)

		return
	}
	
	if CheckBit(int64(buf.typeBuff), int64(BuffMap) ){


		_ , err = sp.File.ReadAt(buf.BufferMap["buffer"][0] , startLine * sp.SizeLine )
		if err != nil {
			log.Println(err)
			return
		}
		
		for startLine < endLine {
			
			startLine += 1

			for val := range buf.BufferMap {

				if val == "buffer" {
					continue
				}

				buf.BufferMap[val] = append(buf.BufferMap[val], bytes.Trim(buf.BufferMap["buffer"][0][  sp.IndexSizeColumns[val][0] : sp.IndexSizeColumns[val][1]  ], " "))
	
				//Activamos PostFormat si existe
				if sp.Hooker !=nil {
					sp.hookerPostFormatMap(buf,val)
				}

			}
		
			//End bucle
			buf.BufferMap["buffer"][0] = buf.BufferMap["buffer"][0][sp.SizeLine:]
		
		}

		delete(buf.BufferMap,"buffer")
		return
	}


	log.Fatalln("Error Grave, Uspace.go ; Rspace.go ; Funcion: MultiColumnSpace ;" +
	"No Hubo coincidencias")
	return
}

func (sp *spaceFile) FullFileSpace(buf *Buffer){

	var err error

	buf.BufferMap["buffer"][0], err = os.ReadFile(sp.Url)

	if err != nil {

		buf.BufferMap["buffer"][0] = []byte{}

	}

}


func (sp *spaceFile) ListBitSpace(buf *Buffer){



	for val, ind := range sp.IndexSizeColumns {

		if ind[0] == 0 {


	var byteLine int64 =  buf.StartLine / 8


	bufferBit := make([]byte , 1 )

	_ , err := sp.File.ReadAt(bufferBit , byteLine)
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
	}

	log.Println("buffer Rspace: ", string(buf.BufferMap["buffer"][0])  )

}

func (sp *spaceFile) ReadEmptyDirSpace(buf *Buffer){

	var err error
	//buf.BufferMap["buffer"][0], err = os.ReadFile(sp.Name + "/" +  strconv.FormatInt( buf.StartLine ,10) + sp.Extension)
	buf.BufferMap["buffer"][0], err = os.ReadFile(sp.Url)
	if err != nil {

		buf.BufferMap["buffer"][0] = []byte{}

	}

}