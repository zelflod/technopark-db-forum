package apiserver

type Config struct {
	BindAddr   string
	SessionKey string
	//DatabaseURL string

}

func NewConfig() *Config {
	return &Config{
		BindAddr:   ":8080",
		SessionKey: "jdfhdfdj",
		//DatabaseURL:	"host=localhost dbname=restapi_dev sslmode=disable port=5432 password=1234 user=d",
	}
}
