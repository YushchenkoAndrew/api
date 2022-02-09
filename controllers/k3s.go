package controllers

import (
	"api/interfaces"

	"github.com/gin-gonic/gin"
)

const (
	SUBSCRIBE_TIME = "00 00 */4 * * *"
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
// @Param model body models.SubscribeDto true "Small info about subscription for k3s"
// @Success 201 {object} models.Success{result=[]models.Link}
// @failure 400 {object} models.Error
// @failure 401 {object} models.Error
// @failure 422 {object} models.Error
// @failure 429 {object} models.Error
// @failure 500 {object} models.Error
// @Router /subscribe [post]
func (*k3sController) Subscribe(c *gin.Context) {
	// var err error
	// var body []byte

	// hasher := md5.New()
	// hasher.Write([]byte(strconv.Itoa(rand.Intn(1000000) + 5000)))
	// token := hex.EncodeToString(hasher.Sum(nil))

	// if body, err = json.Marshal(models.CronCreateDto{
	// 	CronTime: SUBSCRIBE_TIME,
	// 	URL:      fmt.Sprintf("%s/k3s/pods/metrics", config.ENV.URL),
	// 	Method:   "post",
	// 	Key:      token,
	// }); err != nil {
	// 	fmt.Println("Ohh noo; Anyway")
	// 	return
	// }

	// hasher = md5.New()
	// var salt = strconv.Itoa(rand.Intn(1000000))
	// hasher.Write([]byte(salt + config.ENV.BotKey))

	// var req *http.Request
	// if req, err = http.NewRequest("POST", config.ENV.BotUrl+"/cron?key="+hex.EncodeToString(hasher.Sum(nil)), bytes.NewBuffer(body)); err != nil {
	// 	fmt.Println("Ohh noo; Anyway")
	// 	return
	// }

	// req.Header.Set("X-Custom-Header", salt)
	// req.Header.Set("Content-Type", "application/json")

	// var res *http.Response
	// client := &http.Client{}
	// res, err = client.Do(req)
	// if err != nil {
	// 	fmt.Println("Ohh noo; Anyway")
	// 	return
	// }

	// // TODO: Save Request response into db
	// defer res.Body.Close()

	// hasher = md5.New()
	// hasher.Write([]byte(token))

	// ctx := context.Background()
	// db.Redis.Set(ctx, "TOKEN:"+hex.EncodeToString(hasher.Sum(nil)), "OK", 0)
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
	// TODO: Finish the thing above !!!
}
