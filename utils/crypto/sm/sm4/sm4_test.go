package sm4

import (
	"fmt"
	"testing"
)

func TestSM4ECB(t *testing.T) {

	key := "your-secret-key1"
	// 加密
	encrypted, err := EncryptECB("Hello, SM4!", key)
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

func TestSM4CBCFixedIV(t *testing.T) {

	key := "your-secret-key1"
	// 加密
	encrypted, err := EncryptCBCFixedIV("Hello, SM4!", key, key[:16])
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

func TestCBC(t *testing.T) {

	key := "your-secret-key1"
	// 加密
	encrypted, err := EncryptCBC("Hello, SM4!", key)
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
