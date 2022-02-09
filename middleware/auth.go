package middleware

import (
	"api/config"
	"api/db"
	"api/helper"
	"api/logs"
	"api/models"
	"context"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
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

	access := jwt.NewWithClaims(jwt.SigningMethodHS256, &jwt.MapClaims{
		"authorized":  true,
		"access_uuid": token.AccessUUID,
		"user_id":     config.ENV.ID,
		"expire":      token.AccessExpire,
	})

	token.AccessToken, err = access.SignedString([]byte(config.ENV.AccessSecret))
	if err != nil {
		return
	}

	refresh := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"refresh_uuid": token.RefreshUUID,
		"user_id":      config.ENV.ID,
		"expire":       token.RefreshExpire,
	})

	token.RefreshToken, err = refresh.SignedString([]byte(config.ENV.RefreshSecret))
	return
}

func Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		bearToken := strings.Split(c.Request.Header.Get("Authorization"), " ")

		if len(bearToken) != 2 {
			helper.ErrHandler(c, http.StatusUnauthorized, "Invalid token inforamation")
			go logs.SendLogs(&models.LogMessage{
				Stat:    "ERR",
				Name:    "API",
				Url:     "/api/refresh",
				File:    "/controllers/index.go",
				Message: "It's mine first time so please be gentle",
			})
			return
		}

		token, err := jwt.Parse(bearToken[1], func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf(("invalid signing method"))
			}

			if _, ok := t.Claims.(jwt.Claims); !ok && !t.Valid {
				return nil, fmt.Errorf(("expired token"))
			}

			return []byte(config.ENV.AccessSecret), nil
		})

		if err != nil {
			helper.ErrHandler(c, http.StatusUnauthorized, err.Error())
			go logs.SendLogs(&models.LogMessage{
				Stat:    "ERR",
				Name:    "API",
				Url:     "/api/refresh",
				File:    "/controllers/index.go",
				Message: "It's mine first time so please be gentle",
			})
			return
		}

		var userUUID string
		var accessUUID string

		claims, ok := token.Claims.(jwt.MapClaims)
		if ok && token.Valid {
			var ok bool
			if accessUUID, ok = claims["access_uuid"].(string); !ok {
				helper.ErrHandler(c, http.StatusUnprocessableEntity, "Invalid token inforamation")
				return
			}

			if userUUID, ok = claims["user_id"].(string); !ok {
				helper.ErrHandler(c, http.StatusUnprocessableEntity, "Invalid token inforamation")
				return
			}
		}

		// Final check with cache
		ctx := context.Background()
		if cacheUUID, err := db.Redis.Get(ctx, accessUUID).Result(); err != nil || cacheUUID != userUUID {
			helper.ErrHandler(c, http.StatusUnauthorized, "Invalid token inforamation")
			return
		}

		// before request

		c.Next()

		// after request

	}
}

func AuthToken() gin.HandlerFunc {
	return func(c *gin.Context) {
		bearToken := strings.Split(c.Request.Header.Get("Authorization"), " ")

		if len(bearToken) != 2 {
			helper.ErrHandler(c, http.StatusUnauthorized, "Invalid token inforamation")
			go logs.SendLogs(&models.LogMessage{
				Stat:    "ERR",
				Name:    "API",
				File:    "/middleware/auth.go",
				Message: "Ohhh nyo your token is inccorrect",
			})
			return
		}

		// NOTE: Simple way for solving redis injection vulnerabilities
		hasher := md5.New()
		hasher.Write([]byte(bearToken[1]))
		token := hex.EncodeToString(hasher.Sum(nil))

		ctx := context.Background()
		if token, err := db.Redis.Get(ctx, "TOKEN:"+token).Result(); err != nil || token != "OK" {
			helper.ErrHandler(c, http.StatusUnauthorized, err.Error())
			go logs.SendLogs(&models.LogMessage{
				Stat:    "ERR",
				Name:    "API",
				File:    "/middleware/auth.go",
				Message: "Ohhh nyo your token is inccorrect",
			})
			return
		}

		// before request

		c.Next()

		// after request

		// Manially refresh token after each use
		db.Redis.Del(ctx, "TOKEN:"+token)
		salt := strconv.Itoa(rand.Intn(1000000) + 5000)

		go helper.RegenerateToken(salt + bearToken[1])
		helper.ResHandler(c, http.StatusOK, salt)
	}
}
