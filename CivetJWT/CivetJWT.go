package CivetJWT

import (
	"github.com/dgrijalva/jwt-go"
	"time"
)

type JwtConf struct {
	jwtSecret []byte
	Issuer    string
}
type Claims struct {
	LogType string `json:"type"`
	UserID  string `json:"userid"`
	jwt.StandardClaims
}

// GenerateToken 产生token的函数
func (c *JwtConf) GenerateToken(logtype string, userid string) (string, error) {
	nowTime := time.Now()
	expireTime := nowTime.Add(3 * time.Hour)

	claims := Claims{
		logtype,
		userid,
		jwt.StandardClaims{
			ExpiresAt: expireTime.Unix(),
			Issuer:    c.Issuer,
		},
	}
	//
	tokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err := tokenClaims.SignedString(c.jwtSecret)
	return token, err
}

// ParseToken 验证token的函数
func (c *JwtConf) ParseToken(token string) (*Claims, error) {
	tokenClaims, err := jwt.ParseWithClaims(token, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return c.jwtSecret, nil
	})

	if tokenClaims != nil {
		if claims, ok := tokenClaims.Claims.(*Claims); ok && tokenClaims.Valid {
			return claims, nil
		}
	}
	return nil, err
}
