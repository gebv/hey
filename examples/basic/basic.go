package main

import (
	"log"

	chronograph "github.com/zhuharev/chronograph"
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
		ThreadID: "notifications:1",
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

	log.Println(note)

	// достаем последние события
	events, err := chrono.RecentActivity(user1.UserID, notify1.ThreadID, 10)
	checkErr(err)
	if len(events) != 1 {
		log.Fatalln("длинна событий != 1", len(events))
	}
	if _, ok := events[0].Data.(*NotificationData); !ok {
		log.Fatalln("ошибка декодирования данных")
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

	if bm, ok := eventObseres[0].RelatedData.Data.(*Bookmark); ok {
		if !bm.Bookmarked {
			log.Fatalln("не добавилось в избранное")
		}
	} else {
		log.Fatalln("ошибка декодинга")
	}

}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}
