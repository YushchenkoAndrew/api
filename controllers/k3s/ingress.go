package k3s

import (
	"api/config"
	"api/helper"
	"api/interfaces"
	"api/models"
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	v1 "k8s.io/api/networking/v1"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type ingressController struct{}

func NewIngressController() interfaces.Default {
	return &ingressController{}
}

func (*ingressController) CreateAll(c *gin.Context) {}

// @Tags Ingress
// @Summary Create Ingress
// @Accept json
// @Produce application/json
// @Produce application/xml
// @Security BearerAuth
// @Param namespace path string true "Namespace name"
// @Param model body v1.Ingress true "Ingress config file"
// @Success 201 {object} models.Success{result=[]v1.Ingress}
// @failure 400 {object} models.Error
// @failure 422 {object} models.Error
// @failure 429 {object} models.Error
// @failure 500 {object} models.Error
// @Router /k3s/ingress/{namespace} [post]
func (*ingressController) CreateOne(c *gin.Context) {
	var namespace = c.Param("namespace")
	if namespace == "" {
		helper.ErrHandler(c, http.StatusBadRequest, "Namespace name shouldn't be empty")
		return
	}

	var body v1.Ingress
	if err := c.ShouldBind(&body); err != nil {
		helper.ErrHandler(c, http.StatusBadRequest, "Incorrect body is not setted")
		return
	}

	ctx := context.Background()
	result, err := config.K3s.NetworkingV1().Ingresses(namespace).Create(ctx, &body, metaV1.CreateOptions{})
	if err != nil {
		helper.ErrHandler(c, http.StatusInternalServerError, err.Error())
		return
	}

	helper.ResHandler(c, http.StatusCreated, models.Success{
		Status: "OK",
		Result: &[1]v1.Ingress{*result},
		Items:  1,
	})
}

// @Tags Ingress
// @Summary Get Ingress
// @Accept json
// @Produce application/json
// @Produce application/xml
// @Security BearerAuth
// @Param name path string true "Specified name of Ingress"
// @Param namespace path string false "Namespace name"
// @Success 200 {object} models.Success{result=[]v1.Ingress}
// @failure 400 {object} models.Error
// @failure 422 {object} models.Error
// @failure 429 {object} models.Error
// @failure 500 {object} models.Error
// @Router /k3s/ingress/{namespace}/{name} [get]
func (*ingressController) ReadOne(c *gin.Context) {
	var name string
	var namespace string

	if name = c.Param("name"); name == "" {
		helper.ErrHandler(c, http.StatusBadRequest, "Name shouldn't be empty")
		return
	}

	if namespace = c.Param("namespace"); namespace == "" {
		helper.ErrHandler(c, http.StatusBadRequest, "Namespace shouldn't be empty")
		return
	}

	ctx := context.Background()
	result, err := config.K3s.NetworkingV1().Ingresses(namespace).Get(ctx, name, metaV1.GetOptions{})

	if err != nil {
		helper.ErrHandler(c, http.StatusInternalServerError, err.Error())
		return
	}

	helper.ResHandler(c, http.StatusOK, models.Success{
		Status: "OK",
		Result: &[1]v1.Ingress{*result},
		Items:  1,
	})
}

// @Tags Ingress
// @Summary Get Ingress
// @Accept json
// @Produce application/json
// @Produce application/xml
// @Security BearerAuth
// @Param namespace path string false "Namespace name"
// @Success 200 {object} models.Success{result=[]v1.Ingress}
// @failure 422 {object} models.Error
// @failure 429 {object} models.Error
// @failure 500 {object} models.Error
// @Router /k3s/ingress/{namespace} [get]
func (*ingressController) ReadAll(c *gin.Context) {
	ctx := context.Background()
	result, err := config.K3s.NetworkingV1().Ingresses(c.Param("namespace")).List(ctx, metaV1.ListOptions{})

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

func (*ingressController) UpdateAll(c *gin.Context) {}
func (*ingressController) UpdateOne(c *gin.Context) {}
func (*ingressController) DeleteAll(c *gin.Context) {}
func (*ingressController) DeleteOne(c *gin.Context) {}
