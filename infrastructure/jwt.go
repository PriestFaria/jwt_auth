package infrastructure

import (
	"auth_backend/domain"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"time"
)

// TODO Проверка истечения срока действия токена
// TODO Переделать enum ролей на строку (ооочень не нрав типизация)
type JwtTokenGenerator struct{ secretKey []byte }

func NewJwtTokenGenerator(secretKey string) *JwtTokenGenerator {
	return &JwtTokenGenerator{secretKey: []byte(secretKey)}
}

func (j *JwtTokenGenerator) GenerateToken(user *domain.User) (string, error) {
	claims := jwt.MapClaims{
		"user_id":  user.ID,
		"email":    user.Email,
		"role":     user.Role,
		"exp":      time.Now().Add(time.Hour).Unix(),
		"issuedAt": time.Now().Unix(),
	}
	fmt.Println(claims["role"])
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString(j.secretKey)
}
