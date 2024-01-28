package server

import (
	"log"
	"net/http"
	"time"

	"random-projects.net/crayos-backend/game"
	"random-projects.net/crayos-backend/meta"

	"github.com/gorilla/websocket"
)

func serveApi(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "api.html")
}

func serveHome(w http.ResponseWriter, r *http.Request) {
	log.Println(r.URL)
	if r.URL.Path != "/" {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	http.ServeFile(w, r, "home.html")
}

var websocketUpgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(*http.Request) bool {
		return true
	},
}

func acceptPlayerWebsocket(w http.ResponseWriter, r *http.Request) {

	conn, err := websocketUpgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	game.CreatePlayer(conn)
}

var server *http.Server

func Setup() {
	http.HandleFunc("/", serveHome)
	http.HandleFunc("/api", serveApi)
	http.HandleFunc("/ws", acceptPlayerWebsocket)

	server = &http.Server{
		Addr:              *meta.FLAG_ADDR,
		ReadHeaderTimeout: 3 * time.Second,
	}
}

func Run() error {
	return server.ListenAndServe()
}
