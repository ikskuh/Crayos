package game

import "time"

const (
	GAME_ROUND_TIME = 60 * time.Second
)

var AVAILABLE_BACKGROUNDS = []string{
	"arctic",
	"graveyard",
	"pirate_ship",
	"theater_stage1",
}

var AVAILABLE_PROMPTS = []string{
	"a vampire who is afraid of garlic bread",
	"a kangaroo making a clever use of its pouch",
	"a princess breaking the stereotype",
	"a t-rex trying to paint its toe nails",
	"a viking tripping over his beard",
}

var (
	VOTE_PROMPT_PROMPT string = "Select a prompt"
)
