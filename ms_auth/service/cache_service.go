package service

import (
	"errors"
	"ms_auth/infrastructure"
	"ms_auth/model"
)

type cacheService struct{}

func (s *cacheService) CheckInMem(key string) (*model.User, error) {
	dataCache := *infrastructure.GetCache()
	if dataCache == nil {
		return nil, errors.New("missing cache key")
	}

	user, ok := dataCache[key]
	if !ok {
		return nil, errors.New("no email in memory")
	}
	return &user, nil
}

func (s *cacheService) AddInMem(user *model.User) {
	data_cache := *infrastructure.GetCache()
	data_cache[user.Email] = *user
}
