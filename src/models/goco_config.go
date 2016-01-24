// Code generated. DO NOT EDIT.
package models

import (
	"github.com/golang/glog"
	"time"
)

// Config
func NewConfig() *Config {
	model := new(Config)
	return model
}

type Config struct {
	// ServiceSettings
	ServiceSettings ServiceSettings `json:"ServiceSettings" `
	// StorageSettings
	StorageSettings StorageSettings `json:"StorageSettings" `
}

func (model Config) TransformTo(out interface{}) error {
	switch out.(type) {
	default:
		glog.Errorf("Not supported type %v", out)
		return ErrNotSupported
	}
	return nil
}

func (model *Config) TransformFrom(in interface{}) error {
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

func (c Config) Maps() map[string]interface{} {
	return map[string]interface{}{
		// ServiceSettings
		"service_settings": &c.ServiceSettings,
		// StorageSettings
		"storage_settings": &c.StorageSettings,
	}
}

// Fields extract of fields from map
func (c Config) Fields(fields ...string) ([]string, []interface{}) {
	return ExtractFieldsFromMap(c.Maps(), fields...)
}

// FromJson data as []byte or io.Reader
func (c *Config) FromJson(data interface{}) error {
	return FromJson(c, data)
}

// ServiceSettings
func NewServiceSettings() *ServiceSettings {
	model := new(ServiceSettings)
	return model
}

type ServiceSettings struct {
	// ListenAddress
	ListenAddress string `json:"ListenAddress" `
	// Mode
	Mode string `json:"Mode" `
	// TimeoutRequest
	TimeoutRequest time.Duration `json:"TimeoutRequest" `
}

func (model ServiceSettings) TransformTo(out interface{}) error {
	switch out.(type) {
	default:
		glog.Errorf("Not supported type %v", out)
		return ErrNotSupported
	}
	return nil
}

func (model *ServiceSettings) TransformFrom(in interface{}) error {
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

func (s ServiceSettings) Maps() map[string]interface{} {
	return map[string]interface{}{
		// ListenAddress
		"listen_address": &s.ListenAddress,
		// Mode
		"mode": &s.Mode,
		// TimeoutRequest
		"timeout_request": &s.TimeoutRequest,
	}
}

// Fields extract of fields from map
func (s ServiceSettings) Fields(fields ...string) ([]string, []interface{}) {
	return ExtractFieldsFromMap(s.Maps(), fields...)
}

// FromJson data as []byte or io.Reader
func (s *ServiceSettings) FromJson(data interface{}) error {
	return FromJson(s, data)
}

// StorageSettings
func NewStorageSettings() *StorageSettings {
	model := new(StorageSettings)
	return model
}

type StorageSettings struct {
	// Network
	Network string `json:"Network" `
	// Host
	Host string `json:"Host" `
	// Port
	Port string `json:"Port" `
	// User
	User string `json:"User" `
	// Password
	Password string `json:"Password" `
	// Database
	Database string `json:"Database" `
}

func (model StorageSettings) TransformTo(out interface{}) error {
	switch out.(type) {
	default:
		glog.Errorf("Not supported type %v", out)
		return ErrNotSupported
	}
	return nil
}

func (model *StorageSettings) TransformFrom(in interface{}) error {
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

func (s StorageSettings) Maps() map[string]interface{} {
	return map[string]interface{}{
		// Network
		"network": &s.Network,
		// Host
		"host": &s.Host,
		// Port
		"port": &s.Port,
		// User
		"user": &s.User,
		// Password
		"password": &s.Password,
		// Database
		"database": &s.Database,
	}
}

// Fields extract of fields from map
func (s StorageSettings) Fields(fields ...string) ([]string, []interface{}) {
	return ExtractFieldsFromMap(s.Maps(), fields...)
}

// FromJson data as []byte or io.Reader
func (s *StorageSettings) FromJson(data interface{}) error {
	return FromJson(s, data)
}
