package controllers

import (
	"api/db"
	"api/helper"
	"api/models"
	"context"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type BotController struct{}

// @Tags Bot
// @Summary Execute redis Command from request
// @Accept json
// @Produce application/json
// @Produce application/xml
// @Param model body models.BotRedis true "Redis Command"
// @Success 200 {object} models.DefultRes
// @failure 400 {object} models.Error
// @failure 429 {object} models.Error
// @failure 500 {object} models.Error
// @Router /bot/redis [post]
func (*BotController) Redis(c *gin.Context) {
	var body models.BotRedis
	if err := c.ShouldBind(&body); err != nil {
		helper.ErrHandler(c, http.StatusBadRequest, "Incorrect body")
		return
	}

	words := strings.Split(body.Command, " ")
	var command = make([]interface{}, len(words))
	for i := 0; i < len(words); i++ {
		command[i] = words[i]
	}

	ctx := context.Background()
	res, err := db.Redis.Do(ctx, command...).StringSlice()

	if err != nil {
		c.JSON(http.StatusOK, models.DefultRes{
			Status:  "ERR",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.DefultRes{
		Status:  "ERR",
		Message: "Success",
		Result:  res,
	})
}
