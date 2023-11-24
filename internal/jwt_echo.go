package internal

import (
	"net/http"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
)

// jwtCustomClaims are custom claims extending default ones.
type JwtCustomClaims struct {
	Name  string `json:"name"`
	Admin bool   `json:"admin"`
	jwt.StandardClaims
}

func Login(c echo.Context) error {
	username := c.FormValue("username")
	password := c.FormValue("password")

	if username == "infer" && password == "!@#$1234" {

		// Set custom claims
		claims := &JwtCustomClaims{
			username,
			true,
			jwt.StandardClaims{
				ExpiresAt: time.Now().Add(time.Hour * 72).Unix(),
			},
		}

		// Create token with claims
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

		// Generate encoded token and send it as response.
		t, err := token.SignedString([]byte("secret"))
		if err != nil {
			return err
		}
		return c.JSON(http.StatusOK, echo.Map{
			"token": t,
		})
	}

	return echo.ErrUnauthorized
}
