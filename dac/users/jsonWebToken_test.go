package http

import (
	. "dac"
	"strconv"
	"testing"
	"time"
)


var encTestOk [][]byte = [][]byte{
	[]byte("1234567890"),
	[]byte("1234567890123456"),
	[]byte("12345678901234567890"),
	
}
var encTestFail [][]byte = [][]byte{
	nil,
	[]byte(""),
}

var JWTTestOk [][]string = [][]string{
	{"test1"},
	{"test1","test2"},
	{"test1","test2","test3"},
	{"test1","test2","test3","test4"},
}

var JWTTestFail [][]string = [][]string{
	nil,
	{},
	{""},
	{"test1",""},
	{"test1","","test3"},
	{"","test2","test3"},
	{"test1","test2",""},
}

func stringArray(n int)(randArrayString []string){

	for x:= 0 ; x  <= n; x++ {

		randArrayString = append(randArrayString , "User" + strconv.Itoa(x))

	}
	return 
}


// go test -run TestEncripter -v -benchmem
func TestEncripter(t *testing.T) {

	NewBasicDac("/media/franky/GOLANG/tiuviData")
	initToken()

	for _ , text := range encTestOk {

		encripter(&text)
		if text == nil {
			t.FailNow()
		}
	}

	for _ , text := range encTestFail {

		encripter(&text)
		if text != nil {
			t.FailNow()
		}
	}

}


// go test -run TestDencripter -v -benchmem
func TestDencripter(t *testing.T) {

	NewBasicDac("/media/franky/GOLANG/tiuviData")
	initToken()

	for _ , text := range encTestOk {

		encripter(&text)
		if text == nil {
			t.FailNow()
		}
		dencripter(&text)
		if text == nil {
			t.FailNow()
		}
	}
	encTestFail = append(encTestFail , []byte("falso token"))
	for _ , text := range encTestFail {
		dencripter(&text)
		if text != nil {
			t.FailNow()
		}
	}

}



// go test -bench BenchmarkEncripter -benchtime 1000x -benchmem
func BenchmarkEncripter(b *testing.B) {
  
	arrString := stringArray(10000)
	b.ResetTimer()
   
    for i := 0; i < b.N; i++ {

		b.StopTimer()
		textBytes := []byte(arrString[i])
		b.StartTimer()

		encripter(&textBytes)
	
    }
}

// go test -bench BenchmarkEncripter -benchtime 1000x -benchmem
func BenchmarkDencripter(b *testing.B) {
  
	arrString := stringArray(10000)
	b.ResetTimer()
   
    for i := 0; i < b.N; i++ {

		b.StopTimer()
		textBytes := []byte(arrString[i])
		encripter(&textBytes)
		b.StartTimer()

		dencripter(&textBytes)

    }
}


// go test -run TestNewJWT -v -benchmem
func TestNewJWT(t *testing.T) {

	NewBasicDac("/media/franky/GOLANG/tiuviData")

	for _ , text := range JWTTestOk {

		cipherJWT := NewJWT(text...)
		if cipherJWT == "" {
			t.FailNow()
		}
	}

	for _ , text := range JWTTestFail {

		cipherJWT := NewJWT(text...)
		if cipherJWT != "" {
			t.FailNow()
		}

	}

}

// go test -run TestDecodeJWT -v -benchmem

func TestDecodeJWT(t *testing.T) {

	NewBasicDac("/media/franky/GOLANG/tiuviData")

	for _ , text := range JWTTestOk {

		cipherJWT := NewJWT(text...)
		if cipherJWT == "" {
			t.FailNow()
		}
		params := DecodeJWT(cipherJWT)
		if len(params) == 0 {
			t.FailNow()
		}
	}

	params := DecodeJWT("")
	if len(params) != 0 {
		t.FailNow()
	}

	params = DecodeJWT("falsacadena")
	if len(params) != 0 {
		t.FailNow()
	}

}


// go test -bench BenchmarkNewJWT -benchtime 1000x -benchmem
func BenchmarkNewJWT(b *testing.B) {
  

	arrString := stringArray(1000)
	b.ResetTimer()
   
    for i := 0; i < b.N; i++ {

		NewJWT(strconv.FormatInt(int64(i) , 10) ,
		arrString[i] ,
		"0.0.0.0" ,
		strconv.FormatInt(time.Now().Unix() , 10) )

    }
}

// go test -bench BenchmarkDecodeJWT -benchtime 1000x -benchmem
func BenchmarkDecodeJWT(b *testing.B) {
  
	arrString := stringArray(1000)
	b.ResetTimer()
   
    for i := 0; i < b.N; i++ {

		b.StopTimer()
		cipherJWT  := NewJWT(strconv.FormatInt(int64(i) , 10) ,
		arrString[i] ,
		"0.0.0.0" ,
		strconv.FormatInt(time.Now().Unix() , 10) )
		b.StartTimer()
		
		DecodeJWT(cipherJWT  )

    }
}