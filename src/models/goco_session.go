// Code generated. DO NOT EDIT.
package models

import (
	"github.com/golang/glog"
)

// Session
func NewSession() *Session {
	model := new(Session)
	// Custom factory code
	model.Props = make(map[string]interface{})
	return model
}

type Session struct {
	// Id
	Id string `json:"id" `
	// Client
	Client Client `json:"client" `
	// Flags
	Flags StringArray `json:"flags" `
	// Props
	Props map[string]interface{} `json:"props" `
}

func (model Session) TransformTo(out interface{}) error {
	switch out.(type) {
	default:
		glog.Errorf("Not supported type %v", out)
		return ErrNotSupported
	}
	return nil
}

func (model *Session) TransformFrom(in interface{}) error {
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

func (s *Session) Maps() map[string]interface{} {
	return map[string]interface{}{
		// Id
		"id": &s.Id,
		// Client
		"client": &s.Client,
		// Flags
		"flags": &s.Flags,
		// Props
		"props": &s.Props,
	}
}

// Fields extract of fields from map
func (s *Session) Fields(fields ...string) ([]string, []interface{}) {
	return ExtractFieldsFromMap(s.Maps(), fields...)
}

// FromJson data as []byte or io.Reader
func (s *Session) FromJson(data interface{}) error {
	return FromJson(s, data)
}
