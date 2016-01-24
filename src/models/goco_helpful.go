// Code generated. DO NOT EDIT.
package models

import (
	"github.com/golang/glog"
	"github.com/lib/pq"
	"time"
)

// ModelAbstract
func NewModelAbstract() *ModelAbstract {
	model := new(ModelAbstract)
	return model
}

type ModelAbstract struct {
	// IsRemoved
	IsRemoved bool `json:"is_removed" sql:"type:boolean;default:false" `
	// AtCreated
	AtCreated time.Time `json:"at_created" sql:"type:timestamp;default:null" `
	// AtUpdated
	AtUpdated time.Time `json:"at_updated" sql:"type:timestamp;default:null" `
	// AtRemoved
	AtRemoved pq.NullTime `json:"at_removed" sql:"type:timestamp;default:null" `
}

func (model ModelAbstract) TransformTo(out interface{}) error {
	switch out.(type) {
	default:
		glog.Errorf("Not supported type %v", out)
		return ErrNotSupported
	}
	return nil
}

func (model *ModelAbstract) TransformFrom(in interface{}) error {
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

func (m ModelAbstract) Maps() map[string]interface{} {
	return map[string]interface{}{
		// IsRemoved
		"is_removed": &m.IsRemoved,
		// AtCreated
		"at_created": &m.AtCreated,
		// AtUpdated
		"at_updated": &m.AtUpdated,
		// AtRemoved
		"at_removed": &m.AtRemoved,
	}
}

// Fields extract of fields from map
func (m ModelAbstract) Fields(fields ...string) ([]string, []interface{}) {
	return ExtractFieldsFromMap(m.Maps(), fields...)
}

// FromJson data as []byte or io.Reader
func (m *ModelAbstract) FromJson(data interface{}) error {
	return FromJson(m, data)
}

// Trace
func NewTrace() *Trace {
	model := new(Trace)
	return model
}

type Trace struct {
	// RequestId
	RequestId string `json:"request_id" `
	// Path
	Path string `json:"path" `
	// Ip
	Ip string `json:"ip" `
	// ClientId
	ClientId string `json:"client_id" `
}

func (model Trace) TransformTo(out interface{}) error {
	switch out.(type) {
	default:
		glog.Errorf("Not supported type %v", out)
		return ErrNotSupported
	}
	return nil
}

func (model *Trace) TransformFrom(in interface{}) error {
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

func (t Trace) Maps() map[string]interface{} {
	return map[string]interface{}{
		// RequestId
		"request_id": &t.RequestId,
		// Path
		"path": &t.Path,
		// Ip
		"ip": &t.Ip,
		// ClientId
		"client_id": &t.ClientId,
	}
}

// Fields extract of fields from map
func (t Trace) Fields(fields ...string) ([]string, []interface{}) {
	return ExtractFieldsFromMap(t.Maps(), fields...)
}

// FromJson data as []byte or io.Reader
func (t *Trace) FromJson(data interface{}) error {
	return FromJson(t, data)
}

// AppError
func NewAppError() *AppError {
	model := new(AppError)
	return model
}

type AppError struct {
	// Message
	Message string `json:"message" `
	// DevMessage
	DevMessage string `json:"dev_message" `
	// StatusCode
	StatusCode int `json:"status_code" `
	// T
	T Trace `json:"t" `
}

func (model AppError) TransformTo(out interface{}) error {
	switch out.(type) {
	default:
		glog.Errorf("Not supported type %v", out)
		return ErrNotSupported
	}
	return nil
}

func (model *AppError) TransformFrom(in interface{}) error {
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

func (a AppError) Maps() map[string]interface{} {
	return map[string]interface{}{
		// Message
		"message": &a.Message,
		// DevMessage
		"dev_message": &a.DevMessage,
		// StatusCode
		"status_code": &a.StatusCode,
		// T
		"t": &a.T,
	}
}

// Fields extract of fields from map
func (a AppError) Fields(fields ...string) ([]string, []interface{}) {
	return ExtractFieldsFromMap(a.Maps(), fields...)
}

// FromJson data as []byte or io.Reader
func (a *AppError) FromJson(data interface{}) error {
	return FromJson(a, data)
}

// ResponseDTO
func NewResponseDTO() *ResponseDTO {
	model := new(ResponseDTO)
	// Custom factory code
	model.StatusCode = 400 // http.StatusBadRequest
	return model
}

type ResponseDTO struct {
	// StatusCode
	StatusCode int `json:"status_code" `
	// Message
	Message string `json:"message,omitempty" `
	// DevMessage
	DevMessage string `json:"dev_message,omitempty" `
	// Data
	Data interface{} `json:"data,omitempty" `
}

func (model ResponseDTO) TransformTo(out interface{}) error {
	switch out.(type) {
	default:
		glog.Errorf("Not supported type %v", out)
		return ErrNotSupported
	}
	return nil
}

func (model *ResponseDTO) TransformFrom(in interface{}) error {
	switch in.(type) {
	case *AppError:
		dto := in.(*AppError)
		model.StatusCode = dto.StatusCode
		model.Message = dto.Message
		model.DevMessage = dto.DevMessage
	default:
		glog.Errorf("Not supported type %v", in)
		return ErrNotSupported
	}
	return nil

}

//
// Helpful functions
//

func (r ResponseDTO) Maps() map[string]interface{} {
	return map[string]interface{}{
		// StatusCode
		"status_code": &r.StatusCode,
		// Message
		"message": &r.Message,
		// DevMessage
		"dev_message": &r.DevMessage,
		// Data
		"data": &r.Data,
	}
}

// Fields extract of fields from map
func (r ResponseDTO) Fields(fields ...string) ([]string, []interface{}) {
	return ExtractFieldsFromMap(r.Maps(), fields...)
}

// FromJson data as []byte or io.Reader
func (r *ResponseDTO) FromJson(data interface{}) error {
	return FromJson(r, data)
}
