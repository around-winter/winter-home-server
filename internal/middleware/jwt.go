package middleware

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"
)

var jwtSecret = []byte(getJWTSecret())

func getJWTSecret() string {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "winter-jwt-secret-key-2024"
	}
	return secret
}

type Claims struct {
	UserID   int    `json:"user_id"`
	Username string `json:"username"`
	Exp      int64  `json:"exp"`
}

func GenerateToken(userID int, username string) (string, error) {
	expirationTime := time.Now().Add(24 * time.Hour) // 24小时有效期

	claims := &Claims{
		UserID:   userID,
		Username: username,
		Exp:      expirationTime.Unix(),
	}

	header := base64.RawURLEncoding.EncodeToString([]byte(`{"alg":"HS256","typ":"JWT"}`))
	claimsJSON, err := json.Marshal(claims)
	if err != nil {
		return "", err
	}
	payload := base64.RawURLEncoding.EncodeToString(claimsJSON)

	h := hmac.New(sha256.New, jwtSecret)
	h.Write([]byte(header + "." + payload))
	signature := base64.RawURLEncoding.EncodeToString(h.Sum(nil))

	return header + "." + payload + "." + signature, nil
}

func ValidateToken(tokenString string) (*Claims, error) {
	parts := strings.Split(tokenString, ".")
	if len(parts) != 3 {
		return nil, fmt.Errorf("token 格式错误")
	}

	h := hmac.New(sha256.New, jwtSecret)
	h.Write([]byte(parts[0] + "." + parts[1]))
	expectedMAC := h.Sum(nil)
	actualMAC, err := base64.RawURLEncoding.DecodeString(parts[2])
	if err != nil || !hmac.Equal(expectedMAC, actualMAC) {
		return nil, fmt.Errorf("token 签名验证失败")
	}

	payload, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return nil, err
	}

	var claims Claims
	if err := json.Unmarshal(payload, &claims); err != nil {
		return nil, err
	}

	if time.Now().Unix() > claims.Exp {
		return nil, fmt.Errorf("token 已过期")
	}

	return &claims, nil
}
