package dac

import(
	"strings"
	"os"
)

func (SP *space ) OSpace( name string ,folder... string)*spaceFile  {


	if EDAC &&
	SP.ECSD( !SP.compilation ,"Este espacio no se ha copilado") ||
	SP.ECSD( len(name) == 0 ,"No se puede enviar un nombre de archivo vacio."){}


	if checkFileNativeType(SP.fileNativeType, disk ){

		return SP.ospaceDisk(name, folder)

	} 

	if checkFileNativeType(SP.fileNativeType, deferDisk ){

		return SP.ospaceDeferDisk(name, folder)

	} 

	if checkFileNativeType(SP.fileNativeType, permDisk ){

		return SP.ospacePermDisk(name, folder)

	}

	if checkFileNativeType(SP.fileNativeType, openFile ){

		return SP.ospaceOpenFile(name, folder)

	}
	

	if EDAC &&
	SP.ECSD( true ,"Error grave sin coincidencias FileNativeType."){}

	return nil
}

func (SF *spaceFile) SetPublicSpaceFile()*PublicSpaceFile{

	if EDAC &&
	SF.ECSFD( !SF.compilation ,"Este espacio no se ha copilado"){}


	return &PublicSpaceFile{
		spaceFile: SF,
	}

}



func (SP *spaceFile) CheckDirSF()bool {

	urlDir    := strings.SplitAfter(SP.url, "/")
	urlDirStr := strings.Join(urlDir[:len(urlDir)-1], "")

	fInfo, err := os.Stat(urlDirStr)
	if err != nil && 
	err == os.ErrExist && 
	EDAC && 
	SP.ECSD( true,"Error al leer el directorio \n\r" +  err.Error() ){}
	
	if fInfo.IsDir() {
		return true
	}

	return false
}

func (SP *spaceFile) CheckFileSF()bool {

	fInfo, err := os.Stat(SP.url)
	if err != nil && 
	err == os.ErrExist && 
	EDAC && 
	SP.ECSD( true,"Error al leer el directorio \n\r" +  err.Error() ){}
	
	if fInfo.IsDir() {
		return false
	}

	return true
}


func (SP *spaceFile) GetUrl()string {
 
	return SP.url
}
