package config

type JWT struct {
	Secret   string `yaml:"Secret"`
	ExpInDay int    `yaml:"ExpInDay"`
}

type JWTConfig struct {
	Secret   string `yaml:"Secret"`
	ExpInDay int    `yaml:"ExpInDay"`
}
