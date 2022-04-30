package bd

import "log"

func (LDAC *lDAC) NewSpace() *space {

	if len(LDAC.globalDACFolder) == 0 {
		//faTAL
	}

	lSpace := new(space)
	lSpace.lDAC = LDAC

	return lSpace
}

func (SP *space) NewTimeFileDisk() {

	if SP.fileNativeType > 0 {

		log.Fatal("Propiedad filenativetype ya establecida e inmutable.")

	}

	SP.fileNativeType = disk
}

func (SP *space) NewTimeFileDeferDisk() {

	if SP.fileNativeType > 0 {

		log.Fatal("Propiedad filenativetype ya establecida e inmutable.")

	}

	SP.fileNativeType = deferDisk
}

func (SP *space) NewTimeFilePermDisk() {

	if SP.fileNativeType > 0 {

		log.Fatal("Propiedad filenativetype ya establecida e inmutable.")

	}

	SP.fileNativeType = permDisk
}

func (SP *space) SetDir(dir string) {

	if len(SP.dir) > 0 {

		log.Fatal("Propiedad dir ya establecida e inmutable.")

	}

	if len(dir) == 0 {

		log.Fatal("Estas enviando una cadena vacia.", "\r\n",
			"¿Correcto? Revisa - bd.SetGlobalDACFolder")

	}
	dir = regexPathGlobal(dir)

	//Añadimos / barra al final si no la lleva
	if dir[len(dir)-1:] != "/" {

		dir += "/"

	}

	SP.dir = SP.globalDACFolder + dir

}

func (SP *space) NewDacByte() {

	if len(SP.extension) > 0 {

		log.Fatal("Propiedad extension ya establecida e inmutable.")

	}
	SP.extension = dacByte

}

func (SP *space) NewDacBit() {

	if len(SP.extension) > 0 {

		log.Fatal("Propiedad extension ya establecida e inmutable.")

	}
	SP.extension = dacBit

}

func (SP *space) NewField(name string, size int64) {

	if SP.indexSizeFields != nil {

		log.Fatal("Propiedad indexSizeFields ya establecida e inmutable.")

	}

	if size < 1 {

		log.Fatal("Size no puede ser menor o igual a 0.")

	}

	if SP.indexSizeFieldsArray == nil {

		SP.indexSizeFieldsArray = make([]spaceLen, 0)

	}

	SP.indexSizeFieldsArray = append(SP.indexSizeFieldsArray, spaceLen{name, size})

}

func (SP *space) GetNameField() (fields []string) {

	for _, spaceLen := range SP.indexSizeFieldsArray {

		fields = append(fields, spaceLen.name)

	}
	return fields
}

func (SP *space) FieldSizeTotal() (lenField int64) {

	for _, spaceLen := range SP.indexSizeFieldsArray {

		lenField += spaceLen.len

	}
	return lenField

}

func (SP *space) NewColumnByte(name string, size int64) {

	if SP.indexSizeColumns != nil {

		log.Fatal("Propiedad indexSizeColumns ya establecida e inmutable.")

	}

	if size < 1 {

		log.Fatal("Size no puede ser menor o igual a 0.")

	}

	if SP.extension != dacByte {

		log.Fatal("Extension incompatible o indefinida, Activa: space.NewDacByte()")

	}

	if SP.indexSizeColumnsArray == nil {

		SP.indexSizeColumnsArray = make([]spaceLen, 0)

	}

	SP.indexSizeColumnsArray = append(SP.indexSizeColumnsArray, spaceLen{name, size})

}

func (SP *space) NewColumnBit(name string) {

	if SP.indexSizeColumns != nil {

		log.Fatal("Propiedad indexSizeColumns ya establecida e inmutable.")

	}

	if SP.extension != dacBit {

		log.Fatal("Extension incompatible o indefinida, Activa: space.NewDacBit()")

	}

	if SP.indexSizeColumnsArray == nil {

		SP.indexSizeColumnsArray = make([]spaceLen, 0)

	}

	SP.indexSizeColumnsArray = append(SP.indexSizeColumnsArray, spaceLen{name, 1})

}

func (SP *space) GetNameColumn() (columns []string) {

	for _, spaceLen := range SP.indexSizeColumnsArray {

		columns = append(columns, spaceLen.name)

	}
	return columns
}

func (SP *space) ColumnSizeTotal() (lenColumn int64) {

	for _, spaceLen := range SP.indexSizeColumnsArray {

		lenColumn += spaceLen.len

	}
	return lenColumn

}

func (SP *space) PreformatDefault(function func(*[]byte)) {

	if SP.hooker == nil {

		SP.hooker = make(map[string]func(*[]byte))

	}

	if _, found := SP.hooker[preformat]; found {

		log.Fatal("Ya declaraste el preformat predeterminado.")

	}

	SP.hooker[preformat] = function

}

func (SP *space) PreformatGlobal(name string, function func(*[]byte)) {

	if SP.hooker == nil {

		SP.hooker = make(map[string]func(*[]byte))

	}

	if _, found := SP.hooker[preformat+name]; found {

		log.Fatal("Ya declaraste esta columna o field preformat.")

	}

	SP.hooker[preformat+name] = function

}

func (SP *space) PostformatDefault(function func(*[]byte)) {

	if SP.hooker == nil {

		SP.hooker = make(map[string]func(*[]byte))

	}

	if _, found := SP.hooker[postformat]; found {

		log.Fatal("Ya declaraste el preformat predeterminado.")

	}

	SP.hooker[postformat] = function

}

func (SP *space) PostformatGlobal(name string, function func(*[]byte)) {

	if SP.hooker == nil {

		SP.hooker = make(map[string]func(*[]byte))

	}

	if _, found := SP.hooker[postformat+name]; found {

		log.Fatal("Ya declaraste esta columna o field preformat.")

	}

	SP.hooker[postformat+name] = function

}

func (obj *space ) OSpaceInit()bool  {

	return obj.ospaceCompilationFile()

}
