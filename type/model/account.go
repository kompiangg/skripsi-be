package model

import (
	"skripsi-be/pkg/errors"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type Account struct {
	ID        uuid.UUID `db:"id"`
	Username  string    `db:"username"`
	Password  string    `db:"password"`
	CreatedAt time.Time `db:"created_at"`
}

func (a *Account) AssignUUID() bool {
	if a.ID != uuid.Nil {
		return false
	}

	a.ID = uuid.New()
	return true
}

func (a *Account) HashPassword() error {
	_, err := bcrypt.Cost([]byte(a.Password))
	if err == nil {
		return nil
	}

	bytes, err := bcrypt.GenerateFromPassword([]byte(a.Password), bcrypt.DefaultCost)
	if err != nil {
		return errors.New(err)
	}

	a.Password = string(bytes)
	return nil
}

func (a *Account) ComparePassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(a.Password), []byte(password))
	return err == nil
}
