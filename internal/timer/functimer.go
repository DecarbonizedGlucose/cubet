package timer

import (
	"time"

	"github.com/DecarbonizedGlucose/cubet/internal/user"
)

type FunctionalTimer struct {
	solvingTimer     *BasicTimer
	preparationTimer *BasicTimer

	UserState         user.UserState
	enablePreparation bool
	alarmCh           chan time.Duration
	operationCh       chan user.UserAction
	resultCh          chan time.Duration
	punishmentCh      chan user.Punishment
}

func NewFunctionalTimer(
	lastTime time.Duration,
	enablePreparation bool,
	opch chan user.UserAction,
	rech chan time.Duration,
	puch chan user.Punishment,
) *FunctionalTimer {
	ft := &FunctionalTimer{
		solvingTimer:      NewBasicTimer(lastTime),
		preparationTimer:  NewBasicTimer(0),
		UserState:         user.SttInitialOrStopped,
		enablePreparation: enablePreparation,
		alarmCh:           make(chan time.Duration),
		operationCh:       opch,
		resultCh:          rech,
		punishmentCh:      puch,
	}
	go ft.ListenUserAction()
	return ft
}

func (ft *FunctionalTimer) ListenUserAction() {
	for {
		action := <-ft.operationCh
		if action != user.ActExitTimer {
			ft.DoOperation(action)
		} else {
			ft.preparationTimer.Stop()
			ft.solvingTimer.Stop()
			return
		}
	}
}

func (ft *FunctionalTimer) DoOperation(action user.UserAction) {
	switch action {
	case user.ActReadyToPrepare:
		ft.UserState = user.SttToLaunch
	case user.ActStartToPrepare:
		ft.UserState = user.SttPreparing
		ft.preparationTimer.Start(
			ft.alarmCh,
			8*time.Second,
			12*time.Second,
			15*time.Second,
			17*time.Second,
		)
	case user.ActCancelBeforeSolve:
		ft.UserState = user.SttInitialOrStopped
		ft.preparationTimer.Stop()
		ft.punishmentCh <- user.PnsDidNotStart
	case user.ActReadyToSolve:
		ft.UserState = user.SttToStart
	case user.ActStartSolving:
		ft.UserState = user.SttSolving
		ft.preparationTimer.Stop()
		ft.solvingTimer.Start(nil)
	case user.ActCancelDuringSolve:
		ft.UserState = user.SttInitialOrStopped
		elapsed := ft.solvingTimer.Stop()
		ft.resultCh <- elapsed
		ft.punishmentCh <- user.PnsDidNotFinish
	case user.ActStopSolving:
		ft.UserState = user.SttInitialOrStopped
		elapsed := ft.solvingTimer.Stop()
		ft.resultCh <- elapsed
	default:
	}
}
