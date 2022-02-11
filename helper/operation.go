package helper

import (
	"api/config"
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
)

func FormPathFromHandler(c *gin.Context, handler *config.Handler) (string, error) {
	var path = handler.Path

	for _, param := range handler.Required {
		if value := c.DefaultQuery(param, ""); value != "" {
			path = strings.ReplaceAll(path, fmt.Sprintf("{%s}", param), value)
			continue
		}

		return "", fmt.Errorf("Reqired parameter '%s'", param)
	}

	return path, nil
}
