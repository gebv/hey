package store

import (
	"testing"
	"models"
	"github.com/satori/go.uuid"
)

func TestCreateEvent(t *testing.T) {
	// Создать канал и создать событие в канал на уровень выше

	channel := createChannel(t)

	dto := models.NewEventDTO()
	dto.ClientId = channel.ClientId.String()
	dto.Creator = "demo"
	dto.Thread = channel.ExtId + ":specialthreadformessage"
	dto.ExtFlags.Add("flag1")

	_event, err := _s.Get("event").(*EventStore).Create(dto)

	if err != nil {
		t.Error(err)
		return
	}

	// Связанный поток
	relatedThreadDTO := models.NewThreadDTO()
	relatedThreadDTO.ClientId = channel.ClientId.String()
	relatedThreadDTO.ExtId = channel.ExtId + ":specialthreadformessage"
	_relatedThread, err := _s.Get("thread").(*ThreadStore).GetOne(relatedThreadDTO)

	if err != nil {
		t.Error(err)
		return
	}

	relatedTHread := _relatedThread.(*models.Thread)
	relatedThreadId_DeptTwo := relatedTHread.ThreadId

	event := _event.(*models.Event)
	eventId_DeptTwo := event.EventId

	if !uuid.Equal(relatedTHread.ParentThreadId, channel.RootThreadId) {
		t.Error("threads ParentThreadId is not correct")
		return
	}

	if !uuid.Equal(event.ThreadId, relatedTHread.ThreadId) {
		t.Error("threads ThreadId is not correct")
		return
	}

	if relatedTHread.Depth != 2 {
		t.Error("threads Depth is not correct")
		return	
	}

	// dep 3

	dto = models.NewEventDTO()
	dto.ClientId = "b4c8dd5b-852c-460a-9b4a-26109f9162a2"
	dto.Creator = "demo"
	dto.ParentEventId = event.EventId.String()
	dto.Thread = channel.ExtId + ":specialthreadformessage:moredepth"
	dto.ExtFlags.Add("flag1")

	_event, err = _s.Get("event").(*EventStore).Create(dto)

	if err != nil {
		t.Error(err)
		return
	}

	// Связанный поток
	relatedThreadDTO = models.NewThreadDTO()
	relatedThreadDTO.ClientId = channel.ClientId.String()
	relatedThreadDTO.ExtId = channel.ExtId + ":specialthreadformessage:moredepth"
	_relatedThread, err = _s.Get("thread").(*ThreadStore).GetOne(relatedThreadDTO)

	if err != nil {
		t.Error(err)
		return
	}

	relatedTHread = _relatedThread.(*models.Thread)
	relatedThreadId_DeptTree := relatedTHread.ThreadId

	event = _event.(*models.Event)
	// eventId_DeptTree := event.EventId

	if !uuid.Equal(relatedTHread.ParentThreadId, relatedThreadId_DeptTwo) {
		t.Error("threads ParentThreadId is not correct")
		return
	}

	if !uuid.Equal(event.ThreadId, relatedTHread.ThreadId) {
		t.Error("threads ThreadId is not correct")
		return
	}

	if !uuid.Equal(event.ParentThreadId, relatedThreadId_DeptTwo) {
		t.Error("threads ParentThreadId is not correct")
		return
	}

	if !uuid.Equal(event.ParentEventId, eventId_DeptTwo) {
		t.Error("threads ParentEventId is not correct")
		return
	}

	if relatedTHread.Depth != 3 {
		t.Error("threads Depth is not correct")
		return	
	}

	// Обновленный поток relatedThreadId_DeptTwo

	relatedThreadDTO = models.NewThreadDTO()
	relatedThreadDTO.ClientId = channel.ClientId.String()
	relatedThreadDTO.ExtId = channel.ExtId + ":specialthreadformessage:moredepth"
	_relatedThread, err = _s.Get("thread").(*ThreadStore).GetOne(relatedThreadDTO)

	if err != nil {
		t.Error(err)
		return
	}

	relatedTHread = _relatedThread.(*models.Thread)

	if !uuid.Equal(relatedTHread.RelatedEventId, eventId_DeptTwo) {
		t.Errorf("events RelatedEventId is not correct, %v", relatedTHread.RelatedEventId)
		return	
	}

	// Обновленные события
	dtoEvent := models.NewEventDTO()
	dtoEvent.EventId = eventId_DeptTwo.String()
	dtoEvent.ClientId = channel.ClientId.String()

	_event, err = _s.Get("event").(*EventStore).GetOne(dtoEvent)

	if err != nil {
		t.Error(err)
		return
	}

	event = _event.(*models.Event)

	if !uuid.Equal(event.BranchThreadId, relatedThreadId_DeptTree) {
		t.Errorf("events BranchThreadId is not correct")
		return	
	} 

	
}