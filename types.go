package pomogo

import (
	"fmt"
	"time"
)

type Request struct {
	Action   `json:"action"`
	TaskID   string `json:"task_id"`
	Username string `json:"username"`
}

type Action int

const (
	ActionStart Action = iota
	ActionStop
)

type User struct {
	Username string `json:"username"`
	Session  `json:"-"`
	Task     `json:"task"`
}

type Task struct {
	Description string `json:"description"`
	End         string `json:"end"`   // don't care about value, just if exists
	Start       string `json:"start"` // don't care about value, just if exists
	Status      string `json:"status"`
	UUID        string `json:"uuid"`
}

type Hook func(t Task) bool

type SessionStatus int

const (
	SessionStarted SessionStatus = iota
	SessionInterrupted
	SessionCompleted
)

type Session struct {
	Start          time.Time
	End            time.Time
	Status         SessionStatus
	Duration       time.Duration
	StartHooks     []Hook
	InterruptHooks []Hook
	CompleteHooks  []Hook
	Interrupt      chan bool `json:"-"`
}

func (session *Session) RunStartHooks(task Task) {
	runHooks(task, session.StartHooks)
}

func (session *Session) RunCompleteHooks(task Task) {
	runHooks(task, session.CompleteHooks)
}

func (session *Session) RunInterruptHooks(task Task) {
	runHooks(task, session.InterruptHooks)
}

func runHooks(task Task, hooks []Hook) {
	for _, hook := range hooks {
		if ok := hook(task); !ok {
			return
		}
	}
}

func (user *User) Begin() {
	user.Session.RunStartHooks(user.Task)
	user.Session.Status = SessionStarted
	user.Session.Start = time.Now()
	timer := time.NewTimer(user.Session.Duration)

	select {
	case doneTime := <-timer.C:
		user.Session.Status = SessionCompleted
		user.Session.End = doneTime
		user.Session.RunCompleteHooks(user.Task)
	case <-user.Session.Interrupt:
		user.Session.Status = SessionInterrupted
		user.Session.RunInterruptHooks(user.Task)
	}
}

var exampleSession User = User{
	Username: "user",
	Session: Session{
		Start:    time.Now(),
		Duration: time.Second * 2,
		StartHooks: []Hook{
			func(t Task) bool {
				fmt.Println("start session")
				fmt.Println("starting task...")
				fmt.Println("notifying user...")
				return true
			},
		},
		InterruptHooks: []Hook{
			func(t Task) bool {
				fmt.Println("interrupt session")
				fmt.Println("stopping session...")
				return true
			},
		},
		CompleteHooks: []Hook{
			func(t Task) bool {
				fmt.Println("complete session")
				fmt.Println("stopping task...")
				fmt.Println("notifying user...")
				return true
			},
		},
	},
}

func Started(oldTask, newTask *Task) bool {
	if newTask.Start != "" && oldTask.Start == "" {
		return true
	}
	return false
}

func Stopped(oldTask, newTask *Task) bool {
	if newTask.Start == "" && oldTask.Start != "" {
		return true
	}
	return false
}
