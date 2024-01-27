package game

import (
	"fmt"
)

// ------------------------------------------------------------
type Lobby struct{}

func (g Lobby) Enter() {
	fmt.Println("sm: enter Lobby")
}
func (g Lobby) Exit() {
	fmt.Println("sm: exit Lobby")
}
func (g Lobby) Update(sm *StateMachine, pmsg *PlayerMessage) {
	fmt.Println("sm: update Lobby")

	var startPressed bool = true
	if startPressed {
		sm.setState(&SelectPrompt{})
	}
}

// ------------------------------------------------------------
type SelectPrompt struct{}

func (g SelectPrompt) Enter() {
	fmt.Println("sm: enter SelectPrompt")
}
func (g SelectPrompt) Exit() {
	fmt.Println("sm: exit SelectPrompt")
}
func (g SelectPrompt) Update(sm *StateMachine, pmsg *PlayerMessage) {
	fmt.Println("sm: update SelectPrompt")
	var voteDone bool = true
	if voteDone {
		sm.setState(&Painting{})
	}
}

// ------------------------------------------------------------
type Painting struct{}

func (g Painting) Enter() {
	fmt.Println("sm: enter Painting")
}
func (g Painting) Exit() {
	fmt.Println("sm: exit Painting")
}
func (g Painting) Update(sm *StateMachine, pmsg *PlayerMessage) {
	fmt.Println("sm: update Painting")

	var state State = SelectStickers{}
	switch state {
	case PaintingAndVoting{}:
		sm.setState(&PaintingAndVoting{})
	case SelectStickers{}:
		sm.setState(&SelectStickers{})
	default:
		sm.setState(&SelectStickers{})
	}
}

// ------------------------------------------------------------
type PaintingAndVoting struct{}

func (g PaintingAndVoting) Enter() {
	fmt.Println("sm: enter PaintingAndVoting")
}
func (g PaintingAndVoting) Exit() {
	fmt.Println("sm: exit PaintingAndVoting")
}
func (g PaintingAndVoting) Update(sm *StateMachine, pmsg *PlayerMessage) {
	fmt.Println("sm: update PaintingAndVoting")
	sm.setState(&Painting{})
}

// ------------------------------------------------------------
type SelectStickers struct{}

func (g SelectStickers) Enter() {
	fmt.Println("sm: enter SelectStickers")
}
func (g SelectStickers) Exit() {
	fmt.Println("sm: exit SelectStickers")
}
func (g SelectStickers) Update(sm *StateMachine, pmsg *PlayerMessage) {
	fmt.Println("sm: update SelectStickers")
	var places bool = true
	if places {
		sm.setState(&Showcase{})
	}
}

// ------------------------------------------------------------
type Showcase struct{}

func (g Showcase) Enter() {
	fmt.Println("sm: enter Showcase")
}
func (g Showcase) Exit() {
	fmt.Println("sm: exit Showcase")
}
func (g Showcase) Update(sm *StateMachine, pmsg *PlayerMessage) {
	fmt.Println("sm: update Showcase")
	var allSeen bool = true
	if allSeen {
		var allDone bool = true
		if allDone {
			sm.setState(&VoteForBest{})
		} else {
			sm.setState(&Lobby{})
		}
	}
}

// ------------------------------------------------------------
type VoteForBest struct{}

func (g VoteForBest) Enter() {
	fmt.Println("sm: enter VoteForBest")
}
func (g VoteForBest) Exit() {
	fmt.Println("sm: exit VoteForBest")
}
func (g VoteForBest) Update(sm *StateMachine, pmsg *PlayerMessage) {
	fmt.Println("sm: update VoteForBest")
	var allVotes bool = true
	if allVotes {
		sm.setState(&ShowBest{})
	}
}

// ------------------------------------------------------------
type ShowBest struct{}

func (g ShowBest) Enter() {
	fmt.Println("sm: enter ShowBest")
}
func (g ShowBest) Exit() {
	fmt.Println("sm: exit ShowBest")
}
func (g ShowBest) Update(sm *StateMachine, pmsg *PlayerMessage) {
	fmt.Println("sm: update ShowBest")
	var allSeen bool = true
	if allSeen {
		sm.setState(&Lobby{})
	}
}
