package models

type CronCreateDto struct {
	CronTime string `json:"cron_time" xml:"cron_time" example:"00 00 00 */1 * *"`
	URL      string `json:"url" xml:"url" example:"http://127.0.0.1:8000/ping"`
	Method   string `json:"method" xml:"method" example:"post"`
	Token    string `json:"token" xml:"token" example:"HELLO_WORLD"`
	Data     string `json:"data,omitempty" xml:"data" example:"{'data' : 'Hello world'}"`
}

type CronRes struct {
	ID        string        `json:"id" xml:"id" example:"d266389ebf09e1e8a95a5b4286b504b2"`
	CreatedAt string        `json:"created_at" xml:"created_at" example:"Mon Jan 31 2022 00:00:00 GMT+0000 (Coordinated Universal Time)"`
	Exec      CronCreateDto `json:"exec" xml:"exec"`
}
