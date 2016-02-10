package store

import (
	"models"
	"strings"
	"github.com/satori/go.uuid"
	"github.com/jackc/pgx"
	// "github.com/golang/glog"
	"time"
	"strconv"
	"fmt"
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

func (s *EventStore) LoadThreadline(dto *models.LoadThreadlineDTO) (*models.EventLoadResult, error) {
	if err := models.V.Struct(dto); err != nil {
		return nil, err
	}

	var queryWhere string
	whereFirstPage := `WHERE client_id = ? AND channel_id = ? AND thread_id = ?
					ORDER BY created_at DESC, event_id LIMIT ?`
	queryWhere = whereFirstPage

	whereNextPage := `WHERE client_id = ? AND channel_id = ? AND thread_id = ?
							AND (created_at, event_id) < (?, ?)
						ORDER BY created_at DESC, event_id LIMIT ?`

	var lastEventId uuid.UUID
	var lastTime time.Time

	var user *models.User
	var thread *models.Thread
	var args = []interface{}{dto.ClientId}

	var result = models.NewEventLoadResult()

	// Загрузка пользователя
	_userDto := models.NewUserDTO()
	_userDto.ClientId = dto.ClientId
	_userDto.ExtId = dto.User
	_userDto.Tx = dto.Tx
	_user, err := s.StoreManager.Get("user").(*UserStore).GetOne(_userDto)

	if err != nil {
		s.ErrorLog(
			"action", "загрузка сообщений из потока", 
			"subaction", "загрузка связанного пользователя", 
			"thread", dto.Thread, 
			"user", dto.User, 
			"cid", dto.ClientId, 
			"err", err)
		return nil, err
	}

	user = _user.(*models.User)

	if !user.IsEnabled {
		s.ErrorLog(
			"action", "загрузка сообщений из потока", 
			"subaction", "связанный пользователь", 
			"thread", dto.Thread, 
			"user", dto.User, 
			"cid", dto.ClientId, 
			"err", "user not enabled")
		return nil, models.ErrNotAllowed
	}

	// Загрузка потока
	_threadDTO := models.NewThreadDTO()
	_threadDTO.ClientId = dto.ClientId
	_threadDTO.ExtId = dto.Thread
	_threadDTO.Tx = dto.Tx
	_thread, err := s.StoreManager.Get("thread").(*ThreadStore).GetOne(_threadDTO)

	if err != nil {
		s.ErrorLog(
			"action", "загрузка сообщений из потока", 
			"subaction", "загрузка связанного потока", 
			"thread", dto.Thread, 
			"user", dto.User, 
			"cid", dto.ClientId, 
			"err", err)
		return nil, err
	}

	thread = _thread.(*models.Thread)
	args = append(args, thread.ChannelId, thread.ThreadId)

	//TODO: Что проверять? 

	// Анализ курсора
	if len(strings.TrimSpace(dto.Cursor)) > 0 {
		// uuid+time
		_params := strings.Split(strings.TrimSpace(dto.Cursor), ",")
		if len(_params) != 2 {
goto skipcursor
		} 
		
		lastEventId = uuid.FromStringOrNil(_params[0])

		if uuid.Equal(lastEventId, uuid.Nil) {
goto skipcursor
		}

		timestamp, err := strconv.ParseInt(_params[1], 10, 64)

		if err != nil {
goto skipcursor	
		}

		lastTime = time.Unix(0, timestamp)

		queryWhere = whereNextPage
		args = append(args, lastTime, lastEventId)
	}

skipcursor: 

	args = append(args, dto.Limit+1)
	
	var threadline = models.NewThreadline()
	var threadlines []models.Threadline

	fields, threadlineFields := threadline.Fields("event_id", "created_at")
	query := SqlSelect(threadline.TableName(), fields) + " " + queryWhere
	query = FormateToPQuery(query)

	var rows *pgx.Rows
	var eventCounter = 0

	if dto.Tx != nil {
		rows, err = dto.Tx.Query(query, args...)
	} else {
		rows, err = s.db.Query(query, args...)
	}

	defer rows.Close()

	for rows.Next() {
		if err := rows.Scan(threadlineFields...); err != nil {
			s.ErrorLog("action", "загрузка событий из threadline",
				"subaction", "scan threadline",
				"err", err,
				"thread", dto.Thread, 
				"user", dto.User, 
				"cid", dto.ClientId, 
				)
			return nil, err
		}

		eventCounter++
		threadlines = append(threadlines, *threadline)
	}

	rows.Close()

	for _, _threadline := range threadlines {
		event := models.NewEvent()
		event.ClientId = thread.ClientId
		event.EventId = _threadline.EventId
		if err := FindModel(event, s.db, dto.Tx); err != nil {
			s.ErrorLog("action", "загрузка событий из threadline",
				"subaction", "загрузка связанного события",
				"event_id", _threadline.EventId.String(),
				"err", err,
				"thread", dto.Thread, 
				"user", dto.User, 
				"cid", dto.ClientId, 
				)
			return nil, err
		}

		// TODO: Виртуальные поля для события
		// TODO: Преобразовать событие в спец объект

		result.Events = append(result.Events, event)
	}

	// Формирование курсора

	if eventCounter == dto.Limit+1 {
		result.Events = result.Events[:len(result.Events)-1]
		lastEvent := result.Events[len(result.Events)-1]
		result.Cursor = fmt.Sprintf("%v,%v", lastEvent.EventId.String(), lastEvent.CreatedAt.UnixNano())
		result.HasNext = true
	}

	return result, nil
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

	defer tx.Rollback()

	if err != nil {
		s.ErrorLog("action", "cоздание события", "subaction", "поиск пользователья-создателя", "err", err.Error())
		return nil, err
	}

	if !user.(*models.User).IsEnabled {
		s.ErrorLog("action", "cоздание события", "msg", "пользователь не активный", "err", models.ErrNotAllowed)
		return nil, models.ErrNotAllowed
	}

	// Канал
	dtoChannel := models.NewChannelDTO()
	dto.TransformTo(dtoChannel)
	dtoChannel.Tx = tx

	channel, err :=  s.StoreManager.Get("channel").(*ChannelStore).GetOne(dtoChannel)

	if err != nil {
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
			s.ErrorLog("action", "создание события", "subaction", "обновление у ParentEventId BranchThreadId", "err", err.Error())
			return nil, err
		}

		// TODO: вынести в отдельный поток       
		branchThread := models.NewThread()
		branchThread.ThreadId = model.ThreadId
		branchThread.RelatedEventId = parentEvent.EventId

		if err := UpdateModel(branchThread, s.db, tx, "related_event_id"); err != nil {
			s.ErrorLog("action", "создание события", "subaction", "обновление у BranchThreadId RelatedEventId", "err", err.Error())
			return nil, err
		}
	}

	fields, _ := model.Fields()

	if err := CreateModel(model, s.db, tx, fields...); err != nil {
		return nil, err
	}

	// Threadline

	// TODO: вынести в отдельный поток, после проверки сообщение становится доступным публике
	threadline := models.NewThreadline()
	fields, _ = threadline.Fields()
	threadline.ClientId = model.ClientId
	threadline.ThreadId = model.ThreadId
	threadline.ChannelId = model.ChannelId
	threadline.EventId = model.EventId
	threadline.CreatedAt = model.CreatedAt

	if err := CreateModel(threadline, s.db, tx, fields...); err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return model, nil
}