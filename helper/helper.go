package helper

import (
	"api/config"
	"regexp"
	"strconv"
	"strings"

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

	limit = config.ENV.Limit
	return
}

func GetID(c *gin.Context, id *int) bool {
	var err error
	if *id, err = strconv.Atoi(c.Param("id")); err != nil || *id <= 0 {
		return false
	}
	return true
}

var matchFirstCap = regexp.MustCompile("(.)([A-Z][a-z]+)")
var matchAllCap = regexp.MustCompile("([a-z0-9])([A-Z])")

func ToSnakeCase(str string) string {
	snake := matchFirstCap.ReplaceAllString(str, "${1}_${2}")
	snake = matchAllCap.ReplaceAllString(snake, "${1}_${2}")
	return strings.ToLower(snake)
}
