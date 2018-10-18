package main

import (
	"bufio"
	"encoding/json"
	"log"
	"fmt"
	"os"
	"os/user"

	"github.com/tskinn/pomogo"
)

func main() {
	reader := bufio.NewReader(os.Stdin)

	// get old task
	oldRaw, _, _ := reader.ReadLine()
	oldTask := pomogo.Task{}
	err := json.Unmarshal(oldRaw, &oldTask)
	if err != nil {
		log.Println(err)
		return
	}

	// get new task
	newRaw, _, _ := reader.ReadLine()
	newTask := pomogo.Task{}
	err = json.Unmarshal(newRaw, &newTask)
	if err != nil {
		log.Println(err)
		return
	}

	fmt.Println(string(newRaw))

	// create user
	username := "unknown"
	osUser, err := user.Current()
	if err == nil {
		username = osUser.Username
	} else {
		//log it
	}

	client := pomogo.Client{}
	err = client.Connect()
	if err != nil {
		// do soemthing
		return
	}

	err = client.SessionStart(username, newTask.UUID)
	if err != nil {
		return
	}
	// user has session and newTask
	// send start session rpc?

	log.Printf("%v\n", newTask)
}
