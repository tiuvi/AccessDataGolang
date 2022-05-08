package dac

import ("strings")

func (LDAC *lDAC) NewSpace() *space {

	if EDAC && 
	LDAC.ELDAC(len(LDAC.globalDACFolder) == 0,"La ruta DAC no esta establecida."){}
	
	lSpace := new(space)
	lSpace.lDAC = LDAC

	return lSpace
}

func (SP *space) OnErrorFile() {

	if EDAC && 
	SP.ECSD(SP.isErrorFile == true, "Este espacio ya fue definido como un archivo de error."){}
	
	SP.isErrorFile = true

}

func (SP *space) NewTimeFileDisk() {

	if EDAC && 
	SP.ECSD(SP.fileNativeType > 0 , "Propiedad filenativetype ya establecida e inmutable."){}
	
	SP.fileNativeType = disk
}

func (SP *space) NewTimeFileDeferDisk() {

	if EDAC && 
	SP.ECSD(SP.fileNativeType > 0 , "Propiedad filenativetype ya establecida e inmutable."){}
	

	SP.fileNativeType = deferDisk
}

func (SP *space) NewTimeFilePermDisk() {

	if EDAC && 
	SP.ECSD(SP.fileNativeType > 0 , "Propiedad filenativetype ya establecida e inmutable."){}
	

	SP.fileNativeType = permDisk
}

func (SP *space) NewTimeFileOpenFile() {

	if EDAC && 
	SP.ECSD(SP.fileNativeType > 0 , "Propiedad filenativetype ya establecida e inmutable."){}
	

	SP.fileNativeType = openFile
}



func (SP *space) NewDacByte() {

	if EDAC && 
	SP.ECSD( len(SP.extension) > 0, "Propiedad extension ya establecida e inmutable."){}

	SP.extension  = dacByte
	SP.fileCoding = bytes
}


func (SP *space) NewDacBit() {

	if EDAC && 
	SP.ECSD( len(SP.extension) > 0, "Propiedad extension ya establecida e inmutable."){}

	SP.extension = dacBit
	SP.fileCoding = bit
}

//Unicamente para otras extensiones que no sean dacbyte o dacbit
func (SP *space) SetFileCodgingBit() {

	if EDAC && 
	SP.ECSD( SP.fileCoding > 0, "Propiedad filecoding ya establecida e inmutable."){}

	SP.fileCoding = bit
}

//Unicamente para otras extensiones que no sean dacbyte o dacbit
func (SP *space) SetFileCodgingByte() {

	if EDAC && 
	SP.ECSD( SP.fileCoding > 0, "Propiedad filecoding ya establecida e inmutable."){}

	SP.fileCoding = bytes
}

func (SP *space) SetExtensionHtml() {

	if EDAC && 
	SP.ECSD( len(SP.extension) > 0, "Propiedad extension ya establecida e inmutable."){}

	SP.extension = html

}

func (SP *space) SetExtension(extension string) {

	if EDAC && 
	SP.ECSD( len(SP.extension) > 0, "Propiedad extension ya establecida e inmutable."){}

	SP.extension = extension

}

func (SP *space) SetDir(dirName string) {

	if EDAC && 
	SP.ECSD( len(SP.dir) > 0 , "Propiedad dir ya establecida e inmutable.") ||
	SP.ECSD( len(dirName) == 0 , "Estas enviando una cadena vacia."){}
	

	SP.dir =  strings.Join([]string{SP.globalDACFolder , dirName , "/" } , "")

}

func (SP *space) SetSubDir(dir ...string) {

	if EDAC && 
	SP.ECSD( len(SP.dir) > 0 , "Propiedad dir ya establecida e inmutable.") ||
	SP.ECSD( len(dir) == 0 , "Estas enviando una cadena vacia."){}
	
	SP.dir = SP.globalDACFolder
	
	for _ , dirNameStr := range dir {
			
		if EDAC && 
		SP.ECSD( len(dirNameStr) == 0 , "Estas enviando una cadena vacia en un array."){}

		SP.dir = strings.Join([]string{SP.dir ,dirNameStr , "/"}, "" )
	
	}


}




func (SP *space) NewField(name string, size int64) {

	if EDAC && 
	SP.ECSD( SP.indexSizeFields != nil , "Propiedad indexSizeFields ya establecida e inmutable.") ||
	SP.ECSD( size < 1, "Size no puede ser menor o igual a 0."){}


	if SP.indexSizeFieldsArray == nil {

		SP.indexSizeFieldsArray = make([]spaceLen, 0)

	}
	
	SP.indexSizeFieldsArray = append(SP.indexSizeFieldsArray, spaceLen{name, size})

}

func (SP *space) GetNameField() (fields []string) {

	if EDAC && 
	SP.ECSD( SP.indexSizeFields == nil , "No se han iniciado campos en este espacio."){}

	for _, spaceLen := range SP.indexSizeFieldsArray {

		fields = append(fields, spaceLen.name)

	}
	return fields
}

