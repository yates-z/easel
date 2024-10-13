package aes

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
)

// PKCS7 填充
func pkcs7Padding(data []byte, blockSize int) []byte {
	padding := blockSize - len(data)%blockSize
	padText := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(data, padText...)
}

// PKCS7 去除填充
func pkcs7UnPadding(data []byte) []byte {
	length := len(data)
	unPadding := int(data[length-1])
	return data[:(length - unPadding)]
}

// 生成32字节的密钥（AES-256）
func generateKey(key string) []byte {
	hash := sha256.Sum256([]byte(key))
	return hash[:]
}

// EncryptECB AES Encrypt (ECB Mode).
func EncryptECB(text, key string) (string, error) {
	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return "", err
	}

	// 填充明文
	plaintext := pkcs7Padding([]byte(text), block.BlockSize())

	ciphertext := make([]byte, len(plaintext))

	// 逐块加密
	for i := 0; i < len(plaintext); i += block.BlockSize() {
		block.Encrypt(ciphertext[i:i+block.BlockSize()], plaintext[i:i+block.BlockSize()])
	}

	return hex.EncodeToString(ciphertext), nil
}

// DecryptECB AES Decrypt (ECB Mode)
func DecryptECB(ciphertext, key string) (string, error) {
	data, err := hex.DecodeString(ciphertext)
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return "", err
	}

	if len(data)%block.BlockSize() != 0 {
		return "", fmt.Errorf("ciphertext is not a multiple of the block size")
	}

	plaintext := make([]byte, len(data))

	// 逐块解密
	for i := 0; i < len(data); i += block.BlockSize() {
		block.Decrypt(plaintext[i:i+block.BlockSize()], data[i:i+block.BlockSize()])
	}

	// 去除填充
	plaintext = pkcs7UnPadding(plaintext)
	return string(plaintext), nil
}

// EncryptCBCFixedIV : AES Encrypt (CBC Mode).
// key must be 16, 24, or 32 bytes.
// iv must equal 16.
func EncryptCBCFixedIV(text, key, iv string) (string, error) {
	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return "", err
	}

	blockSize := block.BlockSize()
	plaintext := pkcs7Padding([]byte(text), blockSize)
	ciphertext := make([]byte, len(plaintext))

	// Encrypt the data
	mode := cipher.NewCBCEncrypter(block, []byte(iv))
	mode.CryptBlocks(ciphertext, plaintext)

	return hex.EncodeToString(ciphertext), nil
}

// DecryptCBCFixedIV : AES Decrypt (CBC Mode)
func DecryptCBCFixedIV(ciphertext, key, iv string) (string, error) {
	data, err := hex.DecodeString(ciphertext)
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return "", err
	}

	blockSize := block.BlockSize()
	if len(data) < blockSize {
		return "", fmt.Errorf("ciphertext too short")
	}

	// Decrypt the data
	mode := cipher.NewCBCDecrypter(block, []byte(iv))
	mode.CryptBlocks(data, data)

	// Unpad the decrypted data
	data = pkcs7UnPadding(data)
	return string(data), nil
}

// EncryptCBC : AES Encrypt (CBC Mode).
// key must be 16, 24, or 32 bytes.
func EncryptCBC(text, key string) (string, error) {
	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return "", err
	}

	blockSize := block.BlockSize()
	plaintext := pkcs7Padding([]byte(text), blockSize)

	// Create IV
	ciphertext := make([]byte, blockSize+len(plaintext))
	iv := ciphertext[:blockSize]

	// Generate a random IV
	if _, err = io.ReadFull(rand.Reader, iv); err != nil {
		return "", err
	}

	// Encrypt the data
	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(ciphertext[blockSize:], plaintext)

	return hex.EncodeToString(ciphertext), nil
}

// DecryptCBC : AES Decrypt (CBC Mode)
func DecryptCBC(ciphertext, key string) (string, error) {
	data, err := hex.DecodeString(ciphertext)
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return "", err
	}

	blockSize := block.BlockSize()
	if len(data) < blockSize {
		return "", fmt.Errorf("ciphertext too short")
	}

	iv := data[:blockSize]
	data = data[blockSize:]

	if len(data)%blockSize != 0 {
		return "", fmt.Errorf("ciphertext is not a multiple of the block size")
	}
	// Decrypt the data
	mode := cipher.NewCBCDecrypter(block, iv)
	mode.CryptBlocks(data, data)

	// Unpad the decrypted data
	data = pkcs7UnPadding(data)
	return string(data), nil
}

// EncryptGCM AES Encrypt (GCM Mode).
func EncryptGCM(text, key string) (string, error) {
	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return "", err
	}

	// 使用GCM模式
	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonce := make([]byte, aesGCM.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	ciphertext := aesGCM.Seal(nonce, nonce, []byte(text), nil)
	return hex.EncodeToString(ciphertext), nil
}

// DecryptGCM AES Decrypt (GCM Mode).
func DecryptGCM(ciphertext, key string) (string, error) {
	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return "", err
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	data, err := hex.DecodeString(ciphertext)
	if err != nil {
		return "", err
	}

	nonceSize := aesGCM.NonceSize()
	nonce, cipherBytes := data[:nonceSize], data[nonceSize:]

	plaintext, err := aesGCM.Open(nil, nonce, cipherBytes, nil)
	if err != nil {
		return "", err
	}

	return string(plaintext), nil
}
