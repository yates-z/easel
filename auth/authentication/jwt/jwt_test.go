package jwt

import (
	"fmt"
	"testing"
	"time"
)

func Test_JWT(t *testing.T) {
	// 定义 secret 和用户信息
	key := []byte("mysecretkey")
	payload := Payload{
		Sub: "123456",
		Iat: time.Now().Unix(),
		Exp: time.Now().Add(time.Hour).Unix(),
		Iss: "myapp",
	}

	// 生成 JWT
	token := NewToken(NewMethodHS256, payload)
	tokenStr, err := token.Generate(key)
	if err != nil {
		fmt.Println("Failed to generate JWT:", err)
		return
	}
	fmt.Println("Generated JWT:", tokenStr)

	// 验证 JWT
	token, err = FromToken(tokenStr)
	if err != nil {
		fmt.Println("Parse JWT token failed:", err)
	}
	valid := token.Verify(key)
	if valid {
		fmt.Println("JWT is valid. Payload:", token.Payload)
	} else {
		fmt.Println("Invalid JWT")
	}
}
