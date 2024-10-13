package sm2

import (
	"encoding/hex"
	"fmt"
	"testing"
)

func TestSM2(t *testing.T) {
	privKey, pubKey, err := GenerateKeyWithSalt([]byte("hello world"))
	if err != nil {
		t.Fatal(err)
	}
	ciphertext, err := EncryptString("Hello SM2!", hex.EncodeToString(pubKey.GetRawBytes()), C1C3C2)
	if err != nil {
		t.Fatal(err)
	}
	plaintext, err := DecryptString(ciphertext, hex.EncodeToString(privKey.GetRawBytes()), C1C3C2)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(plaintext)
}
