// Code generated. DO NOT EDIT.
package models

import (
	"github.com/golang/glog"
	"github.com/satori/go.uuid"
	"strings"
	"time"
)

// ThreadCounter
func NewThreadCounter() *ThreadCounter {
	model := new(ThreadCounter)
	return model
}

type ThreadCounter struct {
	// ClientId
	ClientId uuid.UUID `json:"client_id" `
	// ThreadId
	ThreadId uuid.UUID `json:"thread_id" `
	// CounterEvents
	CounterEvents int64 `json:"counter_events" `
}

func (model ThreadCounter) TransformTo(out interface{}) error {
	switch out.(type) {
	default:
		glog.Errorf("Not supported type %v", out)
		return ErrNotSupported
	}
	return nil
}

func (model *ThreadCounter) TransformFrom(in interface{}) error {
	switch in.(type) {
	default:
		glog.Errorf("Not supported type %v", in)
		return ErrNotSupported
	}
	return nil

}

//
// Helpful functions
//

func (t *ThreadCounter) Maps() map[string]interface{} {
	return map[string]interface{}{
		// ClientId
		"client_id": &t.ClientId,
		// ThreadId
		"thread_id": &t.ThreadId,
		// CounterEvents
		"counter_events": &t.CounterEvents,
	}
}

// Fields extract of fields from map
func (t *ThreadCounter) Fields(fields ...string) ([]string, []interface{}) {
	return ExtractFieldsFromMap(t.Maps(), fields...)
}

// FromJson data as []byte or io.Reader
func (t *ThreadCounter) FromJson(data interface{}) error {
	return FromJson(t, data)
}

func (ThreadCounter) TableName() string {
	return "thread_counters"
}

// model
// ThreadWatcher
func NewThreadWatcher() *ThreadWatcher {
	model := new(ThreadWatcher)
	return model
}

type ThreadWatcher struct {
	// ClientId
	ClientId uuid.UUID `json:"client_id" `
	// ThreadId
	ThreadId uuid.UUID `json:"thread_id" `
	// UserId
	UserId uuid.UUID `json:"user_id" `
	// Unread
	Unread int64 `json:"unread" `
}

func (model ThreadWatcher) TransformTo(out interface{}) error {
	switch out.(type) {
	default:
		glog.Errorf("Not supported type %v", out)
		return ErrNotSupported
	}
	return nil
}

func (model *ThreadWatcher) TransformFrom(in interface{}) error {
	switch in.(type) {
	default:
		glog.Errorf("Not supported type %v", in)
		return ErrNotSupported
	}
	return nil

}

//
// Helpful functions
//

func (t *ThreadWatcher) Maps() map[string]interface{} {
	return map[string]interface{}{
		// ClientId
		"client_id": &t.ClientId,
		// ThreadId
		"thread_id": &t.ThreadId,
		// UserId
		"user_id": &t.UserId,
		// Unread
		"unread": &t.Unread,
	}
}

// Fields extract of fields from map
func (t *ThreadWatcher) Fields(fields ...string) ([]string, []interface{}) {
	return ExtractFieldsFromMap(t.Maps(), fields...)
}

// FromJson data as []byte or io.Reader
func (t *ThreadWatcher) FromJson(data interface{}) error {
	return FromJson(t, data)
}

func (ThreadWatcher) TableName() string {
	return "thread_watchers"
}

// model
// Threadline
func NewThreadline() *Threadline {
	model := new(Threadline)
	return model
}

type Threadline struct {
	// ClientId
	ClientId uuid.UUID `json:"client_id" `
	// ChannelId
	ChannelId uuid.UUID `json:"channel_id" `
	// ThreadId
	ThreadId uuid.UUID `json:"thread_id" `
	// EventId
	EventId uuid.UUID `json:"event_id" `
	// CreatedAt
	CreatedAt time.Time `json:"created_at" sql:"type:timestamp;default:null" `
}

func (model Threadline) TransformTo(out interface{}) error {
	switch out.(type) {
	default:
		glog.Errorf("Not supported type %v", out)
		return ErrNotSupported
	}
	return nil
}

