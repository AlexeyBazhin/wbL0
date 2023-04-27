package service

import "github.com/AlexeyBazhin/wbL0/internal/domain"

type service struct {
	repo domain.Repository
}

func NewService(store domain.Repository) *service {
	return &service{
		repo: store,
	}
}
