// Code generated. DO NOT EDIT.
package models

import (
	"encoding/base64"
	"github.com/golang/glog"
	"github.com/satori/go.uuid"
	"strings"
)

// EventLoadResult
func NewEventLoadResult() *EventLoadResult {
	model := new(EventLoadResult)
	return model
}

type EventLoadResult struct {
	// Events
	Events []*Event `json:"events" `
	// HasNext
	HasNext bool `json:"has_next" `
	// ThreadTotalCount
	ThreadTotalCount string `json:"thread_total_count" `
	// Cursor
	Cursor string `json:"cursor" `
}

func (model EventLoadResult) TransformTo(out interface{}) error {
	switch out.(type) {
	default:
		glog.Errorf("Not supported type %v", out)
		return ErrNotSupported
	}
	return nil
}

func (model *EventLoadResult) TransformFrom(in interface{}) error {
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

func (e *EventLoadResult) Maps() map[string]interface{} {
	return map[string]interface{}{
		// Events
		"events": &e.Events,
		// HasNext
		"has_next": &e.HasNext,
		// ThreadTotalCount
		"thread_total_count": &e.ThreadTotalCount,
		// Cursor
		"cursor": &e.Cursor,
	}
}

// Fields extract of fields from map
func (e *EventLoadResult) Fields(fields ...string) ([]string, []interface{}) {
	return ExtractFieldsFromMap(e.Maps(), fields...)
}

// FromJson data as []byte or io.Reader
func (e *EventLoadResult) FromJson(data interface{}) error {
	return FromJson(e, data)
}

// EventDTO
func NewEventDTO() *EventDTO {
	model := new(EventDTO)
	return model
}

type EventDTO struct {
	DTOAbstract
	// ClientId
	ClientId string `json:"-" `
	// EventId
	EventId string `json:"-" `
	// Thread	uri channelname:threadname:etc...
	Thread string `json:"thread" `
	// Creator
	Creator string `json:"creator" `
	// ParentEventId
	ParentEventId string `json:"parent_event_id" `
	// DataBase64	base64
	DataBase64 string `json:"data_base" `
	// ExtFlags	Type default: {m: ''}
	ExtFlags StringArray `json:"ext_flags" `
}

func (model EventDTO) TransformTo(out interface{}) error {
	switch out.(type) {
	case *ChannelDTO:
		dto := out.(*ChannelDTO)
		dto.ClientId = model.ClientId
		dto.ExtId = strings.Split(model.Thread, ":")[0]
	case *ThreadDTO:
		dto := out.(*ThreadDTO)
		dto.ClientId = model.ClientId
		dto.EventCreator = model.Creator
		dto.ExtId = model.Thread
	case *UserDTO:
		dto := out.(*UserDTO)
		dto.ClientId = model.ClientId
		dto.ExtId = model.Creator
	default:
		glog.Errorf("Not supported type %v", out)
		return ErrNotSupported
	}
	return nil
}

func (model *EventDTO) TransformFrom(in interface{}) error {
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

func (e *EventDTO) Maps() map[string]interface{} {
	return map[string]interface{}{
		// ClientId
		"client_id": &e.ClientId,
		// EventId
		"event_id": &e.EventId,
		// Thread	uri channelname:threadname:etc...
		"thread": &e.Thread,
		// Creator
		"creator": &e.Creator,
		// ParentEventId
		"parent_event_id": &e.ParentEventId,
		// DataBase64	base64
		"data_base": &e.DataBase64,
		// ExtFlags	Type default: {m: ''}
		"ext_flags": &e.ExtFlags,
	}
}

// Fields extract of fields from map
func (e *EventDTO) Fields(fields ...string) ([]string, []interface{}) {
	return ExtractFieldsFromMap(e.Maps(), fields...)
}

// FromJson data as []byte or io.Reader
func (e *EventDTO) FromJson(data interface{}) error {
	return FromJson(e, data)
}

// Event
func NewEvent() *Event {
	model := new(Event)
	return model
}

type Event struct {
	ModelAbstract
	// EventId
	EventId uuid.UUID `json:"event_id" `
	// ClientId
	ClientId uuid.UUID `json:"client_id" `
	// ChannelId
	ChannelId uuid.UUID `json:"channel_id" `
	// ThreadId
	ThreadId uuid.UUID `json:"thread_id" `
	// Creator
	Creator uuid.UUID `json:"creator" `
	// ParentEventId
	ParentEventId uuid.UUID `json:"parent_event_id" `
	// ParentThreadId
	ParentThreadId uuid.UUID `json:"parent_thread_id" `
	// BranchThreadId
	BranchThreadId uuid.UUID `json:"branch_thread_id" `
	// Data
	Data []byte `json:"data" `
	// ExtFlags
	ExtFlags []string `json:"ext_flags" `
	// Flags
	Flags []string `json:"flags" `
}

func (model Event) TransformTo(out interface{}) error {
	switch out.(type) {
	default:
		glog.Errorf("Not supported type %v", out)
		return ErrNotSupported
	}
	return nil
}

func (model *Event) TransformFrom(in interface{}) error {
	switch in.(type) {
	case *EventDTO:
		dto := in.(*EventDTO)
		model.ExtFlags = []string(dto.ExtFlags)
		model.ClientId = uuid.FromStringOrNil(dto.ClientId)
		model.EventId = uuid.FromStringOrNil(dto.EventId)
		model.Creator = uuid.FromStringOrNil(dto.Creator)
		model.ParentEventId = uuid.FromStringOrNil(dto.ParentEventId)
		var err error
		model.Data, err = base64.StdEncoding.DecodeString(dto.DataBase64)
		if err != nil {
			return err
		}
	default:
		glog.Errorf("Not supported type %v", in)
		return ErrNotSupported
	}
	return nil

}

//
// Helpful functions
//

func (e *Event) Maps() map[string]interface{} {
	maps := e.ModelAbstract.Maps()
	// EventId
	maps["event_id"] = &e.EventId
	// ClientId
	maps["client_id"] = &e.ClientId
	// ChannelId
	maps["channel_id"] = &e.ChannelId
	// ThreadId
	maps["thread_id"] = &e.ThreadId
	// Creator
	maps["creator"] = &e.Creator
	// ParentEventId
	maps["parent_event_id"] = &e.ParentEventId
	// ParentThreadId
	maps["parent_thread_id"] = &e.ParentThreadId
	// BranchThreadId
	maps["branch_thread_id"] = &e.BranchThreadId
	// Data
	maps["data"] = &e.Data
	// ExtFlags
	maps["ext_flags"] = &e.ExtFlags
	// Flags
	maps["flags"] = &e.Flags
	return maps
}

// Fields extract of fields from map
func (e *Event) Fields(fields ...string) ([]string, []interface{}) {
	return ExtractFieldsFromMap(e.Maps(), fields...)
}

// FromJson data as []byte or io.Reader
func (e *Event) FromJson(data interface{}) error {
	return FromJson(e, data)
}

func (Event) TableName() string {
	return "events"
}

// PrimaryName primary field name
func (Event) PrimaryName() string {
	return "event_id"
}

// PrimaryValue primary value
func (e Event) PrimaryValue() uuid.UUID {
	return e.EventId
}

// model
