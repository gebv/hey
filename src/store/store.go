package store

import (
	"database/sql"
	"github.com/golang/glog"

	_ "github.com/lib/pq"
	// "gopkg.in/pg.v3"

	"models"
	"utils"

	"fmt"
	"github.com/RangelReale/osin"
	"github.com/ory-am/osin-storage/storage/postgres"
	"strings"
	"github.com/jackc/pgx"
)

// var db *sql.DB
var db *pgx.Conn

var registredStores = make(map[string]StoreBuilder)

func registrationOfStoreBuilder(name string, builder StoreBuilder) {
	registredStores[name] = builder
}

type StoreBuilder func(*StoreManager) Store

func NewStore() *StoreManager {
	_sm := &StoreManager{}
	var err error
	var _db *sql.DB
	if _db, err = setupConnection(utils.Cfg.StorageAppSettings); err != nil {
		panic("error setup connection, err=" + err.Error())
	}

	_sm.db, err = setupConnectionPGX(utils.Cfg.StorageAppSettings)
	if err != nil {
		panic("error setup connection, err="+err.Error())
	}
	db = _sm.db

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
	osinStorage := postgres.New(_db)
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
	// db   *sql.DB
	// db   *pg.DB
	db   *pgx.Conn
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

func setupConnectionPGX(c models.StorageSettings) (*pgx.Conn, error) {
	return pgx.Connect(extractPGXStorageConfig(c))
}

type databaseLogger struct {}

func (l databaseLogger) Debug(msg string, ctx ...interface{}) {
	glog.Infof("\tSQL[debug]: msg='%s', %v", msg, ctx)
}

func (l databaseLogger) Info(msg string, ctx ...interface{}) {
	glog.Infof("\tSQL: msg='%s', %v", msg, ctx)
}

func (l databaseLogger) Warn(msg string, ctx ...interface{}) {
	glog.Warningf("\tSQL] msg='%s', %v", msg, ctx)
}

func (l databaseLogger) Error(msg string, ctx ...interface{}) {
	glog.Errorf("\tSQL] msg='%s', %v", msg, ctx)
}

func extractPGXStorageConfig(c models.StorageSettings) pgx.ConnConfig {
	var config pgx.ConnConfig

	config.Host = c.Host
	config.User = c.User
	config.Password = c.Password
	config.Database = c.Database
	config.Logger = databaseLogger{}
	config.LogLevel = pgx.LogLevelDebug

	return config
}

// func setupConnectionPg(c models.StorageSettings) (*pg.DB) {

// 	return pg.Connect(&pg.Options{
// 		Host:     c.Host,
// 		Database: c.Database,
// 		User:     c.User,
// 		Password: c.Password,
// 		SSL:      false,
// 	})
// }

func setupConnection(c models.StorageSettings) (*sql.DB, error) {

	addrs := fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable", c.User, c.Password, c.Host, c.Database)

	return sql.Open("postgres", addrs)
}
