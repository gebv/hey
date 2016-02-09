// Code generated. DO NOT EDIT.
package models

import (
	"github.com/golang/glog"
	"github.com/satori/go.uuid"
)

// ChannelCounter
func NewChannelCounter() *ChannelCounter {
	model := new(ChannelCounter)
	return model
}

type ChannelCounter struct {
	// ClientId
	ClientId uuid.UUID `json:"client_id" `
	// ChannelId
	ChannelId uuid.UUID `json:"channel_id" `
	// CounterEvents
	CounterEvents int64 `json:"counter_events" `
}

func (model ChannelCounter) TransformTo(out interface{}) error {
	switch out.(type) {
	default:
		glog.Errorf("Not supported type %v", out)
		return ErrNotSupported
	}
	return nil
}

func (model *ChannelCounter) TransformFrom(in interface{}) error {
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

func (c *ChannelCounter) Maps() map[string]interface{} {
	return map[string]interface{}{
		// ClientId
		"client_id": &c.ClientId,
		// ChannelId
		"channel_id": &c.ChannelId,
		// CounterEvents
		"counter_events": &c.CounterEvents,
	}
}

// Fields extract of fields from map
func (c *ChannelCounter) Fields(fields ...string) ([]string, []interface{}) {
	return ExtractFieldsFromMap(c.Maps(), fields...)
}

// FromJson data as []byte or io.Reader
func (c *ChannelCounter) FromJson(data interface{}) error {
	return FromJson(c, data)
}

func (ChannelCounter) TableName() string {
	return "channel_counters"
}

// model
// ChannelWatcher
func NewChannelWatcher() *ChannelWatcher {
	model := new(ChannelWatcher)
	return model
}

type ChannelWatcher struct {
	// ClientId
	ClientId uuid.UUID `json:"client_id" `
	// ChannelId
	ChannelId uuid.UUID `json:"channel_id" `
	// UserId
	UserId uuid.UUID `json:"user_id" `
	// Unread
	Unread int64 `json:"unread" `
}

func (model ChannelWatcher) TransformTo(out interface{}) error {
	switch out.(type) {
	default:
		glog.Errorf("Not supported type %v", out)
		return ErrNotSupported
	}
	return nil
}

func (model *ChannelWatcher) TransformFrom(in interface{}) error {
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

func (c *ChannelWatcher) Maps() map[string]interface{} {
	return map[string]interface{}{
		// ClientId
		"client_id": &c.ClientId,
		// ChannelId
		"channel_id": &c.ChannelId,
		// UserId
		"user_id": &c.UserId,
		// Unread
		"unread": &c.Unread,
	}
}

// Fields extract of fields from map
func (c *ChannelWatcher) Fields(fields ...string) ([]string, []interface{}) {
	return ExtractFieldsFromMap(c.Maps(), fields...)
}

// FromJson data as []byte or io.Reader
func (c *ChannelWatcher) FromJson(data interface{}) error {
	return FromJson(c, data)
}

func (ChannelWatcher) TableName() string {
	return "channel_watchers"
}

// model
// Channel
func NewChannel() *Channel {
	model := new(Channel)
	// Custom factory code
	model.ExtProps = NewStringMap()
	return model
}

type Channel struct {
	ModelAbstract
	// ChannelId
	ChannelId uuid.UUID `json:"channel_id" `
	// ClientId
	ClientId uuid.UUID `json:"client_id" `
	// ExtId
	ExtId string `json:"ext_id" `
	// ExtIdHash
	ExtIdHash string `json:"ext_id_hash" `
	// ExtProps
	ExtProps StringMap `json:"ext_props" `
	// ExtFlags
	ExtFlags []string `json:"ext_flags" `
	// Owners
	Owners UUIDArray `json:"owners" `
	// RootThreadId
	RootThreadId uuid.UUID `json:"root_thread_id" `
}

func (model Channel) TransformTo(out interface{}) error {
	switch out.(type) {
	default:
		glog.Errorf("Not supported type %v", out)
		return ErrNotSupported
	}
	return nil
}

func (model *Channel) TransformFrom(in interface{}) error {
	switch in.(type) {
	case *ChannelDTO:
		dto := in.(*ChannelDTO)
		model.ExtId = dto.ExtId
		model.ExtProps = dto.ExtProps
		model.ExtFlags = dto.ExtFlags
		model.ClientId = uuid.FromStringOrNil(dto.ClientId)
		model.ExtIdHash = HashText(dto.ExtId)
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

func (c *Channel) Maps() map[string]interface{} {
	maps := c.ModelAbstract.Maps()
	// ChannelId
	maps["channel_id"] = &c.ChannelId
	// ClientId
	maps["client_id"] = &c.ClientId
	// ExtId
	maps["ext_id"] = &c.ExtId
	// ExtIdHash
	maps["ext_id_hash"] = &c.ExtIdHash
	// ExtProps
	maps["ext_props"] = &c.ExtProps
	// ExtFlags
	maps["ext_flags"] = &c.ExtFlags
	// Owners
	maps["owners"] = &c.Owners
	// RootThreadId
	maps["root_thread_id"] = &c.RootThreadId
	return maps
}

// Fields extract of fields from map
func (c *Channel) Fields(fields ...string) ([]string, []interface{}) {
	return ExtractFieldsFromMap(c.Maps(), fields...)
}

// FromJson data as []byte or io.Reader
func (c *Channel) FromJson(data interface{}) error {
	return FromJson(c, data)
}

func (Channel) TableName() string {
	return "channels"
}

// PrimaryName primary field name
func (Channel) PrimaryName() string {
	return "channel_id"
}

// PrimaryValue primary value
func (c Channel) PrimaryValue() uuid.UUID {
	return c.ChannelId
}

// model
// ChannelDTO
func NewChannelDTO() *ChannelDTO {
	model := new(ChannelDTO)
	// Custom factory code
	model.ExtProps = NewStringMap()
	return model
}

type ChannelDTO struct {
	DTOAbstract
	// ClientId
	ClientId string `json:"-" v:"required" `
	// ExtId
	ExtId string `json:"ext_id" v:"required" `
	// ExtProps
	ExtProps StringMap `json:"ext_props" `
	// ExtFlags
	ExtFlags StringArray `json:"ext_flags" `
	// Owners
	Owners StringArray `json:"owners" v:"gt=1" `
}

func (model ChannelDTO) TransformTo(out interface{}) error {
	switch out.(type) {
	default:
		glog.Errorf("Not supported type %v", out)
		return ErrNotSupported
	}
	return nil
}

func (model *ChannelDTO) TransformFrom(in interface{}) error {
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

func (c *ChannelDTO) Maps() map[string]interface{} {
	return map[string]interface{}{
		// ClientId
		"client_id": &c.ClientId,
		// ExtId
		"ext_id": &c.ExtId,
		// ExtProps
		"ext_props": &c.ExtProps,
		// ExtFlags
		"ext_flags": &c.ExtFlags,
		// Owners
		"owners": &c.Owners,
	}
}

// Fields extract of fields from map
func (c *ChannelDTO) Fields(fields ...string) ([]string, []interface{}) {
	return ExtractFieldsFromMap(c.Maps(), fields...)
}

// FromJson data as []byte or io.Reader
func (c *ChannelDTO) FromJson(data interface{}) error {
	return FromJson(c, data)
}
