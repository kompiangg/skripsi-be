package params

import (
	"skripsi-be/type/model"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/google/uuid"
)

type ServiceRegisterAdmin struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Name     string `json:"name"`

	RequestBy     uuid.UUID `json:"-"`
	RequestByRole string    `json:"-"`
}

func (s ServiceRegisterAdmin) Validate() error {
	return validation.ValidateStruct(&s,
		validation.Field(&s.Username, validation.Required),
		validation.Field(&s.Password, validation.Required, validation.Length(8, 72)),
		validation.Field(&s.Name, validation.Required),
	)
}

func (s ServiceRegisterAdmin) ToAccountModel() model.Account {
	return model.Account{
		Username: s.Username,
		Password: s.Password,
	}
}

func (s ServiceRegisterAdmin) ToAdminModel(accountID uuid.UUID) model.Admin {
	return model.Admin{
		Name:      s.Name,
		AccountID: accountID,
	}
}

type ServiceLogin struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (s ServiceLogin) Validate() error {
	return validation.ValidateStruct(&s,
		validation.Field(&s.Username, validation.Required),
		validation.Field(&s.Password, validation.Required),
	)
}

type ServiceRegisterCashier struct {
	Username string `json:"username"`
	Password string `json:"password"`
	StoreID  string `json:"store_id"`
	Name     string `json:"name"`

	RequestBy uuid.UUID `json:"-"`
}

func (s ServiceRegisterCashier) Validate() error {
	return validation.ValidateStruct(&s,
		validation.Field(&s.Username, validation.Required),
		validation.Field(&s.Password, validation.Required, validation.Length(8, 72)),
		validation.Field(&s.Name, validation.Required),
	)
}

func (s ServiceRegisterCashier) ToAccountModel() model.Account {
	return model.Account{
		Username: s.Username,
		Password: s.Password,
	}
}

func (s ServiceRegisterCashier) ToCashierModel(accountID uuid.UUID) model.Cashier {
	return model.Cashier{
		Name:      s.Name,
		AccountID: accountID,
		StoreID:   s.StoreID,
	}
}
