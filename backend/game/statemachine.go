package game

type State interface {
	Enter()
	Exit()
	Update(*StateMachine, *PlayerMessage)
}

type StateMachine struct {
	currentState State
	states       map[string]State
}

func (sm *StateMachine) setState(s State) {
	sm.currentState = s
	sm.currentState.Enter()
}

func (sm *StateMachine) Transition(pmsg *PlayerMessage) {
	sm.currentState.Update(sm, pmsg)
}

func NewStateMachine(initialState State) *StateMachine {
	sm := &StateMachine{
		currentState: initialState,
		states:       make(map[string]State),
	}

	sm.currentState.Enter()
	return sm
}
