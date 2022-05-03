package bd

import "log"


func (SP *space ) OSpace( name string ,folder... string)*spaceFile  {


	if EDAC &&
	SP.ECSD( !SP.compilation ,"Este espacio no se ha copilado"){}



	if len(name) == 0 {

		log.Fatalln("Nombre de archivo vacio en: ", SP.dir, SP.extension)

	}


	if checkFileNativeType(SP.fileNativeType, disk ){

		return SP.ospaceDisk(name, folder)

	} 

	if checkFileNativeType(SP.fileNativeType, deferDisk ){

		return SP.ospaceDeferDisk(name, folder)

	} 

	if checkFileNativeType(SP.fileNativeType, permDisk ){

		return SP.ospacePermDisk(name, folder)

	}

	log.Fatalln("Es obligatorio definir el FileNativeType de la estructura. ", SP.dir,SP.extension)
	return nil
}
