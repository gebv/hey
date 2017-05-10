package main

import chronograph "github.com/zhuharev/chronograph"

const (
	NotificationDataType chronograph.DataType = 1

	BookmarkDataType chronograph.DataType = 2
)

func init() {
	chronograph.RegDataType(NotificationDataType, func() interface{} {
		return &NotificationData{}
	})

	chronograph.RegDataType(BookmarkDataType, func() interface{} {
		return &Bookmark{}
	})
}

func NewNotification(
	title string,
	body string,
) chronograph.Event {
	return chronograph.Event{
		DataType: NotificationDataType,
		Data: &NotificationData{
			Title: title,
			Body:  body,
		},
	}
}

type NotificationData struct {
	Title string
	Body  string
}

type Bookmark struct {
	Bookmarked bool
}

func NewBookmark(userID string, eventID string, bookmarked bool) chronograph.RelatedData {
	return chronograph.RelatedData{
		UserID:   userID,
		EventID:  eventID,
		DataType: BookmarkDataType,
		Data:     &Bookmark{Bookmarked: bookmarked},
	}
}
