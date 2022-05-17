package http

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	. "dac"
	"encoding/base64"
	"strconv"
	"time"
	"strings"
	."fmt"
)

func Falsefunction(){
	Println("")
}

type spaceCipher struct {
	key   []byte
	block cipher.Block
}

var globalCipher *spaceCipher

func initToken() {

	var err error

	ST := new(spaceCipher)

	gCipher := NewSfopenFile(map[string]int64{"keyPrivate": 32}, nil, "users", "globalCipher")

	if found := gCipher.ExistReadfield("keyPrivate"); *found == false {

		bufferBytes := make([]byte, 32)
		_, err := rand.Read(bufferBytes)
		if err != nil &&
			gCipher.NRESM(err != nil, err.Error(), "users", "tokenErrors") {
			return
		}

		ST.key = bufferBytes
		gCipher.SetOneFieldRaw("keyPrivate", &bufferBytes)

	} else {

		ST.key = gCipher.GetOneFieldBytesRaw("keyPrivate")
	}

	//Crea una nueva clave de cifrado
	ST.block, err = aes.NewCipher(ST.key)
	if err != nil &&
		gCipher.NRESM(err != nil, err.Error(), "users", "tokenErrors") {
		return
	}

	globalCipher = ST
}

func encripter(text string) (ciphertext []byte, count int) {

	if len(text) == 0 {
		return nil, 0
	}

	plaintext := []byte(text)

	//El texto debe ser cifrado en bloques de 16 bytes
	//Añadimos padding hasta completar.
	if count = len(plaintext) % aes.BlockSize; count != 0 {

		count = aes.BlockSize - count
		padding := len(plaintext) + count
		SpacePaddingPointer(&plaintext, [2]int64{0, int64(padding)})

	}

	//Creamos un array con el tamaño de bloque + el texto
	ciphertext = make([]byte, aes.BlockSize+len(plaintext))

	//Añadimos bytes aleatorios al principio del texto
	_, err := rand.Read(ciphertext[:aes.BlockSize])
	if err != nil &&
		NRELDACG(err != nil, err.Error(), "users", "tokenErrors") {
		return nil, 0
	}

	mode := cipher.NewCBCEncrypter(globalCipher.block, ciphertext[:aes.BlockSize])
	mode.CryptBlocks(ciphertext[aes.BlockSize:], plaintext)

	return ciphertext, count
}

func dencripter(ciphertext *[]byte, count int) {

	if len(*ciphertext) < aes.BlockSize {
		*ciphertext = nil
		return
	}

	vectorInit := (*ciphertext)[:aes.BlockSize]
	*ciphertext = (*ciphertext)[aes.BlockSize:]

	if len(*ciphertext)%aes.BlockSize != 0 {

		*ciphertext = nil
		return
	}

	mode := cipher.NewCBCDecrypter(globalCipher.block, vectorInit)

	mode.CryptBlocks(*ciphertext, *ciphertext)

	*ciphertext = (*ciphertext)[:len(*ciphertext)-count]

}



func NewJWT(line int64, UserName string, ip string , timeNow time.Time )[]byte{

	if globalCipher == nil {
		initToken()
	}

	 //JWT := strconv.FormatInt(line , 10) + UserName + ip + strconv.FormatInt(timeNow.Unix() , 10)
	

	//	 Println("Original: ",line , UserName , ip  , timeNow.Unix() )

	JWT := strings.Join( []string{
		base64.StdEncoding.EncodeToString( []byte(strconv.FormatInt(line , 10))),
		base64.StdEncoding.EncodeToString( []byte( UserName )),
		base64.StdEncoding.EncodeToString( []byte( ip )),
		base64.StdEncoding.EncodeToString( []byte( strconv.FormatInt(timeNow.Unix() , 10))),
	} , ".")
	
	//	Println("Base 64: ",JWT , "tamaño: " ,len(JWT))

	ciphertext, count := encripter(JWT)

	//Println("encripter: ", string(ciphertext) )

	dencripter(&ciphertext , count)
	//Println("dencripter: ",string(ciphertext), "tamaño: " ,len(ciphertext))

	tokenRaw := strings.Split(string(ciphertext) , ".")
	var decode []byte
	for _ , str := range tokenRaw{

		decode, _ = base64.StdEncoding.DecodeString(str)

		
	

	}
	

	return decode
}