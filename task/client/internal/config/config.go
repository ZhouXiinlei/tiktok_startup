package config

type Config struct {
	Redis struct {
		Addr string
		Pass string
		DB   int
	}
}
