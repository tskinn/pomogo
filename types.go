package main

import (
	"fmt"
	"time"
)

type Task struct {
	Status      string `json:"status"`
	UUID        string `json:"uuid"`
	Description string `json:"description"`
	Start       string `json:"start"`
	End         string `json:"end"`
}

type Hook func(t Task) bool

type SessionStatus int

const (
	SessionStarted SessionStatus = iota
	SessionInterrupted
	SessionCompleted
)

type Session struct {
	Task
	Start          time.Time
	End            time.Time
	Status         SessionStatus
	Duration       time.Duration
	StartHooks     []Hook
	InterruptHooks []Hook
	CompleteHooks  []Hook
	Interrupt      chan bool
}

func (session *Session) RunStartHooks() {
	runHooks(session.Task, session.StartHooks)
}

func (session *Session) RunCompleteHooks() {
	runHooks(session.Task, session.CompleteHooks)
}

func (session *Session) RunInterruptHooks() {
	runHooks(session.Task, session.InterruptHooks)
}

func runHooks(task Task, hooks []Hook) {
	for _, hook := range hooks {
		if ok := hook(task); !ok {
			return
		}
	}
}

func (session *Session) Begin() {
	session.RunStartHooks()
	session.Status = SessionStarted
	session.Start = time.Now()
	timer := time.NewTimer(session.Duration)

	select {
	case doneTime := <-timer.C:
		session.Status = SessionCompleted
		session.End = doneTime
		session.RunCompleteHooks()
	case <-session.Interrupt:
		session.Status = SessionInterrupted
		session.RunInterruptHooks()
	}
}

var exampleSession Session = Session{
	Task:     Task{},
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
}
