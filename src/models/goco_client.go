// Code generated. DO NOT EDIT.
package models

import (
	"github.com/golang/glog"
)

// Client
func NewClient() *Client {
	model := new(Client)
	return model
}

type Client struct {
	ModelAbstract
	// ClientId
	ClientId string `json:"client_id" `
	// Domain
	Domain string `json:"domain" `
	// IPv4
	IPv4 string `json:"i_pv" `
	// IPv6
	IPv6 string `json:"i_pv" `
	// Licenses
	Licenses StringArray `json:"licenses" `
	// Secret
	Secret string `json:"secret" `
	// Redirect
	Redirect string `json:"redirect" `
	// Scopes
	Scopes StringArray `json:"scopes" `
	// Flags
	Flags StringArray `json:"flags" `
	// ClientProps
	ClientProps ClientProps `json:"client_props" `
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

func (c Client) Maps() map[string]interface{} {
	maps := c.ModelAbstract.Maps()
	// ClientId
	maps["client_id"] = &c.ClientId
	// Domain
	maps["domain"] = &c.Domain
	// IPv4
	maps["i_pv"] = &c.IPv4
	// IPv6
	maps["i_pv"] = &c.IPv6
	// Licenses
	maps["licenses"] = &c.Licenses
	// Secret
	maps["secret"] = &c.Secret
	// Redirect
	maps["redirect"] = &c.Redirect
	// Scopes
	maps["scopes"] = &c.Scopes
	// Flags
	maps["flags"] = &c.Flags
	// ClientProps
	maps["client_props"] = &c.ClientProps
	return maps
}

// Fields extract of fields from map
func (c Client) Fields(fields ...string) ([]string, []interface{}) {
	return ExtractFieldsFromMap(c.Maps(), fields...)
}

// FromJson data as []byte or io.Reader
func (c *Client) FromJson(data interface{}) error {
	return FromJson(c, data)
}

func (Client) TableName() string {
	return "clients"
}

// PrimaryName primary field name
func (Client) PrimaryName() string {
	return "client_id"
}

// PrimaryValue primary value
func (c Client) PrimaryValue() string {
	return c.ClientId
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

func (c ClientProps) Maps() map[string]interface{} {
	return map[string]interface{}{
		// FullName
		"full_name": &c.FullName,
		// Email
		"email": &c.Email,
	}
}

// Fields extract of fields from map
func (c ClientProps) Fields(fields ...string) ([]string, []interface{}) {
	return ExtractFieldsFromMap(c.Maps(), fields...)
}

// FromJson data as []byte or io.Reader
func (c *ClientProps) FromJson(data interface{}) error {
	return FromJson(c, data)
}
