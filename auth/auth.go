package auth

import (
	"crypto/rsa"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/dgrijalva/jwt-go/request"
)

type Auth interface {
	Validate(string) bool
}

type Jwt struct {
	PrivateKeyPath string
	PublicKeyPath  string
}

func (j *Jwt) CreateToken(ID int) (string, error) {
	token := jwt.New(jwt.SigningMethodRS256)
	claims := make(jwt.MapClaims)
	claims["exp"] = time.Now().Add(time.Hour * time.Duration(1)).Unix()
	claims["iat"] = time.Now().Unix()
	claims["id"] = ID
	token.Claims = claims
	return token.SignedString(j.GetPrivateKey())
}

func (j *Jwt) GetPrivateKey() *rsa.PrivateKey {
	signBytes, err := ioutil.ReadFile(j.PrivateKeyPath)
	if err != nil {
		log.Fatal(err)
	}
	signKey, err := jwt.ParseRSAPrivateKeyFromPEM(signBytes)
	if err != nil {
		log.Fatal(err)
	}
	return signKey
}

func (j *Jwt) GetPublicKey() *rsa.PublicKey {
	verifyBytes, err := ioutil.ReadFile(j.PublicKeyPath)
	if err != nil {
		log.Fatal(err)
	}
	verifyKey, err := jwt.ParseRSAPublicKeyFromPEM(verifyBytes)
	if err != nil {
		log.Fatal(err)
	}
	return verifyKey
}

func (j *Jwt) Validate(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {

	token, err := request.ParseFromRequest(r, request.AuthorizationHeaderExtractor,
		func(token *jwt.Token) (interface{}, error) {
			return j.GetPublicKey(), nil
		})

	if err == nil {
		if token.Valid {
			next(w, r)
		} else {
			w.WriteHeader(http.StatusUnauthorized)
			fmt.Fprint(w, "Token is not valid")
		}
	} else {
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprint(w, "Unauthorized access to this resource")
	}

}
