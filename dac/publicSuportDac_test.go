package dac

import (
	"regexp"
	"strings"
	"testing"
)


func AllLetterNoSlash()(cadena string) {

	for x := 0 ;x < 128; x++ {
		if x == 12  || x == 47 {
			continue
		}
		cadena += string(byte(x))
	}
	return
}

func AllLetterSlash(n int) string {

	cadena := AllLetterNoSlash()
	cadena += "/"

	var cadenaArr []string
	for x := 0; x < n; x++ {

		cadenaArr = append(cadenaArr, cadena)
	}

	strFinal := strings.Join(cadenaArr, "")

	return strFinal

}

func TestSanitizeUrlRGP(t *testing.T) {

	RegexpNewReactApp := regexp.MustCompile(`[^a-zA-Z0-9/.]`)

	cadena    := "/ruta1/ruta2/ruta3/ruta4/ruta5/ruta6/ruta7/"
	extension := "html"
	for y := 0 ; y <= 7; y++ {
		patch , extName  := SanitizeUrlRGP(RegexpNewReactApp, uint8(y), 7 , cadena + "." + extension)
		t.Log(patch , extName , "\n\r")
	}

	//Extension vacia
	patch , extName  := SanitizeUrlRGP(RegexpNewReactApp, 0, 7 , cadena + "." + "")
	t.Log(patch , extName , "\n\r")

	//Datos vacios
	patch , extName  = SanitizeUrlRGP(RegexpNewReactApp, 0, 7 , "" + "." + "")
	t.Log(patch , extName , "\n\r")
	
	//barra
	patch , extName  = SanitizeUrlRGP(RegexpNewReactApp, 0, 7 , "/" + "." + extension)
	t.Log(patch , extName , "\n\r")
	
	//Directorio equivocados
	patch , extName  = SanitizeUrlRGP(RegexpNewReactApp, 5, 3 , cadena + "." + extension)
	t.Log(patch , extName , "\n\r")

	//Directorio equivocados
	patch , extName  = SanitizeUrlRGP(RegexpNewReactApp, 0, 0 , cadena + "." + extension)
	t.Log(patch , extName , "\n\r")

	t.Log("Test de inyeccion")
	//Test de inyeccion
	for y := 0 ; y <= 5; y++ {

		for x := 0 ; x <= 10; x++ {

			//Cadenas con diferentes /
			cadena :=  AllLetterSlash(x)

			//todas las extensiones
			for _ , ind :=  range GetExtension() {
				
				patch , extName  := SanitizeUrlRGP(RegexpNewReactApp, 0, uint8(y) , cadena + "." + ind)
	
				t.Log(patch , extName , "\n\r")
			}
		}
	}


}


// go test -bench BenchmarkSanitizeUrlRGP -benchtime 1000x -benchmem
var RegexpNewReactApp = regexp.MustCompile(`[^a-zA-Z0-9/.]`)
func BenchmarkSanitizeUrlRGP(b *testing.B) {

	cadena :=  AllLetterSlash(10)
	ext    := "html"
    for i := 0; i < b.N; i++ {

		SanitizeUrlRGP(RegexpNewReactApp, 0, 5 , cadena + "." + ext)
      
    }
}