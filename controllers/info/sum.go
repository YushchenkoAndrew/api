package info

import (
	"api/db"
	"api/helper"
	"api/interfaces/info"
	"api/logs"
	"api/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

type sumController struct{}

func NewSumController() info.Default {
	return &sumController{}
}

// @Tags Info
// @Summary Get Info Sum
// @Accept json
// @Produce application/json
// @Produce application/xml
// @Success 200 {object} models.Success{result=[]models.StatInfo}
// @failure 429 {object} models.Error
// @failure 500 {object} models.Error
// @Router /info/sum [get]
func (o *sumController) Read(c *gin.Context) {
	var model models.StatInfo

	if err := helper.PrecacheResult("INFO:SUM", db.DB.Table("info").Select("SUM(views) as views, SUM(clicks) AS clicks, SUM(media) as media, SUM(visitors) as visitors"), model); err != nil {
		helper.ErrHandler(c, http.StatusInternalServerError, err.Error())
		go logs.DefaultLog("/controllers/info/sum.go", err.Error())
		return
	}

	// items, err := db.Redis.Get(context.Background(), "nINFO").Int64()
	// if err != nil {
	// 	items = -1
	// 	// go (&models.Info{}).Redis(db.DB, db.Redis)
	// 	go logs.SendLogs(&models.LogMessage{
	// 		Stat:    "ERR",
	// 		Name:    "API",
	// 		File:    "/controllers/info/sum.go",
	// 		Message: "Ohh nooo Cache is broken; Anyway...",
	// 		Desc:    err.Error(),
	// 	})
	// }

	helper.ResHandler(c, http.StatusOK, &models.Success{
		Status: "OK",
		Result: []models.StatInfo{model},
		Items:  1,
		// TotalItems: items,
	})
}
