package store

import (
	"models"
	// "strings"
	// "gopkg.in/pg.v3"
	"github.com/satori/go.uuid"
)

func init() {
	registrationOfStoreBuilder("channel", func(sm *StoreManager) Store {
		return NewChannelStore(sm)
	})
}

type ChannelStore struct {
	*StoreManager
}

func NewChannelStore(_store *StoreManager) *ChannelStore {

	return &ChannelStore{_store}
}

func (_manager ChannelStore) ErrorLog(args ...interface{}) {
	_manager.StoreManager.ErrorLog(_manager.Name(), args...)
}

func (_manager ChannelStore) Name() string {
	return "channel"
}

func (s *ChannelStore) Create(dto *models.ChannelDTO) (models.ModelAbstractInterface, error) {
	if err := models.V.StructPartial(dto, "ExtId", "ClientId", "ChannelId"); err != nil {
		return nil, err
	}

	model := models.NewChannel()
	model.TransformFrom(dto)
	model.ChannelId = uuid.NewV1()



	// ThreadRoot
	threadDto := models.NewThreadDTO()
	threadDto.TransformFrom(model)
	threadDto.Tx = dto.Tx
	threadRoot, err := s.StoreManager.Get("thread").(*ThreadStore).Create(threadDto)

	if err != nil {
		return nil, err
	}

	model.RootThreadId = threadRoot.PrimaryValue()

	fields, _ := model.Fields()

	_fields := models.StringArray(fields)

	if err := CreateModel(model, s.db, dto.Tx, _fields...); err != nil {
		return nil, err
	}

	// Создать channel_counter
	tCounter := models.NewChannelCounter();
	tCounter.ClientId = model.ClientId
	tCounter.ChannelId  =model.ChannelId
	tCounter.CounterEvents = 0

	tCounterFields, tCounterValues := tCounter.Fields()
	query := SqlInsert(tCounter.TableName(), tCounterFields)
	query = FormateToPQuery(query)

	if dto.Tx != nil {
		_, err = dto.Tx.Exec(query, tCounterValues...)
	} else {
		_, err = s.db.Exec(query, tCounterValues...)
	}

	if err != nil {
		return nil, err
	}

	// Создать channel_watchers
	tWatcher := models.NewChannelWatcher()
	tWatcher.ClientId = model.ClientId
	tWatcher.ChannelId = model.ChannelId
	tWatcher.Unread = 0

	tWatcherFields, tWatcherValues := tWatcher.Fields()
	query = SqlInsert(tWatcher.TableName(), tWatcherFields)
	query = FormateToPQuery(query)

	for _, owner := range model.Owners {
		tWatcher.UserId = owner

		if dto.Tx != nil {
			_, err = dto.Tx.Exec(query, tWatcherValues...)
		} else {
			_, err = s.db.Exec(query, tWatcherValues...)
		}

		if err != nil {
			s.ErrorLog("err", err, "sql", query, "args", tWatcherValues)
			return nil, err
		}
	}

	return model, nil
}

func (s *ChannelStore) GetOne(dto *models.ChannelDTO) (models.ModelAbstractInterface, error) {
	var err error
	if err := models.V.StructPartial(dto, "ExtId", "ClientId"); err != nil {
		return nil, err
	}

	model := models.NewChannel()
	model.TransformFrom(dto)

	fields, modelValues := model.Fields()

	query := SqlSelect(model.TableName(), fields)
	query += " WHERE client_id = ? AND ext_id_hash = ? AND is_removed = false"
	query = FormateToPQuery(query)

	if dto.Tx != nil {
		err = dto.Tx.QueryRow(query, model.ClientId, model.ExtIdHash).Scan(modelValues...)
	} else {
		err = s.db.QueryRow(query, model.ClientId, model.ExtIdHash).Scan(modelValues...)
	}

	if err != nil {
		return nil, err
	}
	
	return model, nil
}