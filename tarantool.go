package hey

import (
	"os"
	"time"

	uuid "github.com/satori/go.uuid"
	tarantool "github.com/tarantool/go-tarantool"
)

var (
	DefaultTimeout   time.Duration = 500 * time.Millisecond
	DefaultReconnect time.Duration = 1 * time.Second

	DefaultTarantoolPrefix = "chronograph_"

	threadsSpace    = DefaultTarantoolPrefix + "threads"
	observerSpace   = DefaultTarantoolPrefix + "subscriptions"
	threadLineSpace = DefaultTarantoolPrefix + "threadline"
	eventsSpace     = DefaultTarantoolPrefix + "events"
	usersSpace      = DefaultTarantoolPrefix + "users"
	sourcesSpace    = DefaultTarantoolPrefix + "sources"
)

type TarantoolOpts struct {
	Server        string
	Timeout       time.Duration
	Reconnect     time.Duration
	MaxReconnects uint
	User          string
	Pass          string
}

func (t TarantoolOpts) toOpts() tarantool.Opts {
	return tarantool.Opts{
		Timeout:       t.Timeout,
		Reconnect:     t.Reconnect,
		MaxReconnects: t.MaxReconnects,
		User:          t.User,
		Pass:          t.Pass,
	}
}

func setupFromENV() (opts TarantoolOpts, err error) {
	// 127.0.0.1:3013
	return TarantoolOpts{
		Server:        os.Getenv("TARANTOOL_SERVER"),
		Timeout:       DefaultTimeout,
		Reconnect:     DefaultReconnect,
		MaxReconnects: 3,
		User:          os.Getenv("TARANTOOL_USER_NAME"),
		Pass:          os.Getenv("TARANTOOL_USER_PASSWORD"),
	}, err
}

var _ Manager = &TarantoolManager{}

// TarantoolManager main struct
type TarantoolManager struct {
	conn *tarantool.Connection
}

// NewTarantoolManager return manager setupped from env or optsё
func NewTarantoolManager(opts ...TarantoolOpts) (*TarantoolManager, error) {

	var (
		opt TarantoolOpts
		err error
	)

	if len(opts) > 0 {
		opt = opts[0]
	} else {
		opt, err = setupFromENV()
		if err != nil {
			return nil, err
		}
	}

	client, err := tarantool.Connect(opt.Server, opt.toOpts())
	if err != nil {
		return nil, err
	}

	_, err = client.Ping()

	return &TarantoolManager{conn: client}, err
}

// util
func (m *TarantoolManager) get(space string, keyName string, key interface{}, target interface{}) error {
	return m.conn.SelectTyped(space, keyName, 0, 1, tarantool.IterEq, key, &target)
}

// GetThread return thread by ther id
func (m *TarantoolManager) GetThread(threadID string) (*Thread, error) {
	var (
		thread = new(Thread)
	)
	err := m.conn.SelectTyped(threadsSpace, "primary", 0, 1, tarantool.IterEq,
		makeKey(threadID), thread)
	return thread, err
}

// NewThread create new thread
// it return error is thread already exists
// if id empty, will be generated uuid
func (m *TarantoolManager) NewThread(thread *Thread) (err error) {
	// err = m.conn.SelectTyped(threadsSpace, "primary", 0, 1,
	// 	tarantool.IterEq, makeIndex(threadID), &thread)
	// if err != nil {
	// 	return
	// }
	// // thread already exists
	// if thread.ThreadID == threadID {
	// 	return ErrAlreadyExists
	// }
	if thread.ThreadID == "" {
		thread.ThreadID = uuid.NewV4().String()
	}

	_, err = m.conn.Insert(threadsSpace, thread)
	if err != nil {
		return
	}

	return nil
}

func (m *TarantoolManager) UpdateThread(thread *Thread) error {
	_, err := m.conn.Replace(threadsSpace, thread)
	return err
}

// DeleteThread удаляет трэд, подписки и связанные события
// 1. Удаляем все записи из events
// 2. Удаляем все подписки на трэд
func (m *TarantoolManager) DeleteThread(threadID string) error {
	// delete
	m.conn.Delete(threadsSpace, "primary", makeKey(threadID))
	return nil
}

// AddSource подписывает трэд на другой трэд
func (m *TarantoolManager) AddSource(dstThread, sourceThread string) (err error) {
	src := Sources{SourceThreadID: sourceThread, TargetThreadID: dstThread}
	_, err = m.conn.Insert(sourcesSpace, &src)
	if err != nil {
		return
	}
	return nil
}

func (m *TarantoolManager) GetSources(referThreadID string, offset, limit uint32) (threads []Thread, err error) {
	var sources []Sources
	err = m.conn.SelectTyped(sourcesSpace, "primary", offset, limit, tarantool.IterReq, makeKey(referThreadID), &sources)
	if err != nil {
		return
	}
	for _, src := range sources {
		var thread Thread
		err = m.get(threadsSpace, "primary", makeKey(src.SourceThreadID), &thread)
		if err != nil {
			return
		}
		threads = append(threads, thread)
	}
	return
}

// Трэды, у которого данный трэд в источниках
func (m *TarantoolManager) GetRefers(threadID string, offset, limit uint32) (threads []Thread, err error) {
	var sources []Sources
	err = m.conn.SelectTyped(sourcesSpace, "sources_idx", offset, limit, tarantool.IterReq, makeKey(threadID), &sources)
	if err != nil {
		return
	}
	for _, src := range sources {
		var thread Thread
		err = m.get(threadsSpace, "primary", makeKey(src.SourceThreadID), &thread)
		if err != nil {
			return
		}
		threads = append(threads, thread)
	}
	return
}