func (model *Threadline) TransformFrom(in interface{}) error {
	switch in.(type) {
	default:
		glog.Errorf("Not supported type %v", in)
		return ErrNotSupported
	}
	return nil

}

//
// Helpful functions
//

func (t *Threadline) Maps() map[string]interface{} {
	return map[string]interface{}{
		// ClientId
		"client_id": &t.ClientId,
		// ChannelId
		"channel_id": &t.ChannelId,
		// ThreadId
		"thread_id": &t.ThreadId,
		// EventId
		"event_id": &t.EventId,
		// CreatedAt
		"created_at": &t.CreatedAt,
	}
}

// Fields extract of fields from map
func (t *Threadline) Fields(fields ...string) ([]string, []interface{}) {
	return ExtractFieldsFromMap(t.Maps(), fields...)
}

// FromJson data as []byte or io.Reader
func (t *Threadline) FromJson(data interface{}) error {
	return FromJson(t, data)
}

func (Threadline) TableName() string {
	return "threadline"
}

// PrimaryName primary field name
func (Threadline) PrimaryName() string {
	return "event_id"
}

// PrimaryValue primary value
func (t Threadline) PrimaryValue() uuid.UUID {
	return t.EventId
}

// model
// Thread
func NewThread() *Thread {
	model := new(Thread)
	// Custom factory code
	model.ExtProps = NewStringMap()
	return model
}

type Thread struct {
	ModelAbstract
	// ThreadId
	ThreadId uuid.UUID `json:"thread_id" `
	// ClientId
	ClientId uuid.UUID `json:"client_id" `
	// ChannelId
	ChannelId uuid.UUID `json:"channel_id" `
	// ExtId
	ExtId string `json:"ext_id" `
	// ExtIdHash
	ExtIdHash string `json:"ext_id_hash" `
	// ExtProps
	ExtProps StringMap `json:"ext_props" `
	// ExtFlags
	ExtFlags []string `json:"ext_flags" `
	// Depth
	Depth int64 `json:"depth" `
	// Owners
	Owners UUIDArray `json:"owners" `
	// RelatedEventId
	RelatedEventId uuid.UUID `json:"related_event_id" `
	// ParentThreadId
	ParentThreadId uuid.UUID `json:"parent_thread_id" `
}

func (model Thread) TransformTo(out interface{}) error {
	switch out.(type) {
	default:
		glog.Errorf("Not supported type %v", out)
		return ErrNotSupported
	}
	return nil
}

func (model *Thread) TransformFrom(in interface{}) error {
	switch in.(type) {
	case *ThreadDTO:
		dto := in.(*ThreadDTO)
		model.ExtId = dto.ExtId
		model.ExtProps = dto.ExtProps
		model.ExtFlags = dto.ExtFlags
		model.ParentThreadId = uuid.FromStringOrNil(dto.ParentThreadId)
		model.RelatedEventId = uuid.FromStringOrNil(dto.RelatedEventId)
		model.ClientId = uuid.FromStringOrNil(dto.ClientId)
		model.ChannelId = uuid.FromStringOrNil(dto.ChannelId)
		model.ExtIdHash = HashText(dto.ExtId)
		model.Depth = int64(len(strings.Split(dto.ExtId, ":")))
		model.Owners.FromArray(dto.Owners)
	default:
		glog.Errorf("Not supported type %v", in)
		return ErrNotSupported
	}
	return nil

}

//
// Helpful functions
//

func (t *Thread) Maps() map[string]interface{} {
	maps := t.ModelAbstract.Maps()
	// ThreadId
	maps["thread_id"] = &t.ThreadId
	// ClientId
	maps["client_id"] = &t.ClientId
	// ChannelId
	maps["channel_id"] = &t.ChannelId
	// ExtId
	maps["ext_id"] = &t.ExtId
	// ExtIdHash
	maps["ext_id_hash"] = &t.ExtIdHash
	// ExtProps
	maps["ext_props"] = &t.ExtProps
	// ExtFlags
	maps["ext_flags"] = &t.ExtFlags
	// Depth
	maps["depth"] = &t.Depth
	// Owners
	maps["owners"] = &t.Owners
	// RelatedEventId
	maps["related_event_id"] = &t.RelatedEventId
	// ParentThreadId
	maps["parent_thread_id"] = &t.ParentThreadId
	return maps
}

