package internal

import (
	"errors"
	"infer-microservices/internal/flags"
	"net/http"
	"time"

	jwt "github.com/golang-jwt/jwt/v4"
)

var jwtKey []byte

type Claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

func init() {
	flagFactory := flags.FlagFactory{}
	flagJwt := flagFactory.CreateFlagJwt()
	jwtKey = []byte(*flagJwt.GetJwtKey())
}

func JwtAuthMiddleware(hd http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenString := r.Header.Get("Authorization")
		if tokenString == "" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		claims := &Claims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				err := errors.New("unexpected signing method")
				return nil, err
			}
			return jwtKey, nil
		})

		if err != nil {
			if err == jwt.ErrSignatureInvalid {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		if !token.Valid {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		hd.ServeHTTP(w, r)
	})
}

func AuthHandler(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")
	//password := r.FormValue("password")

	// TODO: check username and password against database

	expirationTime := time.Now().Add(30 * time.Minute)
	claims := &Claims{
		Username: username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Write([]byte(tokenString))
}
