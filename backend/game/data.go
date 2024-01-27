package game

const (
	// Length of a drawing round in seconds
	GAME_ROUND_TIME_S = 90

	/// Time of a "trolling" time slice
	GAME_TROLL_EFFECT_COOLDOWN_S = 10
)

const (
	VOTE_PROMPT_PROMPT string = "Select a prompt"
	VOTE_PROMPT_EFFECT string = "Select a trolling effect"
)

var AVAILABLE_PROMPTS = []string{
	"a vampire who is afraid of garlic bread",
	"a kangaroo making a clever use of its pouch",
	"a princess breaking the stereotype",
	"a t-rex trying to paint its toe nails",
	"a viking tripping over his beard",
}
