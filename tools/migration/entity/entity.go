package entity

import validation "github.com/go-ozzo/ozzo-validation/v4"

type MigrationFlag struct {
	Operation      *string
	ConnectionType *string
	TableName      *string
}

func (m MigrationFlag) Validate() error {
	err := validation.ValidateStruct(&m,
		validation.Field(&m.Operation, validation.Required, validation.In("new", "up", "down", "drop")),
		validation.Field(&m.ConnectionType, validation.Required, validation.In("longterm", "shard", "general")),
	)
	if err != nil {
		return err
	}

	if *m.Operation == "new" {
		err := validation.ValidateStruct(&m,
			validation.Field(&m.TableName, validation.Required),
		)
		if err != nil {
			return err
		}
	}

	return nil
}
