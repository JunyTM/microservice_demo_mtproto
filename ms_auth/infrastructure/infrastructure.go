package infrastructure

import (
	"flag"
	"ms_auth/model"

	"gorm.io/gorm"
)

const (
	msql_dns = "user:pass@tcp(127.0.0.1:3306)/dbname?charset=utf8mb4&parseTime=True&loc=Local"
)

var (
	db       *gorm.DB
	cache_user []model.User
)

func init() {
	// load db connection
	ConnectDatabases()

	// migrate database
	isMigrate := flag.Bool("db", false, "Migrate database")
	MigrateDatabases(*isMigrate)
}

// GetDB return the database connection
func GetDB() *gorm.DB {
	return db
}

func CacheInMem(gorm.Model) {
	cache_user = append(cache_user, model.User{})
}

func GetCache() []model.User {
	return cache_user
}