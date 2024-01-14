package config

import (
	"skripsi-be/type/constant"
	"time"
)

type LongTermDatabase struct {
	URIConnection string `yaml:"URIConnection"`
}

type GeneralDatabase struct {
	URIConnection string `yaml:"URIConnection"`
}

type Shard struct {
	RangeInDay    int    `yaml:"RangeInDay"`
	URIConnection string `yaml:"URIConnection"`
}
type ShardingDatabase struct {
	IsUsingSharding bool   `yaml:"IsUsingSharding"`
	Shards          Shards `yaml:"Shards"`
}

type Shards []Shard

func (s Shards) GetShardIndexByDateTime(date time.Time, now time.Time) (int, error) {
	date = time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())
	now = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())

	diff := now.Sub(date)
	diffInDay := int(diff.Hours() / 24)

	for i, v := range s {
		if diffInDay-v.RangeInDay < 0 {
			return i, nil
		}
	}

	return 0, constant.ErrOutOfShardRange
}
