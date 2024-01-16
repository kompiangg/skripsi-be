package config

type Microservice struct {
	OrderService        OrderService        `yaml:"OrderService"`
	AuthService         AuthService         `yaml:"AuthService"`
	IngestionService    IngestionService    `yaml:"IngestionService"`
	LongTermLoadService LongTermLoadService `yaml:"LongTermLoadService"`
	ShardingLoadService ShardingLoadService `yaml:"ShardingLoadService"`
	ServingService      ServingService      `yaml:"ServingService"`
}

type OrderService HTTPServer
type AuthService HTTPServer
type IngestionService HTTPServer
type LongTermLoadService HTTPServer
type ShardingLoadService HTTPServer
type ServingService HTTPServer

type HTTPServer struct {
	Host                 string   `yaml:"Host"`
	Port                 int      `yaml:"Port"`
	WhiteListAllowOrigin []string `yaml:"WhiteListAllowOrigin"`
}
