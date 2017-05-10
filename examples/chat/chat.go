package main

import (
	"log"

	chronograph "github.com/zhuharev/chronograph"
)

func main() {
	chrono, err := chronograph.New()
	checkErr(err)

	// создаём пользователя 1
	user1 := chronograph.User{}
	err = chrono.NewUser(&user1)
	checkErr(err)

	// создаём пользователя 2
	user2 := chronograph.User{}
	err = chrono.NewUser(&user2)
	checkErr(err)

	// создаём чат
	chat := chronograph.Thread{}
	err = chrono.NewThread(&chat)
	checkErr(err)

	// подписываемя двумя пользователями на трэд
	for _, user := range []chronograph.User{user1, user2} {
		err = chrono.Observe(user.UserID, chat.ThreadID)
		checkErr(err)
	}

	// создаём новое событие в чате
	message := NewMessage("ilon", "hello all!")
	err = chrono.NewEvent(&message)
	checkErr(err)

	// получаем сообщение из чата
	events, err := chrono.RecentActivity(user1.UserID, chat.ThreadID, 10)
	checkErr(err)
	if len(events) != 1 {
		log.Fatalln("len events not 1")
	}

	// помечаем все события прочитанными
	err = chrono.MarkAsDelivered(user1.UserID, chat.ThreadID)
	checkErr(err)

	// получаем непрочитанные сообщение из чата, должно быть пусто
	events, err = chrono.RecentActivity(user1.UserID, chat.ThreadID, 10)
	checkErr(err)
	if len(events) != 0 {
		log.Fatalln("len events not 0")
	}

}

func checkErr(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}
