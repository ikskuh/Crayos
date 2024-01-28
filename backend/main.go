package main

import (
	"flag"
	"log"

	"random-projects.net/crayos-backend/game"
	"random-projects.net/crayos-backend/meta"
	"random-projects.net/crayos-backend/server"
)

func main() {
	flag.Parse()

	server.Setup()

	log.Println("Ready.")

	if *meta.DEBUG_MODE {
		default_session := game.CreateSession(nil)

		game.SetDebugSession(default_session)

		log.Println("DEFAULT SESSION ID:", default_session.Id)
	}

	err := server.Run()

	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
