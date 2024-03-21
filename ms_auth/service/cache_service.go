package service

import (
	"errors"
	"log"
	"ms_auth/infrastructure"
	"ms_auth/model"
)

type cacheService struct{}

func (s *cacheService) CheckInMem(key string) (*model.User, error) {
	dataCache := infrastructure.GetCache()
	if dataCache == nil {
		return nil, errors.New("missing cache key")
	}

	user, ok := dataCache[key]
	if !ok {
		return nil, errors.New("no email in memory")
	}
	db := infrastructure.GetDB()
	if err := db.Model(&model.User{}).Find(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (s *cacheService) AddInMem(user *model.User) {
	data_cache := infrastructure.GetCache()
	data_cache[user.Email] = *user
	db := infrastructure.GetDB()
	if err := db.Model(&model.CacheMem{}).Create(&model.CacheMem{
		Email: user.Email,
	}).Error; err != nil {
		log.Println("=====> Missing in-memory")
	}
	// log.Println("=====> Memory:", data_cache)
}
