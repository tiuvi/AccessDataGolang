package dac

import (
	"strings"
	"testing"
)

// Cadena de inyeccion asci sin slash
func AllLetterNoSlash()(cadena string) {

	for x := 0 ;x < 128; x++ {
		if x == 12  || x == 47 {
			continue
		}
		cadena += string(byte(x))
	}
	return
}

//Crea una cadena url de inyeccion asci
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
	"/ruta./",
	"/.ruta/",
	"/.ruta./",
	"./ruta./",
	"//ruta//",
	"/ruta//",
	"//ruta//",
	"/./ruta/./",
	"/.../ruta/.../",
	".../.../ruta/.../",
	"...../...../ruta/.../",
}

var urlTestOk []string = []string{
	"ruta",
	"/ruta1",
	"/ruta1/ruta2",
	"/ruta1/ruta2/ruta3",
	"/ruta1/ruta2/ruta3/ruta4",
	"/ruta1/ruta2/ruta3/ruta4/ruta5",
	"/ruta1/ruta2/ruta3/ruta4/ruta5/ruta6",
	"/ruta1/ruta2/ruta3/ruta4/ruta5/ruta6/ruta7",
	"/ruta1/ruta2/ruta3/ruta4/ruta5/ruta6/ruta7/ruta8",
}

var boleanTest [][2]bool = [][2]bool{
	{true  ,true },
	{true  ,false },
	{false  ,true },
	{false  ,false },
}


//Genera url de prueba
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
	
			if ext == "" && bolean[1] {
				t.Fatal(BCG("Url: ") ,url, BCG("patch: "),path , BCG("extension: "), ext )	
			}
			t.Log(BCG("Url: ") ,url, BCG("patch: "),path , BCG("extension: "), ext )
		}	
	}
}

// go test -run TestSliceUrlToPath -v
func TestSliceUrlToPath(t *testing.T){

	var levelU uint8
	var maxLevelU uint8
	for levelU = 0; levelU <= 6 ; levelU++{

		for maxLevelU = 0; maxLevelU <= 6 ; maxLevelU++{

			for _, url :=  range urlTestOk {

				path := sliceUrlToPath(levelU ,maxLevelU ,url )

				t.Log(BCG("Url: ") ,url,BCG("levelU: "),levelU , BCG("maxLevelU: "),maxLevelU ,BCG("path: "), path )
			}
		}
	}

}

// go test -run TestSanitizePath -v   //Regex: [^a-zA-Z0-9]
func TestSanitizePath(t *testing.T){

	for _ , bolean := range boleanTest{

		t.Log("Url Valida: ",bolean[0] ,"Extension Valida: ", bolean[1])

		for _ , url := range urlProve(bolean[0] , bolean[1]){
	
			patch := sliceUrlToPath(0 ,10 ,url )
			patch  = sanitizePath( patch)
	
			t.Log(BCG("Url: ") ,url, BCG("patch: "),patch  )
		}	
	}

}


// go test -run TestSanitizePathRGP -v   //Regex: [^a-zA-Z0-9/.]
func TestSanitizePathRGP(t *testing.T){

	for _ , bolean := range boleanTest{

		t.Log("Url Valida: ",bolean[0] ,"Extension Valida: ", bolean[1])

		for _ , url := range urlProve(bolean[0] , bolean[1]){
	
			patch := sliceUrlToPath(0 ,10 ,url )
			patch  = sanitizePathRGP( patch)
	
			t.Log(BCG("Url: ") ,url, BCG("patch: "),patch  )
		}	
	}
}


// go test -run TestSanitizeUrl -v   //Regex: [^a-zA-Z0-9]
func TestSanitizeUrl(t *testing.T){

	for _ , bolean := range boleanTest{

		t.Log("Url Valida: ",bolean[0] ,"Extension Valida: ", bolean[1])

		for _ , url := range urlProve(bolean[0] , bolean[1]){
	
			patch, ext  := SanitizeUrl(0,6, url)
	
			t.Log(BCG("Url: ") ,url, BCG("patch: "),patch , BCG("extension: "), ext )
		}	
	}
}

// go test -run TestSanitizeUrlRGP -v    //Regex: [^a-zA-Z0-9/.]
func TestSanitizeUrlRGP(t *testing.T){

	for _ , bolean := range boleanTest{

		t.Log("Url Valida: ",bolean[0] ,"Extension Valida: ", bolean[1])

		for _ , url := range urlProve(bolean[0] , bolean[1]){
	
			patch, ext  := SanitizeUrlRGP(0,6, url)
	
			t.Log(BCG("Url: ") ,url, BCG("patch: "),patch , BCG("extension: "), ext )
		}	
	}
	
}


// go test -bench BenchmarkSanitizeUrl -benchtime 1000x -benchmem
func BenchmarkSanitizeUrl(b *testing.B) {

    for i := 0; i < b.N; i++ {

		for _ , bolean := range boleanTest{

			for _ , url := range urlProve(bolean[0] , bolean[1]){
				b.ResetTimer()
				SanitizeUrl(0,6, url)
			}	
		}
    }
}

// go test -bench BenchmarkSanitizeUrlRGP -benchtime 1000x -benchmem
func BenchmarkSanitizeUrlRGP(b *testing.B) {

    for i := 0; i < b.N; i++ {

		for _ , bolean := range boleanTest{

			for _ , url := range urlProve(bolean[0] , bolean[1]){
				b.ResetTimer()
				SanitizeUrlRGP(0,6, url)
			}	
		}
    }
}