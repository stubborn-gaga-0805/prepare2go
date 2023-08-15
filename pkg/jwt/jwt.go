package jwt

import (
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"time"
)

type Instance struct {
	secretKey string
}

type CustomClaims struct {
	UserId int64
	UserNo string
	jwt.StandardClaims
}

func NewJwtInstance(c Jwt) *Instance {
	return &Instance{
		secretKey: c.SecretKey,
	}
}

func (j *Instance) GenToken(userId int64) (jwtToken string, err error) {
	maxAge := 60 * 60 * 72
	customClaims := &CustomClaims{
		UserId: userId,
		StandardClaims: jwt.StandardClaims{
			IssuedAt:  time.Now().Unix(),                                          // 签发时间
			NotBefore: time.Now().Unix(),                                          // 生效时间
			ExpiresAt: time.Now().Add(time.Duration(maxAge) * time.Second).Unix(), // 过期时间
			Issuer:    fmt.Sprintf("LoginUser:%d", userId),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, customClaims)
	jwtToken, err = token.SignedString([]byte(j.secretKey))
	if err != nil {
		return "", err
	}
	return jwtToken, nil
}

func (j *Instance) DecryptToken(tokenString string) (claims *CustomClaims, err error) {
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(j.secretKey), nil
	})
	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		return claims, nil
	} else {
		return nil, err
	}
}
