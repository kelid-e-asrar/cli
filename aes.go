package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/rand"
	"encoding/hex"
	"io"
)

// // Encrypt the password with an AES key
// func Encrypt(key []byte, input string) error {
// 	c, err := aes.NewCipher(key)
// 	if err != nil {
// 		log.Fatal("error in encrypting the passphrase")
// 		return err
// 	}
// 	out := make([]byte, len(input))
// 	c.Encrypt(out, []byte(input))

// 	return nil
// }

// // Decrypt the password with an AES key
// func Decrypt(key []byte, ct string) (*string, error) {
// 	cipherText, _ := hex.DecodeString(ct)
// 	c, err := aes.NewCipher(key)
// 	if err != nil {
// 		log.Fatal("error in decrypting the passphrase")
// 		return nil, err
// 	}
// 	plain := make([]byte, len(cipherText))
// 	c.Decrypt(plain, cipherText)
// 	s := string(plain[:])

// 	return &s, nil
// }

func createHash(key string) string {
	hasher := md5.New()
	hasher.Write([]byte(key))
	return hex.EncodeToString(hasher.Sum(nil))
}

func encrypt(data []byte, passphrase string) []byte {
	block, _ := aes.NewCipher([]byte(createHash(passphrase)))
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		panic(err.Error())
	}
	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		panic(err.Error())
	}
	ciphertext := gcm.Seal(nonce, nonce, data, nil)
	return ciphertext
}

func decrypt(data []byte, passphrase string) []byte {
	key := []byte(createHash(passphrase))
	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err.Error())
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		panic(err.Error())
	}
	nonceSize := gcm.NonceSize()
	nonce, ciphertext := data[:nonceSize], data[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		panic(err.Error())
	}
	return plaintext
}
