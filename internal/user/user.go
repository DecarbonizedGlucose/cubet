package user

type UserState int

const (
	// free time
	SttFree UserState = iota
	// ready to solve, about to prepare
	SttToPrepare
	// preparation time
	SttPreparing
	// hands timer, release to start timer
	SttToStart
	// preparation time in (15, 17] seconds
	SttOverPrepared
	// preparation time > 17 seconds
	SttOverOverPrepared
	// solving time
	SttSolving
)

type UserAction int

const (
	// go to SttToPrepare
	ActReadyToPrepare UserAction = iota
	// go to SttPreparing
	ActStartToPrepare
	// go to SttFree
	ActCancelBeforeSolve
	// go to SttToStart
	ActReadyToSolve
	// go to SttSolving
	ActStartSolving
	// go to SttFree
	ActCancelDuringSolve
	// go to SttFree
	ActStopSolving
	// exit all timer goroutines
	ActExitTimer
)

type Punishment int

const (
	PnsNoPunishment Punishment = iota
	PnsPlusTwoSeconds
	PnsDidNotFinish
	PnsDidNotStart
)

func (p Punishment) String() string {
	switch p {
	case PnsNoPunishment:
		return "None"
	case PnsPlusTwoSeconds:
		return "+2"
	case PnsDidNotFinish:
		return "DNF"
	case PnsDidNotStart:
		return "DNS"
	default:
		return "Unknown Punishment"
	}
}
