package http

import (
	_"log"
	"math/rand"
	"strconv"
	"testing"
	"time"
	."dac"
)


func init() {
    rand.Seed(time.Now().UnixNano())
}

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func RandStringRunes(n int) string {
    b := make([]rune, n)
    for i := range b {
        b[i] = letterRunes[rand.Intn(len(letterRunes))]
    }
    return string(b)
}

func stringArray(n int)(randArrayString []string){

	for x:= 0 ; x  <= n; x++ {

		randArrayString = append(randArrayString , "User" + strconv.Itoa(x),RandStringRunes(2))

	}
	return 
}

var keyString string = "6368616e676520746869732070617373"
var text string = ""

// go test -run TestEncripter -v -benchmem
func TestEncripter(t *testing.T) {

	NewBasicDac("/media/franky/GOLANG/tiuviData")
	initToken()
	print("Start: " , text , " tamaño: ",len(text) , "\n")
	enc, count := encripter(text)
	print("Start: " , string(enc) , " tamaño: ",len(text) , "\n")
	dencripter(&enc,count)
	print("end: " , string(enc) , " tamaño: ",len(enc) , "\n")
}

// go test -bench BenchmarkEncripter -benchtime 1000x -benchmem
func BenchmarkEncripter(b *testing.B) {
  

	arrString := stringArray(10000)
	b.ResetTimer()
   
    for i := 0; i < b.N; i++ {

		b.StopTimer()
		print("start: " , string(arrString[i]) ,"tamaño: ",len(arrString[i]) , "\n")
		b.StartTimer()

		enc, count := encripter(arrString[i])
		b.ResetTimer()
		dencripter(&enc,count)

		b.StopTimer()
		print("end: " , string(enc) , "tamaño: ",len(enc) , "\n")
		b.StartTimer()
    }
}

// go test -run TestNewJWT -v -benchmem
func TestNewJWT(t *testing.T) {

	NewBasicDac("/media/franky/GOLANG/tiuviData")

	NewJWT(0 , "franky" , "0.0.0.0" , time.Now() )


}
// go test -bench BenchmarkNewJWT -benchtime 1000x -benchmem
func BenchmarkNewJWT(b *testing.B) {
  

	arrString := stringArray(1000)
	b.ResetTimer()
   
    for i := 0; i < b.N; i++ {

		NewJWT(0 , arrString[i] , "0.0.0.0" , time.Now() )

    }
}