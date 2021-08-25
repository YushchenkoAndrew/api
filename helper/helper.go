package helper

import (
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetStat(flags ...bool) (state uint8) {
	state = 0

	for i := 0; i < len(flags); i++ {
		if state = state << 1; flags[i] {
			state++
		}
	}

	return state
}

func Pagination(c *gin.Context) (page int, limit int) {
	var err error
	if page, err = strconv.Atoi(c.DefaultQuery("page", "0")); err != nil {
		page = 0
	}

	if limit, err = strconv.Atoi(c.DefaultQuery("limit", "20")); err != nil {
		limit = 20
	}
	return
}

func GetID(c *gin.Context, id *int) bool {
	var err error
	if *id, err = strconv.Atoi(c.Param("id")); err != nil || *id <= 0 {
		return false
	}
	return true
}

func ResHandler(c *gin.Context, stat int, data *gin.H) {
	switch c.GetHeader("Accept") {
	case "application/xml":
		c.XML(stat, *data)

	default:
		c.JSON(stat, *data)
	}
}

func ErrHandler(c *gin.Context, stat int, message string) {
	ResHandler(c, stat, &gin.H{
		"status":  "ERR",
		"result":  []string{},
		"message": message,
	})
}
