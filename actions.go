package pomogo

import (
	"log"
	"time"
)

// Action ...
type Action func(user *User) bool

// SessionSetStatusComplete ...
func SessionSetStatusComplete(user *User) bool {
	log.Println("SessionSetStatusComplete")
	user.Session.Status = SessionCompleted
	return true
}

// SessionSetStatusInterrupted ...
func SessionSetStatusInterrupted(user *User) bool {
	log.Println("SessionSetStatusInterrupted")
	user.Session.Status = SessionInterrupted
	return true
}

// SessionSetStatusRunning ...
func SessionSetStatusRunning(user *User) bool {
	log.Println("SessionSetStatusRunning")
	user.Session.Status = SessionStarted
	return true
}

// SessionSetPreviousSession ...
func SessionSetPreviousSession(user *User) bool {
	log.Println("SessionSetPreviousSession")
	user.PreviousSession = user.Session
	return true
}

// SessionSetDuration ...
func SessionSetDuration(user *User) bool {
	log.Println("SessionSetDuration")
	switch user.PreviousSession.Type {
	case SessionTypeWork:
		if user.RunningSessions == 4 { // TODO make 4 a config variable?
			user.Session.Duration = user.Config.DurationRestLong
		} else {
			user.Session.Duration = user.Config.DurationRestShort
		}
	case SessionTypeRest:
		user.Session.Duration = user.Config.DurationWork
	}
	return true
}

// UpdateRunningSession increments the running session count if
// previous session was a work session.
func UpdateRunningSession(user *User) bool {
	log.Println("UpdateRunningSession")
	// only update on completed work session
	if user.PreviousSession.Type == SessionTypeRest {
		return true
	}

	if user.RunningSessions == 4 {
		user.RunningSessions = 0
	} else {
		user.RunningSessions++
	}
	return true
}

// SessionSetType ...
func SessionSetType(user *User) bool {
	log.Println("SessionSetType")
	switch user.PreviousSession.Type {
	case SessionTypeRest:
		user.Session.Type = SessionTypeWork
	case SessionTypeWork:
		user.Session.Type = SessionTypeRest
	}
	return true
}

// SessionStart ...
func SessionStart(user *User) bool {
	log.Println("SessionStart")
	user.Session.Start = time.Now()
	user.Session.Timer = time.AfterFunc(user.Session.Duration, func() {
		user.RunCompleteActions()
	})
	user.Session.End = time.Now().Add(user.Session.Duration)
	return true
}

// SessionStopTimer stops the timer. This is used when a session is interrupted.
func SessionStopTimer(user *User) bool {
	log.Println("SessionStopTimer")
	user.Session.Timer.Stop()
	return true
}

// SessionStopModifyEndTime modifies the end time. Used when a session is interrupted.
func SessionStopModifyEndTime(user *User) bool {
	log.Println("SessionStopModifyEndTime")
	user.Session.End = time.Now()
	return true
}

// RunTaskCommand will either trigger a `task <id> start` or a `task <id> stop`.
// This is meant to be run after a session is complete in order to trigger the
// next session.
func RunTaskCommand(user *User) bool {
	log.Println("RunTaskCommand")
	switch user.PreviousSession.Type {
	case SessionTypeRest:
		StartTask(user.Task.UUID)
	case SessionTypeWork:
		StopTask(user.Task.UUID)
	}
	return true
}

// StartNewSession ...
func StartNewSession(user *User) bool {
	log.Println("StartNewSession")
	user.RunStartActions()
	return true
}

var defaultStartActions = []Action{
	SessionSetStatusRunning,
	SessionSetDuration,
	SessionSetType,
	SessionStart,
}

var defaultCompleteActions = []Action{
	SessionSetStatusComplete,
	SessionSetPreviousSession,
	UpdateRunningSession,
	RunTaskCommand,
	StartNewSession,
}

var defaultInteruptActions = []Action{
	SessionStopTimer,
	SessionStopModifyEndTime,
	SessionSetStatusInterrupted,
	SessionSetPreviousSession,
}
