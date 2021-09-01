package controllers

import (
	"api/config"
	"api/db"
	"api/helper"
	"api/middleware"
	"api/models"
	"context"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

type IndexController struct{}

// @Summary Ping/Pong
// @Accept json
// @Produce application/json
// @Success 200 {object} models.Ping
// @failure 429 {object} models.Error
// @Router /ping [get]
func (*IndexController) Ping(c *gin.Context) {
	c.JSON(http.StatusOK, models.Ping{
		Message: "pong",
	})
}

// @Summary Login
// @Accept json
// @Produce application/json
// @Produce application/xml
// @Param model body models.Login true "Login info"
// @Success 200 {object} models.Tokens
// @failure 400 {object} models.Error
// @failure 401 {object} models.Error
// @failure 429 {object} models.Error
// @failure 500 {object} models.Error
// @Router /login [post]
func (*IndexController) Login(c *gin.Context) {
	var login models.Login
	if err := c.ShouldBind(&login); err != nil {
		helper.ErrHandler(c, http.StatusBadRequest, "Incorrect body")
		return
	}

	pass := strings.Split(login.Pass, "$")

	hasher := md5.New()
	hasher.Write([]byte(config.ENV.Pepper + pass[0] + config.ENV.Pass))

	if len(pass) != 2 || !helper.ValidateStr(login.User, config.ENV.User) ||
		!helper.ValidateStr(hex.EncodeToString(hasher.Sum(nil)), pass[1]) {
		helper.ErrHandler(c, http.StatusUnauthorized, "Invalid login inforamation")
		return
	}

	var token models.Auth
	if err := middleware.CreateToken(&token); err != nil {
		helper.ErrHandler(c, http.StatusInternalServerError, "Server Side error: Something went wrong")
		return
	}

	now := time.Now().Unix()
	ctx := context.Background()
	db.Redis.Set(ctx, token.AccessUUID, config.ENV.ID, time.Duration((token.AccessExpire-now)*int64(time.Second)))
	db.Redis.Set(ctx, token.RefreshUUID, config.ENV.ID, time.Duration((token.RefreshExpire-now)*int64(time.Second)))
	helper.ResHandler(c, http.StatusOK, models.Tokens{
		AccessToken:  token.AccessToken,
		RefreshToken: token.RefreshToken,
	})
}

// @Summary Refresh access token
// @Accept json
// @Produce application/json
// @Produce application/xml
// @Success 200 {object} models.Tokens
// @failure 400 {object} models.Error
// @failure 401 {object} models.Error
// @failure 422 {object} models.Error
// @failure 429 {object} models.Error
// @failure 500 {object} models.Error
// @Router /refresh [post]
func (*IndexController) Refresh(c *gin.Context) {
	var tokens models.Tokens
	if err := c.ShouldBind(&tokens); err != nil {
		helper.ErrHandler(c, http.StatusBadRequest, "Refresh token not specified")
		return
	}

	token, err := jwt.Parse(tokens.RefreshToken, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf(("invalid signing method"))
		}

		if _, ok := t.Claims.(jwt.Claims); !ok && !t.Valid {
			return nil, fmt.Errorf(("expired token"))
		}

		return []byte(config.ENV.RefreshSecret), nil
	})

	if err != nil {
		helper.ErrHandler(c, http.StatusUnauthorized, err.Error())
		return
	}

	var userUUID string
	var refreshUUID string
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		helper.ErrHandler(c, http.StatusUnauthorized, "Unauthorized token")
		return
	}

	if refreshUUID, ok = claims["refresh_uuid"].(string); !ok {
		helper.ErrHandler(c, http.StatusUnprocessableEntity, "Invalid token inforamation")
		return
	}

	if userUUID, ok = claims["user_id"].(string); !ok {
		helper.ErrHandler(c, http.StatusUnprocessableEntity, "Invalid token inforamation")
		return
	}

	ctx := context.Background()

	// Double check if such UUID exist in cache + it's the same user
	// (btw don't need it, I have only one user)
	fmt.Println(refreshUUID)
	fmt.Println(userUUID)
	if cacheUUID, err := db.Redis.Get(ctx, refreshUUID).Result(); err != nil || cacheUUID != userUUID {
		helper.ErrHandler(c, http.StatusUnauthorized, "Invalid token inforamation")
		return
	}

	go db.Redis.Del(ctx, refreshUUID)

	var t models.Auth
	if err := middleware.CreateToken(&t); err != nil {
		helper.ErrHandler(c, http.StatusInternalServerError, "Server Side error: Something went wrong")
		return
	}

	now := time.Now().Unix()
	db.Redis.Set(ctx, t.AccessUUID, config.ENV.ID, time.Duration((t.AccessExpire-now)*int64(time.Second)))
	db.Redis.Set(ctx, t.RefreshUUID, config.ENV.ID, time.Duration((t.RefreshExpire-now)*int64(time.Second)))
	helper.ResHandler(c, http.StatusOK, models.Tokens{
		AccessToken:  t.AccessToken,
		RefreshToken: t.RefreshToken,
	})
}
