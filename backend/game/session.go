package game

import (
	"fmt"
	"log"
)

const (
	SESSION_JOINABLE = (1 << 0)
)

type PlayerMessage struct {
	Player  *Player
	Message Message
}

type Session struct {
	Id string

	Flags int

	HostPlayer *Player

	Players map[*Player]bool

	// Channels:
	InboundDataChan chan PlayerMessage
	JoinChan        chan *Player
	LeaveChan       chan *Player
}

var sessions = map[string]*Session{}

func CreateSession(player *Player) *Session {
	session := &Session{
		HostPlayer: player,
		Players:    make(map[*Player]bool),

		InboundDataChan: make(chan PlayerMessage, 256), // buffered channel
		JoinChan:        make(chan *Player),            // synchronous channels
		LeaveChan:       make(chan *Player),            // synchronous channels

		Flags: SESSION_JOINABLE,
	}
	session.Id = fmt.Sprintf("%p", session)

	session.AddPlayer(player)

	// add player before starting main loop, otherwise it will kill itself automatically

	go session.Run()

	// Register session
	sessions[session.Id] = session

	return session
}

func FindSession(id string) *Session {
	// TODO(fqu): Should be mutex checked
	session, ok := sessions[id]
	if ok {
		return session
	} else {
		return nil
	}
}

func (session *Session) Destroy() {
	delete(sessions, session.Id)
}

func (session *Session) AddPlayer(new *Player) {

	log.Println("Player joins", new)
	if session.Flags&SESSION_JOINABLE != 0 {

		session.Players[new] = true

		new.SendChan <- &EnterSessionEvent{
			SessionId: session.Id,
		}

		{
			nicknames := make([]string, len(session.Players))

			i := 0
			for k := range session.Players {
				nicknames[i] = k.NickName
				i++
			}

			for player := range session.Players {
				player.SendChan <- &PlayersChangedEvent{
					Players: nicknames,
					Joined:  &new.NickName,
				}
			}
		}

		new.SendChan <- &ChangeGameViewEvent{
			ViewID: GAMEVIEW_LOBBY,
		}

	} else {
		new.SendChan <- &JoinSessionFailedEvent{
			Reason: "Session is already running.",
		}
	}
}

func (session *Session) Run() {
	log.Println("Starting ", session.Id, " opened")
	defer log.Println("Session ", session.Id, " closed")
	for len(session.Players) > 0 {

		select {
		case pmsg, ok := <-session.InboundDataChan:

			log.Println("inbound message from player", pmsg, ok)

		case new := <-session.JoinChan:
			session.AddPlayer(new)

		case old := <-session.LeaveChan:
			log.Println("Player leaves", old)

			delete(session.Players, old)

			// case <- ticker.C:
			// 	player.ws.SetWriteDeadline(time.Now().Add(writeWait))
			// 	if err := player.ws.WriteMessage(websocket.PingMessage, nil); err != nil {
			// 		return
			// 	}
		}
	}

}
