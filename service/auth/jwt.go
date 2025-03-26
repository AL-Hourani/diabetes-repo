package auth

import (
	"strconv"
	"time"

	"github.com/AL-Hourani/care-center/config"
	"github.com/golang-jwt/jwt/v5"
)

func CreateJWT(secret []byte, patientID int) (string, error) {
	expiration := time.Second * time.Duration(config.Envs.JWTExpirationInSecond)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256 , jwt.MapClaims{
		"patientID":strconv.Itoa(patientID),
		"expiredAt":time.Now().Add(expiration).Unix(),
	})

	tokenString ,err := token.SignedString(secret)
	if err != nil {
		return "",err
	}

	return tokenString , nil
}

