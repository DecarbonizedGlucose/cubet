package user

type UserState int

const (
	SttInitialOrStopped UserState = iota
	SttToLaunch
	SttPreparing
	SttToStart
	SttOverPrepared
	SttOverOverPrepared
	SttSolving
)

type UserAction int

const (
	ActReadyToPrepare UserAction = iota
	ActStartToPrepare
	ActCancelBeforeSolve
	ActReadyToSolve
	ActStartSolving
	ActCancelDuringSolve
	ActStopSolving
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
