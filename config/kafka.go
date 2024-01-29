package config

type Kafka struct {
	Server string     `yaml:"Server"`
	Topic  string     `yaml:"Topic"`
	Group  KafkaGroup `yaml:"Group"`
}

type KafkaGroup struct {
	Shard    string `yaml:"Shard"`
	LongTerm string `yaml:"LongTerm"`
}
