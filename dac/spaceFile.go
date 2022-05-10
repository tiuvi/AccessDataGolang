package dac

import (
	"fmt"

	"os"
	"strings"
	"sync/atomic"
)

func (obj *space) newSpaceFile(folderString string, name string) *spaceFile {

	url := strings.Join([]string{obj.dir, folderString, name, ".", obj.extension}, "")

	var err error
	//Creamos una nueva referencia a spaceFile
	spacef := new(spaceFile)
	
	//Abrimos el archivo
	spacef.file, err = os.OpenFile(url, os.O_RDWR | os.O_CREATE, os.ModePerm)
	if err != nil {

		if os.IsNotExist(err) {

			err = os.MkdirAll(obj.dir+folderString, os.ModeDir | os.ModePerm)
			if err != nil && EDAC &&
				obj.ECSD(len(name) == 0, "Error al crear la carpeta para ese archivo."+fmt.Sprintln(err)) {
			}

			spacef.file, err = os.OpenFile(url, os.O_RDWR | os.O_CREATE , os.ModePerm)
			if err != nil {

				if err != nil && EDAC &&
					obj.ECSD(true, "Error al crear el archivo despues de crear la carpeta para ese archivo. \n\r"+fmt.Sprintln(err)) {
				}

			}

		} else {

			if EDAC &&
				obj.ECSD(true, "Error grave al abrir el archivo. \n\r"+fmt.Sprintln(err)) {
			}

		}
	}

	spacef.space = obj
	//Url pasada como valor dir + name + extension -> name dinamico
	spacef.url = url
	//Iniciamos un puntero a SizeFileLine manejado atomicamente por
	//un contador atomico
	spacef.sizeFileLine = new(int64)
	//Iniciamos el contador atomico
	atomic.StoreInt64(spacef.sizeFileLine, spacef.ospaceAtomicUpdateSizeFileLine())
	//Guardamos nuestro puntero de estructura space en el mapa global DiskFile

	return spacef
}

//Funcion que obtiene el numero de lineas atomico de un archivo
func (obj *spaceFile) ospaceAtomicUpdateSizeFileLine() (line int64) {

	//Leemos el archivo y movemos el puntero al final para contar su tamaño.
	size, err := obj.file.Seek(0, 2)
	if err != nil && EDAC &&
		obj.ECSD(true, "Error al obtener el numero de lineas del archivo, leyendo el archivo \n\r" + fmt.Sprintln(err)) {
	}

	if obj.sizeLine == 0 {

		return -1

	}
	
	//Si el tamaño es mayor que 0 restomos los fields
	if size > 0 {

		size -= obj.lenFields

	}

	//Si el resto entre el numero de lineas, es cero, obtenemos el numero de lineas.
	if size % obj.sizeLine == 0 {

		line = (size / obj.sizeLine)

	}

	//Si el resto entre el numero de lineas, no es cero, obtenemos el numero de lineas
	//sumando 1 debido a que es una linea que se ha escrito a la mitad.
	if size % obj.sizeLine != 0 {

		line = (size / obj.sizeLine) + 1

	}

	//Si es un archivo de bit y las lineas es mayor que cero lo multiplicamos * 8
	//por temas de hardware es totalmente imposible obtener el numero de linea exacto de bits
	//Debido a que el lector hdd lee y escribe 1 byte que son  bits.
	if checkFileCoding(obj.fileCoding, bit) {

		if line > 0 {

			line *= 8

		}

	}

	//El numero de lineas menos 1 debido a que la lineas empiezan en el numero 0
	return line - 1
}
