package http

import (
	. "dac"
	"encoding/base64"
	"net/http"
)

type messageApp struct {
	fieldsRaw map[string]int64
	fields    map[string]string
}


var messageAppFile map[string]*PublicSpaceFile
func NewMessage() *messageApp {

	if messageAppFile == nil {
		messageAppFile = make(map[string]*PublicSpaceFile)
	}

	return &messageApp{

		fieldsRaw: make(map[string]int64),
		fields:    make(map[string]string),
	}
}

func (MA *messageApp) SetMessage(field string, messageRaw string) {

	message := base64.StdEncoding.EncodeToString([]byte(messageRaw))

	MA.fieldsRaw[field] = int64(len(message))
	MA.fields[field] = message

}


var msgSpain *PublicSpaceFile

func (MA *messageApp) InitMessage(country string) {

	file := NewSfPermBytes(MA.fieldsRaw, nil, "NewMessageApp", country)

	if country == "spain" {
		msgSpain = file
	}
	messageAppFile[country] = file

	for field, message := range MA.fields {

		file.SetOneFieldString(field, message)
	}
}

func GetMessage(country string, field string) string {

	return messageAppFile[country].GetOneFieldString(field)

}

func GetMsgSp(field string) string {

	return msgSpain.GetOneFieldString(field)

}

type Speak struct {
	head     http.Header
	response http.ResponseWriter
}

func InitHeader(response http.ResponseWriter) *Speak {

	SK := new(Speak)
	SK.head = response.Header()
	SK.response = response

	return SK
}

func (SK Speak) SendMessageAppSp(field string) {

	SK.head.Set("message", GetMsgSp(field))

}

func (SK Speak) SendCloseMsgSp(field string , code int) {

	SK.head.Set("message", GetMsgSp(field))
	SK.response.WriteHeader(code)
}

func (SK Speak) SendMessage(key string, message string) {

	SK.head.Set(key,
		base64.StdEncoding.EncodeToString([]byte(message)))

}

func (SK Speak) SendMessageRaw(key string, message string) {

	SK.head.Set(key, message)

}

func (SK Speak) CloseHeader(code int) {

	SK.response.WriteHeader(code)

}
