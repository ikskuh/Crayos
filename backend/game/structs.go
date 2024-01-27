package game

import (
	"encoding/json"
	"errors"
	"reflect"
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
	CHANGE_TOOL_MODIFIER_EVENT_TAG = "change-tool-modifier-event"
	PAINTING_CHANGED_EVENT_TAG = "painting-changed-event"
	PLAYERS_CHANGED_EVENT_TAG = "players-changed-event"
	PLAYER_READY_CHANGED_EVENT_TAG = "player-ready-changed-event"
)

var JSON_TYPE_ID = map[reflect.Type]string{
	reflect.TypeOf(&CreateSessionCommand{}): CREATE_SESSION_COMMAND_TAG,
	reflect.TypeOf(&JoinSessionCommand{}): JOIN_SESSION_COMMAND_TAG,
	reflect.TypeOf(&LeaveSessionCommand{}): LEAVE_SESSION_COMMAND_TAG,
	reflect.TypeOf(&UserCommand{}): USER_COMMAND_TAG,
	reflect.TypeOf(&VoteCommand{}): VOTE_COMMAND_TAG,
	reflect.TypeOf(&PlaceStickerCommand{}): PLACE_STICKER_COMMAND_TAG,
	reflect.TypeOf(&SetPaintingCommand{}): SET_PAINTING_COMMAND_TAG,
	reflect.TypeOf(&EnterSessionEvent{}): ENTER_SESSION_EVENT_TAG,
	reflect.TypeOf(&JoinSessionFailedEvent{}): JOIN_SESSION_FAILED_EVENT_TAG,
	reflect.TypeOf(&KickedEvent{}): KICKED_EVENT_TAG,
	reflect.TypeOf(&ChangeGameViewEvent{}): CHANGE_GAME_VIEW_EVENT_TAG,
	reflect.TypeOf(&ChangeToolModifierEvent{}): CHANGE_TOOL_MODIFIER_EVENT_TAG,
	reflect.TypeOf(&PaintingChangedEvent{}): PAINTING_CHANGED_EVENT_TAG,
	reflect.TypeOf(&PlayersChangedEvent{}): PLAYERS_CHANGED_EVENT_TAG,
	reflect.TypeOf(&PlayerReadyChangedEvent{}): PLAYER_READY_CHANGED_EVENT_TAG,
}


type Message interface {
}

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

	dummy["type"] = JSON_TYPE_ID[reflect.TypeOf(msg)]

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

	var out interface{}

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

const (
	GAME_VIEW_TITLE = "title"
	GAME_VIEW_LOBBY = "lobby"
	GAME_VIEW_PROMPTSELECTION = "promptselection"
	GAME_VIEW_ARTSTUDIO_EMPTY = "artstudio-empty"
	GAME_VIEW_ARTSTUDIO_ACTIVE = "artstudio-active"
	GAME_VIEW_EXHIBITION = "exhibition"
	GAME_VIEW_EXHIBITION_VOTING = "exhibition-voting"
	GAME_VIEW_EXHIBITION_STICKERING = "exhibition-stickering"
	GAME_VIEW_SHOWCASE = "showcase"
	GAME_VIEW_GALLERY = "gallery"
)

const (
	EFFECT_FLASHLIGHT = "flashlight"
	EFFECT_DRUNK = "drunk"
	EFFECT_FLIP = "flip"
	EFFECT_SWAP_TOOL = "swap_tool"
	EFFECT_LOCK_PENCIL = "lock_pencil"
)

const (
	USER_ACTION_SET_READY = "set-ready"
	USER_ACTION_SET_NOT_READY = "set-not-ready"
	USER_ACTION_CONTINUE_GAME = "continue"
)

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
	Action string `json:"action"`
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
	View string `json:"view"`
	Painting interface{} `json:"painting"`
	PaintingPrompt *string `json:"paintingPrompt"`
	PaintingBackdrop *string `json:"paintingBackdrop"`
	PaintingStickers []Sticker `json:"paintingStickers"`
	AvailableStickers []string `json:"availableStickers"`
	VotePrompt *string `json:"votePrompt"`
	VoteOptions []string `json:"voteOptions"`
}

type ChangeToolModifierEvent struct {
	Modifier string `json:"modifier"`
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

