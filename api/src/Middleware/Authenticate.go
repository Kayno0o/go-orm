package middleware

import (
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	entity "go-api-test.kayn.ooo/src/Entity"
	repository "go-api-test.kayn.ooo/src/Repository"
)

func Authenticate(c *fiber.Ctx) error {
	tokenString := c.Get("Authorization")
	if tokenString == "" {
		tokenString = c.Cookies("token")
	}

	if tokenString == "" || len(tokenString) < 10 {
		return c.Next()
	}

	if tokenString[:7] == "Bearer " {
		tokenString = tokenString[7:]
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("SECRET_KEY")), nil
	}, jwt.WithValidMethods([]string{"HS256"}), jwt.WithAudience(os.Getenv("JWT_ISSUER")), jwt.WithIssuer(os.Getenv("JWT_ISSUER")))

	if err != nil || !token.Valid {
		c.ClearCookie("token")
		return c.Next()
	}

	claims := token.Claims.(jwt.MapClaims)

	if claims["exp"] == nil || claims["iat"] == nil {
		c.ClearCookie("token")
		return c.Next()
	}

	if int64(claims["exp"].(float64)) < time.Now().Unix() || int64(claims["iat"].(float64)) > time.Now().Unix() {
		c.ClearCookie("token")
		return c.Next()
	}

	id := claims["id"].(float64)

	user, err := repository.FindOneById[entity.User](uint(id))
	if err == nil {
		c.Locals("user", &user)
	}

	return c.Next()
}
