package config

type Kafka struct {
	Server string     `yaml:"Server"`
	Topic  KafkaTopic `yaml:"Topic"`
	Group  KafkaGroup `yaml:"Group"`
}

type KafkaGroup struct {
	Shard     string `yaml:"Shard"`
	LongTerm  string `yaml:"LongTerm"`
	Transform string `yaml:"Transform"`
}

type KafkaTopic struct {
	TransformOrder string `yaml:"TransformOrder"`
	LoadOrder      string `yaml:"LoadOrder"`
}
