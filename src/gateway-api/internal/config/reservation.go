package config

type Reservation struct {
	IPAdrress string `envconfig:"IP_ADRRESS"`
	BaseURL   string `envconfig:"BASE_URL"`
}
