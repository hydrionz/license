package crypto

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"errors"
)

// AesCbcPkcs7 golang implementation of aes-cbc-pkcs7 encryption/decryption with base64 input/output
type AesCbcPkcs7 struct {
	Key []byte // Allows 16, 24, 32 byte lengths
	Iv  []byte // Only allows 16 byte length
}

func (s AesCbcPkcs7) Encrypt(text []byte) ([]byte, error) {
	if len(text) == 0 {
		return nil, errors.New("text is empty")
	}
	// Generate cipher.Block
	block, err := aes.NewCipher(s.Key)
	if err != nil {
		return nil, err
	}
	// Pad content if less than 16 characters
	blockSize := block.BlockSize()
	originData := s.pad(text, blockSize)
	// Encryption mode
	blockMode := cipher.NewCBCEncrypter(block, s.Iv)
	// Encrypt, output to []byte array
	encrypt := make([]byte, len(originData))
	blockMode.CryptBlocks(encrypt, originData)
	return encrypt, nil
}

func (s AesCbcPkcs7) Decrypt(text string) ([]byte, error) {
	if len(text) == 0 {
		return []byte(text), nil
	}
	decodeData, err := base64.StdEncoding.DecodeString(text)
	if err != nil {
		return []byte(text), err
	}
	if len(decodeData) == 0 {
		return []byte(text), nil
	}
	// Generate cipher.Block
	block, _ := aes.NewCipher(s.Key)
	// Decryption mode
	blockMode := cipher.NewCBCDecrypter(block, s.Iv)
	// Output to []byte array
	originData := make([]byte, len(decodeData))
	blockMode.CryptBlocks(originData, decodeData)
	// Remove padding and return
	return s.unPad(originData), nil
}

func (s AesCbcPkcs7) pad(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padText := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padText...)
}

func (s AesCbcPkcs7) unPad(ciphertext []byte) []byte {
	length := len(ciphertext)
	// Remove the last padding
	unPadding := int(ciphertext[length-1])
	return ciphertext[:(length - unPadding)]
}
