package config

type ShardingDatabase struct {
	IsUsingSharding bool    `yaml:"IsUsingSharding"`
	Shards          []Shard `yaml:"Shards"`
}

type Shard struct {
	RangeInDay    int    `yaml:"RangeInDay"`
	URIConnection string `yaml:"URIConnection"`
}

type LongTermDatabase struct {
	URIConnection string `yaml:"URIConnection"`
}

type GeneralDatabase struct {
	URIConnection string `yaml:"URIConnection"`
}
