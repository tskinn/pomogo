package pomogo

import (
	"sync"
	"time"
)

var state *State

// State ...
type State struct {
	mutex *sync.Mutex
	user  User
}

func init() {
	state = &State{
		&sync.Mutex{},
		User{
			Config: Config{
				DurationRestLong:  time.Second * 25,
				DurationRestShort: time.Second * 5,
				DurationWork:      time.Second * 25,
			},
			StartActions:     defaultStartActions,
			InterruptActions: defaultInteruptActions,
			CompleteActions:  defaultCompleteActions,
		},
	}
}

// StartSession ...
func (state *State) StartSession(taskUUID string) {
	state.mutex.Lock()
	defer state.mutex.Unlock()

	state.user.Task.UUID = taskUUID
	state.user.startSession()
}

// StopSession ...
func (state *State) StopSession() {
	state.mutex.Lock()
	defer state.mutex.Unlock()

	state.user.stopSession()
}

// GetStatusResponse ...
func (state *State) GetStatusResponse() Response {
	state.mutex.Lock()
	defer state.mutex.Unlock()

	response := Response{
		TaskID:  state.user.Task.UUID,
		Start:   state.user.Session.Start,
		End:     state.user.Session.End,
		Running: state.user.RunningSessions,
	}

	switch state.user.Session.Status {
	case SessionStarted:
		response.Status = StatusRunning
	case SessionInterrupted:
		response.Status = StatusStopped
	}

	switch state.user.Session.Type {
	case SessionTypeRest:
		response.Type = "rest"
	case SessionTypeWork:
		response.Type = "work"
	}

	return response
}
