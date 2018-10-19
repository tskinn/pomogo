package pomogo

import (
	"time"
)

type Request struct {
	RequestType `json:"type"`
	TaskID      string `json:"task_id"`
	Username    string `json:"username"`
}

type RequestType int

const (
	RequestTypeStart RequestType = iota
	RequestTypeStop
)

type Config struct {
	Option            string
	DurationRestLong  time.Duration
	DurationRestShort time.Duration
	DurationWork      time.Duration
}

type User struct {
	Username        string  `json:"username"`
	Session         Session `json:"-"`
	PreviousSession Session `json:"-"`
	Task            Task    `json:"task"`
	Config          Config  `json:"config"`
	RunningSessions uint    `json:"running_sesions"` // consecutive work sessions
}

type Task struct {
	Description string `json:"description"`
	End         string `json:"end"`   // don't care about value, just if exists
	Start       string `json:"start"` // don't care about value, just if exists
	Status      string `json:"status"`
	UUID        string `json:"uuid"`
}

type SessionStatus int

const (
	SessionStarted SessionStatus = iota
	SessionInterrupted
	SessionCompleted
)

type SessionType uint8

const (
	SessionTypeWork SessionType = iota
	SessionTypeRest
)

type Session struct {
	Start            time.Time
	End              time.Time
	Timer            *time.Timer
	Status           SessionStatus
	Type             SessionType
	Duration         time.Duration
	StartActions     []Action
	InterruptActions []Action
	CompleteActions  []Action
	Interrupt        chan bool `json:"-"`
}

func (user *User) startSession(task Task, length int, sessionType SessionType) {
	user.Session.Status = SessionStarted
	user.Session.Type = sessionType
	user.Session.Start = time.Now()
	user.Session.Timer = time.AfterFunc(time.Minute*time.Duration(length), func() {
		user.RunCompleteActions()
	})
	user.Session.End = time.Now().Add(time.Minute * time.Duration(length))
}

func (user *User) stopSession() {
	user.Session.Status = SessionInterrupted
	user.Session.Timer.Stop() // TODO do  we need to check result?
	user.PreviousSession = user.Session
	user.RunInterruptActions()
}

func (user *User) RunStartActions() {
	runActions(user, user.Session.StartActions)
}

func (user *User) RunCompleteActions() {
	runActions(user, user.Session.CompleteActions)
}

func (user *User) RunInterruptActions() {
	runActions(user, user.Session.InterruptActions)
}

func runActions(user *User, hooks []Action) {
	for _, hook := range hooks {
		if ok := hook(user); !ok {
			return
		}
	}
}

func TaskStarted(oldTask, newTask *Task) bool {
	if newTask.Start != "" && oldTask.Start == "" {
		return true
	}
	return false
}

func TaskStopped(oldTask, newTask *Task) bool {
	if newTask.Start == "" && oldTask.Start != "" {
		return true
	}
	return false
}

