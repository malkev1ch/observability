package model

type Config struct {
	Address               string `env:"ADDRESS,notEmpty"`
	OtelAddress           string `env:"OTEL_ADDRESS,notEmpty"`
	VoucherServiceAddress string `env:"VOUCHER_SERVICE_ADDRESS,notEmpty"`
}
