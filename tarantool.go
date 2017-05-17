package hey

import (
	"errors"
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
	relatedSpace    = DefaultTarantoolPrefix + "related"

	ErrNotFound = errors.New("not_found")
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

func NewTarantoolManager(conn *tarantool.Connection) *TarantoolManager {
	return &TarantoolManager{conn: conn}
}

// // NewTarantoolManager return manager setupped from env or optsё
// func NewTarantoolManager(opts ...TarantoolOpts) (*TarantoolManager, error) {

// 	var (
// 		opt TarantoolOpts
// 		err error
// 	)

// 	if len(opts) > 0 {
// 		opt = opts[0]
// 	} else {
// 		opt, err = setupFromENV()
// 		if err != nil {
// 			return nil, err
// 		}
// 	}

// 	client, err := tarantool.Connect(opt.Server, opt.toOpts())
// 	if err != nil {
// 		return nil, err
// 	}

// 	_, err = client.Ping()

// 	return &TarantoolManager{conn: client}, err
// }

// util
func (m *TarantoolManager) get(space string, keyName string, key interface{}, target interface{}) error {
	return m.conn.SelectTyped(space, keyName, 0, 1, tarantool.IterEq, key, &target)
}

// GetThread return thread by ther id
func (m *TarantoolManager) GetThread(threadID string) (*Thread, error) {
	var (
		threads []Thread
	)
	err := m.conn.SelectTyped(threadsSpace, "primary", 0, 1, tarantool.IterEq,
		makeKey(threadID), &threads)
	if err != nil {
		return nil, err
	}
	if len(threads) == 0 {
		return nil, ErrNotFound
	}
	return &threads[0], nil
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
		UserID:            userID,
		ThreadID:          threadID,
		LastDeliveredTime: time.Now(),
		JoinTime:          time.Now(),
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
		var u []User
		err = m.get(usersSpace, "primary", makeKey(o.UserID), &u)
		if err != nil {
			return
		}
		users = append(users, u[0])
	}
	return
}

// Observes возвращает подписки юзера
func (m *TarantoolManager) Observes(
	userID string,
	offset,
	limit uint32) (threads []Thread, err error) {
	var obs []Observer
	err = m.conn.SelectTyped(observerSpace, "subs_thread_id_idx", offset, limit,
		tarantool.IterReq, makeKey(userID), &obs)
	if err != nil {
		return
	}
	threads = make([]Thread, len(obs))
	for i, o := range obs {
		var thread []Thread
		err = m.get(threadsSpace, "primary", makeKey(o.ThreadID), &thread)
		if err != nil {
			return
		}
		if len(thread) == 1 {
			threads[i] = thread[0]
		}
	}
	return
}

func (m *TarantoolManager) MarkAsDelivered(userID, threadID string, times ...time.Time) (err error) {
	t := time.Now()
	if len(times) > 0 {
		t = times[0]
	}

	_, err = m.conn.Update(observerSpace, "primary", makeKey(threadID, userID), makeUpdate(newUpdateOp("=", 2, t.UnixNano()))) //makeUpdate(newUpdateOp("=", 3, t.Unix())))

	return
}

// Activity события двигаться по limit,offset что предлагает tnt
func (m *TarantoolManager) activity(threadID string, limit,
	offset uint32) (events []Event, err error) {
	err = m.conn.SelectTyped(eventsSpace, "threadline_idx", limit, offset, tarantool.IterReq, makeKey(threadID), &events)
	if err != nil {
		return
	}
	return
}

// ThreadlineActivity range over threadline in revers order
func (m *TarantoolManager) threadlineActivity(userID, threadID string, limit,
	offset uint32) (events []Event, err error) {
	err = m.conn.Call17Typed("threadline", makeKey(userID, threadID, uint64(limit), uint64(offset)), &events)
	if err != nil {
		return
	}
	return
}

// RecentActivity
func (m *TarantoolManager) RecentActivity(userID, threadID string, limit,
	offset uint32) (events []Event, err error) {
	var thread *Thread
	thread, err = m.GetThread(threadID)
	if err != nil {
		return
	}

	if thread.ThreadlineEnabled {
		return m.threadlineActivity(userID, threadID, limit, offset)
	}
	return m.activity(threadID, limit, offset)
}

