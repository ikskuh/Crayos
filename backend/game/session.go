package game

import (
	"fmt"
	"log"
	"math/rand"
	"time"
	"unsafe"
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
		new.Send(&JoinSessionFailedEvent{
			Reason: "Session is already running.",
		})
		return
	}

	log.Println("Player", new.NickName, "joined session", session.Id)

	new.Session = session
	session.Players[new] = true

	new.Send(&EnterSessionEvent{
		SessionId: session.Id,
	})

	session.BroadcastPlayers(new, nil)

	new.Send(&ChangeGameViewEvent{
		View: GAME_VIEW_LOBBY,
	})

}

func (session *Session) Broadcast(msg Message) {
	for player := range session.Players {
		player.Send(msg)
	}
}

func (session *Session) BroadcastExcept(msg Message, except *Player) {
	for player := range session.Players {
		if player != except {
			player.Send(msg)
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

func (_ *NotifyTimeout) GetJsonType() string {
	return ""
}

func (self *NotifyTimeout) FixNils() Message {
	return self
}

type NotifyPlayerJoined struct {
}

func (_ *NotifyPlayerJoined) GetJsonType() string {
	return ""
}

func (self *NotifyPlayerJoined) FixNils() Message {
	return self
}

type NotifyPlayerLeft struct {
	// PlayerMessage.Player is not in the session anymore!
}

func (_ *NotifyPlayerLeft) GetJsonType() string {
	return ""
}

func (self *NotifyPlayerLeft) FixNils() Message {
	return self
}

func (session *Session) PumpEvents(timeout <-chan time.Time) *PlayerMessage {

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

		// Lobby
		{
			players_ready := createPlayerSetFromMap(session.Players, nil)

			for len(session.Players) < 2 || players_ready.any(false) {
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

				backdrop := ALL_BACKDROP_ITEMS[random_source.Intn(len(ALL_BACKDROP_ITEMS))]

				prompts := make([]string, len(AVAILABLE_PROMPTS))
				copy(prompts, AVAILABLE_PROMPTS)
				random_source.Shuffle(len(prompts), func(i, j int) {
					prompts[i], prompts[j] = prompts[j], prompts[i]
				})
				prompts = prompts[0:3]

				log.Println("selected backdrop:", backdrop)
				log.Println("selected prompts: ", prompts)

				// Create prototypes for the views:
				troll_view := &ChangeGameViewEvent{
					View: GAME_VIEW_PROMPTSELECTION,

					Painting:         nil,
					PaintingBackdrop: backdrop,
					PaintingPrompt:   "",
					PaintingStickers: []Sticker{},

					VotePrompt:  VOTE_PROMPT_PROMPT,
					VoteOptions: prompts,
				}
				painter_view := &ChangeGameViewEvent{
					View: GAME_VIEW_ARTSTUDIO_GENERIC,

					Painting:         nil,
					PaintingBackdrop: backdrop,
					PaintingPrompt:   "",
					PaintingStickers: []Sticker{},

					VotePrompt:  "",
					VoteOptions: []string{},
				}

				// local function to update the roles:
				updateViews := func() {
					for _, player := range players {
						switch player_role[player] {
						case ROLE_PAINTER:
							// log.Println("send view (painter)", player.NickName, painter_view)
							player.Send(painter_view)
						case ROLE_TROLL:
							// log.Println("send view (troll)", player.NickName, troll_view)
							player.Send(troll_view)
						}
					}
				}

				changeBoth := func(handler func(view *ChangeGameViewEvent)) {
					handler(troll_view)
					handler(painter_view)
				}

				// Now update the views for the players
				updateViews()

				// Prepare message for trolls to go into "wait for others" state
				troll_view.View = GAME_VIEW_ARTSTUDIO_GENERIC
				troll_view.VoteOptions = []string{}

				// Phase 1: Trolls vote for a prompt
				log.Println(session.Id, "Prompt voting for trolls starts")
				var selected_painting_prompt string
				{
					prompt_voted := createPlayerSetFromList(players, active_painter)

					votes := make([]float32, len(prompts))

					for !prompt_voted.allTrolls() {
						pmsg := session.PumpEvents(no_timeout)
						if pmsg == nil {
							return
						}

						switch msg := pmsg.Message.(type) {
						case *VoteCommand:
							if pmsg.Player != active_painter {

								log.Println("Player ", pmsg.Player.NickName, "voted for", msg)

								index := -1
								for i, val := range prompts {
									if msg.Option == val {
										index = i
									}
								}

								if index >= 0 {
									votes[index] += 0.95 + 0.01*random_source.Float32()

									prompt_voted.add(pmsg.Player)

									// Hide the options for the troll that voted:
									pmsg.Player.Send(troll_view)

								} else {
									log.Println("troll tried to vote illegaly. BAD BOY")

								}

							} else {
								log.Println("painter tried to vote. BAD BOY")
							}
						}
					}

					best_prompt_index := 0
					best_prompt_level := votes[0]

					for index, level := range votes {
						if level >= best_prompt_level {
							best_prompt_index = index
							best_prompt_level = level
						}
					}

					selected_painting_prompt = prompts[best_prompt_index]

					log.Println("Prompt", selected_painting_prompt, "won with", best_prompt_level, "votes")
				}

				changeBoth(func(view *ChangeGameViewEvent) {
					view.VotePrompt = ""
					view.PaintingPrompt = selected_painting_prompt
				})

				troll_view.View = GAME_VIEW_ARTSTUDIO_GENERIC
				painter_view.View = GAME_VIEW_ARTSTUDIO_ACTIVE

				// TODO(fqu): Implement round robin trolling

				updateViews()

				// Phase 2:
				log.Println(session.Id, "Painter is now being tortured")
				{
					// Setup troll order, current troll is always the first one
					trolls := make([]*Player, len(players)-1)

					{
						i := 0
						for _, player := range players {
							if player == active_painter {
								continue
							}
							trolls[i] = player
							i += 1
						}

						// shuffle troll order:
						random_source.Shuffle(len(trolls), func(i, j int) {
							trolls[i], trolls[j] = trolls[j], trolls[i]
						})
					}

					next_troll_event := 0
					troll_did_effect := false

					// Setup session timing:
					total_time_left := GAME_ROUND_TIME_S
					session.Broadcast(&TimerChangedEvent{
						SecondsLeft: total_time_left,
					})

					second_ticker := time.NewTicker(1 * time.Second)
					for total_time_left > 0 {

						if next_troll_event <= 0 {

							trolls[0].Send(troll_view) // troll view is "generic empty" here

							// select next troll by doing round-robin scheduling:
							trolls = append(trolls[1:], trolls[0])

							vote_effect_view := *troll_view

							vote_effect_view.VotePrompt = VOTE_PROMPT_EFFECT
							vote_effect_view.VoteOptions = *(*[]string)(unsafe.Pointer(&ALL_EFFECT_ITEMS))

							trolls[0].Send(&vote_effect_view) // troll view is "generic empty" here
							troll_did_effect = false

							next_troll_event = GAME_TROLL_EFFECT_COOLDOWN_S
						}

						pmsg := session.PumpEvents(second_ticker.C)
						if pmsg == nil {
							return
						}

						switch msg := pmsg.Message.(type) {

						case *NotifyTimeout:
							total_time_left -= 1
							next_troll_event -= 1
							session.Broadcast(&TimerChangedEvent{
								SecondsLeft: total_time_left,
							})

						case *VoteCommand:
							if pmsg.Player == trolls[0] && !troll_did_effect {
								// TODO(fqu): validate that msg.Option is actually a legal vote!
								active_painter.Send(&ChangeToolModifierEvent{
									Modifier: Effect(msg.Option),
								})
								trolls[0].Send(troll_view) // reset troll to regular view, hide the vote options
								troll_did_effect = true
							} else {
								log.Println("someone else tried to harm the painter. BAD BOY!")
							}

						case *SetPaintingCommand:
							if pmsg.Player == active_painter {

								// Keep the state up to date with the painted image:
								troll_view.Painting = msg.Path
								painter_view.Painting = msg.Path

								// Forward painting actions when the user changes the image.
								session.BroadcastExcept(&PaintingChangedEvent{
									Path: msg.Path,
								}, pmsg.Player)

							} else {
								log.Println("someone else tried to paint. BAD BOY!")
							}
						}
					}
				}

				// Phase 3:
				{
					log.Println(session.Id, "Trolls now select stickers")
					for false {
						//
					}
				}

				// Phase 4:
				{
					log.Println(session.Id, "Players can now gaze upon the art")
					round_end_timer := time.NewTimer(GAME_ROUND_TIME_S)
					timeLeft := true
					players_ready := createPlayerSetFromMap(session.Players, nil)
					for timeLeft && players_ready.any(false) {
						pmsg := session.PumpEvents(round_end_timer.C)
						if pmsg == nil {
							return
						}

						switch msg := pmsg.Message.(type) {
						case *UserCommand:
							switch msg.Action {
							case USER_ACTION_CONTINUE_GAME:
								players_ready.add(pmsg.Player)
							}
						case *VoteCommand:

						case *NotifyTimeout:
							timeLeft = false
						}
						broadcastPlayerReadyState(session, players_ready)
					}
				}
			}
		}

		// Phase 5:
		{
			log.Println(session.Id, "All rounds done, show the gallery and vote for the winner")

			session.Broadcast(&ChangeGameViewEvent{
				View: GAME_VIEW_GALLERY,
			})

			round_end_timer := time.NewTimer(GAME_ROUND_TIME_S)
			timeLeft := true
			players_ready := createPlayerSetFromMap(session.Players, nil)
			for timeLeft && players_ready.any(false) {
				pmsg := session.PumpEvents(round_end_timer.C)
				if pmsg == nil {
					return
				}

				switch msg := pmsg.Message.(type) {
				case *UserCommand:
					switch msg.Action {
					case USER_ACTION_CONTINUE_GAME:
						players_ready.add(pmsg.Player)
					}
				case *VoteCommand:

				case *NotifyTimeout:
					timeLeft = false
				}
				broadcastPlayerReadyState(session, players_ready)
			}

		}

		// Phase 6:
		{
			log.Println(session.Id, "Showcase the winner")

			// TODO set drawing of winner
			session.Broadcast(&ChangeGameViewEvent{
				View: GAME_VIEW_PODIUM,
				// Painting: winner.painting
				// PaintingPrompt: "The winner is" + winner.name

				//// Unchanged stuff
				//paintingBackdrop: Backdrop # artstudio*: the ID of the backdrop
				//paintingStickers: list[Sticker] # artstudio*: the current list of stickers that should be shown

				//votePrompt: str # artstudioGeneric: the prompt that is shown when
				//voteOptions: list[str] # promptselection, artstudioGeneric: list of options that the player can vote for.
			})

			round_end_timer := time.NewTimer(GAME_ROUND_TIME_S)
			timeLeft := true
			players_ready := createPlayerSetFromMap(session.Players, nil)
			for timeLeft && players_ready.any(false) {
				pmsg := session.PumpEvents(round_end_timer.C)
				if pmsg == nil {
					return
				}

				switch msg := pmsg.Message.(type) {
				case *UserCommand:
					switch msg.Action {
					case USER_ACTION_CONTINUE_GAME:
						players_ready.add(pmsg.Player)
					}
				case *NotifyTimeout:
					timeLeft = false
				}
				broadcastPlayerReadyState(session, players_ready)
			}

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

func (set *playerSet) any(predicate bool) bool {
	for _, item := range set.items {
		if item.value == predicate {
			return true
		}
	}
	return false
}

func (set *playerSet) all() bool {
	return !set.any(false)
}

func (set *playerSet) none() bool {
	return !set.any(true)
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
