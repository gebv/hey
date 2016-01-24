package models

import (
	"bytes"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"io"
	"strings"
)

type UUID [16]byte

// create a new uuid v4
func NewUUID() *UUID {
	u := &UUID{}
	_, err := rand.Read(u[:16])
	if err != nil {
		panic(err)
	}

	u[8] = (u[8] | 0x80) & 0xBf
	u[6] = (u[6] | 0x40) & 0x4f
	return u
}

func (u *UUID) String() string {
	return fmt.Sprintf("%x-%x-%x-%x-%x", u[:4], u[4:6], u[6:8], u[8:10], u[10:])
}

type StringMap map[string]string

func NewStringMap() map[string]string {
	return make(map[string]string)
}

type StringArray []string

func (c *StringArray) AddAsArray(in []string) {
	for _, str := range in {
		c.Add(str)
	}
}

func (c *StringArray) AddAsFields(in string) {
	for _, str := range strings.Fields(in) {
		c.Add(str)
	}
}

func (c *StringArray) Del(str string) *StringArray {
	str = strings.TrimSpace(str)

	if !c.IsExist(str) {
		return c
	}

	for index, value := range *c {

		if bytes.Equal([]byte(str), []byte(value)) {
			(*c)[index] = (*c)[len((*c))-1]
			(*c) = (*c)[:len((*c))-1]

			return c
		}
	}

	return c
}

func (c *StringArray) IsExist(str string) bool {
	str = strings.TrimSpace(str)

	for _, value := range *c {

		if bytes.Equal([]byte(str), []byte(value)) {

			return true
		}
	}

	return false
}

func (c *StringArray) Add(str string) *StringArray {
	str = strings.TrimSpace(str)

	if len(str) == 0 {
		return c
	}

	if !c.IsExist(str) {
		*c = append(*c, str)
	}

	return c
}

//

// extractFieldsFromMap if len(without) == 0 return all fileds
func ExtractFieldsFromMap(m map[string]interface{}, without ...string) (keys []string, fields []interface{}) {
	_without := make(map[string]bool)
	var flagAllFields = len(without) == 0

	for _, v := range without {
		_without[v] = true
	}

	for fieldName, field := range m {
		if !flagAllFields {
			if !_without[fieldName] {
				continue
			}
		}

		keys = append(keys, fieldName)
		fields = append(fields, field)
	}

	return
}

// FromJson extract object from data (io.Reader OR []byte)
func FromJson(obj interface{}, data interface{}) error {
	switch data.(type) {
	case io.Reader:
		decoder := json.NewDecoder(data.(io.Reader))
		return decoder.Decode(obj)
	case []byte:
		return json.Unmarshal(data.([]byte), obj)
	}

	return ErrNotSupported
}
