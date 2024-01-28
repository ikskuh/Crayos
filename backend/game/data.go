package game

import (
	_ "embed"
	"strings"
	"time"
)

const (
	LIMIT_MAX_PLAYERS      int = 4  // Maximum number of players per session
	LIMIT_MAX_NICKNAME_LEN int = 20 // Maximum number of "chars" in the player name
)

var (
	// Duration for voting the prompt in seconds
	TIME_GAME_PROMPTVOTE_S = 20

	// Duration of a drawing round in seconds
	TIME_GAME_PAINTING_S = 90

	/// Time of a "trolling" time slice in seconds
	TIME_GAME_NEXT_TROLLEFFECT_S = 10

	// Duration of the stickering phase in seconds
	TIME_GAME_STICKERING_S = 20

	// Duration of the showcase phase in seconds
	TIME_GAME_SHOWCASE_S = 15

	// Duration of the picture rating in seconds
	TIME_GAME_RATING_S = 20

	// Retains the image for 1 second after the last vote.
	TIME_GAME_RATING_SLACK time.Duration = 1 * time.Second

	// Duration of the gallery
	TIME_GAME_GALLERY_S = 20

	/// Time how long a "troll" effect does last in milliseconds
	TIME_GAME_TROLL_EFFECT_DURATION_MS = 5000

	/// Timeout for generic announcements
	TIME_ANNOUNCE_GENERIC time.Duration = 3 * time.Second

	/// Duration of a regular popup
	TIME_POPUP_DURATION_MS = 1500
)

const (
	// Error messages:
	TEXT_ERROR_NICK_EMPTY     string = "Empty nick not allowed"
	TEXT_ERROR_NICK_TOO_LONG  string = "Nickname too long!"
	TEXT_ERROR_SESSION_EMPTY  string = "Empty session id not allowed"
	TEXT_ERROR_BAD_SESSION    string = "Session does not exist"
	TEXT_ERROR_SESSION_ONLINE string = "Session is already running."
	TEXT_ERROR_SESSION_FULL   string = "Lobby is already full."

	// Popup messages:
	TEXT_POPUP_START_PAINTING   string = "Start painting!"
	TEXT_POPUP_STOP_PAINTING    string = "Times up!"
	TEXT_POPUP_START_TROLLING   string = "Start trolling!"
	TEXT_POPUP_MISSED_TROLLING  string = "You sleepyhead!"
	TEXT_POPUP_START_STICKERING string = "Let's make a mess!"
	TEXT_POPUP_STOP_STICKERING  string = "Enough of that!"
	TEXT_POPUP_TIMES_UP         string = "Someone's sleepy!"

	// Vote Prompts:
	TEXT_VOTE_PROMPT   string = "Select a prompt"
	TEXT_VOTE_EFFECT   string = "Select a trolling effect"
	TEXT_VOTE_SHOWCASE string = "Gaze upon this masterpiece"

	// Announcements:
	TEXT_ANNOUNCE_YOU_ARE_TROLL   string = "Chose an image that should be drawn"
	TEXT_ANNOUNCE_YOU_ARE_PAINTER string = "You are the painter. Brace yourself!"
	TEXT_ANNOUNCE_WINNER          string = "And the winner is..."
)

//go:embed drawing_prompts_the_other_kind_of_drawcalls.txt
var fileData []byte
var AVAILABLE_PROMPTS = strings.Split(string(fileData), "\n")
