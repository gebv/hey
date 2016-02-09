package store

import (
	// "testing"
	"utils"
	// "models"
	// "time"
	"flag"
	// "github.com/golang/glog"
	"github.com/jackc/pgx"
)

var _s *StoreManager
var _conn *pgx.Conn

func init() {

	// flag.Parse()
	flag.Set("v", "1")
	flag.Set("stderrthreshold", "ERROR")

	utils.IsTesting = true
	utils.LoadConfig("../../config/config.json")
	_s = NewStore()
	
	_conn, _ = setupConnectionPGX(utils.Cfg.StorageAppSettings)
}

// func zzzTestCreate(t *testing.T) {
// 	model := models.NewClient()
// 	model.Id = models.NewUUID().String()
// 	model.Domain = "http://domain" + models.NewUUID().String()
// 	model.Ips.Add("127.0.0.1")
// 	model.Scopes.Add("full")

// 	fieldNames, _ := model.Fields()

// 	if err := CreateModel(model, _s.db, nil, fieldNames...); err != nil {
// 		t.Error(err)
// 	}

// 	time.Sleep(time.Second*1)

// 	model.Scopes.Add("tree")

// 	if err := UpdateModel(model, _s.db, nil, "scopes"); err != nil {
// 		t.Error(err)
// 	}

// 	time.Sleep(time.Second*1)

// 	if err := DeleteModel(model, _s.db, nil); err != nil {
// 		t.Error(err)	
// 	}
// }
