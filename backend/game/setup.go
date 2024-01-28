package game

import "random-projects.net/crayos-backend/meta"

func Setup() {
	if *meta.DEBUG_MODE {
		TIME_GAME_PROMPTVOTE_S = 10
		TIME_GAME_PAINTING_S = 15
		TIME_GAME_NEXT_TROLLEFFECT_S = 0
		TIME_GAME_STICKERING_S = 10
		TIME_GAME_SHOWCASE_S = 5
		TIME_GAME_RATING_S = 10
		TIME_GAME_GALLERY_S = 10
	}
}
