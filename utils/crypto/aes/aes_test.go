package aes

import (
	"fmt"
	"testing"
)

func TestAESECB(t *testing.T) {

	key := string(generateKey("your-secret-key1"))
	// 加密
	encrypted, err := EncryptECB("Hello, ECB!", key)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Encrypted: %s\n", encrypted)

	// 解密
	decrypted, err := DecryptECB(encrypted, key)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Decrypted: %s\n", decrypted)
}

func TestAESCBCFixedIV(t *testing.T) {

	key := string(generateKey("your-secret-key1"))
	// 加密
	encrypted, err := EncryptCBCFixedIV("Hello, CBC!", key, key[:16])
	if err != nil {
		panic(err)
	}
	fmt.Printf("Encrypted: %s\n", encrypted)

	// 解密
	decrypted, err := DecryptCBCFixedIV(encrypted, key, key[:16])
	if err != nil {
		panic(err)
	}
	fmt.Printf("Decrypted: %s\n", decrypted)
}

func TestAESCBC(t *testing.T) {

	key := string(generateKey("your-secret-key1"))
	// 加密
	encrypted, err := EncryptCBC("Hello, CBC!", key)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Encrypted: %s\n", encrypted)

	// 解密
	decrypted, err := DecryptCBC(encrypted, key)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Decrypted: %s\n", decrypted)
}

func TestAESGCM(t *testing.T) {

	key := string(generateKey("your-secret-key1"))

	// 加密
	encrypted, err := EncryptGCM("Hello, GCM!", key)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Encrypted: %s\n", encrypted)

	// 解密
	decrypted, err := DecryptGCM(encrypted, key)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Decrypted: %s\n", decrypted)
}
