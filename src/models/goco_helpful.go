// Code generated. DO NOT EDIT.
package models

import (
	"github.com/golang/glog"
	"github.com/jackc/pgx"
	"time"
)

// DTOAbstract
func NewDTOAbstract() *DTOAbstract {
	model := new(DTOAbstract)
	return model
}

type DTOAbstract struct {
	// Tx
	Tx *pgx.Tx `json:"-" `
}

func (model DTOAbstract) TransformTo(out interface{}) error {
	switch out.(type) {
	default:
		glog.Errorf("Not supported type %v", out)
		return ErrNotSupported
	}
	return nil
}

func (model *DTOAbstract) TransformFrom(in interface{}) error {
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

func (d *DTOAbstract) Maps() map[string]interface{} {
	return map[string]interface{}{
		// Tx
		"tx": &d.Tx,
	}
}

// Fields extract of fields from map
func (d *DTOAbstract) Fields(fields ...string) ([]string, []interface{}) {
	return ExtractFieldsFromMap(d.Maps(), fields...)
}

// FromJson data as []byte or io.Reader
func (d *DTOAbstract) FromJson(data interface{}) error {
	return FromJson(d, data)
}

// ModelAbstract
func NewModelAbstract() *ModelAbstract {
	model := new(ModelAbstract)
	return model
}

type ModelAbstract struct {
	// IsRemoved
	IsRemoved bool `json:"is_removed" sql:"type:boolean;default:false" `
	// CreatedAt
	CreatedAt time.Time `json:"created_at" sql:"type:timestamp;default:null" `
	// UpdatedAt
	UpdatedAt time.Time `json:"updated_at" sql:"type:timestamp;default:null" `
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

func (m *ModelAbstract) Maps() map[string]interface{} {
	return map[string]interface{}{
		// IsRemoved
		"is_removed": &m.IsRemoved,
		// CreatedAt
		"created_at": &m.CreatedAt,
		// UpdatedAt
		"updated_at": &m.UpdatedAt,
	}
}

// Fields extract of fields from map
func (m *ModelAbstract) Fields(fields ...string) ([]string, []interface{}) {
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

func (t *Trace) Maps() map[string]interface{} {
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
func (t *Trace) Fields(fields ...string) ([]string, []interface{}) {
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

func (a *AppError) Maps() map[string]interface{} {
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
func (a *AppError) Fields(fields ...string) ([]string, []interface{}) {
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

func (r *ResponseDTO) Maps() map[string]interface{} {
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
func (r *ResponseDTO) Fields(fields ...string) ([]string, []interface{}) {
	return ExtractFieldsFromMap(r.Maps(), fields...)
}

// FromJson data as []byte or io.Reader
func (r *ResponseDTO) FromJson(data interface{}) error {
	return FromJson(r, data)
}
