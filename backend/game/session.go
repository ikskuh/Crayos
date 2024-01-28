package game

import (
	"fmt"
	"log"
	"math/rand"
	"time"
	"unsafe"

	"random-projects.net/crayos-backend/meta"
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

	// Internals:
	startupTime int64
}

type Role int

const (
	ROLE_PAINTER Role = 0
	ROLE_TROLL   Role = 1
)

var sessions = map[string]*Session{}

func SetDebugSession(session *Session) {
	if !*meta.DEBUG_MODE {
		log.Fatalln("Only allowed in debug mode!")
	}
	sessions["0xDEADBEEF"] = session
}

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

		startupTime: meta.Timestamp(),
	}
	session.Id = fmt.Sprintf("%p", session)

	if player != nil {
		session.AddPlayer(player)
	} else if !*meta.DEBUG_MODE {
		log.Fatalln("Invalid parameter: Session requires a player in non-debug mode")
	}

	// add player before starting main loop, otherwise it will kill itself automatically

	go session.Run()

	// Register session
	session.ServerPrint("Created")
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

	session.ServerPrint("Player", new.NickName, "joined")

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

type gameTimer interface {
	GetChannel() <-chan time.Time
	NotifyTick()
}

func (session *Session) PumpEvents(timer gameTimer) *PlayerMessage {

	for *meta.DEBUG_MODE || len(session.Players) > 0 {
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

			session.ServerPrint("Player", old.NickName, "left")
			delete(session.Players, old)

			session.BroadcastPlayers(nil, old)

			return &PlayerMessage{
				Message: &NotifyPlayerLeft{},
				Player:  old,
			}

		case t := <-timer.GetChannel():
			timer.NotifyTick()
			return &PlayerMessage{
				Player:  nil,
				Message: &NotifyTimeout{timestamp: t},
			}
		}
	}
	return nil
}

func broadcastPlayerReadyState(s *Session, m playerSet) {
	s.ServerPrint("Sending PlayerReadyState", m)
	readyMap := make(map[string]bool)
	for p, b := range m.items {
		readyMap[p.NickName] = b.value
	}
	s.Broadcast(&PlayerReadyChangedEvent{
		Players: readyMap,
	})
}

type gameRoundResult struct {
	painting    Painting
	totalPoints int
}

func (evt *ChangeGameViewEvent) RemoveVote() {
	evt.VotePrompt = ""
	evt.VoteOptions = []string{}
}

func (evt *ChangeGameViewEvent) SetVote(prompt string, options []string) {
	evt.VotePrompt = prompt
	evt.VoteOptions = options
}

func (session *Session) ServerPrint(message ...any) {
	timestamp := meta.Timestamp() - session.startupTime
	formatted := fmt.Sprint(message...)

	log.Println(
		fmt.Sprintf("Session[%s, %7d]: %s", session.Id, timestamp, formatted),
	)
}

func (session *Session) DebugPrint(message ...any) {
	timestamp := meta.Timestamp() - session.startupTime
	formatted := fmt.Sprint(message...)

	log.Println(
		fmt.Sprintf("Session[%s, %7d]: %s", session.Id, timestamp, formatted),
	)
	session.Broadcast(&DebugMessageEvent{
		Message: fmt.Sprintf("%d: %s", timestamp, formatted),
	})
}

func (session *Session) Announce(text string, duration time.Duration) {
	session.Broadcast(&ChangeGameViewEvent{
		View:      GAME_VIEW_ANNOUNCER,
		Announcer: text,
	})
	time.Sleep(duration)
}

