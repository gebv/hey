package main

import (
	"encoding/json"
	"log"

	chronograph "github.com/gebv/hey"
)

func main() {
	chrono, err := chronograph.New()
	if err != nil {
		log.Fatalln(err)
	}

	// создаём пользователя 1
	user1 := chronograph.User{}
	err = chrono.NewUser(&user1)
	checkErr(err)

	// создаём трэд нотификаций пользователя
	notify1 := chronograph.Thread{
		ThreadlineEnabled: true,
	}
	checkErr(err)

	// create thread for user 1 notifications
	err = chrono.NewThread(&notify1)
	checkErr(err)

	// трэд можно читать без подписки на него,
	// но если мы хотим знать последнее прочитанное уведомление,
	// должны подписаться
	// подписываем пользователя 1 на уведомления
	err = chrono.Observe(user1.UserID, notify1.ThreadID)
	checkErr(err)

	obs, err := chrono.Observers(notify1.ThreadID, 0, 1)
	checkErr(err)
	if len(obs) != 1 {
		log.Fatalln("observe not work", len(obs))
	}

	// создаём событие в трэде
	note := NewNotification("Событие", "Новое сообщение!")
	note.ThreadID = notify1.ThreadID
	err = chrono.NewEvent(&note)
	checkErr(err)

	// достаем последние события
	events, err := chrono.RecentActivity(user1.UserID, notify1.ThreadID, 10, 0)
	checkErr(err)
	if len(events) != 1 {
		log.Fatalln("длинна событий != 1", len(events))
	}
	var nd NotificationData
	err = json.Unmarshal(events[0].Data, &nd)
	checkErr(err)
	if nd.Title != "Событие" {
		log.Fatalln("ошибка ")
	}

	// добавляем произвольные данные к событию, видные только данному
	// пользователю
	bookmark := NewBookmark(user1.UserID, note.EventID, true)
	err = chrono.SetRelatedData(&bookmark)
	checkErr(err)

	// чтобы достать события с пользовательской информацией нужно сначала
	// достать события, а затем вызвать
	eventObseres, err := chrono.GetRelatedDatas(user1.UserID, note)
	checkErr(err)
	if len(eventObseres) != 1 {
		log.Fatalln("не хватает")
	}

	var bm Bookmark
	err = json.Unmarshal(eventObseres[0].RelatedData.Data, &bm)
	checkErr(err)
	if !bm.Bookmarked {
		log.Fatalln("не добавилось в избранное")
	}

	log.Println("done")
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}