// Fields extract of fields from map
func (t *Thread) Fields(fields ...string) ([]string, []interface{}) {
	return ExtractFieldsFromMap(t.Maps(), fields...)
}

// FromJson data as []byte or io.Reader
func (t *Thread) FromJson(data interface{}) error {
	return FromJson(t, data)
}

func (Thread) TableName() string {
	return "threads"
}

// PrimaryName primary field name
func (Thread) PrimaryName() string {
	return "thread_id"
}

// PrimaryValue primary value
func (t Thread) PrimaryValue() uuid.UUID {
	return t.ThreadId
}

// model
// ThreadDTO
func NewThreadDTO() *ThreadDTO {
	model := new(ThreadDTO)
	// Custom factory code
	model.ExtProps = NewStringMap()
	return model
}

type ThreadDTO struct {
	DTOAbstract
	// ClientId
	ClientId string `json:"-" v:"required" `
	// ChannelId
	ChannelId string `json:"-" v:"required" `
	// ExtId
	ExtId string `json:"ext_id" v:"required" `
	// ExtProps
	ExtProps StringMap `json:"ext_props" `
	// ExtFlags
	ExtFlags StringArray `json:"ext_flags" `
	// Owners
	Owners StringArray `json:"owners" v:"gt=1" `
	// EventCreator
	EventCreator string `json:"-" `
	// CreatedEventId
	CreatedEventId string `json:"-" `
	// RelatedEventId	Связанная с текущим событием событие из потока ParentThreadId
	RelatedEventId string `json:"related_event_id" `
	// ParentThreadId
	ParentThreadId string `json:"parent_thread_id" `
}

func (model ThreadDTO) TransformTo(out interface{}) error {
	switch out.(type) {
	default:
		glog.Errorf("Not supported type %v", out)
		return ErrNotSupported
	}
	return nil
}

func (model *ThreadDTO) TransformFrom(in interface{}) error {
	switch in.(type) {
	case *Channel:
		dto := in.(*Channel)
		model.ExtFlags = dto.ExtFlags
		model.ClientId = dto.ClientId.String()
		model.ChannelId = dto.ChannelId.String()
		model.ExtId = dto.ExtId
		model.ExtProps = dto.ExtProps
		model.Owners.FromArray(dto.Owners)
	default:
		glog.Errorf("Not supported type %v", in)
		return ErrNotSupported
	}
	return nil

}

//
// Helpful functions
//

func (t *ThreadDTO) Maps() map[string]interface{} {
	return map[string]interface{}{
		// ClientId
		"client_id": &t.ClientId,
		// ChannelId
		"channel_id": &t.ChannelId,
		// ExtId
		"ext_id": &t.ExtId,
		// ExtProps
		"ext_props": &t.ExtProps,
		// ExtFlags
		"ext_flags": &t.ExtFlags,
		// Owners
		"owners": &t.Owners,
		// EventCreator
		"event_creator": &t.EventCreator,
		// CreatedEventId
		"created_event_id": &t.CreatedEventId,
		// RelatedEventId	Связанная с текущим событием событие из потока ParentThreadId
		"related_event_id": &t.RelatedEventId,
		// ParentThreadId
		"parent_thread_id": &t.ParentThreadId,
	}
}

// Fields extract of fields from map
func (t *ThreadDTO) Fields(fields ...string) ([]string, []interface{}) {
	return ExtractFieldsFromMap(t.Maps(), fields...)
}

// FromJson data as []byte or io.Reader
func (t *ThreadDTO) FromJson(data interface{}) error {
	return FromJson(t, data)
}
