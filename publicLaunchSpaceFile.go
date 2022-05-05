package dac



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