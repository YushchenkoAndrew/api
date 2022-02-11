package k3s

import (
	"api/config"
	"api/helper"
	"api/interfaces/k3s"
	"api/logs"
	"api/models"
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	v1 "k8s.io/api/core/v1"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/tools/remotecommand"
)

type podsController struct{}

func NewPodsController() k3s.Pods {
	return &podsController{}
}

// @Tags Pods
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
// @Router /k3s/pod/{namespace}/{name} [post]
func (*podsController) Exec(c *gin.Context) {
	var name string
	var namespace string

	name = c.Param("name")
	if name == "" {
		helper.ErrHandler(c, http.StatusBadRequest, "Namespace name shouldn't be empty")
		return
	}

	if namespace = c.Param("namespace"); namespace == "" {
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

	helper.ResHandler(c, http.StatusCreated, &models.Success{
		Status: "OK",
		Result: string(outWriter.Result),
	})
}

// @Tags Pods
// @Summary Get Pod
// @Accept json
// @Produce application/json
// @Produce application/xml
// @Security BearerAuth
// @Param name path string true "Specified name of Pod"
// @Param namespace path string true "Namespace name"
// @Success 200 {object} models.Success{result=[]v1.Pod}
// @failure 400 {object} models.Error
// @failure 422 {object} models.Error
// @failure 429 {object} models.Error
// @failure 500 {object} models.Error
// @Router /k3s/pod/{namespace}/{name} [get]
func (*podsController) ReadOne(c *gin.Context) {
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
	result, err := config.K3s.CoreV1().Pods(namespace).Get(ctx, name, metaV1.GetOptions{})

	if err != nil {
		helper.ErrHandler(c, http.StatusInternalServerError, err.Error())
		return
	}

	helper.ResHandler(c, http.StatusOK, &models.Success{
		Status: "OK",
		Result: []v1.Pod{*result},
		Items:  1,
	})
}

// @Tags Pods
// @Summary Get Pods
// @Accept json
// @Produce application/json
// @Produce application/xml
// @Security BearerAuth
// @Param namespace path string false "Namespace name"
// @Success 200 {object} models.Success{result=[]v1.Pod}
// @failure 422 {object} models.Error
// @failure 429 {object} models.Error
// @failure 500 {object} models.Error
// @Router /k3s/pod/{namespace} [get]
func (*podsController) ReadAll(c *gin.Context) {
	ctx := context.Background()
	result, err := config.K3s.CoreV1().Pods(c.Param("namespace")).List(ctx, metaV1.ListOptions{})

	if err != nil {
		helper.ErrHandler(c, http.StatusInternalServerError, err.Error())
		return
	}

	helper.ResHandler(c, http.StatusOK, &models.Success{
		Status: "OK",
		Result: result,
		Items:  int64(len(result.Items)),
	})
}
