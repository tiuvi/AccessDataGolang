package http

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	. "dac"
	"encoding/base64"
	"strings"

	. "fmt"
)

func Falsefunction() {
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
			EDAC &&
			gCipher.NRESF(err != nil, err.Error(), "users", "tokenErrors") {
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
		EDAC &&
		gCipher.NRESF(err != nil, err.Error(), "users", "tokenErrors") {
		return
	}

	globalCipher = ST
}

func encripter(bufferBytes *[]byte) {

	if len(*bufferBytes) == 0 {
		*bufferBytes = nil
		return 
	}

	//El texto debe ser cifrado en bloques de 16 bytes
	//Añadimos padding hasta completar.
	if count := len(*bufferBytes) % aes.BlockSize; count != 0 {

		count = aes.BlockSize - count
		padding := len(*bufferBytes) + count
		SpacePaddingPointer(bufferBytes, [2]int64{0, int64(padding)})

	}

	//Creamos un array con el tamaño de bloque + el texto
	ciphertext := make([]byte, aes.BlockSize+len(*bufferBytes))

	//Añadimos bytes aleatorios al principio del texto
	_, err := rand.Read(ciphertext[:aes.BlockSize])
	if err != nil &&
		NRELDACG(err != nil, err.Error(), "users", "tokenErrors") {
		*bufferBytes = nil
		return 
	}

	mode := cipher.NewCBCEncrypter(globalCipher.block, ciphertext[:aes.BlockSize])
	mode.CryptBlocks(ciphertext[aes.BlockSize:], *bufferBytes)

	*bufferBytes = ciphertext
	return
}

func dencripter(ciphertext *[]byte) {

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

	SpaceTrimPointer(ciphertext)

}

func NewJWT(params ...string) (cipherJWT string) {

	if len(params) == 0 {
		return
	}

	if globalCipher == nil {
		initToken()
	}

	for ind, param := range params {
		if param == ""{
			return
		}
		params[ind] = base64.StdEncoding.EncodeToString([]byte(param))
	}

	var JWT []byte
	if len(params) > 1 {

		JWT = []byte(strings.Join(params, "."))
	} else {

		JWT = []byte(params[0])
	}

	encripter(&JWT)
	if JWT == nil {
		return
	}

	return base64.StdEncoding.EncodeToString(JWT)
}

func DecodeJWT(JWT string) (params []string) {

	if len(JWT) == 0 {
		return
	}
	
	cipherJWT , err := base64.StdEncoding.DecodeString(JWT)
	if err != nil {
		return
	}
	dencripter(&cipherJWT)
	if cipherJWT == nil {
		return 
	}
	params = strings.Split(string(cipherJWT), ".")

	for ind, param := range params {

		if param == ""{
			return []string{}
		}
		decode, err := base64.StdEncoding.DecodeString(param)
		if err != nil {
			return []string{}
		}

		params[ind] = string(decode)

	}

	return params
}
