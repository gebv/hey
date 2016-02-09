package models

import (
	"encoding/json"
	"github.com/golang/glog"
	"github.com/satori/go.uuid"
	"time"
)

type ModelFields interface {
	Maps() map[string]interface{}
	Fields(fields ...string) ([]string, []interface{})
	FromJson(data interface{}) error
}

type Transformer interface {
	TransformTo(interface{}) error
	TransformFrom(interface{}) error
}

// Все модели являются транформерами и обладают вспомогательными функциями ModelFields
type ModelAbstractInterface interface {
	ModelFields
	Transformer

	BeforeCreate()
	BeforeSave()
	BeforeDelete()

	PrimaryName() string
	PrimaryValue() uuid.UUID
	TableName() string
}

func (c *ModelAbstract) BeforeCreate() {
	c.CreatedAt = time.Now()
}

func (c *ModelAbstract) BeforeSave() {
	c.UpdatedAt = time.Now()
}

func (c *ModelAbstract) BeforeDelete() {
	c.IsRemoved = true
}

// ResponseDTO

func (c *ResponseDTO) ToJson() []byte {
	_b, err := json.Marshal(c)

	if err != nil {
		glog.Errorf("Marshal ResponseDTO error, %s", err)
		return []byte(`{}`)
	}

	return _b
}
