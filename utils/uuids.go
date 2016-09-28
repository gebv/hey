package utils

import (
	"bytes"
	"database/sql/driver"
	"encoding/json"
	"log"
	"strings"

	"github.com/satori/go.uuid"
)

// NewUUID4 returns uuid version 4
func NewUUID4() string {
	return strings.Replace(uuid.NewV4().String(), "-", "", -1)
}

func UUIDSFrom(arr []uuid.UUID) UUIDS {
	return UUIDS(arr)
}

type UUIDS []uuid.UUID

func (m *UUIDS) Scan(value interface{}) error {
	_value := value.([]byte)
	_value = _value[1 : len(_value)-1]

	m.FromArray(strings.Split(string(_value), ","))

	return nil
}

func (m UUIDS) Value() (driver.Value, error) {
	if len(m) == 0 {
		return string("{}"), nil
	}

	b, err := json.Marshal(m)
	if err != nil {
		return string("{}"), err
	}

	return "{" + string(b)[1:len(b)-1] + "}", nil
}

func (c *UUIDS) FromArray(v interface{}) *UUIDS {
	switch v.(type) {
	case []uuid.UUID:
		for _, _v := range v.([]uuid.UUID) {
			c.Add(_v)
		}
	case []string:
		for _, _v := range v.([]string) {
			c.Add(_v)
		}
	default:
		log.Printf("[WRN] Not supported type %T", v)
	}

	return c
}

func (c *UUIDS) IsExist(v uuid.UUID) bool {
	for _, value := range *c {

		if bytes.Equal(v.Bytes(), value.Bytes()) {

			return true
		}
	}

	return false
}

func (c *UUIDS) Add(v interface{}) *UUIDS {
	switch v.(type) {
	case uuid.UUID:
		if c.IsExist(v.(uuid.UUID)) {
			return c
		}

		*c = append(*c, v.(uuid.UUID))
	case string:
		if c.IsExist(uuid.FromStringOrNil(v.(string))) {
			return c
		}

		*c = append(*c, uuid.FromStringOrNil(v.(string)))
	default:
		log.Printf("[WRN] Not supported type %T", v)
	}

	return c
}
