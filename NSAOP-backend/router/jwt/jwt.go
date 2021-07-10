package jwt

import (
	"errors"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"

	"nsaop/config"
)

type JWT struct {
	SigningKey []byte
}

type UserClaims struct {
	UserId   uint
	UserRole string
	jwt.StandardClaims
}

func NewJWT(which string) *JWT {
	return &JWT{
		[]byte(config.Secret[which]),
	}
}

func (j *JWT) GenerateToken(claims UserClaims) (string, error) { // encode
	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(j.SigningKey)
}

func (j *JWT) ParseToken(t string) (*UserClaims, error) { // decode
	token, err := jwt.ParseWithClaims(t, &UserClaims{}, func(token *jwt.Token) (interface{}, error) { return j.SigningKey, nil })
	if err != nil {
		if v, ok := err.(*jwt.ValidationError); ok {
			switch {
			case v.Errors&jwt.ValidationErrorMalformed != 0:
				return nil, errors.New("That's not even a token")
			case v.Errors&jwt.ValidationErrorExpired != 0:
				return nil, errors.New("Token is expired")
			case v.Errors&jwt.ValidationErrorNotValidYet != 0:
				return nil, errors.New("Token not active yet")
			default:
				return nil, errors.New("Couldn't handle this token")
			}
		}
	}
	if token != nil && token.Valid {
		if claims, ok := token.Claims.(*UserClaims); ok {
			return claims, nil
		}
	}
	return nil, errors.New("Couldn't handle this token")
}

func GenerateTokenForUser(c *gin.Context, userId uint, userRole string, which string) string {
	claims := UserClaims{
		UserId:   userId,
		UserRole: userRole,
		StandardClaims: jwt.StandardClaims{
			NotBefore: time.Now().Unix() - 1000,
			ExpiresAt: time.Now().Unix() + config.Router.GetInt64("jwt."+which+"Time"),
			Issuer:    "nsaop",
		},
	}
	token, _ := NewJWT(which).GenerateToken(claims)
	return token
}
