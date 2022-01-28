package k3s

import (
	"api/config"
	"api/helper"
	"api/logs"
	"api/models"
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	v1 "k8s.io/api/core/v1"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/tools/remotecommand"
	metricsV1 "k8s.io/metrics/pkg/apis/metrics/v1beta1"
)

type PodsController struct{}

// @Tags K3s
// @Summary Exec command inside Pod
// @Accept json
// @Produce application/json
// @Produce application/xml
// @Security BearerAuth
// @Param namespace path string true "Namespace name"
// @Param model body string true "Deployment config file"
// @Success 201 {object} models.Success{result=[]v1.Service}
// @failure 400 {object} models.Error
// @failure 422 {object} models.Error
// @failure 429 {object} models.Error
// @failure 500 {object} models.Error
// @Router /k3s/service/{name} [post]
func (*PodsController) Exec(c *gin.Context) {
	var name string
	var namespace string

	name = c.Param("name")
	if name == "" {
		helper.ErrHandler(c, http.StatusBadRequest, "Namespace name shouldn't be empty")
		return
	}

	if namespace = c.DefaultQuery("namespace", ""); namespace == "" {
		helper.ErrHandler(c, http.StatusBadRequest, "Namespace shouldn't be empty")
		return
	}

	cmd, err := c.GetRawData()
	if err != nil {
		helper.ErrHandler(c, http.StatusBadRequest, err.Error())
		return
	}

	req := config.K3s.CoreV1().RESTClient().Post().Namespace(namespace).Resource("pods").Name(name).SubResource("exec").VersionedParams(&v1.PodExecOptions{
		Command: []string{"sh", "-c", string(cmd)},
		Stdout:  true,
		Stderr:  true,
		TTY:     true,
	}, scheme.ParameterCodec)

	exec, err := remotecommand.NewSPDYExecutor(config.K3sConfig, "POST", req.URL())
	if err != nil {
		helper.ErrHandler(c, http.StatusInternalServerError, err.Error())
		return
	}

	outWriter := helper.StreamWriter{}
	errWriter := helper.StreamWriter{}
	err = exec.Stream(remotecommand.StreamOptions{
		Stdout: &outWriter,
		Stderr: &errWriter,
	})

	if err != nil {
		helper.ErrHandler(c, http.StatusInternalServerError, err.Error())
		return
	}

	if len(errWriter.Result) != 0 {
		logs.DefaultLog("containers/k3s/pods", string(errWriter.Result))
		helper.ErrHandler(c, http.StatusInternalServerError, string(errWriter.Result))
		return
	}

	helper.ResHandler(c, http.StatusCreated, models.Success{
		Status: "OK",
		Result: string(outWriter.Result),
	})
}

// @Tags K3s
// @Summary Get Pod
// @Accept json
// @Produce application/json
// @Produce application/xml
// @Security BearerAuth
// @Param name path string true "Specified name of Service"
// @Param namespace query string true "Namespace name"
// @Success 200 {object} models.Success{result=[]v1.Service}
// @failure 400 {object} models.Error
// @failure 422 {object} models.Error
// @failure 429 {object} models.Error
// @failure 500 {object} models.Error
// @Router /k3s/service/{name} [get]
func (*PodsController) ReadOne(c *gin.Context) {
	var name string
	var namespace string

	if name = c.Param("name"); name == "" {
		helper.ErrHandler(c, http.StatusBadRequest, "Name shouldn't be empty")
		return
	}

	if namespace = c.DefaultQuery("namespace", ""); namespace == "" {
		helper.ErrHandler(c, http.StatusBadRequest, "Namespace shouldn't be empty")
		return
	}

	ctx := context.Background()
	result, err := config.K3s.CoreV1().Services(namespace).Get(ctx, name, metaV1.GetOptions{})

	if err != nil {
		helper.ErrHandler(c, http.StatusInternalServerError, err.Error())
		return
	}

	helper.ResHandler(c, http.StatusOK, models.Success{
		Status: "OK",
		Result: &[1]v1.Service{*result},
		Items:  1,
	})
}

// @Tags K3s
// @Summary Get Pods
// @Accept json
// @Produce application/json
// @Produce application/xml
// @Security BearerAuth
// @Param namespace query string false "Namespace name"
// @Success 200 {object} models.Success{result=[]v1.Pod}
// @failure 422 {object} models.Error
// @failure 429 {object} models.Error
// @failure 500 {object} models.Error
// @Router /k3s/pod [get]
func (*PodsController) ReadAll(c *gin.Context) {
	ctx := context.Background()
	result, err := config.K3s.CoreV1().Pods(c.DefaultQuery("namespace", metaV1.NamespaceAll)).List(ctx, metaV1.ListOptions{})

	if err != nil {
		helper.ErrHandler(c, http.StatusInternalServerError, err.Error())
		return
	}

	// helper.ResHandler(c, http.StatusOK, models.Success{
	// 	Status: "OK",
	// 	Result: result.Items,
	// 	Items:  int64(len(result.Items)),
	// })

	result2, err2 := config.Metrics.MetricsV1beta1().PodMetricses(c.DefaultQuery("namespace", metaV1.NamespaceAll)).List(ctx, metaV1.ListOptions{})
	if err2 != nil {
		helper.ErrHandler(c, http.StatusInternalServerError, err.Error())
		return
	}

	helper.ResHandler(c, http.StatusOK, models.Success{
		Status: "OK",
		Result: &result2,
		Items:  int64(len(result.Items)),
	})
}

// @Tags K3s
// @Summary Get Pod Metrics
// @Accept json
// @Produce application/json
// @Produce application/xml
// @Security BearerAuth
// @Param name path string true "Specified name of Service"
// @Param namespace query string false "Namespace name"
// @Success 200 {object} models.Success{result=[]v1.Pod}
// @failure 422 {object} models.Error
// @failure 429 {object} models.Error
// @failure 500 {object} models.Error
// @Router /k3s/pod/metrics/{name} [get]
func (*PodsController) ReadMetricsOne(c *gin.Context) {
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
		Result: &[1]metricsV1.PodMetrics{*result},
		Items:  1,
	})
}

// @Tags K3s
// @Summary Get Pods Metrics
// @Accept json
// @Produce application/json
// @Produce application/xml
// @Security BearerAuth
// @Param namespace query string false "Namespace name"
// @Success 200 {object} models.Success{result=[]v1.Pod}
// @failure 422 {object} models.Error
// @failure 429 {object} models.Error
// @failure 500 {object} models.Error
// @Router /k3s/pod/metrics [get]
func (*PodsController) ReadMetricsAll(c *gin.Context) {
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
