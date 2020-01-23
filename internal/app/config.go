package apiserver

type Config struct {
	BindAddr    string
	LogLevel    string
	DatabaseURL string
}

func NewConfig() *Config {
	return &Config{
		BindAddr:		":5000",
		LogLevel:		"debug",
		DatabaseURL:	"host=localhost dbname=docker sslmode=disable port=5432 password=docker user=docker",
	}
}

