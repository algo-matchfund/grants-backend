package middlewares

import (
	"crypto/x509"
	"encoding/pem"
	"errors"
	"log"
	"strings"
	"time"

	"github.com/algo-matchfund/grants-backend/gen/models"
	"gopkg.in/square/go-jose.v2/jwt"
)

// NewAuthenticator returns a function that verifies jwt tokens using the provided public key
func NewAuthenticator(publicKey string) (func(string) (*models.Principal, error), error) {
	pemBlock, _ := pem.Decode([]byte(publicKey))
	if pemBlock == nil {
		return nil, errors.New("unable to parse public key")
	}

	k, err := x509.ParsePKIXPublicKey(pemBlock.Bytes)
	if err != nil {
		log.Println("Error parsing keycloak public key")
		log.Println(err)
		return nil, err
	}

	return func(bearerHeader string) (*models.Principal, error) {
		token := strings.Split(bearerHeader, " ")[1]
		tok, err := jwt.ParseSigned(token)
		if err != nil {
			return nil, err
		}

		claims := &struct {
			jwt.Claims
			models.Principal
			Type                       string `json:"typ,omitempty"`
			AuthorizedParty            string `json:"azp,omitempty"`
			AuthenticationContextClass string `json:"acr,omitempty"`
		}{}

		err = tok.Claims(k, &claims)

		if err != nil {
			log.Println("Could not verify JWT signature.")
			log.Println(err)
			return nil, err
		}

		if time.Now().After(claims.Expiry.Time()) {
			return nil, jwt.ErrExpired
		}

		if time.Now().Before(claims.NotBefore.Time()) {
			return nil, jwt.ErrNotValidYet
		}

		claims.Principal.ID = claims.Subject

		return &claims.Principal, nil
	}, nil
}
