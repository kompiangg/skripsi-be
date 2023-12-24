package config

type Redis struct {
	Hostname string `yaml:"Hostname"`
	Username string `yaml:"Username"`
	Password string `yaml:"Password"`
	DB       int    `yaml:"DB"`
}
