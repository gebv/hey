package store

import (
	"models"
	// "strings"
	"github.com/satori/go.uuid"
)

func init() {
	registrationOfStoreBuilder("event", func(sm *StoreManager) Store {
		return NewEventStore(sm)
	})
}

type EventStore struct {
	*StoreManager
}

func NewEventStore(_store *StoreManager) *EventStore {

	return &EventStore{_store}
}

func (_manager EventStore) ErrorLog(args ...interface{}) {
	_manager.StoreManager.ErrorLog(_manager.Name(), args...)
}

func (_manager EventStore) Name() string {
	return "event"
}

func (s *EventStore) GetOne(dto *models.EventDTO) (models.ModelAbstractInterface, error) {
	var err error
	if err := models.V.StructPartial(dto, "EventId", "ClientId"); err != nil {
		return nil, err
	}

	model := models.NewEvent()
	model.TransformFrom(dto)

	fields, modelValues := model.Fields()

	query := SqlSelect(model.TableName(), fields)
	query += " WHERE client_id = ? AND event_id = ? AND is_removed = false"
	query = FormateToPQuery(query)

	if dto.Tx != nil {
		err = dto.Tx.QueryRow(query, model.ClientId, model.EventId).Scan(modelValues...)
	} else {
		err = s.db.QueryRow(query, model.ClientId, model.EventId).Scan(modelValues...)
	}

	if err != nil {
		return nil, err
	}
	
	return model, nil
}

func (s *EventStore) Create(dto *models.EventDTO) (models.ModelAbstractInterface, error) {
	if err := models.V.Struct(dto); err != nil {
		return nil, err
	}

	tx, err := s.db.Begin()

	if err != nil {
		s.ErrorLog("action", "создание события", "subaction", "создание транзакции", "err", err.Error())
		return nil, models.ErrUnknown
	}

	// Пользователь
	dtoUser := models.NewUserDTO()
	dto.TransformTo(dtoUser)
	dtoUser.Tx = tx
	user, err := s.StoreManager.Get("user").(*UserStore).GetOne(dtoUser)

	if err != nil {
		tx.Rollback()
		s.ErrorLog("action", "cоздание события", "subaction", "поиск пользователья-создателя", "err", err.Error())
		return nil, err
	}

	if !user.(*models.User).IsEnabled {
		tx.Rollback()
		s.ErrorLog("action", "cоздание события", "msg", "пользователь не активный", "err", models.ErrNotAllowed)
		return nil, models.ErrNotAllowed
	}

	// Канал
	dtoChannel := models.NewChannelDTO()
	dto.TransformTo(dtoChannel)
	dtoChannel.Tx = tx

	channel, err :=  s.StoreManager.Get("channel").(*ChannelStore).GetOne(dtoChannel)

	if err != nil {
		tx.Rollback()
		s.ErrorLog("action", "создание события", "subaction", "поиск связанного канала", "err", err.Error())
		return nil, err
	}

	// Цепочка потоков
	dtoThread := models.NewThreadDTO()
	dtoThread.ChannelId = channel.PrimaryValue().String()
	dto.TransformTo(dtoThread)
	dtoThread.Tx = tx
	_threads, err := s.StoreManager.Get("thread").(*ThreadStore).FindAllThreadsFromPathAndCreateLast(dtoThread)

	if err != nil {
		tx.Rollback()
		s.ErrorLog("action", "создание события", "subaction", "поиск связанных потоков по цепочке", "err", err.Error())
		return nil, err
	}
	threads := _threads.([]*models.Thread)

	model := models.NewEvent()
	model.TransformFrom(dto)
	model.EventId = uuid.NewV1()
	model.ThreadId = threads[len(threads)-1].ThreadId
	model.ChannelId = channel.PrimaryValue()
	model.ParentEventId = uuid.FromStringOrNil(dto.ParentEventId)
	model.Creator = user.PrimaryValue()

	if len(threads) > 1 {
		model.ParentThreadId = threads[len(threads)-2].ThreadId	
	}

	if model.ParentEventId != uuid.Nil {
		// TODO: вынести в отдельный поток
		parentEvent := models.NewEvent()
		parentEvent.EventId = model.ParentEventId
		parentEvent.BranchThreadId = model.ThreadId

		if err := UpdateModel(parentEvent, s.db, tx, "branch_thread_id"); err != nil {
			tx.Rollback()
			s.ErrorLog("action", "создание события", "subaction", "обновление у ParentEventId BranchThreadId", "err", err.Error())
			return nil, err
		}

		// TODO: вынести в отдельный поток
		branchThread := models.NewThread()
		branchThread.ThreadId = model.ThreadId
		branchThread.RelatedEventId = parentEvent.EventId

		if err := UpdateModel(branchThread, s.db, tx, "related_event_id"); err != nil {
			tx.Rollback()
			s.ErrorLog("action", "создание события", "subaction", "обновление у BranchThreadId RelatedEventId", "err", err.Error())
			return nil, err
		}
	}

	fields, _ := model.Fields()

	if err := CreateModel(model, s.db, tx, fields...); err != nil {
		tx.Rollback()
		return nil, err
	}

	// Threadline

	// TODO: вынести в отдельный поток
	threadline := models.NewThreadline()
	fields, _ = threadline.Fields()
	threadline.ClientId = model.ClientId
	threadline.ThreadId = model.ThreadId
	threadline.ChannelId = model.ChannelId
	threadline.EventId = model.EventId

	if err := CreateModel(threadline, s.db, tx, fields...); err != nil {
		tx.Rollback()
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return model, nil
}