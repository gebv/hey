package store

import (
	"models"
	"github.com/satori/go.uuid"
)

func init() {
	registrationOfStoreBuilder("user", func(sm *StoreManager) Store {
		return NewUserStore(sm)
	})
}

type UserStore struct {
	*StoreManager
}

func NewUserStore(_store *StoreManager) *UserStore {

	return &UserStore{_store}
}

func (_manager UserStore) ErrorLog(args ...interface{}) {
	_manager.StoreManager.ErrorLog(_manager.Name(), args...)
}

func (_manager UserStore) Name() string {
	return "user"
}

func (s *UserStore) GetOne(dto *models.UserDTO) (models.ModelAbstractInterface, error) {
	if err := models.V.StructPartial(dto, "ExtId", "ClientId"); err != nil {
		return nil, err
	}

	model := models.NewUser()
	model.TransformFrom(dto)

	fields, modelValues := model.Fields()

	query := SqlSelect(model.TableName(), fields)
	query += " WHERE client_id = ? AND ext_id_hash = ? AND is_removed = false LIMIT 1"
	query = FormateToPQuery(query)

	if dto.Tx != nil {
		err := dto.Tx.QueryRow(query, model.ClientId, model.ExtIdHash).Scan(modelValues...)
		if err != nil {
			return nil, err
		}
	} else {
		err := s.db.QueryRow(query, model.ClientId, model.ExtIdHash).Scan(modelValues...)
		if err != nil {
			return nil, err
		}
	}
	
	return model, nil
}

func (s *UserStore) Create(dto *models.UserDTO) (models.ModelAbstractInterface, error) {
	if err := models.V.StructPartial(dto, "ExtId", "ClientId"); err != nil {
		return nil, err
	}

	model := models.NewUser()
	model.TransformFrom(dto)
	model.UserId = uuid.NewV1()

	fields, _ := model.Fields()

	if err := CreateModel(model, s.db, nil, fields...); err != nil {
		return nil, err
	}

	return model, nil
}

func (s *UserStore) Delete(dto *models.UserDTO) (error) {
	if err := models.V.StructPartial(dto, "ExtId", "ClientId"); err != nil {
		return err
	}

	model, err := s.GetOne(dto)

	if err != nil {

		return err
	}

	model.BeforeDelete()
	model.BeforeSave()

	return UpdateModel(model, s.db, nil, "update_at", "is_removed")
}

func (s *UserStore) Update(dto *models.UserDTO) (models.ModelAbstractInterface, error) {
	var err error
	if err := models.V.StructPartial(dto, "ExtId", "ClientId"); err != nil {
		return nil, err
	}
	model := models.NewUser()
	model.TransformFrom(dto)
	model.BeforeSave()

	fieldss, modelValues := model.Fields("is_enabled", "ext_props", "updated_at")

	query := SqlUpdate(model.TableName(), fieldss)
	query += " WHERE client_id = ? AND ext_id_hash = ? RETURNING user_id"
	query = FormateToPQuery(query)

	if dto.Tx != nil {
		err = dto.Tx.QueryRow(query, append(modelValues, model.ClientId, model.ExtIdHash)...).Scan(&model.UserId)
	} else {
		err = s.db.QueryRow(query, append(modelValues, model.ClientId, model.ExtIdHash)...).Scan(&model.UserId)
	}

	if err != nil {
		return nil, err
	}
	
	return model, nil
}
