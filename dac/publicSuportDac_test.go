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

var extensionesTestFail []string = []string{
	//Sin extension
	"",
	//Extension vacia
	".",
	//Doble extension primera vacia
	".ext.",
	//Doble extension las dos vacias
	"..",
}
var extensionesTestOk []string = []string{
	".ext",
	".ext.map",
}

var urlTestFail []string = []string{
	"",
	"/",
	"//",
	"//ruta//",
	"/ruta//",
	"//ruta//",
}

var urlTestOk []string = []string{
	"ruta",
	"/ruta1",
	"/ruta1/ruta2",
	"/ruta1/ruta2/ruta3",
	"/ruta1/ruta2/ruta3/ruta4",
}

var boleanTest [][2]bool = [][2]bool{
	{true  ,true },
	{true  ,false },
	{false  ,true },
	{false  ,false },
}



func urlProve (url bool , extension bool)(urls[]string){

	if url  && extension {
		for _ , urlVal := range urlTestOk{

			for _ , extVal := range extensionesTestOk{

				urls = append(urls, urlVal + extVal)
			}
		}
	}

	if !url  && extension {
		for _ , urlVal := range urlTestFail{

			for _ , extVal := range extensionesTestOk{

				urls = append(urls, urlVal + extVal)
			}
		}
	}

	if url  && !extension {
		for _ , urlVal := range urlTestOk{

			for _ , extVal := range extensionesTestFail{

				urls = append(urls, urlVal + extVal)
			}
		}
	}

	if !url  && !extension {
		for _ , urlVal := range urlTestFail{

			for _ , extVal := range extensionesTestFail{

				urls = append(urls, urlVal + extVal)
			}
		}
	}


	return
}

// go test -run TestCutExtensionToPath -v
func TestCutExtensionToPath(t *testing.T){

	for _ , bolean := range boleanTest{

		t.Log("Url Valida: ",bolean[0] ,"Extension Valida: ", bolean[1])

		for _ , url := range urlProve(bolean[0] , bolean[1]){
	
			path , ext := cutExtensionToPath(url)
			if ext == "" && bolean[1]{
				t.Fatal(BCG("Url: ") ,url, BCG("patch: "),path , BCG("extension: "), ext )	
			}
				
			if ext != "" && !bolean[1] {
				t.Fatal(BCG("Url: ") ,url, BCG("patch: "),path , BCG("extension: "), ext )	
			}
			t.Log(BCG("Url: ") ,url, BCG("patch: "),path , BCG("extension: "), ext )
		}	
	}
}







var RegexpNewReactApp = regexp.MustCompile(`[^a-zA-Z0-9/.]`)
func BenchmarkSanitizeUrlRGP(b *testing.B) {

	cadena :=  AllLetterSlash(10)
	ext    := "html"
    for i := 0; i < b.N; i++ {

		SanitizeUrlRGP(RegexpNewReactApp, 0, 5 , cadena + "." + ext)
      
    }
}