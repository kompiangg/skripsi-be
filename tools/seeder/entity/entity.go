package entity

import validation "github.com/go-ozzo/ozzo-validation/v4"

type SeederFlag struct {
	ConnectionType *string
}

func (s SeederFlag) Validate() error {
	err := validation.ValidateStruct(&s,
		validation.Field(&s.ConnectionType, validation.Required, validation.In("else", "general")),
	)
	if err != nil {
		return err
	}

	return nil
}