func (session *Session) Run() {
	random_source := rand.New(rand.NewSource(time.Now().UnixNano()))

	session.ServerPrint("Started")
	defer session.ServerPrint("Stopped")

	no_timeout := &noTimeoutGameTimer{
		channel: make(chan time.Time), // pass when no timeout is required
	}

	for *meta.DEBUG_MODE || len(session.Players) > 0 {

		// Lobby
		session.DebugPrint("Enter lobby")
		{
			// Show lobby
			session.Broadcast(&ChangeGameViewEvent{
				View: GAME_VIEW_LOBBY,
			})

			players_ready := createPlayerSetFromMap(session.Players, nil)

			for len(session.Players) < 2 || players_ready.any(false) {

				broadcastPlayerReadyState(session, players_ready)

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
			}
		}

		session.DebugPrint("Start game")

		// Game
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

			results := make([]gameRoundResult, len(players))

			// Each player gets their turn:
			for index, active_painter := range players {

				round_id := fmt.Sprintf("Round %d: ", index+1)

				session.DebugPrint(round_id, "Initialize")

				// Assign roles:
				player_role := make(map[*Player]Role)
				for _, player := range players {
					if player == active_painter {
						player_role[player] = ROLE_PAINTER
					} else {
						player_role[player] = ROLE_TROLL
					}
				}

				// local function to update the roles:
				splitAnnounce := func(painterText string, trollText string) {
					for _, player := range players {
						var text string
						switch player_role[player] {
						case ROLE_PAINTER:
							text = painterText
						case ROLE_TROLL:
							text = trollText
						}
						player.Send(&ChangeGameViewEvent{
							View:      GAME_VIEW_ANNOUNCER,
							Announcer: text,
						})
					}
					time.Sleep(ANNOUNCE_GENERIC_TIMEOUT)
				}

				// Select one random background:

				backdrop := ALL_BACKDROP_ITEMS[random_source.Intn(len(ALL_BACKDROP_ITEMS))]

				prompts := make([]string, len(AVAILABLE_PROMPTS))
				copy(prompts, AVAILABLE_PROMPTS)
				random_source.Shuffle(len(prompts), func(i, j int) {
					prompts[i], prompts[j] = prompts[j], prompts[i]
				})
				prompts = prompts[0:3]

				session.ServerPrint("selected backdrop:", backdrop)
				session.ServerPrint("selected prompts: ", prompts)

				// Tell them what's happening
				splitAnnounce(
					TEXT_ANNOUNCE_YOU_ARE_PAINTER,
					TEXT_ANNOUNCE_YOU_ARE_TROLL,
				)

				// Create prototypes for the views:
				troll_view := &ChangeGameViewEvent{
					View: GAME_VIEW_PROMPTSELECTION,

					Painting: Painting{
						Graphics: EMPTY_GRAPHICS,
						Backdrop: backdrop,
						Prompt:   "",
						Stickers: []Sticker{},
					},
				}
				painter_view := &ChangeGameViewEvent{
					View: GAME_VIEW_ARTSTUDIO_GENERIC,

					Painting: Painting{
						Graphics: EMPTY_GRAPHICS,
						Backdrop: backdrop,
						Prompt:   "",
						Stickers: []Sticker{},
					},
				}

				// local function to update the roles:
				updateViews := func() {
					for _, player := range players {
						switch player_role[player] {
						case ROLE_PAINTER:
							// session.ServerPrint("send view (painter)", player.NickName, painter_view)
							player.Send(painter_view)
						case ROLE_TROLL:
							// session.ServerPrint("send view (troll)", player.NickName, troll_view)
							player.Send(troll_view)
						}
					}
				}

				changeBoth := func(handler func(view *ChangeGameViewEvent)) {
					handler(troll_view)
					handler(painter_view)
				}

				troll_view.SetVote(TEXT_VOTE_PROMPT, prompts)
				painter_view.RemoveVote()

				// Now update the views for the players
				updateViews()

				// Prepare message for trolls to go into "wait for others" state
				troll_view.View = GAME_VIEW_ARTSTUDIO_GENERIC
				troll_view.RemoveVote()

				// Phase 1: Trolls vote for a prompt
				session.DebugPrint(round_id, "Prompt voting for trolls starts")
				var selected_painting_prompt string
				{
					prompt_voted := createPlayerSetFromList(players, active_painter)

					votes := make([]float32, len(prompts))
					for i := range votes {
						// initialize votes with some basic noise so timeout can happen
						votes[i] = 0.1 * random_source.Float32()
					}

					vote_end_timer := session.createTimer(GAME_TOPIC_VOTE_TIME_S)

					for !vote_end_timer.TimedOut() && !prompt_voted.allTrollsSet() {
						pmsg := session.PumpEvents(vote_end_timer)
						if pmsg == nil {
							return
						}

						switch msg := pmsg.Message.(type) {
						case *VoteCommand:
							if pmsg.Player != active_painter {

								session.ServerPrint("Player ", pmsg.Player.NickName, "voted for", msg)

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
									session.ServerPrint("troll tried to vote illegaly. BAD BOY")

								}

							} else {
								session.ServerPrint("painter tried to vote. BAD BOY")
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

					session.ServerPrint("Prompt", selected_painting_prompt, "won with", best_prompt_level, "votes")
				}

				changeBoth(func(view *ChangeGameViewEvent) {
					view.RemoveVote()
					view.Painting.Prompt = selected_painting_prompt
				})

				troll_view.View = GAME_VIEW_ARTSTUDIO_GENERIC
				painter_view.View = GAME_VIEW_ARTSTUDIO_ACTIVE

				updateViews()

				// Phase 2:
				session.DebugPrint(round_id, "Painter is now being tortured")
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

					round_end_timer := session.createTimer(GAME_ROUND_TIME_S)

					for !round_end_timer.TimedOut() {

						if next_troll_event <= 0 {

							trolls[0].Send(troll_view) // troll view is "generic empty" here

							// select next troll by doing round-robin scheduling:
							trolls = append(trolls[1:], trolls[0])

							vote_effect_view := *troll_view

							vote_effect_view.SetVote(TEXT_VOTE_EFFECT, *(*[]string)(unsafe.Pointer(&ALL_EFFECT_ITEMS)))

							trolls[0].Send(&vote_effect_view) // troll view is "generic empty" here
							troll_did_effect = false

							next_troll_event = GAME_TROLL_EFFECT_COOLDOWN_S
						}

						pmsg := session.PumpEvents(round_end_timer)
						if pmsg == nil {
							return
						}

						switch msg := pmsg.Message.(type) {

						case *NotifyTimeout:
							next_troll_event -= 1

						case *VoteCommand:
							if pmsg.Player == trolls[0] && !troll_did_effect {
								// TODO(fqu): validate that msg.Option is actually a legal vote!
								session.Broadcast(&ChangeToolModifierEvent{
									Modifier: Effect(msg.Option),
									Duration: GAME_TROLL_EFFECT_DURATION_MS,
								})
								trolls[0].Send(troll_view) // reset troll to regular view, hide the vote options
								troll_did_effect = true
							} else {
								session.ServerPrint("someone else tried to harm the painter. BAD BOY!")
							}

						case *SetPaintingCommand:
							if pmsg.Player == active_painter {

								// Keep the state up to date with the painted image:
								troll_view.Painting.Graphics = msg.Graphics
								painter_view.Painting.Graphics = msg.Graphics

								// Forward painting actions when the user changes the image.
								session.BroadcastExcept(&PaintingChangedEvent{
									Graphics: msg.Graphics,
								}, pmsg.Player)

							} else {
								session.ServerPrint("someone else tried to paint. BAD BOY!")
							}
						}
					}

					// Hide timer:
					session.Broadcast(&TimerChangedEvent{
						SecondsLeft: -1,
					})
				}

				updateViews()

				// Disable all active effects
				session.Broadcast(&ChangeToolModifierEvent{
					Modifier: "",
					Duration: 0,
				})

				// Phase 3:
				session.DebugPrint(round_id, "Trolls now select stickers")
				{
					// TODO(fqu): Set stickering mode here

					updateViews()

					for false {
						//
					}
				}

				// Store the result of that round
				results[index] = gameRoundResult{
					painting:    painter_view.Painting,
					totalPoints: 0,
				}

				// Phase 4:
				session.DebugPrint(round_id, "Showcase the artwork")
				{
					round_end_timer := session.createTimer(GALLERY_ROUND_TIME_S)
					players_ready := createPlayerSetFromMap(session.Players, nil)

					changeBoth(func(view *ChangeGameViewEvent) {
						view.View = GAME_VIEW_ARTSTUDIO_GENERIC
						view.SetVote(TEXT_VOTE_SHOWCASE, []string{"", "", "", "", "continue"})
					})

					updateViews()

					// Remove the vote so we can hide it if the player hits the button
					changeBoth(func(view *ChangeGameViewEvent) {
						view.RemoveVote()
					})

					for !round_end_timer.TimedOut() && !players_ready.allSet() {
						pmsg := session.PumpEvents(round_end_timer)
						if pmsg == nil {
							return
						}

						switch msg := pmsg.Message.(type) {
						case *VoteCommand:
							if msg.Option != "continue" {
								session.ServerPrint("User sent bad continue option, BAD BOY")
							} else {
								players_ready.add(pmsg.Player)
								pmsg.Player.Send(troll_view) // it doesn't matter, they should be equal
							}
						}
					}
					round_end_timer.Hide()
				}
			} // end of inner loop over players

			session.Announce("Vote for the winner now!", ANNOUNCE_GENERIC_TIMEOUT)

			// Phase 5:
			{
				// TODO: Loop through all results and let the players vote for the pictures

				for index, result := range results {
					round_id := fmt.Sprintf("Showcase %d: ", index+1)

					session.DebugPrint(round_id, "Vote for image")

					vote_view := ChangeGameViewEvent{
						View:     GAME_VIEW_ARTSTUDIO_GENERIC,
						Painting: result.painting,
					}
					vote_view.SetVote(TEXT_VOTE_SHOWCASE, []string{
						"star1",
						"star2",
						"star3",
						"star4",
						"star5",
					})
					session.Broadcast(&vote_view)

					// Hide the vote for later sending:
					vote_view.RemoveVote()

					round_end_timer := session.createTimer(GALLERY_ROUND_TIME_S)
					players_ready := createPlayerSetFromMap(session.Players, nil)
					for !round_end_timer.TimedOut() && !players_ready.allSet() {
						pmsg := session.PumpEvents(round_end_timer)
						if pmsg == nil {
							return
						}

						switch msg := pmsg.Message.(type) {
						case *VoteCommand:

							if !players_ready.isSet(pmsg.Player) {

								ok := true
								switch msg.Option {
								case "star1":
									results[index].totalPoints += 1
								case "star2":
									results[index].totalPoints += 2
								case "star3":
									results[index].totalPoints += 3
								case "star4":
									results[index].totalPoints += 4
								case "star5":
									results[index].totalPoints += 5

								default:
									ok = false
								}

								if ok {
									players_ready.add(pmsg.Player)
									pmsg.Player.Send(&vote_view)
								}
							} else {
								session.ServerPrint("don't wont twice my friend. BAD BOY!")
							}

						}
					}
					round_end_timer.Hide()
				}
			}

			session.Announce(TEXT_ANNOUNCE_WINNER, ANNOUNCE_GENERIC_TIMEOUT)

			// Determine winner:
			{
				best_painting_score := 0
				best_painting_index := 0

				for i := range results {
					results[i].painting.Winner = false

					if results[i].totalPoints >= best_painting_score {
						best_painting_score = results[i].totalPoints
						best_painting_index = i
					}
				}

				results[best_painting_index].painting.Winner = true
			}

			// Phase 6:
			session.DebugPrint("Showcase the winner")
			{
				view_cmd := ChangeGameViewEvent{
					View:    GAME_VIEW_GALLERY,
					Results: make([]Painting, len(results)),
				}

				for i := range view_cmd.Results {
					view_cmd.Results[i] = results[i].painting
				}

				// TODO set drawing of winner
				session.Broadcast(&view_cmd)

				round_end_timer := session.createTimer(GAME_ROUND_TIME_S)
				players_ready := createPlayerSetFromMap(session.Players, nil)
				for !round_end_timer.TimedOut() && players_ready.any(false) {
					pmsg := session.PumpEvents(round_end_timer)
					if pmsg == nil {
						return
					}

					switch msg := pmsg.Message.(type) {
					case *UserCommand:
						switch msg.Action {
						case USER_ACTION_LEAVE_GALLERY:
							players_ready.add(pmsg.Player)
						}
					}
				}
				round_end_timer.Hide()
			}

			session.DebugPrint("Round done. Back to lobby!")
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

func (set *playerSet) allSet() bool {
	return !set.any(false)
}

func (set *playerSet) noneSet() bool {
	return !set.any(true)
}

func (set *playerSet) allTrollsSet() bool {
	for _, item := range set.items {
		if item.role == ROLE_TROLL && !item.value {
			return false
		}
	}
	return true
}

func (set *playerSet) painterSet() bool {
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

func (set *playerSet) isSet(p *Player) bool {
	return set.items[p].value
}

type autoGameTimer struct {
	session *Session
	ticker  *time.Ticker

	timeLeft int
}

func (session *Session) createTimer(timeout_secs int) *autoGameTimer {
	session.Broadcast(&TimerChangedEvent{
		SecondsLeft: timeout_secs,
	})
	return &autoGameTimer{
		session:  session,
		ticker:   time.NewTicker(1 * time.Second),
		timeLeft: timeout_secs,
	}
}

func (timer *autoGameTimer) TimedOut() bool {
	return timer.timeLeft <= 0
}

func (timer *autoGameTimer) NotifyTick() {
	timer.timeLeft -= 1
	timer.session.Broadcast(&TimerChangedEvent{
		SecondsLeft: timer.timeLeft,
	})
}

func (timer *autoGameTimer) Hide() {
	timer.session.Broadcast(&TimerChangedEvent{
		SecondsLeft: -1,
	})
}

func (timer *autoGameTimer) GetChannel() <-chan time.Time {
	return timer.ticker.C
}

type noTimeoutGameTimer struct {
	channel <-chan time.Time
}

func (timer *noTimeoutGameTimer) GetChannel() <-chan time.Time {
	return timer.channel
}

func (timer *noTimeoutGameTimer) NotifyTick() {

}
