package jwt

import (
	"encoding/json"
	"strings"
	"time"
)

// Header definition.
type Header struct {
	Alg string `json:"alg"`
	Typ string `json:"typ"`
}

// Payload definition.
type Payload struct {
	Sub string `json:"sub"` // 用户标识
	Iat int64  `json:"iat"` // 签发时间
	Exp int64  `json:"exp"` // 过期时间
	Iss string `json:"iss"` // 签发者
}

// Token is a JWT token.
type Token struct {
	// Raw is the raw token.
	Raw string
	// Header is the first part of the token.
	Header Header
	// Payload is the second part of the token.
	Payload Payload
	// Signature is the third part of the token.
	Signature []byte
	// Method is the signing method used.
	Method SigningMethod
}

func NewToken(m func() SigningMethod, payload Payload) *Token {
	method := m()
	return &Token{
		Header: Header{
			Alg: method.Alg(),
			Typ: "JWT",
		},
		Payload: payload,
		Method:  method,
	}
}

// FromToken generates a Token from a string.
func FromToken(token string) (*Token, error) {
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return nil, ErrTokenInvalid
	}

	headerEncoded, payloadEncoded, signatureEncoded := parts[0], parts[1], parts[2]
	// Decode and unmarshal header.
	headerJSON, err := Base64URLDecode(headerEncoded)
	if err != nil {
		return nil, ErrTokenHeaderInvalid
	}
	header := Header{}
	err = json.Unmarshal(headerJSON, &header)
	if err != nil {
		return nil, ErrTokenHeaderInvalid
	}

	// Get signing method.
	method := GetSigningMethod(header.Alg)
	if method == nil {
		return nil, ErrMethodNotFound
	}

	// Decode and unmarshal payload.
	payloadJSON, err := Base64URLDecode(payloadEncoded)
	if err != nil {
		return nil, ErrTokenPayloadInvalid
	}
	payload := Payload{}
	err = json.Unmarshal(payloadJSON, &payload)
	if err != nil {
		return nil, ErrTokenPayloadInvalid
	}

	signature, err := Base64URLDecode(signatureEncoded)
	if err != nil {
		return nil, ErrSignatureInvalid
	}

	return &Token{
		Raw:       token,
		Header:    header,
		Payload:   payload,
		Signature: signature,
		Method:    method,
	}, nil
}

// Generate generates a new JWT token.
func (t *Token) Generate(key []byte) (string, error) {

	// 编码 Header 和 Payload
	headerJSON, err := json.Marshal(t.Header)
	if err != nil {
		return "", err
	}
	payloadJSON, err := json.Marshal(t.Payload)
	if err != nil {
		return "", err
	}

	headerEncoded := Base64URLEncode(headerJSON)
	payloadEncoded := Base64URLEncode(payloadJSON)

	// 生成签名
	t.Signature, err = t.Method.Sign(headerEncoded+"."+payloadEncoded, key)
	if err != nil {
		return "", err
	}
	signatureEncoded := Base64URLEncode(t.Signature)
	// 拼接成 JWT
	t.Raw = headerEncoded + "." + payloadEncoded + "." + signatureEncoded
	return t.Raw, nil
}

// Verify verifies the JWT token.
func (t *Token) Verify(key []byte) bool {
	parts := strings.Split(t.Raw, ".")

	err := t.Method.Verify(parts[0]+"."+parts[1], t.Signature, key)
	if err != nil {
		return false
	}

	if time.Now().Unix() > t.Payload.Exp {
		return false
	}

	return true
}
