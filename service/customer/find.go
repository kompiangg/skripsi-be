package customer

import (
	"context"
	"skripsi-be/type/result"
)

func (s service) FindLikeOneOfAllColumn(ctx context.Context, req string) ([]result.Customer, error) {
	customers, err := s.customerRepo.FindByIDLikeOrNameLikeOrContactLike(ctx, req)
	if err != nil {
		return nil, err
	}

	res := make([]result.Customer, len(customers))
	for i, customer := range customers {
		res[i].FromModel(customer)
	}

	return res, nil
}
