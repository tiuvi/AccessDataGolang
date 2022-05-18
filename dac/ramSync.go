package dac

import (
//	"log"
)

type SpaceRamSync struct {
	*spaceFile
	colName string
	size    [2]int64
	Map     map[string]int64
}


func (sF *spaceFile) InitSync(colName string) *SpaceRamSync {

	size, found := sF.indexSizeColumns[colName]

	if EDAC &&
		sF.ECSD(!found, "Columna no encontrada.") {
	}

	if found {

		//Creamos un puntero a la estructura.
		SGMS := &SpaceRamSync{
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
			SpaceTrimPointer(&mapColumn.BufferMap[colName][x])

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
func (SGMS *SpaceRamSync) GetLine(bufferBytes *[]byte)*int64 {

	//Preformateamos los datos para que mantenga la concordancia con el archivo.
	if SGMS.hooker != nil {

		SGMS.hookerPreFormatPointer(bufferBytes, SGMS.colName)

	}

	//Aplicamos los padding del archivo.
	SpacePaddingPointer(bufferBytes, SGMS.size)

	//Borramos los espacios a la derecha
	SpaceTrimPointer(bufferBytes)

	//Transformamos a string
	valueStr := string(*bufferBytes)

	//Añadimos un bloqueo de lectura
	SGMS.RLock()
	defer SGMS.RUnlock()

	//Si existe devolvemos true si no false.
	line , found := SGMS.Map[valueStr]
	if found {

		return &line

	}

	return nil
}

func (SGMS *SpaceRamSync) GetLineString(str string)*int64 {

	RBuffer := []byte(str)
	return SGMS.GetLine(&RBuffer)
}



//SGMapStringLineUpd: Actualiza el archivo y el mapa superglobal, unicamente campos unicos
//si no retorna false la funcion.
//#bd/ramSync/SGMapStringLineUpd
func (SGMS *SpaceRamSync) SetLine(line int64, bufferBytes *[]byte)*int64 {

	if EDAC && 
	SGMS.ECSFD( len(*bufferBytes) == 0 , "No se puede enviar un buffer de bytes vacio.") ||
	SGMS.ECSFD( line < 0, "Line no puede ser inferior a 0."){}

	//preformat para el mapa
	if SGMS.hooker != nil {

		SGMS.hookerPreFormatPointer(bufferBytes, SGMS.colName)

	}

	//Aplicamos los padding del archivo.
	SpacePaddingPointer(bufferBytes, SGMS.size)

	//Borramos los espacios a la derecha
	SpaceTrimPointer(bufferBytes)

	//Transformamos a string
	strBuffer := string(*bufferBytes)

	//Creamos el buffer de escritura

	//creamos los bloqueos
	SGMS.Lock()
	defer SGMS.Unlock()


	//Primero leemos el fichero y borramos esa linea del mapa
	//BuffBytes := SGMS.BRspace( BuffBytes, false ,  line, line  , SGMS.colName)
	//BuffBytes.Rspace()

	BuffBytes := SGMS.GetOneLine(SGMS.colName, line)

	//Borramos los espacios a la derecha
	SpaceTrimPointer(BuffBytes.Buffer)

	delete(SGMS.Map, string(*BuffBytes.Buffer))

	//Despues escribimos la nueva linea en el archivo
	lineNew := SGMS.SetOneLine(SGMS.colName, line, bufferBytes)

	//Añadimos esa linea al mapa tambien
	SGMS.Map[strBuffer] = *lineNew

	return lineNew

}


func (SGMS *SpaceRamSync) NewLine(bufferBytes *[]byte)*int64 {

	//preformat para el mapa
	if SGMS.hooker != nil {

		SGMS.hookerPreFormatPointer(bufferBytes, SGMS.colName)

	}

	//Aplicamos los padding del archivo.
	SpacePaddingPointer(bufferBytes, SGMS.size)

	//Borramos los espacios a la derecha
	SpaceTrimPointer(bufferBytes)

	//Transformamos a string
	strBuffer := string(*bufferBytes)

	//Creamos el buffer de escritura

	//creamos los bloqueos
	SGMS.Lock()
	defer SGMS.Unlock()

	//Si la linea es -1  buscamos en el mapa, creamos una nueva entrada y una nueva linea en el archivo.
	_, found := SGMS.Map[strBuffer]
	if !found {

		line := *SGMS.NewOneLine(SGMS.colName, bufferBytes)

		SGMS.Map[strBuffer] = line

		return &line

	}
	return nil
}

func (SGMS *SpaceRamSync) NewLineString(str string)*int64 {

	WBuffer := []byte(str)
	return SGMS.NewLine( &WBuffer)
}


//SGMapStringLineDel: Borramos un elemento del archivo y del mapa
//#bd/ramSync/SGMapStringLineDel
func (SGMS *SpaceRamSync) DeleteLine(line int64)(linePointer *int64) {

	//Creamos un buffer de lectura y leemos el numero de linea
	//BuffBytes := SGMS.BRspace( BuffBytes, false ,  line, line  , SGMS.colName)
	//BuffBytes.Rspace()
	BuffBytes := SGMS.GetOneLine(SGMS.colName, line)

	//Borramos los espacios a la derecha
	SpaceTrimPointer(BuffBytes.Buffer)

	//Borramos ese numero de linea del mapa
	delete(SGMS.Map, string(*BuffBytes.Buffer))

	//borramos la linea enviando un byte null de escritura.
	linePointer = SGMS.SetOneLine(SGMS.colName, line, &[]byte{})

	return linePointer
}
