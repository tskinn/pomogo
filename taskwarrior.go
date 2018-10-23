package pomogo

import (
	"bytes"
	"log"
	"os/exec"
)

// StopTask ...
func StopTask(taskUUID string) ([]byte, error) {
	log.Println("StopTask")
	cmd := exec.Command("task", taskUUID, "stop")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	return out.Bytes(), err
}

// StartTask ...
func StartTask(taskUUID string) ([]byte, error) {
	log.Println("StartTask")
	cmd := exec.Command("task", taskUUID, "start")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	return out.Bytes(), err
}
