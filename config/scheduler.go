package config

type Scheduler struct {
	MoveShardingData BaseScheduler `yaml:"MoveShardingData"`
}

type BaseScheduler struct {
	Enable   bool   `yaml:"Enable"`
	Duration string `yaml:"Duration"`
}
