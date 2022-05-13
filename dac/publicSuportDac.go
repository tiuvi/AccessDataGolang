package dac

import (
	"fmt"

	"os"
	"regexp"
	"strings"
)

/**********************************************************************************************/
/* Columnas */
/**********************************************************************************************/

//Verifica si el string, es una columna en ese espacio.
func (SF *space) IsColumnMap(column string) *bool {

	if SF.indexSizeColumns != nil {

		_, found := SF.indexSizeColumns[column]
		if found {

			exist := true
			return &exist
		}

		exist := false
		return &exist
	}

	return nil
}

//Verifica si el string, es una columna en ese espacio.
func (SF *space) IsColumn(column string) bool {

	if SF.indexSizeColumns != nil {

		_, found := SF.indexSizeColumns[column]
		if found {

			return true
		}

		return false
	}

	return false
}

//Verifica si el string, es una columna en ese espacio.
func (SF *space) IsNotColumn(column string) bool {

	if SF.indexSizeColumns != nil {

		_, found := SF.indexSizeColumns[column]
		if found {

			return false
		}

		return true
	}

	return true
}

//CAlcula el tamaño de una columna
func (SP *space) CalcSizeColumnBWspace(field string) (sizeTotal int64) {

	if EDAC &&
		SP.ECSD(SP.indexSizeColumns == nil, "Este espacio no tiene columnas") {
	}

	size, found := SP.indexSizeColumns[field]
	if EDAC &&
		SP.ECSD(!found, "Este espacio no es un columnas") {
	}
	if found {
		sizeTotal = size[1] - size[0]
	}

	return sizeTotal
}



/**********************************************************************************************/
/* Fields */
/**********************************************************************************************/

//Verifica si el string, es un campo en ese espacio.
func (SF *space) IsFieldMap(field string) *bool {

	if SF.indexSizeFields != nil {

		_, found := SF.indexSizeFields[field]
		if found {

			exist := true
			return &exist
		}

		exist := false
		return &exist
	}

	return nil
}

//Verifica si el string, es un campo en ese espacio.
func (SF *space) IsField(field string) bool {

	if SF.indexSizeFields != nil {

		_, found := SF.indexSizeFields[field]
		if found {

			return true
		}

		return false
	}

	return false
}

//Verifica si el string, es un campo en ese espacio.
func (SF *space) IsNotField(field string) bool {

	if SF.indexSizeFields != nil {

		_, found := SF.indexSizeFields[field]
		if found {

			return false
		}

		return true
	}

	return true
}

//CAlcula el tamaño de un campo
func (SP *space) CalcSizeField(field string) (sizeTotal int64) {

	if EDAC &&
		SP.ECSD(SP.indexSizeFields == nil, "Este espacio no tiene fields") {
	}

	size, found := SP.indexSizeFields[field]

	if EDAC &&
		SP.ECSD(!found, "Este espacio no es un fields") {
	}

	if found {
		sizeTotal = size[1] - size[0]
	}

	return sizeTotal
}



/**********************************************************************************************/
/* Extension */
/**********************************************************************************************/

//Verifica si el string, es un campo en ese espacio.
func IsExtensionContent(extName string) (value string, found bool) {

	if extName == "dacByte" || extName ==  "dacBit" {
		return 
	}

	value, found = extensionFile[extName]
	if found {

		return 
	}
	return 
}

func IsRagesExtension(extName string) bool {

	value, found := allowedRange[extName]
	if found {
		return value
	}
	return false
}


func GetExtension()(extension []string){

	for extName := range extensionFile {

		if extName == "dacByte" || extName ==  "dacBit" {
			continue
		}
		
		extension = append(extension, extName)
	}

	return
}

func (SP *space) GetExtension()(string){

	return SP.extension
}

/**********************************************************************************************/
/* Columnas y fields */
/**********************************************************************************************/

func (sP *space) IsColFil(data ...string) bool {

	for _, name := range data {

		if sP.IsField(name) {
			continue
		}

		if sP.IsColumn(name) {
			continue
		}
		return false
	}

	return true
}

func (sP *space) IsNotColFil(data ...string) bool {

	for _, name := range data {

		if sP.IsField(name) {
			continue
		}

		if sP.IsColumn(name) {
			continue
		}
		return true
	}

	return false
}

/**********************************************************************************************/
/* RAngos */
/**********************************************************************************************/

