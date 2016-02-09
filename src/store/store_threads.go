package store

import (
	"models"
	"strings"
	// "fmt"
	"github.com/jackc/pgx"
	"github.com/satori/go.uuid"
)

func init() {
	registrationOfStoreBuilder("thread", func(sm *StoreManager) Store {
		return NewThreadStore(sm)
	})
}

type ThreadStore struct {
	*StoreManager
}

func NewThreadStore(_store *StoreManager) *ThreadStore {

	return &ThreadStore{_store}
}

func (_manager ThreadStore) ErrorLog(args ...interface{}) {
	_manager.StoreManager.ErrorLog(_manager.Name(), args...)
}

func (_manager ThreadStore) Name() string {
	return "thread"
}

func (s *ThreadStore) GetOne(dto *models.ThreadDTO) (models.ModelAbstractInterface, error) {
	var err error
	if err := models.V.StructPartial(dto, "ExtId", "ClientId"); err != nil {
		return nil, err
	}

	model := models.NewThread()
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

func (s *ThreadStore) Create(dto *models.ThreadDTO) (models.ModelAbstractInterface, error) {
	var err error
	if err = models.V.StructPartial(dto, "ExtId", "ClientId", "ChannelId"); err != nil {
		return nil, err
	}

	// Создание модели
	model := models.NewThread()
	model.TransformFrom(dto)
	model.ThreadId = uuid.NewV1()

	fields, _ := model.Fields()

	_fields := models.StringArray(fields)

	// if len(model.ParentThreadId) == 0 {
	// 	_fields.Del("parent_thread_id")
	// }

	// if len(model.RelatedEventId) == 0 {
	// 	_fields.Del("related_event_id")
	// }

	// _fields.Del("owners")

	if err := CreateModel(model, s.db, dto.Tx, _fields...); err != nil {
		return nil, err
	}

	// Создать thread_counter
	tCounter := models.NewThreadCounter();
	tCounter.ClientId = model.ClientId
	tCounter.ThreadId = model.ThreadId
	tCounter.CounterEvents = 0
	tCounterFields, tCounterValues := tCounter.Fields() //"client_id", "thread_id", "counter_events"
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

	// Создать thread_watchers
	tWatcher := models.NewThreadWatcher()
	tWatcher.ClientId = model.ClientId
	tWatcher.ThreadId = model.ThreadId
	tWatcher.Unread = 0

	tWatcherFields, tWatcherValues := tWatcher.Fields() // "client_id", "thread_id", "user_id", "unread"
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

// func (s *ThreadStore) UpdateThreadWatchers(dto *models.ThreadDTO) error {
// 	// client_id uuid,
// 	// thread_id uuid,

// 	// counter_events int8,
// }

// func (s *ThreadStore) UpdateThreadCounter(dto *models.ThreadDTO) error {
// 	// client_id uuid,
// 	// thread_id uuid,
// 	// user_id uuid,
	
// 	// unread int4,
// }

// Вернуть все потоки. В случае если последний отсутствует, создать его
func (s *ThreadStore) FindAllThreadsFromPathAndCreateLast(dto *models.ThreadDTO) (interface{}, error) {
	var currentThread models.ModelAbstractInterface
	currentThread = models.NewThread()

	fields, modelValues := currentThread.Fields()
	// _fields := models.StringArray(fields)
	// _fields.Del("owners")
	// _fields.Del("thread_id")
	query := SqlSelect(currentThread.TableName(), fields)
	query += " WHERE client_id = ? AND ext_id = ?"
	query = FormateToPQuery(query)

	var threads []*models.Thread
	var parentThread *models.Thread
	var path string
	var err error

	var threadNodes = strings.Split(dto.ExtId, ":")

	for index, _path := range threadNodes {
		if index == 0 {
			path += _path
		} else {
			path += ":"+_path
		}

		if dto.Tx != nil {
			err = dto.Tx.QueryRow(query, dto.ClientId, path).Scan(modelValues...)
		} else {
			err = s.db.QueryRow(query, dto.ClientId, path).Scan(modelValues...)
		}

		if err == pgx.ErrNoRows && index == len(threadNodes)-1 {
			// Если предпоследний отсутствует, создаем

			if parentThread != nil {
				// Если не корневой
				dto.ParentThreadId = parentThread.ThreadId.String()
			}

			dto.RelatedEventId = dto.CreatedEventId
			dto.Owners.FromArray(parentThread.Owners) // Владелец потока остается быть прежним
			
			if currentThread, err = s.Create(dto); err != nil {
				s.ErrorLog("action", "поиск потоков", "subaction", "создание по горячему поток", "err", err.Error())
				return nil, models.ErrUnknown
			}
		} else if err != nil {
			s.ErrorLog("action", "поиск и создание потока на основе пути, потомок отсутствует", "ext_id", path, "err", err, "index", index, "query", query)
			return nil, models.ErrNotValid
		} else if err == nil && index == len(threadNodes)-2 {
			if uuid.FromStringOrNil(dto.CreatedEventId) != uuid.Nil {
				currentThread.(*models.Thread).RelatedEventId = 	uuid.FromStringOrNil(dto.CreatedEventId)
				if err := UpdateModel(currentThread, s.db, dto.Tx, "related_event_id") ; err != nil {
					s.ErrorLog("action", "обновление RelatedEvent", "ext_id", path, "err", err, "index", index)
					return nil, err;
				}
			}	
		}

		threads = append(threads, currentThread.(*models.Thread))
		parentThread = currentThread.(*models.Thread)
	}

	return threads, nil
}