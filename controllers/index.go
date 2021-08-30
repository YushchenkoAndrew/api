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
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type IndexController struct{}

// @Summary Ping/Pong
// @Accept json
// @Produce application/json
// @Success 200 {object} models.Ping
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
	db.Redis.Set(ctx, token.RefreshToken, config.ENV.ID, time.Duration((token.RefreshExpire-now)*int64(time.Second)))
	helper.ResHandler(c, http.StatusOK, models.Tokens{
		AccessToken:  token.AccessToken,
		RefreshToken: token.RefreshToken,
	})
}

// @Summary Refresh
// @Accept json
// @Produce application/json
// @Produce application/xml
// @Success 200 {object} models.Tokens
// @failure 500 {object} models.Error
// @Router /refresh [post]
func (*IndexController) Refresh(c *gin.Context) {
	// TODO:
}
