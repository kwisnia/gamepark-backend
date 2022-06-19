package crypto

import (
	"log"
	"time"

	"github.com/golang-jwt/jwt/v4"
	config_loader "github.com/kwisnia/inzynierka-backend/internal/pkg/config"
	"golang.org/x/crypto/bcrypt"
)

const oneYear = time.Hour * 24 * 365

func HashAndSalt(pwd string) string {
	bytePwd := []byte(pwd)
	hash, err := bcrypt.GenerateFromPassword(bytePwd, bcrypt.MinCost)
	if err != nil {
		log.Println(err)
	}
	return string(hash)
}

func ComparePasswords(hashedPwd string, plainPwd string) bool {
	byteHash := []byte(hashedPwd)
	bytePwd := []byte(plainPwd)
	err := bcrypt.CompareHashAndPassword(byteHash, bytePwd)
	return err == nil
}

func CreateToken(email string) (string, error) {
	config := config_loader.Config

	var err error
	accessTokenClaims := jwt.MapClaims{}
	accessTokenClaims["email"] = email
	accessTokenClaims["exp"] = time.Now().Add(oneYear).Unix()
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS512, accessTokenClaims)
	token, err := accessToken.SignedString([]byte(config.Server.Secret))
	if err != nil {
		return "error creating token", err
	}
	return token, nil
}

func ValidateToken(tokenString string) bool {
	config := config_loader.Config
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.Server.Secret), nil
	})
	if err != nil {
		return false
	}
	return token.Valid
}
