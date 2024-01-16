package config

type JWT struct {
	Admin   AdminJWTConfig   `yaml:"Admin"`
	Cashier CashierJWTConfig `yaml:"Cashier"`
}

type AdminJWTConfig JWTConfig
type CashierJWTConfig JWTConfig

type JWTConfig struct {
	Secret   string `yaml:"Secret"`
	ExpInDay int    `yaml:"ExpInDay"`
}
