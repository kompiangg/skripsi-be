package config

type ShardingDatabase struct {
	Shards []Shard `yaml:"Shards"`
}

type Shard struct {
	RangeInDay    int    `yaml:"RangeInDay"`
	URIConnection string `yaml:"URIConnection"`
}

type LongTermDatabase struct {
	URIConnection string `yaml:"URIConnection"`
}
