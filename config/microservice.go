package config

type Microservice struct {
	OrderService OrderService `yaml:"OrderService"`
}

type OrderService HTTPServer

type HTTPServer struct {
	Host                 string   `yaml:"Host"`
	Port                 int      `yaml:"Port"`
	WhiteListAllowOrigin []string `yaml:"WhiteListAllowOrigin"`
}
