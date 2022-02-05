package k3s

import (
	"api/config"
	"api/helper"
	"api/interfaces"
	"api/models"
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	v1 "k8s.io/api/apps/v1"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type deploymentController struct{}

func NewDeploymentController() interfaces.Default {
	return &deploymentController{}
}

func (*deploymentController) CreateAll(c *gin.Context) {}

// @Tags K3s/Deployment
// @Summary Create Deployment
// @Accept json
// @Produce application/json
// @Produce application/xml
// @Security BearerAuth
// @Param namespace path string true "Namespace name"
// @Param model body v1.Deployment true "Deployment config file"
// @Success 201 {object} models.Success{result=[]v1.Deployment}
// @failure 400 {object} models.Error
// @failure 422 {object} models.Error
// @failure 429 {object} models.Error
// @failure 500 {object} models.Error
// @Router /k3s/deployment/{namespace} [post]
func (*deploymentController) CreateOne(c *gin.Context) {
	var namespace = c.Param("namespace")
	if namespace == "" {
		helper.ErrHandler(c, http.StatusBadRequest, "Namespace name shouldn't be empty")
		return
	}

	var body v1.Deployment
	if err := c.ShouldBind(&body); err != nil {
		helper.ErrHandler(c, http.StatusBadRequest, "Incorrect body is not setted")
		return
	}

	ctx := context.Background()
	result, err := config.K3s.AppsV1().Deployments(namespace).Create(ctx, &body, metaV1.CreateOptions{})
	if err != nil {
		helper.ErrHandler(c, http.StatusInternalServerError, err.Error())
		return
	}

	helper.ResHandler(c, http.StatusCreated, models.Success{
		Status: "OK",
		Result: &[1]v1.Deployment{*result},
		Items:  1,
	})
}

// @Tags K3s
// @Summary Get Deployments
// @Accept json
// @Produce application/json
// @Produce application/xml
// @Security BearerAuth
// @Param name path string true "Specified name of deployment"
// @Param namespace query string false "Namespace name"
// @Success 200 {object} models.Success{result=[]v1.Deployment}
// @failure 400 {object} models.Error
// @failure 422 {object} models.Error
// @failure 429 {object} models.Error
// @failure 500 {object} models.Error
// @Router /k3s/deployment/{name} [get]
func (*deploymentController) ReadOne(c *gin.Context) {
	var name string

	if name = c.Param("name"); name == "" {
		helper.ErrHandler(c, http.StatusBadRequest, "Name shouldn't be empty")
		return
	}

	ctx := context.Background()
	result, err := config.K3s.AppsV1().Deployments(c.DefaultQuery("namespace", metaV1.NamespaceDefault)).Get(ctx, name, metaV1.GetOptions{})

	if err != nil {
		helper.ErrHandler(c, http.StatusInternalServerError, err.Error())
		return
	}

	helper.ResHandler(c, http.StatusOK, models.Success{
		Status: "OK",
		Result: &[1]v1.Deployment{*result},
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
// @Success 200 {object} models.Success{result=[]v1.Deployment}
// @failure 422 {object} models.Error
// @failure 429 {object} models.Error
// @failure 500 {object} models.Error
// @Router /k3s/deployment [get]
func (*deploymentController) ReadAll(c *gin.Context) {
	ctx := context.Background()
	result, err := config.K3s.AppsV1().Deployments(c.DefaultQuery("namespace", metaV1.NamespaceAll)).List(ctx, metaV1.ListOptions{})

	if err != nil {
		helper.ErrHandler(c, http.StatusInternalServerError, err.Error())
		return
	}

	helper.ResHandler(c, http.StatusOK, models.Success{
		Status: "OK",
		Result: &result.Items,
		Items:  int64(len(result.Items)),
	})
}

func (*deploymentController) UpdateAll(c *gin.Context) {}
func (*deploymentController) UpdateOne(c *gin.Context) {}
func (*deploymentController) DeleteAll(c *gin.Context) {}
func (*deploymentController) DeleteOne(c *gin.Context) {}
