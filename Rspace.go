package bd

import (
	"bytes"
	"log"
	"os"
	"strconv"


)


func (sp *Space) OneColumnSpace(buf Buffer){
	
	startLine := buf.StartLine
	endLine   := buf.EndLine

	_ , sp.err = sp.File.ReadAt(buf.Buffer["buffer"][0] , startLine * sp.Size_line )

	if (endLine - startLine ) == 1 {

		for val := range buf.Buffer {

			if val == "buffer" {
				continue
			}

			buf.Buffer[val] = append(buf.Buffer[val], bytes.Trim(buf.Buffer["buffer"][0] , " "))
			
			//Postformat por columnas
			function, exist := sp.Hooker[Postformat + val]
			if exist{

				buf.Buffer[val][len(buf.Buffer[val])-1] = function(buf.Buffer[val][len(buf.Buffer[val])-1])

			} else {

				//Postformat global
				function, exist = sp.Hooker[Postformat]
				if exist {

					buf.Buffer[val][len(buf.Buffer[val])-1] = function(buf.Buffer[val][len(buf.Buffer[val])-1])

				}
			}
		}
	}

	if (endLine - startLine) > 1 {


		for startLine < endLine {
			
			startLine += 1
	
			for val := range buf.Buffer {

				if val == "buffer" {
						continue
					}
	
				buf.Buffer[val] = append(buf.Buffer[val], bytes.Trim(buf.Buffer["buffer"][0][:sp.Size_line] , " "))

				//Postformat por columnas
				function, exist := sp.Hooker[Postformat + val]
				if exist{

					buf.Buffer[val][len(buf.Buffer[val])-1] = function(buf.Buffer[val][len(buf.Buffer[val])-1])

				} else {

					//Postformat global
					function, exist = sp.Hooker[Postformat]
					if exist {

						buf.Buffer[val][len(buf.Buffer[val])-1] = function(buf.Buffer[val][len(buf.Buffer[val])-1])

					}
				}

				
			}
	
			//End bucle
			buf.Buffer["buffer"][0] = buf.Buffer["buffer"][0][sp.Size_line:]
		
		}
		
	}

	delete(buf.Buffer,"buffer")
	return
	
}



func (sp *Space) MultiColumnSpace(buf Buffer){

	startLine := buf.StartLine
	endLine   := buf.EndLine

	_ , sp.err = sp.File.ReadAt(buf.Buffer["buffer"][0] , startLine * sp.Size_line )
	
	
	for startLine < endLine {
		
		startLine += 1

		for val := range buf.Buffer {

			if val == "buffer" {
				continue
			}

			buf.Buffer[val] = append(buf.Buffer[val], bytes.Trim(buf.Buffer["buffer"][0][  sp.IndexSizeColumns[val][0] : sp.IndexSizeColumns[val][1]  ], " "))
		
			//Postformat por columnas
			function, exist := sp.Hooker[Postformat + val]
			if exist{

				buf.Buffer[val][len(buf.Buffer[val])-1] = function(buf.Buffer[val][len(buf.Buffer[val])-1])

			} else {

				//Postformat global
				function, exist = sp.Hooker[Postformat]
				if exist {

					buf.Buffer[val][len(buf.Buffer[val])-1] = function(buf.Buffer[val][len(buf.Buffer[val])-1])

				}
			}
		}
	
		//End bucle
		buf.Buffer["buffer"][0] = buf.Buffer["buffer"][0][sp.Size_line:]
	
	}

	delete(buf.Buffer,"buffer")

	return
}

func (sp *Space) FullFileSpace(buf Buffer){


	buf.Buffer["buffer"][0], sp.err = os.ReadFile(sp.Url + "/" +  strconv.FormatInt( buf.StartLine ,10) + sp.Extension)

	if sp.err != nil {

		buf.Buffer["buffer"][0] = []byte{}

	}

}


func (sp *Space) ListBitSpace(buf Buffer){



	for val, ind := range sp.IndexSizeColumns {

		if ind[0] == 0 {


	var byteLine int64 =  buf.StartLine / 8


	bufferBit := make([]byte , 1 )

	_ , err := sp.File.ReadAt(bufferBit , byteLine)
	if err !=nil {

		buf.Buffer[val] = append(buf.Buffer[val] , []byte("off"))
		return
	}

	var bitLine int64 =  buf.StartLine % 8 

	turn := sp.readBit(bitLine,bufferBit)

			if turn {

				buf.Buffer[val] = append(buf.Buffer[val] , []byte("on"))

			}else{

				buf.Buffer[val] = append(buf.Buffer[val] , []byte("off"))

			}
		}
	}

	log.Println("buffer Rspace: ", string(buf.Buffer["buffer"][0])  )

}

func (sp *Space) ReadEmptyDirSpace(buf Buffer){

	buf.Buffer["buffer"][0], sp.err = os.ReadFile(sp.Name + "/" +  strconv.FormatInt( buf.StartLine ,10) + sp.Extension)

	if sp.err != nil {

		buf.Buffer["buffer"][0] = []byte{}

	}

}