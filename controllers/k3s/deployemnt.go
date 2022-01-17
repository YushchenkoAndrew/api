package k3s

import (
	"api/config"
	"api/helper"
	"api/models"
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	appsV1 "k8s.io/api/apps/v1"
	coreV1 "k8s.io/api/core/v1"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type DeploymentController struct{}

// @Tags K3s
// @Summary Create
// @Accept json
// @Produce application/json
// @Produce application/xml
// @Security BearerAuth
// @Param id path int true "Project primaray id"
// @Param model body coreV1.Namespace true "Link info"
// @Success 201 {object} models.Success{result=[]models.Link}
// @failure 400 {object} models.Error
// @failure 422 {object} models.Error
// @failure 429 {object} models.Error
// @failure 500 {object} models.Error
// @Router /link/{id} [post]
func (*DeploymentController) Create(c *gin.Context) {
	var body coreV1.Namespace
	if err := c.ShouldBind(&body); err != nil {
		helper.ErrHandler(c, http.StatusBadRequest, "Incorrect body is not setted")
		return
	}

	ctx := context.Background()
	result, err := config.K3s.CoreV1().Namespaces().Create(ctx, &body, metaV1.CreateOptions{})
	if err != nil {
		helper.ErrHandler(c, http.StatusInternalServerError, err.Error())
	}

	helper.ResHandler(c, http.StatusCreated, models.Success{
		Status: "OK",
		Result: [1]coreV1.Namespace{*result},
		Items:  1,
	})
}

// @Tags K3s
// @Summary Get Deployments
// @Accept json
// @Produce application/json
// @Produce application/xml
// @Security BearerAuth
// @Param name path string false "Specified name of deployment"
// @Param namespace query string false "Namespace name"
// @Success 200 {object} models.Success{result=[]appsV1.Deployment}
// @failure 400 {object} models.Error
// @failure 422 {object} models.Error
// @failure 429 {object} models.Error
// @failure 500 {object} models.Error
// @Router /k3s/deployment/{name} [get]
func (*DeploymentController) ReadOne(c *gin.Context) {
	var name = c.Param("name")
	if name == "" {
		helper.ErrHandler(c, http.StatusBadRequest, "Name shouldn't be empty")
		return
	}

	ctx := context.Background()
	result, err := config.K3s.AppsV1().Deployments(c.DefaultQuery("namespace", "")).Get(ctx, name, metaV1.GetOptions{})

	if err != nil {
		helper.ErrHandler(c, http.StatusInternalServerError, err.Error())
		return
	}

	helper.ResHandler(c, http.StatusOK, models.Success{
		Status: "OK",
		Result: [1]appsV1.Deployment{*result},
		Items:  1,
	})
}

// @Tags K3s
// @Summary Get Deployments
// @Accept json
// @Produce application/json
// @Produce application/xml
// @Security BearerAuth
// @Param namespace query string false "Namespace name"
// @Success 200 {object} models.Success{result=[]appsV1.Deployment}
// @failure 422 {object} models.Error
// @failure 429 {object} models.Error
// @failure 500 {object} models.Error
// @Router /k3s/deployment [get]
func (*DeploymentController) ReadAll(c *gin.Context) {
	ctx := context.Background()
	result, err := config.K3s.AppsV1().Deployments(c.DefaultQuery("namespace", "")).List(ctx, metaV1.ListOptions{})

	if err != nil {
		helper.ErrHandler(c, http.StatusInternalServerError, err.Error())
		return
	}

	helper.ResHandler(c, http.StatusOK, models.Success{
		Status: "OK",
		Result: result.Items,
		Items:  int64(len(result.Items)),
	})
}
