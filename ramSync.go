package dac

import (
//	"log"
)

type SpaceSGMStr struct {
	*spaceFile
	colName string
	size    [2]int64
	Map     map[string]int64
}

//SGMapStringLineInit: Sincroniza un mapa string/int64 superglobal con un archivo.
//#bd/ramSync/SGMapStringLineInit
func (sF *spaceFile) SGMapStringLineInit(colName string) *SpaceSGMStr {

	size, found := sF.indexSizeColumns[colName]

	if EDAC &&
		sF.ECSD(!found, "Columna no encontrada.") {
	}

	if found {

		//Creamos un puntero a la estructura.
		SGMS := &SpaceSGMStr{
			spaceFile: sF,
			colName:   colName,
			size:      size,
			Map:       make(map[string]int64),
		}

		//Activamos candados de lectura y escritura.
		sF.Lock()
		defer sF.Unlock()

		//Leemos el fichero completo y desactivamos postformat.
		//Estos datos son preformateados.
		mapColumn := SGMS.GetAllLines(colName)

		//Guardamos el fichero en un mapa
		var x int64
		for x = 0; x <= *sF.sizeFileLine; x++ {

			//Borramos los espacios a la derecha
			SGMS.spaceTrimPointer(&mapColumn.BufferMap[colName][x])

			mapString := string(mapColumn.BufferMap[colName][x])

			if mapString == "" {

				continue

			}

			SGMS.Map[mapString] = x

		}

		return SGMS
	}
	return nil
}

//SGMapStringLineRead: Lee un elemento y devuelve si existe o no con un boleano.
//#bd/ramSync/SGMapStringLineRead
func (SGMS *SpaceSGMStr) SGMapStringLineRead(bufferBytes *[]byte) bool {

	//Preformateamos los datos para que mantenga la concordancia con el archivo.
	if SGMS.hooker != nil {

		SGMS.hookerPreFormatPointer(bufferBytes, SGMS.colName)

	}

	//Aplicamos los padding del archivo.
	SGMS.spacePaddingPointer(bufferBytes, SGMS.size)

	//Borramos los espacios a la derecha
	SGMS.spaceTrimPointer(bufferBytes)

	//Transformamos a string
	valueStr := string(*bufferBytes)

	//Añadimos un bloqueo de lectura
	SGMS.RLock()
	defer SGMS.RUnlock()

	//Si existe devolvemos true si no false.
	_, found := SGMS.Map[valueStr]
	if found {

		return true

	}

	return false
}

//SGMapStringLineUpd: Actualiza el archivo y el mapa superglobal, unicamente campos unicos
//si no retorna false la funcion.
//#bd/ramSync/SGMapStringLineUpd
func (SGMS *SpaceSGMStr) SGMapStringLineUpd(line int64, bufferBytes *[]byte) (*int64, bool) {

	//preformat para el mapa
	if SGMS.hooker != nil {

		SGMS.hookerPreFormatPointer(bufferBytes, SGMS.colName)

	}

	//Aplicamos los padding del archivo.
	SGMS.spacePaddingPointer(bufferBytes, SGMS.size)

	//Borramos los espacios a la derecha
	SGMS.spaceTrimPointer(bufferBytes)

	//Transformamos a string
	strBuffer := string(*bufferBytes)

	//Creamos el buffer de escritura

	//creamos los bloqueos
	SGMS.Lock()
	defer SGMS.Unlock()

	//Si la linea es -1  buscamos en el mapa, creamos una nueva entrada y una nueva linea en el archivo.
	if line == -1 {

		_, found := SGMS.Map[strBuffer]
		if !found {

			line := *SGMS.NewOneLine(SGMS.colName, bufferBytes)

			SGMS.Map[strBuffer] = line

			return &line, true

		}
		return nil , false
	}

	//Si la linea es > -1 actualizamos la linea y el mapa
	if line > -1 {

		//Primero leemos el fichero y borramos esa linea del mapa
		//BuffBytes := SGMS.BRspace( BuffBytes, false ,  line, line  , SGMS.colName)
		//BuffBytes.Rspace()

		BuffBytes := SGMS.GetOneLine(SGMS.colName, line)

		//Borramos los espacios a la derecha
		SGMS.spaceTrimPointer(BuffBytes.Buffer)

		delete(SGMS.Map, string(*BuffBytes.Buffer))

		//Despues escribimos la nueva linea en el archivo
		line := SGMS.SetOneLine(SGMS.colName, line, bufferBytes)
		//Añadimos esa linea al mapa tambien
		SGMS.Map[strBuffer] = *line

		return line, true
	}

	return nil , false
}

//SGMapStringLineDel: Borramos un elemento del archivo y del mapa
//#bd/ramSync/SGMapStringLineDel
func (SGMS *SpaceSGMStr) SGMapStringLineDel(line int64)(linePointer *int64) {

	//Creamos un buffer de lectura y leemos el numero de linea
	//BuffBytes := SGMS.BRspace( BuffBytes, false ,  line, line  , SGMS.colName)
	//BuffBytes.Rspace()
	BuffBytes := SGMS.GetOneLine(SGMS.colName, line)

	//Borramos los espacios a la derecha
	SGMS.spaceTrimPointer(BuffBytes.Buffer)

	//Borramos ese numero de linea del mapa
	delete(SGMS.Map, string(*BuffBytes.Buffer))

	//borramos la linea enviando un byte null de escritura.
	linePointer = SGMS.SetOneLine(SGMS.colName, line, &[]byte{})

	return linePointer
}