// CountEvents возвращает количество событий после даты t
func (m *TarantoolManager) CountEvents(userID, threadID string, t time.Time) (int, error) {
	var counts []int
	err := m.conn.Call17Typed("count_events", makeKey(userID, threadID, t.Nanosecond()), &counts)
	if err != nil {
		return 0, err
	}
	return int(counts[0]), nil
}

// NewEvent cerate new event. if id empty, wiil be generated uuid.
// If CreatedAt zero, it will be setted to time.Now()
func (m *TarantoolManager) NewEvent(ev *Event) (err error) {
	if ev.ThreadID == "" {
		return ErrEmptyThreadID
	}
	if ev.EventID == "" {
		ev.EventID = uuid.NewV4().String()
	}
	if ev.CreatedAt.IsZero() {
		ev.CreatedAt = time.Now()
	}
	ev.UpdatedAt = time.Now()

	_, err = m.conn.Insert(eventsSpace, ev)
	if err != nil {
		return
	}

	// check is trhreadline enabled
	var thread *Thread
	thread, err = m.GetThread(ev.ThreadID)
	if err != nil {
		return
	}
	if !thread.ThreadlineEnabled {
		return
	}

	_, err = m.conn.Call17("new_event_in_threadline", makeKey(thread.ThreadID, ev.CreatedAt.UnixNano(), ev.EventID))

	// recieve all thread observers
	// var (
	// 	obs    []User
	// 	tmpObs []User
	// 	offset uint32 = 1000
	// 	limit  uint32
	// )
	// for {
	// 	tmpObs, err = m.Observers(thread.ThreadID, offset, limit)
	// 	if err != nil {
	// 		return
	// 	}
	// 	obs = append(obs, tmpObs...)
	// 	if uint32(len(tmpObs)) != offset {
	// 		break
	// 	}
	// }
	//
	// for _,o:=range obs {
	//
	// }

	return
}

// GetEvent return event by their ID
func (m *TarantoolManager) GetEvent(eventID string) (ev *Event, err error) {
	ev = new(Event)
	err = m.conn.SelectTyped(eventsSpace, "primary", 0, 1, tarantool.IterEq,
		makeKey(eventID), ev)
	return
}

// GetEvents returns events by their ids
func (m *TarantoolManager) GetEvents(ids ...string) (events []Event, err error) {
	for _, id := range ids {
		var event *Event
		event, err = m.GetEvent(id)
		if err != nil {
			return
		}
		events = append(events, *event)
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

func (m *TarantoolManager) updateEventUpdatedAt(eventID string) error {
	_, err := m.conn.Update(threadLineSpace, "primary", makeKey(eventID), makeUpdate(newUpdateOp("=", 4, time.Now().UnixNano())))
	if err != nil {
		return err
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

// SetRelatedData set data to tarantool
func (m *TarantoolManager) SetRelatedData(rel *RelatedData) (err error) {
	_, err = m.conn.Replace(relatedSpace, rel)
	if err != nil {
		return err
	}
	err = m.updateEventUpdatedAt(rel.EventID)
	return
}

// GetRelatedDatas возвращает события с кастомными данными юзера
func (m *TarantoolManager) GetRelatedDatas(userID string, events ...Event) (obs []EventObserver, err error) {
	for _, ev := range events {
		var rel []RelatedData
		err = m.get(relatedSpace, "primary", makeKey(userID, ev.EventID), &rel)
		if err != nil {
			return
		}
		obs = append(obs, EventObserver{Event: ev, RelatedData: rel[0]})
	}

	return
}

// NewUser create new user. If u.UserID empty, will be generated uuid
func (m *TarantoolManager) NewUser(u *User) (err error) {
	if u.UserID == "" {
		u.UserID = uuid.NewV4().String()
	}

	_, err = m.conn.Insert(usersSpace, u)
	if err != nil {
		return
	}
	return nil
}

// GetUser return user by their ID
func (m *TarantoolManager) GetUser(userID string) (u *User, err error) {
	u = new(User)
	err = m.get(usersSpace, "primary", makeKey(userID), u)
	return
}

// DeleteUser delete user from db
func (m *TarantoolManager) DeleteUser(userID string) (err error) {
	_, err = m.conn.Delete(usersSpace, "primary", makeKey(userID))
	return
}
