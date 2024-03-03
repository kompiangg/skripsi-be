package item

import (
	"context"
	"skripsi-be/type/result"
)

func (s service) FindLikeNameOrID(ctx context.Context, nameOrID string) ([]result.Item, error) {
	items, err := s.itemRepo.FindLikeNameOrID(ctx, nameOrID)
	if err != nil {
		return nil, err
	}

	res := make([]result.Item, len(items))
	for i, item := range items {
		res[i].FromModel(item)
	}

	return res, nil
}
