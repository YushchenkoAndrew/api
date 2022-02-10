package controllers

import (
	"api/config"
	"api/db"
	"api/helper"
	"api/interfaces"
	"api/logs"
	"api/models"
	"bytes"
	"context"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type subscriptionController struct{}

func NewSubscriptionController() interfaces.Default {
	return &subscriptionController{}
}

func (*subscriptionController) CreateAll(c *gin.Context) {}

// @Tags Subscription
// @Summary Create Subscription to run operation
// @Accept json
// @Produce application/json
// @Produce application/xml
// @Security BearerAuth
// @Param model body models.SubscribeDto true "Small info about subscription for k3s"
// @Success 201 {object} models.Success{result=[]models.DefultRes}
// @failure 400 {object} models.Error
// @failure 401 {object} models.Error
// @failure 422 {object} models.Error
// @failure 429 {object} models.Error
// @failure 500 {object} models.Error
// @Router /subscription [post]
func (*subscriptionController) CreateOne(c *gin.Context) {
	var err error
	var body models.SubscribeDto
	if err = c.ShouldBind(&body); err != nil || body.CronTime == "" {
		helper.ErrHandler(c, http.StatusBadRequest, "Incorrect body params")
		return
	}

	handler, ok := config.GetOperation(body.Operation)
	if !ok {
		helper.ErrHandler(c, http.StatusNotFound, fmt.Sprintf("Operation '%s' not founded", body.Operation))
		return
	}

	hasher := md5.New()
	hasher.Write([]byte(strconv.Itoa(rand.Intn(1000000) + 5000)))
	token := hex.EncodeToString(hasher.Sum(nil))

	var reqBody []byte
	if reqBody, err = json.Marshal(&models.CronCreateDto{
		CronTime: body.CronTime,
		URL:      config.ENV.URL + handler.Path,
		Method:   handler.Method,
		Token:    token,
	}); err != nil {
		fmt.Printf("Ohh noo; Anyway: %v", err)
		return
	}

	hasher = md5.New()
	salt := strconv.Itoa(rand.Intn(1000000) + 5000)
	hasher.Write([]byte(salt + config.ENV.BotKey))

	var req *http.Request
	if req, err = http.NewRequest("POST", config.ENV.BotUrl+"/cron/subscribe?key="+hex.EncodeToString(hasher.Sum(nil)), bytes.NewBuffer(reqBody)); err != nil {
		fmt.Printf("Ohh noo; Anyway: %v", err)
		return
	}

	req.Header.Set("X-Custom-Header", salt)
	req.Header.Set("Content-Type", "application/json")

	var res *http.Response
	client := &http.Client{}
	res, err = client.Do(req)
	if err != nil {
		helper.ErrHandler(c, http.StatusInternalServerError, "Server side error: Something went wrong response")
		return
	}

	defer res.Body.Close()

	if res.StatusCode != http.StatusCreated {
		helper.ErrHandler(c, res.StatusCode, "Bot request error")
		return
	}

	var cron = models.CronRes{}
	if err = json.NewDecoder(res.Body).Decode(&cron); err != nil {
		helper.ErrHandler(c, http.StatusInternalServerError, "Server side error: Something went wrong response")
		return
	}

	result := db.DB.Create(&models.Subscription{
		CronID:    cron.ID,
		CronTime:  cron.Exec.CronTime,
		Operation: body.Operation,
	})

	if result.Error != nil || result.RowsAffected == 0 {
		helper.ErrHandler(c, http.StatusInternalServerError, "Something unexpected happend")
		go logs.DefaultLog("/controllers/subscription.go", result.Error)
		return
	}

	go func() {
		hasher = md5.New()
		hasher.Write([]byte(token))

		ctx := context.Background()
		db.Redis.Set(ctx, "TOKEN:"+hex.EncodeToString(hasher.Sum(nil)), "OK", 0)
	}()

	helper.ResHandler(c, http.StatusOK, models.DefultRes{
		Status:  "OK",
		Message: "Success",
		Result:  []string{},
	})
}

func (*subscriptionController) ReadOne(c *gin.Context)   {}
func (*subscriptionController) ReadAll(c *gin.Context)   {}
func (*subscriptionController) UpdateOne(c *gin.Context) {}
func (*subscriptionController) UpdateAll(c *gin.Context) {}

// @Tags Subscription
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
func (*subscriptionController) DeleteOne(c *gin.Context) {
	// TODO: Finish the thing above !!!
}

func (*subscriptionController) DeleteAll(c *gin.Context) {}
