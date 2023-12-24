package service

import (
	"skripsi-be/config"
	"skripsi-be/repository"
)

type Service struct {
}

type Config struct {
}

func New(
	repository repository.Repository,
	config config.Config,
) (Service, error) {
	return Service{}, nil
}
