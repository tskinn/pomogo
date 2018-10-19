package pomogo

import (
  "time"
)

type Action func(user *User) bool

func SessionSetStatusComplete(user *User) bool {
	user.Session.Status = SessionCompleted
	return true
}

func SessionSetStatusInterrupted(user *User) bool {
	user.Session.Status = SessionInterrupted
	return true
}

func SessionSetStatusRunning(user *User) bool {
	user.Session.Status = SessionStarted
	return true
}

func SessionSetPreviousSession(user *User) bool {
  user.PreviousSession = user.Session
  return true
}

func SessionSetDuration(user *User) bool {
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

func UpdateRunningSession(user *User) bool {
  // only update on completed work session
  if user.PreviousSession.Type == SessionTypeRest {
    return true
  }

	if user.RunningSessions == 4 {
		user.RunningSessions = 0
	} else {
		user.RunningSessions += 1
	}
	return true
}

func SessionSetType(user *User) bool {
	switch user.PreviousSession.Type {
	case SessionTypeRest:
		user.Session.Type = SessionTypeWork
	case SessionTypeWork:
		user.Session.Type = SessionTypeRest
	}
	return true
}

func SessionStart(user *User) bool {
	user.Session.Start = time.Now()
	user.Session.Timer = time.AfterFunc(user.Session.Duration, func() {
		user.RunCompleteActions()
	})
	user.Session.End = time.Now().Add(user.Session.Duration)
	return true
}

// RunTaskCommand will either trigger a `task <id> start` or a `task <id> stop`.
// This is meant to be run after a session is complete in order to trigger the
// next session.
func RunTaskCommand(user *User) bool {
  switch user.PreviousSession.Type {
    case SessionTypeRest:
      TaskStart(user.Task.UUID)
    case SessionTypeWork:
      TaskStop(user.Task.UUID)
  }
  return true
}

var defaultStartActions = []Action{
	SessionSetStatusRunning,
	SessionSetDuration,
	SessionSetType,
	SessionStart,
}

var defaultCompleteActions = []Action{
  SessionSetPreviousSession,
  SessionSetStatusComplete,
  UpdateRunningSession,
  RunTaskCommand,
}


