package config

type Microservice struct {
	OrderService HTTPServer `yaml:"OrderService"`
}

type HTTPServer struct {
	Host                 string   `yaml:"Host"`
	Port                 int      `yaml:"Port"`
	WhiteListAllowOrigin []string `yaml:"WhiteListAllowOrigin"`
}
