package infrastructure

import (
	"flag"
	"log"
	"ms_auth/model"

	"gorm.io/gorm"
)

const (
	msql_dns = "user:pass@tcp(127.0.0.1:3306)/dbname?charset=utf8mb4&parseTime=True&loc=Local"
)

var (
	db         *gorm.DB
	cache_user map[string]model.User
)

func init() {
	// load db connection
	ConnectDatabases()

	isMigrate := flag.Bool("db", false, "Migrate database")
	isLoadCache := flag.Bool("cache", false, "Load cache")
	flag.Parse()

	// migrate database
	MigrateDatabases(*isMigrate)

	// load cache
	loadMemoryCache(*isLoadCache)
}

// GetDB return the database connection
func GetDB() *gorm.DB {
	return db
}

func loadMemoryCache(isLoadCache bool) {
	if !isLoadCache {
		return
	}

	var users []model.User
	err := db.Model(&model.User{}).Find(&users).Order("id").Error
	if err != nil {
		log.Println("=> Warrning: Cannot load MemCache")
	}

	for i := range users {
		cache_user[users[i].Email] = users[i]
	}
}

func CacheInMem(key string, value model.User) {
	cache_user[key] = value
}

func GetCache() map[string]model.User {
	return cache_user
}
