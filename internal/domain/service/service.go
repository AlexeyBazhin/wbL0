package service

import "github.com/AlexeyBazhin/wbL0/internal/domain"

type service struct {
	repo  domain.Repository
	cache domain.Cache
}

func NewService(store domain.Repository, cache domain.Cache) *service {
	return &service{
		repo:  store,
		cache: cache,
	}
}
