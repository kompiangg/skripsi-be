package constant

import "github.com/go-errors/errors"

var (
	ErrOutOfShardRange = errors.New("date is out of shard range")
)
