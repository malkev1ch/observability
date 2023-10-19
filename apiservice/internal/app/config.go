package app

type config struct {
	Address            string `env:"ADDRESS" envDefault:"0.0.0.0:8080"`
	UserServiceAddress string `env:"USER_SERVICE_ADDRESS" envDefault:"0.0.0.0:8081"`
}
