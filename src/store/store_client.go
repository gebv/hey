package store

import (
	"models"
)

func init() {
	registrationOfStoreBuilder("client", func(sm *StoreManager) Store {
		return NewClientStore(sm)
	})
}

type ClientStore struct {
	*StoreManager
}

func NewClientStore(_store *StoreManager) *ClientStore {

	return &ClientStore{_store}
}

func (_manager ClientStore) ErrorLog(args ...interface{}) {
	_manager.StoreManager.ErrorLog(_manager.Name(), args...)
}

func (_manager ClientStore) Name() string {
	return "client"
}

func (_manager ClientStore) Get(modelId string) (models.ModelAbstractInterface, error) {
	_model := models.NewClient()
	// _model.ClientId = modelId

	// fields, _ := _model.Fields()

	// if err := FindModel(_model, _manager.db, nil, fields...); err != nil {
	// 	return nil, err
	// }

	return _model, nil
}

// CreateAnonymUser создание пользователя на основе email
func (_manager ClientStore) CreateAnonymUser(email, target string) (models.ModelAbstractInterface, error) {
	_model := models.NewClient()
	// _model.ClientId = models.NewUUID().String()

	// fields, _ := _model.Fields("user_id", "name", "email", "licenses", "sys_tags")

	// if err := CreateModel(_model, db, nil, fields...); err != nil {
	// 	return nil, err
	// }

	return _model, nil
}

// GetByEmail Найти пользователя по email
func (_manager ClientStore) GetByEmail(email string) (models.ModelAbstractInterface, error) {
	_model := models.NewClient()

	// fields, _ := _model.Fields()
	// query := SqlSelect(_model.TableName(), fields)
	// query = FormateToPG(query, fields)
	// query += " WHERE email = ?"
	// args := []interface{}{email}

	// if err := SelectOne(_model, db, nil, query, args); err != nil {
	// 	return nil, err
	// }

	return _model, nil
}
