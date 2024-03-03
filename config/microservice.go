package config

type Microservice struct {
	GeneralService   GeneralService   `yaml:"GeneralService"`
	AuthService      AuthService      `yaml:"AuthService"`
	IngestionService IngestionService `yaml:"IngestionService"`
	ServingService   ServingService   `yaml:"ServingService"`
}

type GeneralService HTTPServer
type AuthService HTTPServer
type IngestionService HTTPServer
type ServingService HTTPServer

type HTTPServer struct {
	Host                 string   `yaml:"Host"`
	Port                 int      `yaml:"Port"`
	WhiteListAllowOrigin []string `yaml:"WhiteListAllowOrigin"`
}
