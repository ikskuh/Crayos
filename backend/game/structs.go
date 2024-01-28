package game

import (
	"encoding/json"
	"errors"
)

const (
	CREATE_SESSION_COMMAND_TAG = "create-session-command"
	JOIN_SESSION_COMMAND_TAG = "join-session-command"
	LEAVE_SESSION_COMMAND_TAG = "leave-session-command"
	USER_COMMAND_TAG = "user-command"
	VOTE_COMMAND_TAG = "vote-command"
	PLACE_STICKER_COMMAND_TAG = "place-sticker-command"
	SET_PAINTING_COMMAND_TAG = "set-painting-command"
	ENTER_SESSION_EVENT_TAG = "enter-session-event"
	JOIN_SESSION_FAILED_EVENT_TAG = "join-session-failed-event"
	KICKED_EVENT_TAG = "kicked-event"
	CHANGE_GAME_VIEW_EVENT_TAG = "change-game-view-event"
	TIMER_CHANGED_EVENT_TAG = "timer-changed-event"
	CHANGE_TOOL_MODIFIER_EVENT_TAG = "change-tool-modifier-event"
	PAINTING_CHANGED_EVENT_TAG = "painting-changed-event"
	PLAYERS_CHANGED_EVENT_TAG = "players-changed-event"
	PLAYER_READY_CHANGED_EVENT_TAG = "player-ready-changed-event"
	POP_UP_EVENT_TAG = "pop-up-event"
	DEBUG_MESSAGE_EVENT_TAG = "debug-message-event"
)

func DeserializeMessage(data []byte) (Message, error) {

	var raw_map map[string]interface{} // must be an object

	err := json.Unmarshal(data, &raw_map)
	if err != nil {
		return nil, err
	}

	type_tag, ok := raw_map["type"]
	if !ok {
		return nil, errors.New("Invalid json")
	}

	var out Message

	switch type_tag {

	case CREATE_SESSION_COMMAND_TAG:
		out = &CreateSessionCommand{}
	case JOIN_SESSION_COMMAND_TAG:
		out = &JoinSessionCommand{}
	case LEAVE_SESSION_COMMAND_TAG:
		out = &LeaveSessionCommand{}
	case USER_COMMAND_TAG:
		out = &UserCommand{}
	case VOTE_COMMAND_TAG:
		out = &VoteCommand{}
	case PLACE_STICKER_COMMAND_TAG:
		out = &PlaceStickerCommand{}
	case SET_PAINTING_COMMAND_TAG:
		out = &SetPaintingCommand{}
	case ENTER_SESSION_EVENT_TAG:
		out = &EnterSessionEvent{}
	case JOIN_SESSION_FAILED_EVENT_TAG:
		out = &JoinSessionFailedEvent{}
	case KICKED_EVENT_TAG:
		out = &KickedEvent{}
	case CHANGE_GAME_VIEW_EVENT_TAG:
		out = &ChangeGameViewEvent{}
	case TIMER_CHANGED_EVENT_TAG:
		out = &TimerChangedEvent{}
	case CHANGE_TOOL_MODIFIER_EVENT_TAG:
		out = &ChangeToolModifierEvent{}
	case PAINTING_CHANGED_EVENT_TAG:
		out = &PaintingChangedEvent{}
	case PLAYERS_CHANGED_EVENT_TAG:
		out = &PlayersChangedEvent{}
	case PLAYER_READY_CHANGED_EVENT_TAG:
		out = &PlayerReadyChangedEvent{}
	case POP_UP_EVENT_TAG:
		out = &PopUpEvent{}
	case DEBUG_MESSAGE_EVENT_TAG:
		out = &DebugMessageEvent{}

	default:
		return nil, errors.New("Invalid type")
	}

	err = json.Unmarshal(data, &out)
	if err != nil {
		return nil, err
	}

	return out, nil
}

type Sticker struct {
	Id string `json:"id"`
	X float32 `json:"x"`
	Y float32 `json:"y"`
}

