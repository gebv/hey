// Code generated. DO NOT EDIT.
package models

import (
	"github.com/golang/glog"
)

// Context
func NewContext() *Context {
	model := new(Context)
	// Custom factory code
	model.T = NewTrace()
	return model
}

type Context struct {
	// T	Информация о текущем действии
	T *Trace `json:"t" `
	// Err
	Err *AppError `json:"err" `
}

func (model Context) TransformTo(out interface{}) error {
	switch out.(type) {
	default:
		glog.Errorf("Not supported type %v", out)
		return ErrNotSupported
	}
	return nil
}

func (model *Context) TransformFrom(in interface{}) error {
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

func (c Context) Maps() map[string]interface{} {
	return map[string]interface{}{
		// T	Информация о текущем действии
		"t": &c.T,
		// Err
		"err": &c.Err,
	}
}

// Fields extract of fields from map
func (c Context) Fields(fields ...string) ([]string, []interface{}) {
	return ExtractFieldsFromMap(c.Maps(), fields...)
}

// FromJson data as []byte or io.Reader
func (c *Context) FromJson(data interface{}) error {
	return FromJson(c, data)
}
