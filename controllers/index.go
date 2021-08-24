package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type IndexController struct{}

func (*IndexController) Ping(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "pong",
	})
}

func (*IndexController) Navigation(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"current_url":                 "/api",
		"authorization_url":           "/api/login",
		"check_token_url":             "/api/token",
		"delete_element_by_id_url":    "/api/{table}/:id",
		"delete_element_by_param_url": "/api/{table}?{query}",
		"documentation_url":           "/api/doc",
		"get_table_list_url":          "/api/tables",
		"get_elements_url":            "/api/{table}",
		"get_element_by_id_url":       "/api/{table}/:id",
		"get_element_by_param_url":    "/api/{table}?{query}",
		"ping_url":                    "/api/ping",
		"post_element_url":            "/api/{table}",
		"put_element_by_id_url":       "/api/{table}/:id",
		"put_element_by_param_url":    "/api/{table}?{query}",
	})
}

func (*IndexController) TableList(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"tables": [2]string{
			"Info",
			"World"},
	})
}
