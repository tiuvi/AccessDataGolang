package dac

import (
	"os"
	"strings"
)

func NewContentRead(extension string, dirName ...string) *PublicSpaceFile {

	DAC := GetGlobalDac()
	dir := DAC.globalDACFolder
	lastElement := len(dirName) - 1
	var lastElementName string

	for ind, dirNameStr := range dirName {


		if EDAC &&
			DAC.ELDACF(len(dirNameStr) == 0, "Estas enviando una cadena vacia en un array.") {
		}

		if lastElement == ind {

			lastElementName = dirNameStr 

		} else {

			dir = strings.Join([]string{dir, dirNameStr, "/"}, "")

		}

	}



	//Inicio directorio de mapas -> structuras
	var PSC *PublicSpaceCache
	PSC = GetCache(dir);
	if PSC == nil {

		PSC = InsertCache(dir)

	}


	dir = strings.Join([]string{dir, lastElementName,".",extension}, "")

	
	if PublicSF := PSC.GetFileCache(dir); PublicSF != nil {

		return PublicSF
	}
	

	info, err := os.Stat(dir)
	if err != nil {
		return nil
	}

	PSF := NewSf(openFile, bytes, extension, map[string]int64{extension: info.Size()}, nil, dirName...)

	PSC.InsertFileCache(PSF)

	return PSF
}

func NewContentWrite(extension string, lenContent int64, dirName ...string) *PublicSpaceFile {

	DAC := GetGlobalDac()
	dir := DAC.globalDACFolder
	lastElement := len(dirName) - 1
	var lastElementName string

	for ind, dirNameStr := range dirName {

	
		if EDAC &&
			DAC.ELDACF(len(dirNameStr) == 0, "Estas enviando una cadena vacia en un array.") {
		}

		if lastElement == ind {

			lastElementName = dirNameStr 
			
		} else {

			dir = strings.Join([]string{dir, dirNameStr, "/"}, "")

		}

	}

	//Inicio directorio de mapas -> structuras
	var PSC *PublicSpaceCache
	PSC = GetCache(dir);
	if PSC == nil {

		PSC = InsertCache(dir)

	}


	dir = strings.Join([]string{dir,lastElementName, ".",extension}, "")

	if PSF := PSC.GetFileCache(dir); PSF != nil {

		if size , _ := PSF.indexSizeFields[extension]; size[1] == lenContent{
		
			
			return PSF

		} else {

			PSF := NewSf(openFile, bytes, extension, map[string]int64{extension: lenContent}, nil, dirName...)
			PSC.UpdateFileCache(PSF)

		

			return PSF
		}
	}
	
	PSF:= NewSf(openFile, bytes, extension, map[string]int64{extension: lenContent}, nil, dirName...)

	PSC.InsertFileCache(PSF)
	
	return PSF
}





