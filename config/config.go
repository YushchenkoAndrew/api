package config

func Init() {
	LoadEnv("./")
	LoadK3s()
}
