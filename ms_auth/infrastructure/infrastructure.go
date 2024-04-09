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

	authKey         string = "ACA2547BAA3B4B3E6AD9129ADEE796042FC4ACAAABA40DCE7A1196640595829E2BA85C5E68E5D2B7B713B28720CCD0E65DA3B80CD8A9282BDE895755B924F4CBBCED119DEE057D3066D7095C25771A53D9AD6EC3464EE0E7FE1A4D9851F17B6CC7D6E636E0BDEA7D66153CF37005835A528E65673C1F5ADAA27465964BA89AA6FD045108B25B1DC9AB2D979864858C9EACBFA4A399CDBF2154D2CC6743C5CAA40683343D80A36C7F1289DC5AC5EE585DEA1E1CD296444EA6CD8F6C3564B2D355B310A3C02608EEA5EC36D38E941BF811489AB7A24C069E44FF54714CB4A45EEE2F310FDDC0A1B08BD3DFF8801E27C8508EFD7AD51EF3C13EF1D9A02EE4601741"
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
