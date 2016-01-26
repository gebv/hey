package store

import (
	"database/sql"
	"github.com/golang/glog"

	_ "github.com/lib/pq"

	"models"
	"utils"

	"fmt"
	"github.com/RangelReale/osin"
	"github.com/ory-am/osin-storage/storage/postgres"
	"strings"
)

var db *sql.DB

var registredStores = make(map[string]StoreBuilder)

func registrationOfStoreBuilder(name string, builder StoreBuilder) {
	registredStores[name] = builder
}

type StoreBuilder func(*StoreManager) Store

func NewStore() *StoreManager {
	_sm := &StoreManager{}
	var err error
	if _sm.db, err = setupConnection(utils.Cfg.StorageSettings); err != nil {
		panic("error setup connection, err=" + err.Error())
	}

	_sm.stores = make(map[string]Store)

	for _name, _builder := range registredStores {
		_sm.stores[_name] = _builder(_sm)
	}

	config := osin.NewServerConfig()
	// config.AuthorizationExpiration = 5
	// config.AccessExpiration = 30
	config.ErrorStatusCode = 400
	config.AllowGetAccessRequest = false
	// config.AllowedAuthorizeTypes = osin.AllowedAuthorizeType{osin.CODE, osin.TOKEN}
	config.AllowedAuthorizeTypes = osin.AllowedAuthorizeType{osin.CODE}
	config.AllowedAccessTypes = osin.AllowedAccessType{osin.AUTHORIZATION_CODE, osin.REFRESH_TOKEN}

	// TODO: reserve tables client, authorize, access, refresh
	osinStorage := postgres.New(_sm.db)
	if err := osinStorage.CreateSchemas(); err != nil {
		panic("osin storage: err="+err.Error())
	}
	_sm.Osin = osin.NewServer(config, osinStorage)

	return _sm
}

type Store interface {
	Name() string
}

type StoreManager struct {
	db   *sql.DB
	Osin *osin.Server

	stores map[string]Store
}

// Get получить Store по имени
func (sm *StoreManager) Get(name string) Store {
	return sm.stores[name]
}

func (sm StoreManager) ErrorLog(prefix string, args ...interface{}) {
	if len(args)%2 == 0 {
		glog.Errorf(prefix+": "+strings.Repeat("%v='%v', ", len(args)/2), args...)
		return
	}

	glog.Errorf(prefix+": "+strings.Repeat("%v, ", len(args)), args...)
}

func setupConnection(c models.StorageSettings) (*sql.DB, error) {

	addrs := fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable", c.User, c.Password, c.Host, c.Database)

	return sql.Open("postgres", addrs)
}
