package app

type config struct {
	Address string `env:"ADDRESS" envDefault:"0.0.0.0:8080"`
}
