package main

import (
	"encoding/json"

	chronograph "github.com/gebv/hey"
)

const (
	NotificationDataType = "1"

	BookmarkDataType = "2"
)

func NewNotification(
	title string,
	body string,
) chronograph.Event {
	return chronograph.Event{
		DataType: NotificationDataType,
		Data: NotificationData{
			Title: title,
			Body:  body,
		}.Marshal(),
	}
}

type NotificationData struct {
	Title string
	Body  string
}

func (n NotificationData) Marshal() []byte {
	bts, _ := json.Marshal(n)
	return bts
}

type Bookmark struct {
	Bookmarked bool
}

func (n Bookmark) Marshal() []byte {
	bts, err := json.Marshal(n)
	if err != nil {
		panic(err)
	}
	return bts
}

func NewBookmark(userID string, eventID string, bookmarked bool) chronograph.RelatedData {
	return chronograph.RelatedData{
		UserID:   userID,
		EventID:  eventID,
		DataType: BookmarkDataType,
		Data:     Bookmark{Bookmarked: bookmarked}.Marshal(),
	}
}
