package order

import (
	"context"
	"skripsi-be/pkg/errors"
	"skripsi-be/type/params"
	"skripsi-be/type/result"

	"github.com/google/uuid"
)

func (s service) IngestOrder(ctx context.Context, param []params.ServiceIngestionOrder) ([]result.ServiceIngestOrder, error) {
	res := make([]result.ServiceIngestOrder, len(param))
	repoParam := make([]params.RepositoryPublishTransformOrderEvent, len(param))

	storeIDMap := make(map[string]bool)
	cashierIDMap := make(map[uuid.UUID]bool)
	customerIDMap := make(map[string]bool)
	itemIDMap := make(map[string]bool)

	for i, v := range param {
		err := v.Validate()
		if err != nil {
			return nil, errors.Wrap(err)
		}

		if v.CreatedAt.After(s.config.Date.Now()) {
			return nil, errors.Wrap(errors.ErrDataParamMustNotAFterCurrentTime)
		}

		if _, ok := storeIDMap[v.StoreID]; !ok {
			_, err = s.storeRepo.FindByID(ctx, v.StoreID)
			if errors.Is(err, errors.ErrRecordNotFound) {
				return nil, errors.Wrap(errors.ErrStoreNotFound)
			} else if err != nil {
				return nil, errors.Wrap(err)
			}

			storeIDMap[v.StoreID] = true
		}

		var cashierStoreID string
		if _, ok := cashierIDMap[v.CashierID]; !ok {
			cashier, err := s.cashierRepo.FindByID(ctx, v.CashierID)
			if errors.Is(err, errors.ErrRecordNotFound) {
				return nil, errors.Wrap(errors.ErrCashierNotFound)
			} else if err != nil {
				return nil, errors.Wrap(err)
			}

			cashierStoreID = cashier.StoreID
			cashierIDMap[v.CashierID] = true
		}

		if cashierStoreID != v.StoreID {
			return nil, errors.Wrap(errors.ErrCustomerCashierNotMatch)
		}

		if _, ok := customerIDMap[v.CustomerID]; !ok {
			_, err = s.customerRepo.FindByID(ctx, v.CustomerID)
			if errors.Is(err, errors.ErrRecordNotFound) {
				return nil, errors.Wrap(errors.ErrCustomerNotFound)
			} else if err != nil {
				return nil, errors.Wrap(err)
			}

			customerIDMap[v.CustomerID] = true
		}

		for _, item := range v.OrderDetails {
			if _, ok := itemIDMap[item.ItemID]; !ok {
				_, err = s.itemRepo.FindByID(ctx, item.ItemID)
				if errors.Is(err, errors.ErrRecordNotFound) {
					return nil, errors.Wrap(errors.ErrItemNotFound)
				} else if err != nil {
					return nil, errors.Wrap(err)
				}
			}
		}

		repoParam[i] = v.ToRepositoryPublishTransformOrderEvent()
		res[i].FromParamServiceIngestionOrder(v, repoParam[i].ID)
	}

	err := s.publisherRepo.PublishTransformOrderEvent(ctx, repoParam)
	if err != nil {
		return nil, errors.Wrap(err)
	}

	return res, nil
}