type GameView string
const (
	GAME_VIEW_TITLE GameView = "title"
	GAME_VIEW_LOBBY GameView = "lobby"
	GAME_VIEW_PROMPTSELECTION GameView = "promptselection"
	GAME_VIEW_ARTSTUDIO_GENERIC GameView = "artstudio-generic"
	GAME_VIEW_ARTSTUDIO_ACTIVE GameView = "artstudio-active"
	GAME_VIEW_ARTSTUDIO_STICKER GameView = "artstudio-sticker"
	GAME_VIEW_GALLERY GameView = "gallery"
	GAME_VIEW_ANNOUNCER GameView = "announcer"
)
var ALL_GAME_VIEW_ITEMS = []GameView{
	"title",
	"lobby",
	"promptselection",
	"artstudio-generic",
	"artstudio-active",
	"artstudio-sticker",
	"gallery",
	"announcer",
}

type Effect string
const (
	EFFECT_FLASHLIGHT Effect = "flashlight"
	EFFECT_DRUNK Effect = "drunk"
	EFFECT_FLIP Effect = "flip"
	EFFECT_SWAP_TOOL Effect = "swap_tool"
	EFFECT_LOCK_PENCIL Effect = "lock_pencil"
)
var ALL_EFFECT_ITEMS = []Effect{
	"flashlight",
	"drunk",
	"flip",
	"swap_tool",
	"lock_pencil",
}

type UserAction string
const (
	USER_ACTION_SET_READY UserAction = "set-ready"
	USER_ACTION_SET_NOT_READY UserAction = "set-not-ready"
	USER_ACTION_LEAVE_GALLERY UserAction = "leave-gallery"
)
var ALL_USER_ACTION_ITEMS = []UserAction{
	"set-ready",
	"set-not-ready",
	"leave-gallery",
}

type Backdrop string
const (
	BACKDROP_ARCTIC Backdrop = "arctic"
	BACKDROP_GRAVEYARD Backdrop = "graveyard"
	BACKDROP_PIRATE_SHIP Backdrop = "pirate_ship"
	BACKDROP_THEATER_STAGE1 Backdrop = "theater_stage1"
	BACKDROP_DESERT Backdrop = "desert"
)
var ALL_BACKDROP_ITEMS = []Backdrop{
	"arctic",
	"graveyard",
	"pirate_ship",
	"theater_stage1",
	"desert",
}

type CreateSessionCommand struct {
	NickName string `json:"nickName"`
}

type JoinSessionCommand struct {
	NickName string `json:"nickName"`
	SessionId string `json:"sessionId"`
}

type LeaveSessionCommand struct {
}

type UserCommand struct {
	Action UserAction `json:"action"`
}

type VoteCommand struct {
	Option string `json:"option"`
}

type PlaceStickerCommand struct {
	Sticker string `json:"sticker"`
	X float32 `json:"x"`
	Y float32 `json:"y"`
}

type SetPaintingCommand struct {
	Graphics Graphics `json:"graphics"`
}

type EnterSessionEvent struct {
	SessionId string `json:"sessionId"`
}

type JoinSessionFailedEvent struct {
	Reason string `json:"reason"`
}

type KickedEvent struct {
	Reason string `json:"reason"`
}

type Painting struct {
	Prompt string `json:"prompt"`
	Graphics Graphics `json:"graphics"`
	Backdrop Backdrop `json:"backdrop"`
	Stickers []Sticker `json:"stickers"`
	Winner bool `json:"winner"`
}

type ChangeGameViewEvent struct {
	View GameView `json:"view"`
	Painting Painting `json:"painting"`
	Results []Painting `json:"results"`
	VotePrompt string `json:"votePrompt"`
	VoteOptions []string `json:"voteOptions"`
	Announcer string `json:"announcer"`
}

type TimerChangedEvent struct {
	SecondsLeft int `json:"secondsLeft"`
}

type ChangeToolModifierEvent struct {
	Modifier Effect `json:"modifier"`
	Duration int `json:"duration"`
}

type PaintingChangedEvent struct {
	Graphics Graphics `json:"graphics"`
}

type PlayersChangedEvent struct {
	Players []string `json:"players"`
	AddedPlayer *string `json:"addedPlayer"`
	RemovedPlayer *string `json:"removedPlayer"`
}

