package infrastructure

import (
	"crypto/rsa"
	"flag"
	"log"
	"ms_auth/model"
	"sync"

	"gorm.io/gorm"
)

const (
	// msql_dns = "users:147563@tcp(authen_db:3306)/app_auth?charset=utf8mb4&parseTime=True&loc=Local"
	msql_dns = "users:147563@tcp(localhost:3306)/app_auth?charset=utf8mb4&parseTime=True&loc=Local"
)

var (
	db        *gorm.DB
	cacheUser map[string]model.User
	cachMutex sync.Mutex

	publicKey  string
	privateKey *rsa.PrivateKey

	authKey         string = "QmJlMzUxNTAyMzEyNDY4OTc5MTM1OWM0NzExNTIwNDUzNjg3NDMxMTQzMjU3MjY5NTU0NTI1NzI3OTY4NDUyMQ=="
	clientPublicKey string
)

func init() {
	// Make memcache
	cacheUser = make(map[string]model.User)

	// load db connection
	ConnectDatabases()

	isMigrate := flag.Bool("db", false, "Migrate database")
	// isLoadCache := flag.Bool("cache", true, "Load cache")
	flag.Parse()

	// migrate database
	MigrateDatabases(*isMigrate)

	// load cache
	// loadMemoryCache(*isLoadCache)
	loadKeyPemParam()
}

// GetDB return the database connection
func GetDB() *gorm.DB {
	return db
}

func GetCache() *map[string]model.User {
	cachMutex.Lock()
	defer cachMutex.Unlock()
	return &cacheUser
}

func LoadMemoryCache(isLoadCache bool) {
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

func GetAuthKey() string {
	return authKey
}

func GetServerPublicKey() string {
	return publicKey
}

func GetServerPrivateKey() *rsa.PrivateKey {
	return privateKey
}

func SetClientPublicKey(key string) {
	clientPublicKey = key
}
