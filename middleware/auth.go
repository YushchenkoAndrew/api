package middleware

import (
	"api/config"
	"api/models"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/twinj/uuid"
)

// type Middleware struct{}

func CreateToken(token *models.Auth) (err error) {
	token.AccessUUID = uuid.NewV4().String()
	token.RefreshUUID = uuid.NewV4().String()

	token.AccessExpire = time.Now().Add(time.Minute * 15).Unix()
	token.RefreshExpire = time.Now().Add(time.Hour * 24 * 7).Unix()

	access := jwt.NewWithClaims(jwt.SigningMethodES256, jwt.MapClaims{
		"authorized": true,
		"user_id":    config.ENV.ID,
		"expire":     token.AccessExpire,
	})

	token.AccessToken, err = access.SignedString([]byte(config.ENV.AccessSecret))
	if err != nil {
		return
	}

	refresh := jwt.NewWithClaims(jwt.SigningMethodES256, jwt.MapClaims{
		"refresh_uuid": token.RefreshUUID,
		"user_id":      config.ENV.ID,
		"expire":       token.RefreshExpire,
	})

	token.RefreshToken, err = refresh.SignedString([]byte(config.ENV.RefreshSecret))
	return
}

func Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// TODO:

		// before request

		c.Next()

		// after request

	}
}
