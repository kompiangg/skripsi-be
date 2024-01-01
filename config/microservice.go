package config

type Microservice struct {
	WebService HTTPServer `yaml:"WebService"`
}

type HTTPServer struct {
	Host                 string   `yaml:"Host"`
	Port                 int      `yaml:"Port"`
	WhiteListAllowOrigin []string `yaml:"WhiteListAllowOrigin"`
}