type PlayerReadyChangedEvent struct {
	Players map[string]bool `json:"players"`
}

type PopUpEvent struct {
	Message string `json:"message"`
	Duration int `json:"duration"`
}

type DebugMessageEvent struct {
	Message string `json:"message"`
}


func (item *CreateSessionCommand) GetJsonType() string {
	return "create-session-command"
}
func (item *CreateSessionCommand) FixNils() Message {
	copy := *item
	return &copy
}

func (item *JoinSessionCommand) GetJsonType() string {
	return "join-session-command"
}
func (item *JoinSessionCommand) FixNils() Message {
	copy := *item
	return &copy
}

func (item *LeaveSessionCommand) GetJsonType() string {
	return "leave-session-command"
}
func (item *LeaveSessionCommand) FixNils() Message {
	copy := *item
	return &copy
}

func (item *UserCommand) GetJsonType() string {
	return "user-command"
}
func (item *UserCommand) FixNils() Message {
	copy := *item
	return &copy
}

func (item *VoteCommand) GetJsonType() string {
	return "vote-command"
}
func (item *VoteCommand) FixNils() Message {
	copy := *item
	return &copy
}

func (item *PlaceStickerCommand) GetJsonType() string {
	return "place-sticker-command"
}
func (item *PlaceStickerCommand) FixNils() Message {
	copy := *item
	return &copy
}

func (item *SetPaintingCommand) GetJsonType() string {
	return "set-painting-command"
}
func (item *SetPaintingCommand) FixNils() Message {
	copy := *item
	return &copy
}

func (item *EnterSessionEvent) GetJsonType() string {
	return "enter-session-event"
}
func (item *EnterSessionEvent) FixNils() Message {
	copy := *item
	return &copy
}

func (item *JoinSessionFailedEvent) GetJsonType() string {
	return "join-session-failed-event"
}
func (item *JoinSessionFailedEvent) FixNils() Message {
	copy := *item
	return &copy
}

func (item *KickedEvent) GetJsonType() string {
	return "kicked-event"
}
func (item *KickedEvent) FixNils() Message {
	copy := *item
	return &copy
}

func (item *ChangeGameViewEvent) GetJsonType() string {
	return "change-game-view-event"
}
func (item *ChangeGameViewEvent) FixNils() Message {
	copy := *item
	if copy.Results == nil {
		copy.Results = []Painting{}
	}
	if copy.VoteOptions == nil {
		copy.VoteOptions = []string{}
	}
	return &copy
}

func (item *TimerChangedEvent) GetJsonType() string {
	return "timer-changed-event"
}
func (item *TimerChangedEvent) FixNils() Message {
	copy := *item
	return &copy
}

func (item *ChangeToolModifierEvent) GetJsonType() string {
	return "change-tool-modifier-event"
}
func (item *ChangeToolModifierEvent) FixNils() Message {
	copy := *item
	return &copy
}

func (item *PaintingChangedEvent) GetJsonType() string {
	return "painting-changed-event"
}
func (item *PaintingChangedEvent) FixNils() Message {
	copy := *item
	return &copy
}

func (item *PlayersChangedEvent) GetJsonType() string {
	return "players-changed-event"
}
func (item *PlayersChangedEvent) FixNils() Message {
	copy := *item
	if copy.Players == nil {
		copy.Players = []string{}
	}
	return &copy
}

func (item *PlayerReadyChangedEvent) GetJsonType() string {
	return "player-ready-changed-event"
}
func (item *PlayerReadyChangedEvent) FixNils() Message {
	copy := *item
	if copy.Players == nil {
		copy.Players = map[string]bool{}
	}
	return &copy
}

func (item *PopUpEvent) GetJsonType() string {
	return "pop-up-event"
}
func (item *PopUpEvent) FixNils() Message {
	copy := *item
	return &copy
}

func (item *DebugMessageEvent) GetJsonType() string {
	return "debug-message-event"
}
func (item *DebugMessageEvent) FixNils() Message {
	copy := *item
	return &copy
}

