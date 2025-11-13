package middleware

import (
	"net/http"
	"strings"

	"github.com/MicahParks/keyfunc"
	"github.com/golang-jwt/jwt/v4"
	"github.com/labstack/echo/v4"
)

var jwks *keyfunc.JWKS

func InitJWKS(authURL string) error {
	jwksURL := authURL + "/api/auth/jwks"

	var err error
	jwks, err = keyfunc.Get(jwksURL, keyfunc.Options{
		RefreshUnknownKID: true,
	})

	return err
}

func JWTMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		header := c.Request().Header.Get("Authorization")
		if header == "" {
			return echo.NewHTTPError(http.StatusUnauthorized, "missing Authorization header")
		}

		parts := strings.Split(header, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			return echo.NewHTTPError(http.StatusUnauthorized, "invalid Authorization header")
		}

		tokenStr := parts[1]

		// Do not force RS256 here: Better Auth defaults to EdDSA (Ed25519)
		token, err := jwt.Parse(tokenStr, jwks.Keyfunc)
		if err != nil || !token.Valid {
			return echo.NewHTTPError(http.StatusUnauthorized, "invalid token")
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			return echo.NewHTTPError(http.StatusUnauthorized, "invalid token claims")
		}

		// store claims in echo context
		c.Set("claims", claims)

		return next(c)
	}
}
