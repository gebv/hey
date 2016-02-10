package store

import (
	"testing"
	"models"
	"github.com/satori/go.uuid"
)


func createChannel(t *testing.T) *models.Channel{
	tx, _ := _s.db.Begin()

	dto := models.NewChannelDTO()
	dto.ExtId = "channel"+uuid.NewV1().String()
	dto.ExtFlags.Add("type:demo")
	dto.ClientId = "b4c8dd5b-852c-460a-9b4a-26109f9162a2"
	dto.Owners.Add("193b7a9c-42ad-456e-8886-aa6ae8ebcf17")
	dto.ExtFlags.Add("flag1")
	dto.ExtFlags.Add("flag2")
	dto.ExtProps["key"] = "value"
	dto.Tx = tx
	channel, err := _s.Get("channel").(*ChannelStore).Create(dto)

	defer tx.Rollback()

	if err != nil {
		t.Error(err)
		return nil
	} else if err := tx.Commit(); err != nil {
		t.Error(err)
		return nil
	}

	return channel.(*models.Channel)
}
