package timer

import (
	"context"
	"time"
)

type FunctionalTimer struct {
	solvingTimer     *BasicTimer
	preparationTimer *BasicTimer

	UserState         UserState
	enablePreparation bool
	alarmCh           chan time.Duration
	operationCh       chan UserAction
	resultCh          chan time.Duration
	punishmentCh      chan Punishment

	ctx        context.Context
	cancelFunc context.CancelFunc
}

func NewFunctionalTimer(
	lastTime time.Duration,
	enablePreparation bool,
	opch chan UserAction,
	rech chan time.Duration,
	puch chan Punishment,
) *FunctionalTimer {
	ft := &FunctionalTimer{
		solvingTimer:      NewBasicTimer(lastTime),
		preparationTimer:  NewBasicTimer(0),
		UserState:         SttInitialOrStopped,
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
		if action != ActExitTimer {
			ft.DoOperation(action)
		} else {
			ft.solvingTimer.Stop()
			return
		}
	}
}

func (ft *FunctionalTimer) DoOperation(action UserAction) {
	switch action {
	case ActReadyToPrepare:
		ft.UserState = SttToLaunch
	case ActStartToPrepare:
		ft.UserState = SttPreparing
		ft.preparationTimer.Start(
			ft.alarmCh,
			8*time.Second,
			12*time.Second,
			15*time.Second,
			17*time.Second,
		)
	case ActCancelBeforeSolve:
		ft.UserState = SttInitialOrStopped
		ft.preparationTimer.Stop()
		ft.punishmentCh <- PnsDidNotStart
	case ActReadyToSolve:
		ft.UserState = SttToStart
	case ActStartSolving:
		ft.UserState = SttSolving
		ft.preparationTimer.Stop()
		ft.solvingTimer.Start(nil)
	case ActCancelDuringSolve:
		ft.UserState = SttInitialOrStopped
		elapsed := ft.solvingTimer.Stop()
		ft.resultCh <- elapsed
		ft.punishmentCh <- PnsDidNotFinish
	case ActStopSolving:
		ft.UserState = SttInitialOrStopped
		elapsed := ft.solvingTimer.Stop()
		ft.resultCh <- elapsed
	default:
	}
}
