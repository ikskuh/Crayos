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
	if session.Flags&SESSION_JOINABLE == 0 {
		new.SendChan <- &JoinSessionFailedEvent{
			Reason: "Session is already running.",
		}
		return
	}

	new.Session = session
	session.Players[new] = true

	new.SendChan <- &EnterSessionEvent{
		SessionId: session.Id,
	}

	session.BroadcastPlayers(new, nil)

	new.SendChan <- &ChangeGameViewEvent{
		View: GAME_VIEW_LOBBY,
	}

}

func (session *Session) Broadcast(msg Message) {
	for player := range session.Players {
		player.SendChan <- msg
	}
}

func (session *Session) BroadcastPlayers(added_player *Player, removed_player *Player) {
	nicknames := make([]string, len(session.Players))

	i := 0
	for k := range session.Players {
		nicknames[i] = k.NickName
		i++
	}

	evt := PlayersChangedEvent{
		Players:       nicknames,
		AddedPlayer:   nil,
		RemovedPlayer: nil,
	}

	if added_player != nil {
		evt.AddedPlayer = &added_player.NickName
	}
	if removed_player != nil {
		evt.RemovedPlayer = &removed_player.NickName
	}

	session.Broadcast(&evt)
}

func (session *Session) PumpEvents() *PlayerMessage {
	for len(session.Players) > 0 {
		select {
		case pmsg := <-session.InboundDataChan:
			return &pmsg

		case new := <-session.JoinChan:
			session.AddPlayer(new)

		case old := <-session.LeaveChan:

			// TODO(fqu): Handle dropping players out of active session!
			// In the Lobby, it's totally fine to join/leave all the time

			log.Println("Player leaves", old)
			delete(session.Players, old)

			session.BroadcastPlayers(nil, old)

			// case <- ticker.C:
			// 	player.ws.SetWriteDeadline(time.Now().Add(writeWait))
			// 	if err := player.ws.WriteMessage(websocket.PingMessage, nil); err != nil {
			// 		return
			// 	}
		}
	}
	return nil
}

func (session *Session) Run() {
	log.Println("Starting ", session.Id, " opened")
	defer log.Println("Session ", session.Id, " closed")

	for len(session.Players) > 0 {

		// show lobby
		var startGame = false
		for startGame {
			pmsg := session.PumpEvents()
			if pmsg == nil {
				return
			}

			switch msg := pmsg.Message.(type) {
			case *UserCommand:
				startGame = (msg.Action == USER_ACTION_START_GAME) && len(session.Players) >= 2
			}
		}

		for current_player := range session.Players {

			// determine painter
			_ = current_player

			// change view for all, clear current painting
			var painting_time_not_up = false
			for painting_time_not_up {
				pmsg := session.PumpEvents()
				if pmsg == nil {
					return
				}

				switch msg := pmsg.Message.(type) {
				// case *Timeout:
				// 	break
				default:
					_ = msg
				}

			}

			// enter sticker stage
			var stickers_not_placed = false
			for stickers_not_placed {
				pmsg := session.PumpEvents()
				if pmsg == nil {
					return
				}

				switch msg := pmsg.Message.(type) {
				//case Timeout:
				//	break
				default:
					_ = msg
				}

			}

			// show picture/showcase
			for painting_time_not_up {
				pmsg := session.PumpEvents()
				if pmsg == nil {
					return
				}

				switch msg := pmsg.Message.(type) {
				// case Timeout:
				// 	break
				default:
					_ = msg
				}

			}
		}

		// show art gallery with voting
		var gallery_time_not_up_and_players_not_finished = false
		for gallery_time_not_up_and_players_not_finished {
			pmsg := session.PumpEvents()
			if pmsg == nil {
				return
			}

			switch msg := pmsg.Message.(type) {
			// case Timeout:
			//	break
			default:
				_ = msg
			}
		}

		// higjlight winner

		// show art gallery with winner
		for gallery_time_not_up_and_players_not_finished {
			pmsg := session.PumpEvents()
			if pmsg == nil {
				return
			}

			switch msg := pmsg.Message.(type) {
			// case Timeout:
			//	break
			default:
				_ = msg
			}
		}
		pmsg := session.PumpEvents()
		if pmsg == nil {
			return
		}

		log.Println("Handle message from [", pmsg.Player.NickName, "]: ", pmsg.Message)

	}
}
