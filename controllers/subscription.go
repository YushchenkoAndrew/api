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
// @Param _ query string false "For more info about query see request: 'GET /operations'"
// @Success 201 {object} models.Success{result=[]models.Subscription}
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

	var path string
	if path, err = helper.FormPathFromHandler(c, handler); err != nil {
		helper.ErrHandler(c, http.StatusNotFound, err.Error())
		return
	}

	hasher := md5.New()
	hasher.Write([]byte(strconv.Itoa(rand.Intn(1000000) + 5000)))
	token := hex.EncodeToString(hasher.Sum(nil))

	var reqBody []byte
	if reqBody, err = json.Marshal(&models.CronCreateDto{
		CronTime: body.CronTime,
		URL:      config.ENV.URL + path,
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

	var cron = models.CronEntity{}
	if err = json.NewDecoder(res.Body).Decode(&cron); err != nil {
		helper.ErrHandler(c, http.StatusInternalServerError, "Server side error: Something went wrong response")
		return
	}

	model := models.Subscription{
		CronID:   cron.ID,
		CronTime: cron.Exec.CronTime,
		Method:   handler.Method,
		Path:     path,
	}
	result := db.DB.Create(&model)

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

	go db.FlushValue("SUBSCRIPTION")
	helper.ResHandler(c, http.StatusCreated, &models.Success{
		Status: "OK",
		Result: []models.Subscription{model},
		Items:  1,

		// TODO: Maybe on day I'll add this ....
		// TotalItems: items,
	})
}

func (*subscriptionController) ReadAll(c *gin.Context) {}

// @Tags Subscription
// @Summary Read subscription by id/cron_id
// @Accept json
// @Produce application/json
// @Produce application/xml
// @Security BearerAuth
// @Param id path string true "This id can be a ID (Primary Key) or a CronID"
// @Success 200 {object} models.Success{result=[]models.Subscription}
// @failure 400 {object} models.Error
// @failure 401 {object} models.Error
// @failure 422 {object} models.Error
// @failure 429 {object} models.Error
// @failure 500 {object} models.Error
// @Router /subscription/{id} [get]
func (*subscriptionController) ReadOne(c *gin.Context) {
	var id string
	var model []models.Subscription

	if id = c.Param("id"); id == "" {
		helper.ErrHandler(c, http.StatusBadRequest, "Incorrect id value")
		return
	}

	var query = "cron_id = ?"
	if _, err := strconv.Atoi(id); err == nil {
		query = "id = ?"
	}

	hasher := md5.New()
	hasher.Write([]byte(fmt.Sprintf("ID=%s", id)))
	if err := helper.PrecacheResult(fmt.Sprintf("SUBSCRIPTION:%s", hex.EncodeToString(hasher.Sum(nil))), db.DB.Where(query, id), &model); err != nil {
		helper.ErrHandler(c, http.StatusInternalServerError, err.Error())
		go logs.DefaultLog("/controllers/subscription.go", err.Error())
		return
	}

	helper.ResHandler(c, http.StatusOK, &models.Success{
		Status: "OK",
		Result: model,
		Items:  1,

		// TODO: Maybe one day ....
		// TotalItems: items,
	})
}

func (*subscriptionController) UpdateOne(c *gin.Context) {}
func (*subscriptionController) UpdateAll(c *gin.Context) {}

// @Tags Subscription
// @Summary Delete subscription by id/cron_id
// @Accept json
// @Produce application/json
// @Produce application/xml
// @Security BearerAuth
// @Param id path string true "This id can be a ID (Primary Key) or a CronID"
// @Success 200 {object} models.Success{result=[]string{}}
// @failure 400 {object} models.Error
// @failure 401 {object} models.Error
// @failure 422 {object} models.Error
// @failure 429 {object} models.Error
// @failure 500 {object} models.Error
// @Router /subscription/{id} [delete]
func (*subscriptionController) DeleteOne(c *gin.Context) {
	var id string
	var model []models.Subscription

	if id = c.Param("id"); id == "" {
		helper.ErrHandler(c, http.StatusBadRequest, "Incorrect id value")
		return
	}

	var query = "cron_id = ?"
	if _, err := strconv.Atoi(id); err == nil {
		query = "id = ?"
	}

	hasher := md5.New()
	hasher.Write([]byte(fmt.Sprintf("ID=%s", id)))
	if err := helper.PrecacheResult(fmt.Sprintf("SUBSCRIPTION:%s", hex.EncodeToString(hasher.Sum(nil))), db.DB.Where(query, id), &model); err != nil || len(model) == 0 {
		helper.ErrHandler(c, http.StatusInternalServerError, err.Error())
		go logs.DefaultLog("/controllers/subscription.go", err.Error())
		return
	}

	hasher = md5.New()
	salt := strconv.Itoa(rand.Intn(1000000) + 5000)
	hasher.Write([]byte(salt + config.ENV.BotKey))

	var req *http.Request
	var err error
	if req, err = http.NewRequest("DELETE", fmt.Sprintf("%s/cron/subscribe?key=%s&id=%s", config.ENV.BotUrl, hex.EncodeToString(hasher.Sum(nil)), model[0].CronID), nil); err != nil {
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

	if res.StatusCode != http.StatusOK {
		helper.ErrHandler(c, res.StatusCode, "Bot request error")
		return
	}

	db.DB.Where("id = ?", model[0].ID).Delete(&models.Subscription{})
	go db.FlushValue("SUBSCRIPTION")
	helper.ResHandler(c, http.StatusOK, &models.Success{
		Status: "OK",
		Result: []string{},
	})
}

func (*subscriptionController) DeleteAll(c *gin.Context) {}
