package game

import (
	"fmt"
	"log"
	"math/rand"
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

type Role int

const (
	ROLE_PAINTER Role = 0
	ROLE_TROLL   Role = 1
)

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

func broadcastPlayerReadyState(s *Session, m playerSet) {
	log.Println("Sending PlayerReadyState", m)
	readyMap := make(map[string]bool)
	for p, b := range m.items {
		readyMap[p.NickName] = b.value
	}
	s.Broadcast(&PlayerReadyChangedEvent{
		Players: readyMap,
	})
}

func (session *Session) Run() {
	random_source := rand.New(rand.NewSource(time.Now().UnixNano()))

	log.Println("Starting ", session.Id, " opened")
	defer log.Println("Session ", session.Id, " closed")

	no_timeout := make(chan time.Time) // pass when no timeout is required

	for len(session.Players) > 0 {

		// show lobby

		{
			players_ready := createPlayerSetFromMap(session.Players, nil)

			for len(session.Players) < 2 || !players_ready.all() {
				pmsg := session.PumpEvents(no_timeout)
				if pmsg == nil {
					return
				}

				switch msg := pmsg.Message.(type) {
				case *UserCommand:
					switch msg.Action {
					case USER_ACTION_SET_READY:
						players_ready.add(pmsg.Player)
					case USER_ACTION_SET_NOT_READY:
						players_ready.remove(pmsg.Player)
					}
				case *NotifyPlayerJoined:
					players_ready.insertNewPlayer(pmsg.Player, false)

				case *NotifyPlayerLeft:
					players_ready.removePlayer(pmsg.Player)

				}
				broadcastPlayerReadyState(session, players_ready)

			}
		}

		log.Println(session.Id, "Start game")

		{
			// Create a list of players:
			players := make([]*Player, len(session.Players))
			{
				i := 0
				for p := range session.Players {
					players[i] = p
					i += 1
				}
			}

			// create random player order which we will use this round:
			random_source.Shuffle(len(players), func(i, j int) {
				players[i], players[j] = players[j], players[i]
			})

			// Each player gets their turn:
			for index, active_painter := range players {

				log.Println(session.Id, "Start round", index+1)

				// Assign roles:
				player_role := make(map[*Player]Role)
				for _, player := range players {
					if player == active_painter {
						player_role[player] = ROLE_PAINTER
					} else {
						player_role[player] = ROLE_TROLL
					}
				}

				// Select one random background:

				backdrop := AVAILABLE_BACKGROUNDS[random_source.Intn(len(AVAILABLE_BACKGROUNDS))]

				prompts := make([]string, len(AVAILABLE_PROMPTS))
				random_source.Shuffle(len(prompts), func(i, j int) {
					prompts[i], prompts[j] = prompts[j], prompts[i]
				})
				prompts = prompts[0:3]

				log.Println("selected backdrop:", backdrop)
				log.Println("selected prompts: ", prompts)

				// Create prototypes for the views:
				troll_view := ChangeGameViewEvent{
					View: GAME_VIEW_EXHIBITION,

					Painting:         nil,
					PaintingBackdrop: &backdrop,
					PaintingPrompt:   nil,
					PaintingStickers: []Sticker{},

					AvailableStickers: []string{},

					VotePrompt:  &VOTE_PROMPT_PROMPT,
					VoteOptions: prompts,
				}
				painter_view := ChangeGameViewEvent{
					View: GAME_VIEW_EXHIBITION,

					Painting:         nil,
					PaintingBackdrop: &backdrop,
					PaintingPrompt:   nil,
					PaintingStickers: []Sticker{},

					AvailableStickers: []string{},

					VotePrompt:  nil,
					VoteOptions: []string{},
				}

				// local function to update the roles:
				updateViews := func() {
					for _, player := range players {
						switch player_role[player] {
						case ROLE_PAINTER:
							player.SendChan <- &painter_view
						case ROLE_TROLL:
							player.SendChan <- &troll_view
						}
					}
				}

				// Now update the views for the players
				updateViews()

				// Phase 1: Trolls vote for a prompt
				log.Println(session.Id, "Voting for trolls starts")
				{
					prompt_voted := createPlayerSetFromList(players, active_painter)
					for !prompt_voted.allTrolls() {
						pmsg := session.PumpEvents(no_timeout)
						if pmsg == nil {
							return
						}

						switch msg := pmsg.Message.(type) {
						case *VoteCommand:
							log.Println("Player ", pmsg.Player.NickName, "voted for", msg)
							prompt_voted.add(pmsg.Player)

							// TODO(fqu): add vote to election
						}
					}

					// TODO(fqu): evaluate votes, select prompt
				}

				// Phase 2:
				log.Println(session.Id, "Painter is now being tortured")
				for {
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

				// Phase 3:
				log.Println(session.Id, "Trolls now select stickers")
				for {
					//
				}

				// Phase 4:
				log.Println(session.Id, "Players can now gaze upon the art")
				for {
					//
				}
			}
		}

		// Phase 5:
		log.Println(session.Id, "All rounds done, show the gallery")

		for {
			//
		}

		// Phase 6:
		log.Println(session.Id, "Showcase the winner")

		for {
			//
		}

	}
}

type playerSetItem struct {
	value bool
	role  Role
}

type playerSet struct {
	items map[*Player]*playerSetItem
}

func createPlayerSetFromMap(players map[*Player]bool, painter *Player) playerSet {

	items := make(map[*Player]*playerSetItem)

	for p := range players {
		item := playerSetItem{
			value: false,
			role:  ROLE_TROLL,
		}
		if p == painter {
			item.role = ROLE_PAINTER
		}
		items[p] = &item
	}

	return playerSet{
		items: items,
	}
}

func createPlayerSetFromList(players []*Player, painter *Player) playerSet {

	items := make(map[*Player]*playerSetItem)

	for _, p := range players {
		item := playerSetItem{
			value: false,
			role:  ROLE_TROLL,
		}
		if p == painter {
			item.role = ROLE_PAINTER
		}
		items[p] = &item
	}

	return playerSet{
		items: items,
	}
}

func (set *playerSet) add(p *Player) {
	set.items[p].value = true
}

func (set *playerSet) remove(p *Player) {
	set.items[p].value = false
}

func (set *playerSet) all() bool {
	for _, item := range set.items {
		if !item.value {
			return false
		}
	}
	return true
}

func (set *playerSet) none() bool {
	for _, item := range set.items {
		if item.value {
			return false
		}
	}
	return true
}

func (set *playerSet) allTrolls() bool {
	for _, item := range set.items {
		if item.role == ROLE_TROLL && !item.value {
			return false
		}
	}
	return true
}

func (set *playerSet) painter() bool {
	for _, item := range set.items {
		if item.role == ROLE_PAINTER {
			return item.value
		}
	}
	return false
}

func (set *playerSet) insertNewPlayer(p *Player, inital bool) {
	set.items[p] = &playerSetItem{
		value: inital,
	}
}

func (set *playerSet) removePlayer(p *Player) {
	delete(set.items, p)
}
