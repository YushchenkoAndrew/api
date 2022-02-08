package models

type CronCreateDto struct {
	CronTime string `json:"cronTime" xml:"cronTime" example:"00 00 00 */1 * *"`
	URL      string `json:"url" xml:"url" example:"http://127.0.0.1:8000/ping"`
	Method   string `json:"method" xml:"method" example:"post"`
	ApiKey   string `json:"apiKey" xml:"apiKey" example:"post"`
	Data     string `json:"data,omitempty" xml:"data" example:""`
}

type CronRes struct {
	ID        string        `json:"id" xml:"id"`
	CreatedAt string        `json:"createdAt" xml:"createdAt" example:"Mon Jan 31 2022 00:00:00 GMT+0000 (Coordinated Universal Time)"`
	Exec      CronCreateDto `json:"exec" xml:"exec"`
}
