package crypto

import (
	"github.com/golang-jwt/jwt/v4"
	config_loader "github.com/kwisnia/inzynierka-backend/internal/pkg/config"
	"golang.org/x/crypto/bcrypt"
	"log"
	"strings"
	"time"
)

const oneYear = time.Hour * 24 * 365

type JwtClaims struct {
	jwt.RegisteredClaims
	Username string
	UserID   uint
	Exp      int64
}

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

func CreateToken(username string, userID uint) (string, error) {
	secret := config_loader.GetEnv("JWT_SECRET")
	var err error
	accessTokenClaims := JwtClaims{
		Username: username,
		Exp:      time.Now().Add(oneYear).Unix(),
		UserID:   userID,
	}
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS512, accessTokenClaims)
	token, err := accessToken.SignedString([]byte(secret))
	if err != nil {
		return "error creating token", err
	}
	return token, nil
}

func ValidateToken(tokenString string) (*string, *uint, bool) {
	secret := config_loader.GetEnv("JWT_SECRET")
	jwtString := strings.Split(tokenString, "Bearer ")[1]
	token, err := jwt.ParseWithClaims(jwtString, &JwtClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
	if err != nil {
		return nil, nil, false
	}
	return &token.Claims.(*JwtClaims).Username, &token.Claims.(*JwtClaims).UserID, token.Valid
}
