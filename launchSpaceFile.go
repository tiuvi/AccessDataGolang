package bd

import "log"


func (obj *space ) OSpace( name string ,folder... string)*spaceFile  {


	if !obj.compilation {

		log.Fatalln("Este espacio no se ha copilado, copilar mejora la seguridad.", obj.dir)

	}


	if len(name) == 0 {

		log.Fatalln("Nombre de archivo vacio en: ", obj.dir, obj.extension)

	}


	if checkFileNativeType(obj.fileNativeType, disk ){

		return obj.ospaceDisk(name, folder)

	} 

	if checkFileNativeType(obj.fileNativeType, deferDisk ){

		return obj.ospaceDeferDisk(name, folder)

	} 

	if checkFileNativeType(obj.fileNativeType, permDisk ){

		return obj.ospacePermDisk(name, folder)

	}

	log.Fatalln("Es obligatorio definir el FileNativeType de la estructura. ", obj.dir,obj.extension)
	return nil
}