//CAlcula el numero de rangos de un field o nil si no existe.
func (SP *space) CalcRangeField(field string, RangeBytes int64) int64 {

	var sizeTotal int64

	if EDAC &&
		SP.ECSD(SP.indexSizeFields == nil, "Este espacio no tiene fields") ||
		SP.ECSD(RangeBytes < 0, "Rango inferior a 0.") {
	}

	size, found := SP.indexSizeFields[field]

	if EDAC &&
		SP.ECSD(!found, "Ese field no existe") {
	}

	if found {
		sizeTotal = size[1] - size[0]
	} else {
		return 0
	}

	if EDAC &&
		SP.ECSD(RangeBytes > sizeTotal, "El rango es mayor que el tamaño total del campo.") {
	}

	TotalRangue := sizeTotal / RangeBytes
	restoRangue := sizeTotal % RangeBytes

	if restoRangue != 0 {

		TotalRangue += 1
	}

	return TotalRangue
}

//Funcion de soporte para calcular los rangos
func (SP *space) CalcRangesBytes(lenBuffer int64, RangeBytes int64) int64 {

	if EDAC &&
	SP.ECSD(RangeBytes <= 0, "El rango de bytes no puede ser inferior o igual a cero") ||
	SP.ECSD(lenBuffer <= 0, "El buffer de bytes no puede ser inferior o igual a cero") ||
	SP.ECSD(lenBuffer < RangeBytes, "El tamaño del buffer no puede ser inferior al tamaño del rango."){}

	TotalRangue := lenBuffer / RangeBytes
	restoRangue := lenBuffer % RangeBytes

	if restoRangue != 0 {

		TotalRangue += 1
	}
	
	return TotalRangue
}


/**********************************************************************************************/
/* Operaciones con directorios */
/**********************************************************************************************/

func (SP *space) CheckDirSP()bool {

	fInfo, err := os.Stat(SP.dir)
	if err != nil && 
	err != os.ErrNotExist && 
	EDAC && 
	SP.ECSD( true,"Error al leer el directorio \n\r" +  fmt.Sprintln(err) ){}
	
	if err != nil && err == os.ErrNotExist {

		return false
	} 

	if fInfo.IsDir() {

		return true
	}

	return false
}









func cutExtensionToPath(url string)(patch string, extension string){

	position := strings.LastIndex(url, ".")
	if position == -1 {
		return
	}
	extension    = url[position +1:]
	patch        = url[:position]

	return
}

func sliceUrlToPath(levelU uint8,maxLevelU uint8,url string)(patch []string ){

	if maxLevelU == 0 || maxLevelU < levelU {
		return
	}

	if len(url) <= 1 {
		return
	}

	if url[:1] == "/" {

		url =  url[1:]
	} 

	level    := int(levelU)
	maxLevel := int(maxLevelU)

	patch   = strings.Split(url , "/")
	if len(patch) < level || len(patch) == 0{
		return patch[:0]
	}

	if len(patch) >= level {

		patch  = patch[level:]
	} 

	if len(patch) > maxLevel {

		patch  = patch[:maxLevel]
	}
	
	return
}

var SanitizeUrlContentRegexp = regexp.MustCompile(`[^a-zA-Z0-9]`)
func sanitizePath( patch []string)[]string {

	var x int
	for ind, value :=  range patch {

		ind = ind - x 

			if value := SanitizeUrlContentRegexp.ReplaceAllString(value , ""); len(value) > 0 {

					patch[ind] = value
					continue
			}

		x++	
	}

	return patch[:len(patch)-x]
}


func SanitizeUrl(levelU uint8,maxLevelU uint8,url string)(patch []string, extension string){


	url , extension = cutExtensionToPath(url)
	if url == "" || extension == ""{
		return
	}

	patch = sliceUrlToPath(levelU ,maxLevelU ,url )
	if len(patch) == 0 {
		return
	}
	
	patch   = sanitizePath( patch )

	return
}

var RegexpNewReactApp = regexp.MustCompile(`[^a-zA-Z0-9/.]`)
func sanitizePathRGP(patch []string)[]string {

	var x int
	for ind, value :=  range patch {

		ind = ind - x 

			if value := RegexpNewReactApp.ReplaceAllString(value , ""); len(value) > 0 {

				if value := strings.ReplaceAll(value, "..", "" ); value != "."  && len(value) > 0  {
					
					patch[ind] = value
					continue
				}
			}

		x++	
	}

	return patch[:len(patch)-x]
}


func SanitizeUrlRGP(levelU uint8,maxLevelU uint8, url string)(patch []string, extension string){


	if maxLevelU == 0 || maxLevelU < levelU  || len(url) <= 1 {
		return
	}

	if url[:1] == "/" {

		url =  url[1:]
	} 

	url , extension = cutExtensionToPath(url)
	if url == "" || extension == ""{
		return
	}

	patch = sliceUrlToPath(levelU ,maxLevelU ,url )
	if len(patch) == 0 {
		return
	}

	patch   = sanitizePathRGP(patch )

	return
}
