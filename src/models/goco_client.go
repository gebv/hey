// Code generated. DO NOT EDIT.
package models

import (
	"github.com/golang/glog"
	"github.com/satori/go.uuid"
)

// Client
func NewClient() *Client {
	model := new(Client)
	return model
}

type Client struct {
	ModelAbstract
	// Id
	Id uuid.UUID `json:"id" `
	// Domain
	Domain string `json:"domain" `
	// Ips
	Ips []string `json:"ips" `
	// Secret
	Secret string `json:"secret" `
	// RedirectUri
	RedirectUri string `json:"redirect_uri" `
	// Scopes
	Scopes []string `json:"scopes" `
	// Flags
	Flags []string `json:"flags" `
	// Props
	Props ClientProps `json:"props" `
	// IsEnabled
	IsEnabled bool `json:"is_enabled" `
}

func (model Client) TransformTo(out interface{}) error {
	switch out.(type) {
	default:
		glog.Errorf("Not supported type %v", out)
		return ErrNotSupported
	}
	return nil
}

func (model *Client) TransformFrom(in interface{}) error {
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

func (c *Client) Maps() map[string]interface{} {
	maps := c.ModelAbstract.Maps()
	// Id
	maps["id"] = &c.Id
	// Domain
	maps["domain"] = &c.Domain
	// Ips
	maps["ips"] = &c.Ips
	// Secret
	maps["secret"] = &c.Secret
	// RedirectUri
	maps["redirect_uri"] = &c.RedirectUri
	// Scopes
	maps["scopes"] = &c.Scopes
	// Flags
	maps["flags"] = &c.Flags
	// Props
	maps["props"] = &c.Props
	// IsEnabled
	maps["is_enabled"] = &c.IsEnabled
	return maps
}

// Fields extract of fields from map
func (c *Client) Fields(fields ...string) ([]string, []interface{}) {
	return ExtractFieldsFromMap(c.Maps(), fields...)
}

// FromJson data as []byte or io.Reader
func (c *Client) FromJson(data interface{}) error {
	return FromJson(c, data)
}

func (Client) TableName() string {
	return "client"
}

// PrimaryName primary field name
func (Client) PrimaryName() string {
	return "id"
}

// PrimaryValue primary value
func (c Client) PrimaryValue() uuid.UUID {
	return c.Id
}

// model
// ClientProps
func NewClientProps() *ClientProps {
	model := new(ClientProps)
	return model
}

type ClientProps struct {
	// FullName
	FullName string `json:"full_name" `
	// Email
	Email string `json:"email" `
}

func (model ClientProps) TransformTo(out interface{}) error {
	switch out.(type) {
	default:
		glog.Errorf("Not supported type %v", out)
		return ErrNotSupported
	}
	return nil
}

func (model *ClientProps) TransformFrom(in interface{}) error {
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

func (c *ClientProps) Maps() map[string]interface{} {
	return map[string]interface{}{
		// FullName
		"full_name": &c.FullName,
		// Email
		"email": &c.Email,
	}
}

// Fields extract of fields from map
func (c *ClientProps) Fields(fields ...string) ([]string, []interface{}) {
	return ExtractFieldsFromMap(c.Maps(), fields...)
}

// FromJson data as []byte or io.Reader
func (c *ClientProps) FromJson(data interface{}) error {
	return FromJson(c, data)
}
