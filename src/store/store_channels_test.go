package store

import (
	"testing"
	"models"
	"github.com/satori/go.uuid"
)


func zzzTestFindAllThreadsFromPathAndCreateLast(t *testing.T) {
	channel := createChannel(t)

	tx, _ := _s.db.Begin()

	dto := models.NewThreadDTO()
	dto.ClientId = channel.ClientId.String()
	dto.ChannelId = channel.PrimaryValue().String()
	dto.ExtId = channel.ExtId + ":thread" + uuid.NewV1().String()
	dto.EventCreator = "193b7a9c-42ad-456e-8886-aa6ae8ebcf17"
	// dto.CreatedEventId = uuid.NewV1().String()
	dto.Owners.Add("193b7a9c-42ad-456e-8886-aa6ae8ebcf17")
	dto.ExtFlags.Add("flag1")
	dto.ExtFlags.Add("flag2")
	dto.ExtProps["key"] = "value"
	dto.Tx = tx

	_, err := _s.Get("thread").(*ThreadStore).FindAllThreadsFromPathAndCreateLast(dto)

	if err != nil {
		t.Error(err)
		tx.Rollback()
		return
	} else if err := tx.Commit(); err != nil {
		t.Error(err)
		return
	}

	// Chech created thread
	dto.Tx = nil
	_thread, err := _s.Get("thread").(*ThreadStore).GetOne(dto)

	if err != nil {
		t.Error(err)
		return
	}

	thread := _thread.(*models.Thread)

	if thread.ExtId != dto.ExtId {
		t.Error("thread ExtId is not correct")
		return
	}

	if !uuid.Equal(thread.ChannelId, uuid.FromStringOrNil(dto.ChannelId)) {
		t.Error("thread ChannelId is not correct")
		return
	}

	if !uuid.Equal(thread.ClientId, uuid.FromStringOrNil(dto.ClientId)) {
		t.Error("thread ClientId is not correct")
		return
	}

	if !uuid.Equal(thread.Owners[0], uuid.FromStringOrNil(dto.Owners[0])) {
		t.Error("thread Owners is not correct")
		return
	}

	_flags := models.StringArray(thread.ExtFlags)
	if !_flags.IsExist("flag1") || !_flags.IsExist("flag2") {
		t.Error("thread ExtFlags is not correct")
		return
	}

	if thread.ExtProps["key"] != "value" {
		t.Error("thread ExtProps is not correct")
		return	
	}

	if thread.Depth != 2 {
		t.Error("thread Depth is not correct")
		return	
	}
}