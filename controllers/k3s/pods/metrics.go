package pods

import (
	"api/config"
	"api/helper"
	"api/interfaces/k3s"
	"api/models"
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	v1 "k8s.io/metrics/pkg/apis/metrics/v1beta1"
)

type metricsController struct{}

func NewMetricsController() k3s.Metrics {
	return &metricsController{}
}

// @Tags Metrics
// @Summary Get Pod Metrics
// @Accept json
// @Produce application/json
// @Produce application/xml
// @Security BearerAuth
// @Param name path string true "Specified name of Service"
// @Param namespace query string false "Namespace name"
// @Success 200 {object} models.Success{result=[]v1.PodMetrics}
// @failure 422 {object} models.Error
// @failure 429 {object} models.Error
// @failure 500 {object} models.Error
// @Router /k3s/pod/metrics/{name} [get]
func (*metricsController) ReadOne(c *gin.Context) {
	var name string
	if name = c.Param("name"); name == "" {
		helper.ErrHandler(c, http.StatusBadRequest, "Name shouldn't be empty")
		return
	}

	ctx := context.Background()
	result, err := config.Metrics.MetricsV1beta1().PodMetricses(c.DefaultQuery("namespace", metaV1.NamespaceDefault)).Get(ctx, name, metaV1.GetOptions{})
	if err != nil {
		helper.ErrHandler(c, http.StatusInternalServerError, err.Error())
		return
	}

	helper.ResHandler(c, http.StatusOK, models.Success{
		Status: "OK",
		Result: &[1]v1.PodMetrics{*result},
		Items:  1,
	})
}

// @Tags Metrics
// @Summary Get Pods Metrics
// @Accept json
// @Produce application/json
// @Produce application/xml
// @Security BearerAuth
// @Param namespace query string false "Namespace name"
// @Success 200 {object} models.Success{result=[]v1.PodMetrics}
// @failure 422 {object} models.Error
// @failure 429 {object} models.Error
// @failure 500 {object} models.Error
// @Router /k3s/pod/metrics [get]
func (*metricsController) ReadAll(c *gin.Context) {
	ctx := context.Background()
	result, err := config.Metrics.MetricsV1beta1().PodMetricses(c.DefaultQuery("namespace", metaV1.NamespaceAll)).List(ctx, metaV1.ListOptions{})
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
