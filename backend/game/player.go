package game

import (
	"log"
	"reflect"
	"sync"
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
	mu     sync.Mutex
	closed bool

	ws *websocket.Conn

	Session *Session

	NickName string

	// NOTE(fqu):
	// Must be a buffered channel, as we have to be able to send
	// non-blockingly.
	sendChan chan []byte
}

func CreatePlayer(ws *websocket.Conn) *Player {
	player := &Player{
		ws:       ws,
		Session:  nil,
		NickName: "Anonymouse",

		sendChan: make(chan []byte, 256),
	}

	player.Send(&ChangeGameViewEvent{
		View: GAME_VIEW_TITLE,
	})

	go player.writePump()
	go player.readPump()

	return player
}

func (player *Player) Send(msg Message) {
	if player.closed {
		return
	}
	encoded_msg, err := SerializeMessage(msg)
	if err != nil {
		log.Fatalln("failed to serialize message for client: ", err, msg)
	}
	player.sendChan <- encoded_msg
}

func (player *Player) Close() {

	player.mu.Lock()
	defer player.mu.Unlock()

	if player.closed {
		return
	}

	if player.Session != nil {
		player.Session.LeaveChan <- player
		player.Session = nil
	}
	player.ws.Close()
	close(player.sendChan)

	player.closed = true
}

// Pumps messages from websocket to the session or creates/joins a new session.
func (player *Player) readPump() {
	defer player.Close()
	// player.ws.SetReadLimit(maxMessageSize)
	player.ws.SetReadDeadline(time.Now().Add(pongWait))
	player.ws.SetPongHandler(func(string) error {
		player.ws.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})
	for {
		_, raw_message, err := player.ws.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("websocket error: %v", err)
			}
			break
		}

		// log.Println("raw message from websocket", string(raw_message))

		msg, err := DeserializeMessage(raw_message)

		// log.Println("message from websocket", string(raw_message), reflect.TypeOf(msg), msg, err)

		if err != nil {
			log.Println("failed to read message from client: ", err)
			return
		}
		if player.Session != nil {
			// log.Println("Forward message to session ", msg)
			player.Session.InboundDataChan <- PlayerMessage{
				Player:  player,
				Message: msg,
			}
		} else {
			switch v := msg.(type) {
			case *CreateSessionCommand:
				if v.NickName == "" {
					player.Send(&JoinSessionFailedEvent{
						Reason: TEXT_ERROR_NICK_EMPTY,
					})
				} else if len(v.NickName) > LIMIT_MAX_NICKNAME_LEN {
					player.Send(&JoinSessionFailedEvent{
						Reason: TEXT_ERROR_NICK_TOO_LONG,
					})
				} else {
					player.NickName = v.NickName

					session := CreateSession(player)

					_ = session
				}

			case *JoinSessionCommand:
				if v.SessionId == "" {
					player.Send(&JoinSessionFailedEvent{
						Reason: TEXT_ERROR_SESSION_EMPTY,
					})
				} else if v.NickName == "" {
					player.Send(&JoinSessionFailedEvent{
						Reason: TEXT_ERROR_NICK_EMPTY,
					})
				} else if len(v.NickName) > LIMIT_MAX_NICKNAME_LEN {
					player.Send(&JoinSessionFailedEvent{
						Reason: TEXT_ERROR_NICK_TOO_LONG,
					})
				} else {
					player.NickName = v.NickName

					session := FindSession(v.SessionId)

					if session != nil {
						session.JoinChan <- player
					} else {
						log.Println("didn't find session", v.SessionId)
						player.Send(&JoinSessionFailedEvent{
							Reason: TEXT_ERROR_BAD_SESSION,
						})
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
	defer player.Close()

	ticker := time.NewTicker(pingPeriod)
	defer ticker.Stop()

	for {
		select {
		case message, ok := <-player.sendChan:
			player.ws.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// The hub closed the channel.
				player.ws.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			err := player.ws.WriteMessage(websocket.TextMessage, message)
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
