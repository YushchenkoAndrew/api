package logs

import (
	"api/config"
	"api/models"
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"strconv"
)

func SendLogs(message *models.LogMessage) {
	var err error
	var body []byte

	if body, err = json.Marshal(*message); err != nil {
		fmt.Printf("Ohh noo; Anyway: %v", err)
		return
	}

	hasher := md5.New()
	var salt = strconv.Itoa(rand.Intn(1000000))
	hasher.Write([]byte(salt + config.ENV.BotKey))

	var req *http.Request
	if req, err = http.NewRequest("POST", config.ENV.BotUrl+"/logs/alert?key="+hex.EncodeToString(hasher.Sum(nil)), bytes.NewBuffer(body)); err != nil {
		fmt.Printf("Ohh nyo; Anyway: %v", err)
		return
	}

	req.Header.Set("X-Custom-Header", salt)
	req.Header.Set("Content-Type", "application/json")

	var res *http.Response
	client := &http.Client{}
	res, err = client.Do(req)
	if err != nil {
		fmt.Printf("Ohh noo; Anyway: %v", err)
		return
	}

	defer res.Body.Close()
}

func DefaultLog(file string, err interface{}) {
	SendLogs(&models.LogMessage{
		Stat:    "ERR",
		Name:    "API",
		File:    file,
		Message: "Ohh nooo something is broken; Anyway...",
		Desc:    err,
	})
}
