package game

import (
	"log"
	"reflect"
	"time"

	"github.com/gorilla/websocket"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 512
)

// client.hub.register <- client

type Player struct {
	ws *websocket.Conn

	Session *Session

	NickName string

	SendChan chan Message
}

func CreatePlayer(ws *websocket.Conn) *Player {

	player := &Player{
		ws:       ws,
		Session:  nil,
		SendChan: make(chan Message, 256),
		NickName: "Anonymouse",
	}

	player.SendChan <- &ChangeGameViewEvent{
		ViewID: GAMEVIEW_TITLE,
	}

	go player.writePump()
	go player.readPump()

	return player
}

// Pumps messages from websocket to the session or creates/joins a new session.
func (player *Player) readPump() {
	defer func() {
		if player.Session != nil {
			player.Session.LeaveChan <- player
		}
		player.ws.Close()
	}()
	player.ws.SetReadLimit(maxMessageSize)
	player.ws.SetReadDeadline(time.Now().Add(pongWait))
	player.ws.SetPongHandler(func(string) error {
		player.ws.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})
	for {
		_, raw_message, err := player.ws.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}

		var message Message //  = bytes.TrimSpace(bytes.Replace(message, newline, space, -1))

		msg, err := DeserializeMessage(raw_message)

		log.Println("message from websocket", string(raw_message), reflect.TypeOf(msg), msg, err)

		if err != nil {
			log.Println("failed to read message from client: ", err)
			return
		}

		if player.Session != nil {
			log.Println("Forward message to session ", msg)
			player.Session.InboundDataChan <- PlayerMessage{
				Player:  player,
				Message: message,
			}
		} else {
			switch v := msg.(type) {
			case *CreateSessionCommand:
				log.Println("new session?")
				player.NickName = v.NickName

				log.Println("nick set:", v.NickName, player.NickName)

				session := CreateSession(player)

				log.Println("Create new session...", session.Id)

			case *JoinSessionCommand:
				log.Println("join session?")
				player.NickName = v.NickName

				session := FindSession(v.SessionId)

				if session != nil {
					log.Println("found session: ", session.Id)
					session.JoinChan <- player
				} else {
					log.Println("didn't find session", v.SessionId)
					player.SendChan <- &JoinSessionFailedEvent{
						Reason: "Session does not exist",
					}
				}

			default:
				log.Println("Bad command, dropping client, type was ", reflect.TypeOf(msg))
				return
			}
		}
	}
}

// Pumps messages from Player.SendChan to the websocket.
func (player *Player) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		player.ws.Close()
	}()
	for {
		select {
		case message, ok := <-player.SendChan:
			player.ws.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// The hub closed the channel.
				player.ws.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			encoded_msg, err := SerializeMessage(message)
			if err != nil {
				log.Println("failed to serialize message for client: ", err, message)
				return
			}

			err = player.ws.WriteMessage(websocket.TextMessage, encoded_msg)
			if err != nil {
				log.Println("failed to send message to client: ", err)
				return
			}

		case <-ticker.C:
			player.ws.SetWriteDeadline(time.Now().Add(writeWait))
			if err := player.ws.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}
