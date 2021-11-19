package models

type LogMessage struct {
	Stat    string      `json:"stat" xml:"stat" binding:"required"`
	Name    string      `json:"name" xml:"name" binding:"required"`
	Url     string      `json:"url" xml:"url"`
	File    string      `json:"file" xml:"file"`
	Message string      `json:"message" xml:"message" binding:"required"`
	Desc    interface{} `json:"desc" xml:"desc"`
}

type BotRedis struct {
	Command string `json:"command" xml:"command" binding:"required"`
}
