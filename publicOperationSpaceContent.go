package dac

import (
	"os"
	"log"
	"strings"
)

func NewContentReadDiskSf(extension string, dirName ...string) *PublicSpaceFile {

	DAC := GetGlobalDac()
	dir := DAC.globalDACFolder
	lastElement := len(dirName) - 1

	for ind, dirNameStr := range dirName {

		regDirName := regexPathGlobalNoSlash(dirNameStr)
		if EDAC &&
			DAC.ELDACF(len(regDirName) == 0, "Estas enviando una cadena vacia en un array.") {
		}

		if lastElement == ind {

			dir = strings.Join([]string{dir, regDirName, ""}, "")

		} else {

			dir = strings.Join([]string{dir, regDirName, "/"}, "")

		}

	}

	dir = strings.Join([]string{dir, ".html"}, "")

	if spaceFile := getFileDisk(dir); spaceFile != nil {

		return spaceFile.SetPublicSpaceFile()
	}

	info, err := os.Stat(dir)
	if err != nil {
		return nil
	}

	return NewSf(disk, bytes, html, map[string]int64{extension: info.Size()}, nil, dirName...)
}

func NewContentWriteDiskSf(extension string, lenContent int64, dirName ...string) *PublicSpaceFile {

	DAC := GetGlobalDac()
	dir := DAC.globalDACFolder
	lastElement := len(dirName) - 1

	for ind, dirNameStr := range dirName {

		regDirName := regexPathGlobalNoSlash(dirNameStr)
		if EDAC &&
			DAC.ELDACF(len(regDirName) == 0, "Estas enviando una cadena vacia en un array.") {
		}

		if lastElement == ind {

			dir = strings.Join([]string{dir, regDirName, ""}, "")

		} else {

			dir = strings.Join([]string{dir, regDirName, "/"}, "")

		}

	}

	dir = strings.Join([]string{dir, ".html"}, "")

	deleteFileDisk(dir)

	return NewSf(disk, bytes, html, map[string]int64{extension: lenContent}, nil, dirName...)

}











/*
Funciones automaticas

*/
//ESta funcion genera todas las funciones de contenido.
func automaticFunction(){

	for name  :=range extensionFile{

		if name == "dacBit" || name == "DacByte"{
			continue
		}

		nameMayus :=  strings.ToUpper(name[:1]) + name[1:]

log.Println(`
func New`+nameMayus+`ReadDiskSf(dirName ...string) *PublicSpaceFile { 

	return NewContentReadDiskSf( `+name+`  , dirName...)

}


func New`+nameMayus+`WriteDiskSf(lenContent int64,dirName ...string) *PublicSpaceFile {

	return NewContentWriteDiskSf( `+name+`  , lenContent ,dirName...)

}
`)

	}
	
}



func NewCssReadDiskSf(dirName ...string) *PublicSpaceFile {

	return NewContentReadDiskSf(css, dirName...)

}

func NewCssWriteDiskSf(lenContent int64, dirName ...string) *PublicSpaceFile {

	return NewContentWriteDiskSf(css, lenContent, dirName...)

}

func NewGlbReadDiskSf(dirName ...string) *PublicSpaceFile {

	return NewContentReadDiskSf(glb, dirName...)

}

func NewGlbWriteDiskSf(lenContent int64, dirName ...string) *PublicSpaceFile {

	return NewContentWriteDiskSf(glb, lenContent, dirName...)

}

func NewJsonReadDiskSf(dirName ...string) *PublicSpaceFile {

	return NewContentReadDiskSf(json, dirName...)

}

func NewJsonWriteDiskSf(lenContent int64, dirName ...string) *PublicSpaceFile {

	return NewContentWriteDiskSf(json, lenContent, dirName...)

}

func NewJsReadDiskSf(dirName ...string) *PublicSpaceFile {

	return NewContentReadDiskSf(js, dirName...)

}

func NewJsWriteDiskSf(lenContent int64, dirName ...string) *PublicSpaceFile {

	return NewContentWriteDiskSf(js, lenContent, dirName...)

}

func NewSvgReadDiskSf(dirName ...string) *PublicSpaceFile {

	return NewContentReadDiskSf(svg, dirName...)

}

func NewSvgWriteDiskSf(lenContent int64, dirName ...string) *PublicSpaceFile {

	return NewContentWriteDiskSf(svg, lenContent, dirName...)

}


func NewHtmlReadDiskSf(dirName ...string) *PublicSpaceFile {

	return NewContentReadDiskSf(html, dirName...)

}

func NewHtmlWriteDiskSf(lenContent int64, dirName ...string) *PublicSpaceFile {

	return NewContentWriteDiskSf(html, lenContent, dirName...)

}

func NewJpgReadDiskSf(dirName ...string) *PublicSpaceFile {

	return NewContentReadDiskSf(jpg, dirName...)

}

func NewJpgWriteDiskSf(lenContent int64, dirName ...string) *PublicSpaceFile {

	return NewContentWriteDiskSf(jpg, lenContent, dirName...)

}

func NewMp4ReadDiskSf(dirName ...string) *PublicSpaceFile {

	return NewContentReadDiskSf(mp4, dirName...)

}

func NewMp4WriteDiskSf(lenContent int64, dirName ...string) *PublicSpaceFile {

	return NewContentWriteDiskSf(mp4, lenContent, dirName...)

}

func NewTxtReadDiskSf(dirName ...string) *PublicSpaceFile {

	return NewContentReadDiskSf(txt, dirName...)

}

func NewTxtWriteDiskSf(lenContent int64, dirName ...string) *PublicSpaceFile {

	return NewContentWriteDiskSf(txt, lenContent, dirName...)

}

func NewGifReadDiskSf(dirName ...string) *PublicSpaceFile {

	return NewContentReadDiskSf(gif, dirName...)

}

func NewGifWriteDiskSf(lenContent int64, dirName ...string) *PublicSpaceFile {

	return NewContentWriteDiskSf(gif, lenContent, dirName...)

}

func NewPngReadDiskSf(dirName ...string) *PublicSpaceFile {

	return NewContentReadDiskSf(png, dirName...)

}

func NewPngWriteDiskSf(lenContent int64, dirName ...string) *PublicSpaceFile {

	return NewContentWriteDiskSf(png, lenContent, dirName...)

}

func NewMp3ReadDiskSf(dirName ...string) *PublicSpaceFile {

	return NewContentReadDiskSf(mp3, dirName...)

}

func NewMp3WriteDiskSf(lenContent int64, dirName ...string) *PublicSpaceFile {

	return NewContentWriteDiskSf(mp3, lenContent, dirName...)

}

func NewPdfReadDiskSf(dirName ...string) *PublicSpaceFile {

	return NewContentReadDiskSf(pdf, dirName...)

}

func NewPdfWriteDiskSf(lenContent int64, dirName ...string) *PublicSpaceFile {

	return NewContentWriteDiskSf(pdf, lenContent, dirName...)

}
