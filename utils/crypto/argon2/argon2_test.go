package argon2

import (
	"fmt"
	"testing"
)

func TestArgon2(t *testing.T) {
	// 生成密码哈希
	password := "SuperSecretPassword123!"
	hashedPassword, err := HashPassword([]byte(password))
	if err != nil {
		fmt.Println("Error hashing password:", err)
		return
	}
	fmt.Println("Hashed Password:", hashedPassword)

	// 验证密码
	match := VerifyPassword([]byte(password), hashedPassword)
	if match {
		fmt.Println("Password is correct!")
	} else {
		fmt.Println("Invalid password!")
	}
}
