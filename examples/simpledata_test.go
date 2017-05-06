package examples

import (
	"testing"
	"time"

	"github.com/gebv/hey"
	"github.com/stretchr/testify/assert"
	msgpack "gopkg.in/vmihailenco/msgpack.v2"
)

func TestSimpleData(t *testing.T) {
	tick := time.Now()
	obj := NewThreadSimpleData("a", 1234, tick)
	obj.ThreadID = "uri:threads:example"

	network, err := msgpack.Marshal(obj)
	assert.NoError(t, err, "encode")

	got := &hey.Thread{}
	err = msgpack.Unmarshal(network, got)
	assert.NoError(t, err, "decode")

	assert.Equal(t, obj.ThreadID, got.ThreadID)
	assert.EqualValues(t, obj.Data.(*ThreadSimpleData).A, got.Data.(*ThreadSimpleData).A)
	assert.EqualValues(t, obj.Data.(*ThreadSimpleData).B, got.Data.(*ThreadSimpleData).B)
	assert.EqualValues(t, obj.Data.(*ThreadSimpleData).C, got.Data.(*ThreadSimpleData).C)
}
