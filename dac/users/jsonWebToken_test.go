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
//BenchmarkEncripter-8  1000times 1778 ns/op 160 B/op 5 allocs/op
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

// go test -bench BenchmarkDencripter -benchtime 1000x -benchmem
//BenchmarkDencripter-8 1000 time 730.9 ns/op 112 B/op 3 allocs/op
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
//BenchmarkNewJWT-8 1000 times 3635 ns/op 622 B/op 17 allocs/op
func BenchmarkNewJWT(b *testing.B) {
  

	arrString := stringArray(1000)
	b.ResetTimer()
   
    for i := 0; i < b.N; i++ {

		b.StopTimer()
		value1 := strconv.FormatInt(int64(i) , 10) 
		value2 := strconv.FormatInt(time.Now().Unix() , 10)
		b.StartTimer()
		
		NewJWT(value1, arrString[i] ,"0.0.0.0" ,value2)

    }
}

// go test -bench BenchmarkDecodeJWT -benchtime 1000x -benchmem
//BenchmarkDecodeJWT-8 1000 times 1480 ns/op 479 B/op 14 allocs/op
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