package game

import (
	"fmt"
	"log"
	"time"
)

type PlayerMessage struct {
	Player  *Player
	Message Message
}

type SessionFlags struct {
	Joinable bool
}

type Session struct {
	Id string

	Flags SessionFlags

	HostPlayer *Player

	Players map[*Player]bool

	// Channels:
	InboundDataChan chan PlayerMessage
	JoinChan        chan *Player // receives players that have joined the session
	LeaveChan       chan *Player // receives players that have left  the session
}

var sessions = map[string]*Session{}

func CreateSession(player *Player) *Session {
	session := &Session{
		HostPlayer: player,
		Players:    make(map[*Player]bool),

		InboundDataChan: make(chan PlayerMessage, 256), // buffered channel
		JoinChan:        make(chan *Player),            // synchronous channels
		LeaveChan:       make(chan *Player),            // synchronous channels

		Flags: SessionFlags{
			Joinable: true,
		},
	}
	session.Id = fmt.Sprintf("%p", session)

	session.AddPlayer(player)

	// add player before starting main loop, otherwise it will kill itself automatically

	go session.Run()

	// Register session
	log.Println("Created session...", session.Id)
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

	if !session.Flags.Joinable {
		new.SendChan <- &JoinSessionFailedEvent{
			Reason: "Session is already running.",
		}
		return
	}

	log.Println("Player", new.NickName, "joined session", session.Id)

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

func (session *Session) BroadcastExcept(msg Message, except *Player) {
	for player := range session.Players {
		if player != except {
			player.SendChan <- msg
		}
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

type NotifyTimeout struct {
	timestamp time.Time
}

type NotifyPlayerJoined struct {
}

type NotifyPlayerLeft struct {
	// PlayerMessage.Player is not in the session anymore!
}

func (session *Session) PumpEvents(timeout chan time.Time) *PlayerMessage {

	for len(session.Players) > 0 {
		select {
		case pmsg := <-session.InboundDataChan:
			return &pmsg

		case new := <-session.JoinChan:
			session.AddPlayer(new)

			return &PlayerMessage{
				Message: &NotifyPlayerJoined{},
				Player:  new,
			}

		case old := <-session.LeaveChan:

			log.Println("Player", old.NickName, "left session", session.Id)
			delete(session.Players, old)

			session.BroadcastPlayers(nil, old)

			return &PlayerMessage{
				Message: &NotifyPlayerLeft{},
				Player:  old,
			}

		case t := <-timeout:
			return &PlayerMessage{
				Player:  nil,
				Message: &NotifyTimeout{timestamp: t},
			}
		}
	}
	return nil
}

func broadcastPlayerReadyState(s *Session, m map[*Player]bool) {
	log.Println("Sending PlayerReadyState")
	readyMap := make(map[string]bool)
	for p, b := range m {
		readyMap[p.NickName] = b
	}
	s.Broadcast(&PlayerReadyChangedEvent{
		Players: readyMap,
	})
}

func (session *Session) Run() {
	log.Println("Starting ", session.Id, " opened")
	defer log.Println("Session ", session.Id, " closed")

	no_timeout := make(chan time.Time) // pass when no timeout is required

	for len(session.Players) > 0 {

		// show lobby
		// startGame := false
		// playersReady := false
		playersReadyMap := make(map[*Player]bool)
		// len(session.Players) < 2 && !playersReady && !startGame
		for true {
			pmsg := session.PumpEvents(no_timeout)
			if pmsg == nil {
				return
			}

			switch msg := pmsg.Message.(type) {
			case *UserCommand:
				switch msg.Action {
				case USER_ACTION_SET_READY:
					playersReadyMap[pmsg.Player] = true
				case USER_ACTION_SET_NOT_READY:
					playersReadyMap[pmsg.Player] = false
				}
				broadcastPlayerReadyState(session, playersReadyMap)
			}
		}

		for current_player := range session.Players {

			// determine painter
			_ = current_player

			// change view for all, clear current painting
			var painting_time_not_up = false
			for painting_time_not_up {
				pmsg := session.PumpEvents(no_timeout)
				if pmsg == nil {
					return
				}

				switch msg := pmsg.Message.(type) {
				// case *Timeout:
				// 	break

				case *NotifyTimeout:

					log.Println("message timeout received")

				// Forward painting actions
				case *SetPaintingCommand:
					session.BroadcastExcept(&PaintingChangedEvent{
						Path: msg.Path,
					}, pmsg.Player)

				default:
					_ = msg
				}

			}

			// enter sticker stage
			var stickers_not_placed = false
			for stickers_not_placed {
				pmsg := session.PumpEvents(no_timeout)
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
				pmsg := session.PumpEvents(no_timeout)
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
			pmsg := session.PumpEvents(no_timeout)
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
			pmsg := session.PumpEvents(no_timeout)
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
		// pmsg := session.PumpEvents(no_timeout)
		// log.Println("Handle message from [", pmsg.Player.NickName, "]: ", pmsg.Message)
	}
}
