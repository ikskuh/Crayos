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
)


func SerializeMessage(msg Message) ([]byte, error) {

	temp, err := json.Marshal(msg)
	if err != nil {
		return nil, err
	}

	var dummy map[string]interface{}

	err = json.Unmarshal(temp, &dummy)
	if err != nil {
		return nil, err
	}

    dummy["type"] = msg.GetJsonType()

	return json.Marshal(dummy)
}

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
)
var ALL_GAME_VIEW_ITEMS = []GameView{
	"title",
	"lobby",
	"promptselection",
	"artstudio-generic",
	"artstudio-active",
	"artstudio-sticker",
	"gallery",
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
	USER_ACTION_CONTINUE_GAME UserAction = "continue"
)
var ALL_USER_ACTION_ITEMS = []UserAction{
	"set-ready",
	"set-not-ready",
	"continue",
}

type Backdrop string
const (
	BACKDROP_ARCTIC Backdrop = "arctic"
	BACKDROP_GRAVEYARD Backdrop = "graveyard"
	BACKDROP_PIRATE_SHIP Backdrop = "pirate_ship"
	BACKDROP_THEATER_STAGE1 Backdrop = "theater_stage1"
)
var ALL_BACKDROP_ITEMS = []Backdrop{
	"arctic",
	"graveyard",
	"pirate_ship",
	"theater_stage1",
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
	Path interface{} `json:"path"`
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

type ChangeGameViewEvent struct {
	View GameView `json:"view"`
	Painting interface{} `json:"painting"`
	PaintingPrompt string `json:"paintingPrompt"`
	PaintingBackdrop Backdrop `json:"paintingBackdrop"`
	PaintingStickers []Sticker `json:"paintingStickers"`
	VotePrompt string `json:"votePrompt"`
	VoteOptions []string `json:"voteOptions"`
}

type TimerChangedEvent struct {
	SecondsLeft int `json:"secondsLeft"`
}

type ChangeToolModifierEvent struct {
	Modifier Effect `json:"modifier"`
}

type PaintingChangedEvent struct {
	Path interface{} `json:"path"`
}

type PlayersChangedEvent struct {
	Players []string `json:"players"`
	AddedPlayer *string `json:"addedPlayer"`
	RemovedPlayer *string `json:"removedPlayer"`
}

type PlayerReadyChangedEvent struct {
	Players map[string]bool `json:"players"`
}


func (item *CreateSessionCommand) GetJsonType() string {
	return "create-session-command"
}

func (item *JoinSessionCommand) GetJsonType() string {
	return "join-session-command"
}

func (item *LeaveSessionCommand) GetJsonType() string {
	return "leave-session-command"
}

func (item *UserCommand) GetJsonType() string {
	return "user-command"
}

func (item *VoteCommand) GetJsonType() string {
	return "vote-command"
}

func (item *PlaceStickerCommand) GetJsonType() string {
	return "place-sticker-command"
}

func (item *SetPaintingCommand) GetJsonType() string {
	return "set-painting-command"
}

func (item *EnterSessionEvent) GetJsonType() string {
	return "enter-session-event"
}

func (item *JoinSessionFailedEvent) GetJsonType() string {
	return "join-session-failed-event"
}

func (item *KickedEvent) GetJsonType() string {
	return "kicked-event"
}

func (item *ChangeGameViewEvent) GetJsonType() string {
	return "change-game-view-event"
}

func (item *TimerChangedEvent) GetJsonType() string {
	return "timer-changed-event"
}

func (item *ChangeToolModifierEvent) GetJsonType() string {
	return "change-tool-modifier-event"
}

func (item *PaintingChangedEvent) GetJsonType() string {
	return "painting-changed-event"
}

func (item *PlayersChangedEvent) GetJsonType() string {
	return "players-changed-event"
}

func (item *PlayerReadyChangedEvent) GetJsonType() string {
	return "player-ready-changed-event"
}

