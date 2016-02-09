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

	if err != nil {
		t.Error(err)
		tx.Rollback()
		return nil
	} else if err := tx.Commit(); err != nil {
		t.Error(err)
		return nil
	}

	return channel.(*models.Channel)
}

func zzzTestCreateChannel(t *testing.T) {
	channel := createChannel(t)

	// Проверка созданного канала
	channelDTO := models.NewChannelDTO()
	channelDTO.ExtId = channel.ExtId
	channelDTO.ClientId = channel.ClientId.String()
	_createdChannel, err := _s.Get("channel").(*ChannelStore).GetOne(channelDTO)

	if err != nil {
		t.Error(err)
		return
	}

	createdChannel := _createdChannel.(*models.Channel)

	if !uuid.Equal(createdChannel.ClientId, channel.ClientId) {
		t.Error("channel ClientId is not correct")
		return	
	}

	if createdChannel.ExtId != channel.ExtId {
		t.Error("channel ExtId is not correct")
		return	
	}

	_flags := models.StringArray(createdChannel.ExtFlags)
	if !_flags.IsExist("flag1") ||
		!_flags.IsExist("flag2") {
		t.Error("channel ExtFlags is not correct")
		return	
	}

	if len(createdChannel.Owners) == 0 || !uuid.Equal(createdChannel.Owners[0], channel.Owners[0]) {
		t.Error("channel Owners is not correct")
		return	
	}

	if createdChannel.ExtProps["key"] != "value" {
		t.Error("channel ExtFlags is not correct")
		return	
	}

	// Проверка связанного с каналом channel_watchers
	_channelWatchers := models.NewChannelWatcher()
	fields, modelValues := _channelWatchers.Fields()
	query := SqlSelect(_channelWatchers.TableName(), fields)
	query += " WHERE client_id = ? AND channel_id = ?"
	query = FormateToPQuery(query)

	if err := _s.db.QueryRow(query, channel.ClientId, channel.PrimaryValue()).Scan(modelValues...); err != nil {
		t.Error(err)
		return
	}

	if !uuid.Equal(_channelWatchers.UserId, channel.Owners[0]) {
		t.Error("channel_watchers UserId is not correct")
		return	
	}

	if _channelWatchers.Unread != 0 {
		t.Error("channel_watchers Unread is not correct")
		return	
	}

	// Проверка связанного с каналом channel_counter
	_channelCounter := models.NewChannelCounter()
	fields, modelValues = _channelCounter.Fields()
	query = SqlSelect(_channelCounter.TableName(), fields)
	query += " WHERE client_id = ? AND channel_id = ?"
	query = FormateToPQuery(query)

	if err := _s.db.QueryRow(query, channel.ClientId, channel.PrimaryValue().String()).Scan(modelValues...); err != nil {
		t.Error(err)
		return
	}

	if _channelCounter.CounterEvents != 0 {
		t.Error("channel_counter CounterEvents is not correct")
		return	
	}

	// Проверка связанного root_thread
	relatedThreadDTO := models.NewThreadDTO()
	relatedThreadDTO.ClientId = channel.ClientId.String()
	relatedThreadDTO.ExtId = channel.ExtId
	_relatedThread, err := _s.Get("thread").(*ThreadStore).GetOne(relatedThreadDTO)

	if err != nil {
		t.Error(err)
		return
	}

	

	relatedThread := _relatedThread.(*models.Thread)

	// RootThread for channel
	if !uuid.Equal(channel.RootThreadId, relatedThread.ThreadId) {
		t.Error("channels RootThreadId is not correct")
		return	
	}

	if !uuid.Equal(relatedThread.ClientId, channel.ClientId) {
		t.Error("threads ClientId is not correct")
		return	
	}

	if !uuid.Equal(relatedThread.ChannelId, channel.PrimaryValue()) {
		t.Error("threads ChannelId is not correct")
		return	
	}

	if relatedThread.ExtId != channel.ExtId {
		t.Error("threads ExtId is not correct")
		return	
	}

	_flags = models.StringArray(relatedThread.ExtFlags)
	if !_flags.IsExist("flag1") ||
		!_flags.IsExist("flag2") {
		t.Error("threads ExtFlags is not correct")
		return	
	}

	if relatedThread.ExtProps["key"] != "value" {
		t.Error("threads ExtProps is not correct")
		return	
	}

	if len(relatedThread.Owners) == 0 || !uuid.Equal(relatedThread.Owners[0], channel.Owners[0]) {
		t.Error("threads Owners is not correct")
		return		
	}

	if relatedThread.Depth != 1 {
		t.Error("threads Depth is not correct")
		return		
	}

	// Связанный с потоком thread_counter
	_threadCounter := models.NewThreadCounter()
	fields, modelValues = _threadCounter.Fields()
	query = SqlSelect(_threadCounter.TableName(), fields)
	query += " WHERE client_id = ? AND thread_id = ?"
	query = FormateToPQuery(query)

	if err := _s.db.QueryRow(query, channel.ClientId, relatedThread.PrimaryValue().String()).Scan(modelValues...); err != nil {
		t.Error(err)
		return
	}

	if _threadCounter.CounterEvents != 0 {
		t.Error("thread_counter CounterEvents is not correct")
		return
	}

	// Связанный с потоком thread_watchers
	_threadWatchers := models.NewThreadWatcher()
	fields, modelValues = _threadWatchers.Fields()
	query = SqlSelect(_threadWatchers.TableName(), fields)
	query += " WHERE client_id = ? AND thread_id = ?"
	query = FormateToPQuery(query)

	if err := _s.db.QueryRow(query, channel.ClientId, relatedThread.PrimaryValue().String()).Scan(modelValues...); err != nil {
		t.Error(err)
		return
	}

	if !uuid.Equal(_threadWatchers.UserId, relatedThread.Owners[0]) {
		t.Error("thread_watchers UserId is not correct")
		return
	}

	if _threadWatchers.Unread != 0 {
		t.Error("thread_watchers Unread is not correct")
		return
	}
}