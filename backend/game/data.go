package game

import (
	_ "embed"
	"strings"
	"time"
)

const (
	// Length of a drawing round in seconds
	GAME_ROUND_TIME_S = 20 // 90

	// Timeout for the showcase duration
	GALLERY_ROUND_TIME_S = 30

	/// Time of a "trolling" time slice
	GAME_TROLL_EFFECT_COOLDOWN_S = 10

	/// Time how long a "trol
	GAME_TROLL_EFFECT_DURATION_MS = 5000

	// Time for voting
	GAME_TOPIC_VOTE_TIME_S = 20

	/// Timeout for generic announcements
	ANNOUNCE_GENERIC_TIMEOUT time.Duration = 5 * time.Second
)

const (
	TEXT_POPUP_START_PAINTING string = "Start painting!"
	TEXT_POPUP_TIMES_UP       string = "Times up!"

	TEXT_VOTE_PROMPT   string = "Select a prompt"
	TEXT_VOTE_EFFECT   string = "Select a trolling effect"
	TEXT_VOTE_SHOWCASE string = "Gaze upon this masterpiece"

	TEXT_ANNOUNCE_YOU_ARE_TROLL   string = "Chose an image that should be drawn"
	TEXT_ANNOUNCE_YOU_ARE_PAINTER string = "You are the painter. Brace yourself!"
	TEXT_ANNOUNCE_WINNER          string = "And the winner is..."
)

//go:embed drawing_prompts_the_other_kind_of_drawcalls.txt
var fileData []byte
var AVAILABLE_PROMPTS = strings.Split(string(fileData), "\n")
