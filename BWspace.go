package bd

import (
	"sync/atomic"
)

//El canal de escritura linea , nombre de columna y el buffer.
type WChanBuf struct{
	Line int64
	ColName string
	Buffer 	[]byte
}

//Buffer de escritura con tres tipos de buffer
//Tipo buffer unicamente puede escribir en una columna o un field.
//Tipo mapaBuffer puede escribir simultaneametne en columnas y fields.
//Abre un canal que puede actualizar tantno columnas como fields.
type WBuffer struct {
	*spaceFile
	Line int64
	ColumnName string
	typeBuff FileTypeBuffer

	Buffer *[]byte
	BufferMap map[string][]byte
	Channel chan WChanBuf
}

//La funcion crea un nuevo buffer de escritura y lo devuelve su referencia.
//Esta funcion usa una referencia de bytes para escribir en un archivo.
func (sp *spaceFile) BWspaceBuff(line int64 ,columnName string, bufferBytes *[]byte)(*WBuffer){

	sp.check(columnName, "Archivo: BWspace.go ; Funcion: BWspaceBuff")

	return &WBuffer{
		spaceFile: sp,
		typeBuff: BuffBytes,
		Line: line,
		ColumnName: columnName,
		Buffer: bufferBytes,
	}

}
//La funcion crea un buffer de escritura y devuelve su referencia
//La funcion usa un mapa que se pasa por referencia para escribir en un archivo.
func (sp *spaceFile) BWspaceBuffMap(line int64 , bufferMap map[string][]byte)(*WBuffer){

	for columnName := range bufferMap {
	
		sp.check(columnName, "Archivo: BWspace.go ; Funcion: BWspaceBuffMap")

	}

   return &WBuffer{
		spaceFile: sp,
		typeBuff: BuffMap,
	    Line: line,
	    BufferMap: bufferMap,
   }

}

//Esta funcion abre un canal
//Primer uso abrir un canal en tiempo de ejecucion y usarlo despues en las rutas.
//Segundo uso usarlo en una ruta y mandar los valores que se quieran actualizar de forma
//dinamica
func (sp *spaceFile)BWChanBuf()(*WBuffer){

	return &WBuffer{
		spaceFile: sp,
		typeBuff: BuffChan,
		Channel: make(chan WChanBuf, 0),
	}
}

//Envia un buffer por el canal
func (WBuffer *WBuffer)BWspaceSendchan(line int64, columnName string , bufferChan *[]byte)int64{


	WBuffer.check(columnName, "Archivo: BWspace.go ; Funcion: BWspaceSendchan")

	//Actualización de campos sin lineas.
	if line != -2 {

		if line == -1 {
		
			line = atomic.AddInt64(WBuffer.SizeFileLine, 1)

		}

		if line > *WBuffer.SizeFileLine {

			atomic.AddInt64(WBuffer.SizeFileLine, line - *WBuffer.SizeFileLine )
			
		}
	}


	WBuffer.Channel <- WChanBuf{line,columnName, *bufferChan }

	return line
}

//Cierra un canal.
func (WBuffer *WBuffer) BWspaceClosechan(){
	
	close(WBuffer.Channel)
}

