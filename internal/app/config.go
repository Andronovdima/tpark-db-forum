package apiserver

type Config struct {
	BindAddr    string
	LogLevel    string
	DatabaseURL string
	SessionKey  string
	ClientUrl	string
}

func NewConfig() *Config {
	return &Config{
		BindAddr:		":8080",
		LogLevel:		"debug",
		SessionKey:		"jdfhdfdj",
		DatabaseURL:	"host=localhost dbname=db-forum sslmode=disable port=5432 password=1234 user=d",
		//ClientUrl:		"http://127.0.0.1:9000", //"http://localhost:9000",//
		ClientUrl:		"http://89.208.211.100:9000",
	}
}

