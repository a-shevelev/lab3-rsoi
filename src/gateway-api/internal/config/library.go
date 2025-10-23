package config

type Library struct {
	IPAdrress string `envconfig:"IP_ADRRESS"`
	BaseURL   string `envconfig:"BASE_URL"`
}
