package middleware

import (
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

type authError struct {
	Message string `json:"message" example:"error" `
}

func Authorize(c *gin.Context) {

	token := c.Request.Header.Get("Authorization")
	if token == "" {
		c.AbortWithStatusJSON(401, authError{
			Message: "Unauthorized",
		})
		return
	}

	if len(strings.Split(token, " ")) != 2 {
		c.AbortWithStatusJSON(401, authError{
			Message: "Unauthorized",
		})
		return
	}

	token = strings.Split(token, " ")[1]

	claims, err := jwt.ParseWithClaims(token, jwt.MapClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_SECRET")), nil
	})

	if err != nil {
		c.AbortWithStatusJSON(401, authError{
			Message: "Unauthorized",
		})
		return
	}

	if !claims.Valid {
		c.AbortWithStatusJSON(401, authError{
			Message: "Unauthorized",
		})
		return
	}

	if claims.Claims.(jwt.MapClaims)["deviceId"] != nil {
		c.AbortWithStatusJSON(401, authError{
			Message: "Unauthorized",
		})
		return
	}

	if claims.Claims.(jwt.MapClaims)["adminId"] == nil &&
		claims.Claims.(jwt.MapClaims)["phone"] == nil &&
		claims.Claims.(jwt.MapClaims)["username"] == nil &&
		claims.Claims.(jwt.MapClaims)["email"] == nil &&
		claims.Claims.(jwt.MapClaims)["role"] == nil {
		c.AbortWithStatusJSON(401, authError{
			Message: "Unauthorized",
		})
		return
	}

	c.Set("adminId", claims.Claims.(jwt.MapClaims)["adminId"])
	c.Set("phone", claims.Claims.(jwt.MapClaims)["phone"])
	c.Set("username", claims.Claims.(jwt.MapClaims)["username"])
	c.Set("email", claims.Claims.(jwt.MapClaims)["email"])
	c.Set("role", claims.Claims.(jwt.MapClaims)["role"])

	c.Next()
}

func UserAuthorize(c *gin.Context) {

	token := c.Request.Header.Get("Authorization")
	if token == "" {
		c.AbortWithStatusJSON(401, authError{
			Message: "Unauthorized",
		})
		return
	}

	if len(strings.Split(token, " ")) != 2 {
		c.AbortWithStatusJSON(401, authError{
			Message: "Unauthorized",
		})
		return
	}

	token = strings.Split(token, " ")[1]

	claims, err := jwt.ParseWithClaims(token, jwt.MapClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_SECRET_USER")), nil
	})

	if err != nil {
		c.AbortWithStatusJSON(401, authError{
			Message: "Unauthorized",
		})
		return
	}

	if !claims.Valid {
		c.AbortWithStatusJSON(401, authError{
			Message: "Unauthorized",
		})
		return
	}

	if claims.Claims.(jwt.MapClaims)["userId"] == nil &&
		claims.Claims.(jwt.MapClaims)["phone"] == nil &&
		claims.Claims.(jwt.MapClaims)["username"] == nil {
		c.AbortWithStatusJSON(401, authError{
			Message: "Unauthorized",
		})
		return
	}

	c.Set("userId", claims.Claims.(jwt.MapClaims)["userId"])
	c.Set("phone", claims.Claims.(jwt.MapClaims)["phone"])
	c.Set("username", claims.Claims.(jwt.MapClaims)["username"])

	c.Next()
}
