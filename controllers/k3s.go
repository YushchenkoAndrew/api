package controllers

import (
	"api/interfaces"

	"github.com/gin-gonic/gin"
)

type k3sController struct{}

func NewK3sController() interfaces.K3s {
	return &k3sController{}
}

// @Tags K3s
// @Summary Create
// @Accept json
// @Produce application/json
// @Produce application/xml
// @Security BearerAuth
// @Param id path int true "Project primaray id"
// @Param model body models.ReqLink true "Link info"
// @Success 201 {object} models.Success{result=[]models.Link}
// @failure 400 {object} models.Error
// @failure 401 {object} models.Error
// @failure 422 {object} models.Error
// @failure 429 {object} models.Error
// @failure 500 {object} models.Error
// @Router /link/{id} [post]
func (*k3sController) Subscribe(c *gin.Context) {

}

// @Tags K3s
// @Summary Create
// @Accept json
// @Produce application/json
// @Produce application/xml
// @Security BearerAuth
// @Param id path int true "Project primaray id"
// @Param model body models.ReqLink true "Link info"
// @Success 201 {object} models.Success{result=[]models.Link}
// @failure 400 {object} models.Error
// @failure 401 {object} models.Error
// @failure 422 {object} models.Error
// @failure 429 {object} models.Error
// @failure 500 {object} models.Error
// @Router /link/{id} [post]
func (*k3sController) Unsubscribe(c *gin.Context) {

}
