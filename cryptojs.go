package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"encoding/base64"
	"encoding/json"
)

// Cryptojs object
type Crypto struct {
	Iv   string `json:"iv"`
	Salt string `json:"s"`
	Ct   string `json:"ct"`
}

func buildDecryptionKey(pw string, salt string) []byte {
	var key []byte
	cPw := pw + salt
	hash := getMD5(cPw)
	key = hash

	for i := 0; i < 2; i++ {
		hash = getMD5(string(hash) + cPw)
		key = append(key, hash...)
	}

	return key[:32]
}

// Takes json string
func Decrypt(s string, password string) (string, error) {
	var obj Crypto
	if err := json.Unmarshal([]byte(s), &obj); err != nil {
		return "", err
	}

	iv, err := base64.StdEncoding.DecodeString(obj.Iv)
	if err != nil {
		return "", err
	}

	salt, err := base64.StdEncoding.DecodeString(obj.Salt)
	if err != nil {
		return "", err
	}

	ct, err := base64.StdEncoding.DecodeString(obj.Ct)
	if err != nil {
		return "", err
	}

	key := buildDecryptionKey(password, string(salt))

	block, err := aes.NewCipher(key[:32])
	if err != nil {
		panic(err)
	}

	decrypter := cipher.NewCBCDecrypter(block, iv)
	data := make([]byte, len(ct))
	decrypter.CryptBlocks(data, ct)

	return string(data), nil
}

func getMD5(str string) []byte {
	hasher := md5.New()
	hasher.Write([]byte(str))
	return hasher.Sum(nil)
}
