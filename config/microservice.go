package config

type Microservice struct {
	GeneralService      GeneralService      `yaml:"GeneralService"`
	AuthService         AuthService         `yaml:"AuthService"`
	IngestionService    IngestionService    `yaml:"IngestionService"`
	ServingService      ServingService      `yaml:"ServingService"`
	ShardingLoadService ShardingLoadService `yaml:"ShardingLoadService"`
	LongtermLoadService LongtermLoadService `yaml:"LongtermLoadService"`
	TransformService    TransformService    `yaml:"TransformService"`
}

type GeneralService HTTPServer
type AuthService HTTPServer
type IngestionService HTTPServer
type ServingService HTTPServer
type ShardingLoadService HTTPServer
type LongtermLoadService HTTPServer
type TransformService HTTPServer

type HTTPServer struct {
	Host                 string   `yaml:"Host"`
	Port                 int      `yaml:"Port"`
	WhiteListAllowOrigin []string `yaml:"WhiteListAllowOrigin"`
}
