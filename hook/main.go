package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"os/user"

	"github.com/tskinn/pomodorogo"
)

func main() {
	reader := bufio.NewReader(os.Stdin)

	// get old task
	oldRaw, _, _ := reader.ReadLine()
	oldTask := pomodorogo.Task{}
	err := json.Unmarshal(oldRaw, &oldTask)
	if err != nil {
		fmt.Println(err)
		return
	}

  // get new task
	newRaw, _, _ := reader.ReadLine()
	newTask := pomodorogo.Task{}
	err = json.Unmarshal(newRaw, &newTask)
	if err != nil {
  	fmt.Println(err)
  	return
	}

  // create user
	username := "unknown"
	osUser, err := user.Current()
	if err == nil {
		username = osUser.Username
	} else {
		//log it
	}

	user := pomodorogo.User{
		ID:      username,
		Session: pomodorogo.Session{},
		Task:    newTask,
	}

	client := pomodorogo.Client{}
	err = client.Connect()
	if err != nil {
    // do soemthing
    return
  }

  err = client.SessionStart(user)
  if err != nil {
    return
  }
	// user has session and newTask
	// send start session rpc?

	fmt.Printf("%v\n", newTask)
}
