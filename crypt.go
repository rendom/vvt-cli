package main

import "encoding/base64"
import "crypto/aes"

type EncryptionObject struct {
	Iv      string
	Chipper string
	Salt    string
}

func (t EncryptionObject) toJson() string {
	return "{}"
}

func decodeBase64(str string) string {
	data, err := base64.StdEncoding.DecodeString(str)

	if err != nil {
		panic(err)
	}

	return string(data)
}

func NewFromJson(str string) *EncryptionObject {
	// parse json
}

func New(str string, key string) *EncryptionObject {
	obj := EncryptionObject
}

func main() {
	var obj EncryptionObject
	// key: test
	iv := "zHjG8MGnM4AAwJGhBaHQiA=="
	cText := "0hUiCtsahYN9+Xy59Do9FQ=="
	salt := "OmQOJCKze9E="
	key := []byte("test")

	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}
}
