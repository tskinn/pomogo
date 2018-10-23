package pomogo

import (
	"time"
)

// Request is the request made from the client to server
type Request struct {
	RequestType `json:"type"`
	TaskID      string `json:"task_id"`
	Username    string `json:"username"`
}

// RequestType is
type RequestType int

const (
	// RequestTypeStart start
	RequestTypeStart RequestType = iota
	// RequestTypeStop stop
	RequestTypeStop
	// RequestTypeStatus get status
	RequestTypeStatus
)

// Response ...
type Response struct {
	End     time.Time `json:"end"`
	Running uint8     `json:"running"`
	Start   time.Time `json:"start"`
	Status  string    `json:"status"`
	TaskID  string    `json:"task_id"`
	Type    string    `json:"type"`
}

const (
	StatusRunning = "running"
	StatusStopped = "stopped"
)

// Config is user config
type Config struct {
	Option            string
	DurationRestLong  time.Duration
	DurationRestShort time.Duration
	DurationWork      time.Duration
}

// User is garbage. get rid of it
type User struct {
	Username         string  `json:"username"`
	Session          Session `json:"-"`
	PreviousSession  Session `json:"-"`
	Task             Task    `json:"task"`
	Config           Config  `json:"config"`
	RunningSessions  uint8   `json:"running_sesions"` // consecutive work sessions
	StartActions     []Action
	InterruptActions []Action
	CompleteActions  []Action
}

// Task is task from task warrior
type Task struct {
	Description string `json:"description"`
	End         string `json:"end"`   // don't care about value, just if exists
	Start       string `json:"start"` // don't care about value, just if exists
	Status      string `json:"status"`
	UUID        string `json:"uuid"`
}

// SessionStatus ...
type SessionStatus int

const (
	// SessionStarted ...
	SessionStarted SessionStatus = iota
	// SessionInterrupted ...
	SessionInterrupted
	// SessionCompleted ...
	SessionCompleted
)

// SessionType ...
type SessionType uint8

const (
	// SessionTypeWork ...
	SessionTypeWork SessionType = iota
	// SessionTypeRest ...
	SessionTypeRest
)

// Session ...
type Session struct {
	Start    time.Time
	End      time.Time
	Timer    *time.Timer
	Status   SessionStatus
	Type     SessionType
	Duration time.Duration
}

func (user *User) startSession() {
	user.RunStartActions()
}

func (user *User) stopSession() {
	user.RunInterruptActions()
}

// RunStartActions ...
func (user *User) RunStartActions() {
	runActions(user, user.StartActions)
}

// RunCompleteActions ...
func (user *User) RunCompleteActions() {
	runActions(user, user.CompleteActions)
}

// RunInterruptActions ...
func (user *User) RunInterruptActions() {
	runActions(user, user.InterruptActions)
}

func runActions(user *User, hooks []Action) {
	for _, hook := range hooks {
		if ok := hook(user); !ok {
			return
		}
	}
}

// TaskStarted ...
func TaskStarted(oldTask, newTask *Task) bool {
	if newTask.Start != "" && oldTask.Start == "" {
		return true
	}
	return false
}

// TaskStopped ...
func TaskStopped(oldTask, newTask *Task) bool {
	if newTask.Start == "" && oldTask.Start != "" {
		return true
	}
	return false
}
