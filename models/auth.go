package models

type Auth struct {
	AccessToken   string
	RefreshToken  string
	AccessUUID    string
	RefreshUUID   string
	AccessExpire  int64
	RefreshExpire int64
}

type Login struct {
	User string `json:"user" xml:"user" binding:"required"`
	Pass string `json:"pass" xml:"pass" binding:"required"`
}

type Tokens struct {
	AccessToken  string `json:"access_token" xml:"access_token"`
	RefreshToken string `json:"refresh_token" xml:"refresh_token"`
}
