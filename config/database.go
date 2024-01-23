package config

type LongTermDatabase struct {
	URIConnection string `yaml:"URIConnection"`
}

type GeneralDatabase struct {
	URIConnection string `yaml:"URIConnection"`
}

type Shard struct {
	DataRetention int    `yaml:"DataRetention"`
	URIConnection string `yaml:"URIConnection"`
}
type ShardingDatabase struct {
	IsUsingSharding bool   `yaml:"IsUsingSharding"`
	Shards          Shards `yaml:"Shards"`
}

type Shards []Shard