func (SP *space) FieldSizeTotal() (lenField int64) {

	if EDAC && 
	SP.ECSD( SP.indexSizeFields == nil , "No se han iniciado campos en este espacio."){}

	for _, spaceLen := range SP.indexSizeFieldsArray {

		lenField += spaceLen.len

	}
	return lenField

}

func (SP *space) NewColumnByte(name string, size int64) {


	if EDAC && 
	SP.ECSD( SP.indexSizeColumns != nil , "Propiedad indexSizeColumns ya establecida e inmutable.") ||
	SP.ECSD( size < 1 , "Size no puede ser menor o igual a 0.") ||
	SP.ECSD( SP.extension != dacByte , "Extension incompatible o indefinida, Activa: space.NewDacByte()"){}


	if SP.indexSizeColumnsArray == nil {

		SP.indexSizeColumnsArray = make([]spaceLen, 0)

	}

	SP.indexSizeColumnsArray = append(SP.indexSizeColumnsArray, spaceLen{name, size})

}

func (SP *space) NewColumnBit(name string) {

	if EDAC && 
	SP.ECSD( SP.indexSizeColumns != nil , "Propiedad indexSizeColumns ya establecida e inmutable.") ||
	SP.ECSD( SP.extension != dacBit , "Extension incompatible o indefinida, Activa: space.NewDacByte()"){}


	if SP.indexSizeColumnsArray == nil {

		SP.indexSizeColumnsArray = make([]spaceLen, 0)

	}

	SP.indexSizeColumnsArray = append(SP.indexSizeColumnsArray, spaceLen{name, 1})

}

func (SP *space) GetNameColumn() (columns []string) {

	if EDAC && 
	SP.ECSD( SP.indexSizeColumns == nil , "No se han iniciado columnas en este espacio"){}

	for _, spaceLen := range SP.indexSizeColumnsArray {

		columns = append(columns, spaceLen.name)

	}
	return columns
}

func (SP *space) ColumnSizeTotal() (lenColumn int64) {

	if EDAC && 
	SP.ECSD( SP.indexSizeColumns == nil , "No se han iniciado columnas en este espacio"){}

	for _, spaceLen := range SP.indexSizeColumnsArray {

		lenColumn += spaceLen.len

	}
	return lenColumn

}

func (SP *space) PreformatDefault(function func(*[]byte)) {

	if SP.hooker == nil {

		SP.hooker = make(map[string]func(*[]byte))

	}

	if EDAC {
		_, found := SP.hooker[preformat]; 
		SP.ECSD( found , "Ya declaraste el preformat predeterminado.")
	}

	SP.hooker[preformat] = function

}

func (SP *space) PreformatGlobal(name string, function func(*[]byte)) {

	if SP.hooker == nil {

		SP.hooker = make(map[string]func(*[]byte))

	}

	if EDAC {
		_, found := SP.hooker[preformat+name]; 
		SP.ECSD( found , "Ya declaraste el preformat global.")
	}

	SP.hooker[preformat+name] = function

}

func (SP *space) PostformatDefault(function func(*[]byte)) {

	if SP.hooker == nil {

		SP.hooker = make(map[string]func(*[]byte))

	}

	if EDAC {
		_, found := SP.hooker[postformat]; 
		SP.ECSD( found , "Ya declaraste el postformat predeterminado.")
	}

	SP.hooker[postformat] = function

}

func (SP *space) PostformatGlobal(name string, function func(*[]byte)) {

	if SP.hooker == nil {

		SP.hooker = make(map[string]func(*[]byte))

	}

	if EDAC {
		_, found := SP.hooker[postformat+name]; 
		SP.ECSD( found , "Ya declaraste este postformat global.")
	}

	SP.hooker[postformat+name] = function

}

func (SP *space ) OSpaceInit()bool  {

	if EDAC && 
	SP.ECSD( SP.fileNativeType == 0 ,"fileNativeType no definido") ||
	SP.ECSD( len(SP.dir) == 0 , "Directorio no definido") ||
	SP.ECSD( len(SP.extension) == 0 ,"extension no definido") ||
	SP.ECSD( len(SP.indexSizeFieldsArray) == 0 && len(SP.indexSizeColumnsArray) == 0 ,
	"Se ha iniciado un espacio sin fields y sin columnas.") ||
	SP.ECSD( SP.compilation ,"Este espacio ya se ha copilado"){}
	 

	return SP.ospaceCompilationFile()

}

func (SP *space ) OSpaceGlobal(name string)bool  {

	if EDAC {

		SP.ECSD( !SP.compilation ,"Este espacio no se ha copilado")

		_ , found := Space[name]; 
		SP.ECSD( found ,"Este nombre ya existe en el mapa global.")
	}
	
	Space[name] = SP
	
	return true
}

func (SP *space) SetPublicSpace()*PublicSpace{

	if EDAC &&
	SP.ECSD( !SP.compilation ,"Este espacio no se ha copilado"){}


	return &PublicSpace{
		space: SP,
	}

}