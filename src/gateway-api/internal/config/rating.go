package config

type Rating struct {
	IPAdrress string `envconfig:"IP_ADRRESS"`
	BaseURL   string `envconfig:"BASE_URL"`
}
