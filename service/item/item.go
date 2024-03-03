package item

import (
	"context"
	"skripsi-be/repository/item"
	"skripsi-be/type/result"
)

type Service interface {
	FindLikeNameOrID(ctx context.Context, nameOrID string) ([]result.Item, error)
}

type Config struct{}

type service struct {
	config   Config
	itemRepo item.Repository
}

func New(
	config Config,
	itemRepo item.Repository,
) Service {
	return service{
		config:   config,
		itemRepo: itemRepo,
	}
}
