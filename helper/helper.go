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
