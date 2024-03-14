package order

import (
	"context"
	"skripsi-be/pkg/errors"
	"skripsi-be/type/model"
	"skripsi-be/type/params"
	"skripsi-be/type/result"

	"github.com/google/uuid"
)

func (s service) IngestOrder(ctx context.Context, param []params.ServiceIngestionOrder) ([]result.ServiceIngestOrder, error) {
	res := make([]result.ServiceIngestOrder, len(param))
	repoParam := make([]params.RepositoryPublishTransformOrderEvent, len(param))

	cashierIDMap := make(map[uuid.UUID]bool)
	customerIDMap := make(map[string]bool)
	itemIDMap := make(map[string]model.Item)

	for i, v := range param {
		err := v.Validate()
		if err != nil {
			return nil, errors.Wrap(err)
		}

		if v.PaymentDate.After(s.config.Date.Now()) {
			return nil, errors.Wrap(errors.ErrDataParamMustNotAFterCurrentTime)
		}

		if _, ok := cashierIDMap[v.CashierID]; !ok {
			_, err := s.cashierRepo.FindByID(ctx, v.CashierID)
			if errors.Is(err, errors.ErrRecordNotFound) {
				return nil, errors.Wrap(errors.ErrCashierNotFound)
			} else if err != nil {
				return nil, errors.Wrap(err)
			}

			cashierIDMap[v.CashierID] = true
		}

		if v.CustomerID.Valid {
			if _, ok := customerIDMap[v.CustomerID.String]; !ok {
				_, err = s.customerRepo.FindByID(ctx, v.CustomerID.String)
				if errors.Is(err, errors.ErrRecordNotFound) {
					return nil, errors.Wrap(errors.ErrCustomerNotFound)
				} else if err != nil {
					return nil, errors.Wrap(err)
				}

				customerIDMap[v.CustomerID.String] = true
			}
		}

		for idx, item := range v.OrderDetails {
			if _, ok := itemIDMap[item.ItemID]; !ok {
				itemData, err := s.itemRepo.FindByID(ctx, item.ItemID)
				if errors.Is(err, errors.ErrRecordNotFound) {
					return nil, errors.Wrap(errors.ErrItemNotFound)
				} else if err != nil {
					return nil, errors.Wrap(err)
				}

				v.OrderDetails[idx].Price = itemData.Price
				v.OrderDetails[idx].Unit = itemData.Unit
				itemIDMap[item.ItemID] = itemData
			}
		}

		repoParam[i] = v.ToRepositoryPublishTransformOrderEvent()

		for _, item := range v.OrderDetails {
			repoParam[i].OrderDetails = append(repoParam[i].OrderDetails, item.ToRepositoryPublishTransformOrderDetailEvent(itemIDMap[item.ItemID]))
		}

		res[i].FromParamServiceIngestionOrder(v, repoParam[i].ID)
	}

	err := s.publisherRepo.PublishTransformOrderEvent(ctx, repoParam)
	if err != nil {
		return nil, errors.Wrap(err)
	}

	return res, nil
}
