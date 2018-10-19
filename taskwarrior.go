package pomogo

import (
	"bytes"
	"os/exec"
)

func StopTask(taskUUID string) ([]byte, error) {
	cmd := exec.Command("task", taskUUID, "stop")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	return out.Bytes(), err
}

func StartTask(taskUUID string) ([]byte, error) {
	cmd := exec.Command("task", taskUUID, "start")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	return out.Bytes(), err
}


