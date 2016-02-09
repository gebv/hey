// Code generated. DO NOT EDIT.
package models

import (
	"github.com/golang/glog"
	"github.com/satori/go.uuid"
)

// User
func NewUser() *User {
	model := new(User)
	return model
}

type User struct {
	ModelAbstract
	// UserId
	UserId uuid.UUID `json:"user_id" `
	// ClientId
	ClientId uuid.UUID `json:"client_id" `
	// ExtId
	ExtId string `json:"ext_id" `
	// ExtIdHash
	ExtIdHash string `json:"ext_id_hash" `
	// ExtProps
	ExtProps StringMap `json:"ext_props" `
	// IsEnabled
	IsEnabled bool `json:"is_enabled" `
}

func (model User) TransformTo(out interface{}) error {
	switch out.(type) {
	default:
		glog.Errorf("Not supported type %v", out)
		return ErrNotSupported
	}
	return nil
}

func (model *User) TransformFrom(in interface{}) error {
	switch in.(type) {
	case *UserDTO:
		dto := in.(*UserDTO)
		model.ExtId = dto.ExtId
		model.ExtProps = dto.ExtProps
		model.IsEnabled = dto.IsEnabled
		model.ClientId = uuid.FromStringOrNil(dto.ClientId)
		model.ExtIdHash = HashText(dto.ExtId)
	default:
		glog.Errorf("Not supported type %v", in)
		return ErrNotSupported
	}
	return nil

}

//
// Helpful functions
//

func (u *User) Maps() map[string]interface{} {
	maps := u.ModelAbstract.Maps()
	// UserId
	maps["user_id"] = &u.UserId
	// ClientId
	maps["client_id"] = &u.ClientId
	// ExtId
	maps["ext_id"] = &u.ExtId
	// ExtIdHash
	maps["ext_id_hash"] = &u.ExtIdHash
	// ExtProps
	maps["ext_props"] = &u.ExtProps
	// IsEnabled
	maps["is_enabled"] = &u.IsEnabled
	return maps
}

// Fields extract of fields from map
func (u *User) Fields(fields ...string) ([]string, []interface{}) {
	return ExtractFieldsFromMap(u.Maps(), fields...)
}

// FromJson data as []byte or io.Reader
func (u *User) FromJson(data interface{}) error {
	return FromJson(u, data)
}

func (User) TableName() string {
	return "users"
}

// PrimaryName primary field name
func (User) PrimaryName() string {
	return "user_id"
}

// PrimaryValue primary value
func (u User) PrimaryValue() uuid.UUID {
	return u.UserId
}

// model
// UserDTO
func NewUserDTO() *UserDTO {
	model := new(UserDTO)
	return model
}

type UserDTO struct {
	DTOAbstract
	// ClientId
	ClientId string `v:"required" json:"-" `
	// ExtId
	ExtId string `json:"ext_id" v:"required" `
	// ExtProps
	ExtProps StringMap `json:"ext_props" `
	// IsEnabled
	IsEnabled bool `json:"is_enabled" `
}

func (model UserDTO) TransformTo(out interface{}) error {
	switch out.(type) {
	default:
		glog.Errorf("Not supported type %v", out)
		return ErrNotSupported
	}
	return nil
}

func (model *UserDTO) TransformFrom(in interface{}) error {
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

func (u *UserDTO) Maps() map[string]interface{} {
	return map[string]interface{}{
		// ClientId
		"client_id": &u.ClientId,
		// ExtId
		"ext_id": &u.ExtId,
		// ExtProps
		"ext_props": &u.ExtProps,
		// IsEnabled
		"is_enabled": &u.IsEnabled,
	}
}

// Fields extract of fields from map
func (u *UserDTO) Fields(fields ...string) ([]string, []interface{}) {
	return ExtractFieldsFromMap(u.Maps(), fields...)
}

// FromJson data as []byte or io.Reader
func (u *UserDTO) FromJson(data interface{}) error {
	return FromJson(u, data)
}