// RemoveSource отписывает трэд
func (m *TarantoolManager) RemoveSource(dstThread, sourceThread string) (err error) {
	_, err = m.conn.Delete(sourcesSpace, "primary", makeKey(dstThread, sourceThread))
	if err != nil {
		return err
	}
	return nil
}

// Observe - подписка конкретного юзера на трэд
func (m *TarantoolManager) Observe(userID, threadID string) error {
	obs := Observer{
		UserID:   userID,
		ThreadID: threadID,
	}
	_, err := m.conn.Insert(observerSpace, &obs)
	return err
}

// Ignore удаляет подписку на трэд
func (m *TarantoolManager) Ignore(userID, threadID string) (err error) {

	_, err = m.conn.Delete(observerSpace, "primary", makeKey(threadID, userID))
	if err != nil {
		return err
	}

	// удаляем из трэадлайна все записи с targetThreadID

	return err
}

// Observers возвращает подписчиков трэда
func (m *TarantoolManager) Observers(threadID string, offset, limit uint32) (users []User, err error) {
	var obs []Observer
	err = m.conn.SelectTyped(observerSpace, "primary", offset, limit, tarantool.IterReq, makeKey(threadID), &obs)
	if err != nil {
		return
	}
	for _, o := range obs {
		var u User
		err = m.get(usersSpace, "primary", makeKey(o.UserID), &u)
		if err != nil {
			return
		}
		users = append(users, u)
	}
	return
}

// Observes возвращает подписки юзера
func (m *TarantoolManager) Observes(
	userID string,
	offset,
	limit uint32) (threads []Thread, err error) {
	var obs []Observer
	err = m.conn.SelectTyped(
		observerSpace,
		"subscriptions_idx",
		offset,
		limit,
		tarantool.IterReq,
		makeKey(userID),
		&obs)
	if err != nil {
		return
	}
	for _, o := range obs {
		var thread Thread
		err = m.conn.SelectTyped(threadsSpace,
			"primary",
			0,
			1,
			tarantool.IterEq,
			makeKey(o.ThreadID),
			&thread)
		if err != nil {
			return
		}
		threads = append(threads, thread)
	}
	return
}

func (m *TarantoolManager) MarkAsDelivered(userID, threadID string, times ...time.Time) (err error) {
	t := time.Now()
	if len(times) > 0 {
		t = times[0]
	}

	_, err = m.conn.Update(observerSpace, "primary", makeKey(threadID, userID), makeUpdate(newUpdateOp("=", 3, t)))

	return
}

// RecentActivityByLastTS возвращает события позже lastts
func (m *TarantoolManager) RecentActivityByLastTS(threadID string, limit uint32,
	lastts time.Time) (events []Event, err error) {

	var threadline []Threadline
	err = m.conn.SelectTyped(threadLineSpace, "threadline_idx", limit, 0,
		tarantool.IterGt, makeKey(threadID, lastts), &threadline)
	if err != nil {
		return
	}

	for _, obs := range threadline {
		var event Event
		err = m.conn.SelectTyped(
			eventsSpace,
			"primary",
			0,
			1,
			tarantool.IterEq,
			makeKey(obs.EventID),
			&event)
		if err != nil {
			return
		}
		events = append(events, event)
	}
	return
}

// RecentActivity события двигаться по limit,offset что предлагает tnt
func (m *TarantoolManager) RecentActivity(threadID string, limit, offset uint32) (events []Event, err error) {
	var threadline []Threadline
	err = m.conn.SelectTyped(threadLineSpace, "threadline_idx", limit, offset, tarantool.IterReq, makeKey(threadID), &threadline)
	if err != nil {
		return
	}
	for _, obs := range threadline {
		var event Event
		err = m.conn.SelectTyped(eventsSpace, "primary", 0, 1, tarantool.IterEq, makeKey(obs.EventID), &event)
		if err != nil {
			return
		}
		events = append(events, event)
	}
	return
}

// events
// todo
// 1. Достаём всех подписчиков трэда.
// 2. Вставляем в Timeline для всех подписчиков этот eventID
func (m *TarantoolManager) NewEvent(ev *Event) (err error) {
	if ev.EventID == "" {
		ev.EventID = uuid.NewV4().String()
	}

	_, err = m.conn.Insert(eventsSpace, ev)
	if err != nil {
		return
	}

	return nil
}

func (m *TarantoolManager) GetEvent(eventID string) (ev *Event, err error) {
	ev = new(Event)
	err = m.conn.SelectTyped(eventsSpace, "primary", 0, 1, tarantool.IterEq,
		makeKey(eventID), ev)
	return
}

// GetEvents returns events by their ids
func (m *TarantoolManager) GetEvents(ids ...string) (events []Event, err error) {
	for _, id := range ids {
		var event Event
		err = m.conn.SelectTyped(eventsSpace, "primary", 0, 1, tarantool.IterEq,
			makeKey(id), &event)
		if err != nil {
			return
		}
		events = append(events, event)
	}
	return
}

func (m *TarantoolManager) UpdateEvent(ev *Event) (err error) {
	// todo обновить updatedat связаных трэдов
	ev.UpdatedAt = time.Now()
	_, err = m.conn.Replace(eventsSpace, ev)
	if err != nil {
		return
	}
	return nil
}

func (m *TarantoolManager) DeleteEvent(eventID string) (err error) {
	// todo удалить из threadline связанные записи
	_, err = m.conn.Delete(eventsSpace, "primary", makeKey(eventID))
	if err != nil {
		return
	}
	return nil
}
