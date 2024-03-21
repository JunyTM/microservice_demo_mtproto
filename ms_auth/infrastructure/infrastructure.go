package infrastructure

import (
	"flag"
	"log"
	"ms_auth/model"

	"gorm.io/gorm"
)

const (
	msql_dns = "users:147563@tcp(0.0.0.0:3306)/app_auth?charset=utf8mb4&parseTime=True&loc=Local"
)

var (
	db        *gorm.DB
	cacheUser map[string]model.User
)

func init() {
	// Make memcache
	cacheUser = make(map[string]model.User)

	// load db connection
	ConnectDatabases()

	isMigrate := flag.Bool("db", false, "Migrate database")
	isLoadCache := flag.Bool("cache", true, "Load cache")
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

func GetCache() map[string]model.User {
	return cacheUser
}

func loadMemoryCache(isLoadCache bool) {
	if !isLoadCache {
		return
	}
	var caches []model.CacheMem
	err := db.Model(&model.CacheMem{}).Find(&caches).Order("id").Error
	if err != nil {
		log.Println("=> Warrning: Cannot load MemCache")
		return
	}

	for i := range caches {
		cacheUser[caches[i].Email] = model.User{}
	}
}
