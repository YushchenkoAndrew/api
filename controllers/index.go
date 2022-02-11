package controllers

import (
	"api/config"
	"api/db"
	"api/helper"
	"api/interfaces"
	"api/logs"
	"api/middleware"
	"api/models"
	"context"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

type indexController struct{}

func NewIndexController() interfaces.Index {
	return &indexController{}
}

// @Summary Ping/Pong
// @Accept json
// @Produce application/json
// @Success 200 {object} models.Ping
// @failure 429 {object} models.Error
// @Router /ping [get]
func (*indexController) Ping(c *gin.Context) {
	c.JSON(http.StatusOK, models.Ping{
		Status:  "OK",
		Message: "pong",
	})
}

// @Summary Trace Ip :ip
// @Accept json
// @Produce application/json
// @Produce application/xml
// @Param ip path string true "Client IP"
// @Success 200 {object} models.Success{result=[]models.GeoIpLocations}
// @failure 429 {object} models.Error
// @failure 400 {object} models.Error
// @failure 500 {object} models.Error
// @Router /trace/{ip} [get]
func (*indexController) TraceIp(c *gin.Context) {
	var ip string
	if ip = c.Param("ip"); ip == "" {
		helper.ErrHandler(c, http.StatusBadRequest, "Incorrect ip value")
		return
	}

	key := "ip:" + ip
	ctx := context.Background()
	var location []models.GeoIpLocations

	// Check if cache have requested data
	if data, err := db.Redis.Get(ctx, key).Result(); err == nil {
		json.Unmarshal([]byte(data), &location)
		go db.Redis.Expire(ctx, key, time.Duration(config.ENV.LiveTime)*time.Second)
	} else {

		var block []models.GeoIpBlocks
		if result := db.DB.Where("network >>= ?::inet", ip).Find(&block); result.Error != nil || len(block) == 0 {
			helper.ErrHandler(c, http.StatusInternalServerError, "Server side error: Something went wrong")
			go logs.DefaultLog("/controllers/index.go", result.Error)
			return
		}

		if result := db.DB.Where("geoname_id = ?", block[0].GeonameId).Find(&location); result.Error != nil {
			helper.ErrHandler(c, http.StatusInternalServerError, "Server side error: Something went wrong")
			go logs.DefaultLog("/controllers/index.go", result.Error)
			return
		}

		// Encode json to str
		if str, err := json.Marshal(&location); err == nil {
			go db.Redis.Set(ctx, key, str, time.Duration(config.ENV.LiveTime)*time.Second)
		}
	}

	helper.ResHandler(c, http.StatusOK, &models.Success{
		Status: "OK",
		Result: &location,
		Items:  int64(len(location)),
	})
}

// @Summary Login
// @Accept json
// @Produce application/json
// @Produce application/xml
// @Param model body models.LoginDto true "Login info"
// @Success 200 {object} models.TokenEntity
// @failure 400 {object} models.Error
// @failure 401 {object} models.Error
// @failure 429 {object} models.Error
// @failure 500 {object} models.Error
// @Router /login [post]
func (*indexController) Login(c *gin.Context) {
	var login models.LoginDto
	if err := c.ShouldBind(&login); err != nil {
		helper.ErrHandler(c, http.StatusBadRequest, "Incorrect body")
		return
	}

	pass := strings.Split(login.Pass, "$")

	hasher := md5.New()
	hasher.Write([]byte(pass[0] + config.ENV.Pepper + config.ENV.Pass))

	if len(pass) != 2 || !helper.ValidateStr(login.User, config.ENV.User) ||
		!helper.ValidateStr(hex.EncodeToString(hasher.Sum(nil)), pass[1]) {
		helper.ErrHandler(c, http.StatusUnauthorized, "Invalid login inforamation")
		go logs.SendLogs(&models.LogMessage{
			Stat:    "ERR",
			Name:    "API",
			Url:     "/api/refresh",
			File:    "/controllers/index.go",
			Message: "First rule of User Validation; It's not to talk about Users",
		})
		return
	}

	var token models.Auth
	if err := middleware.CreateToken(&token); err != nil {
		helper.ErrHandler(c, http.StatusInternalServerError, "Server Side error: Something went wrong")
		go logs.DefaultLog("/controllers/index.go", err.Error())
		return
	}

	now := time.Now().Unix()
	ctx := context.Background()
	db.Redis.Set(ctx, token.AccessUUID, config.ENV.ID, time.Duration((token.AccessExpire-now)*int64(time.Second)))
	db.Redis.Set(ctx, token.RefreshUUID, config.ENV.ID, time.Duration((token.RefreshExpire-now)*int64(time.Second)))
	helper.ResHandler(c, http.StatusOK, models.TokenEntity{
		Status:       "OK",
		AccessToken:  token.AccessToken,
		RefreshToken: token.RefreshToken,
	})
}

// @Summary Refresh access token
// @Accept json
// @Produce application/json
// @Produce application/xml
// @Success 200 {object} models.TokenEntity
// @failure 400 {object} models.Error
// @failure 401 {object} models.Error
// @failure 422 {object} models.Error
// @failure 429 {object} models.Error
// @failure 500 {object} models.Error
// @Router /refresh [post]
func (*indexController) Refresh(c *gin.Context) {
	var tokens models.TokenEntity
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
	if cacheUUID, err := db.Redis.Get(ctx, refreshUUID).Result(); err != nil || cacheUUID != userUUID {
		helper.ErrHandler(c, http.StatusUnauthorized, "Invalid token inforamation")
		return
	}

	go db.Redis.Del(ctx, refreshUUID)

	var t models.Auth
	if err := middleware.CreateToken(&t); err != nil {
		helper.ErrHandler(c, http.StatusInternalServerError, "Server Side error: Something went wrong")
		go logs.DefaultLog("/controllers/index.go", err.Error())
		return
	}

	now := time.Now().Unix()
	db.Redis.Set(ctx, t.AccessUUID, config.ENV.ID, time.Duration((t.AccessExpire-now)*int64(time.Second)))
	db.Redis.Set(ctx, t.RefreshUUID, config.ENV.ID, time.Duration((t.RefreshExpire-now)*int64(time.Second)))
	helper.ResHandler(c, http.StatusOK, models.TokenEntity{
		Status:       "OK",
		AccessToken:  t.AccessToken,
		RefreshToken: t.RefreshToken,
	})
}
