package pomodorogo

import (
	"encoding/json"
	"log"
	"net"
	"time"
)

type Client struct {
	conn net.Conn
}

func (client *Client) Connect() error {
	c, err := net.Dial("unix", defaultSock)
	if err != nil {
		// do something
		log.Fatal("Dial error", err)
		return err
	}
	client.conn = c
	return nil
}

func (client *Client) read() ([]byte, error) {
	buf := make([]byte, 1024)
	client.conn.SetReadDeadline(time.Now().Add(time.Millisecond * 50)) // do something about the error it returns
	n, err := client.conn.Read(buf[:])
	if err != nil {
		return []byte{}, err
	}
	log.Println("Client got:", string(buf[:n]))
	return buf[:n], nil
}

func (client *Client) write(msg []byte) error {
	_, err := client.conn.Write(msg)
	return err
}

func (client *Client) SessionStart(user User, taskID string) error {
	return client.sendRequest(ActionStart, user, taskID)
}

func (client *Client) SessionStop(user User, taskID string) error {
	return client.sendRequest(ActionStop, user, taskID)
}

func (client *Client) sendRequest(action Action, username, taskID string) error {
	request := Request{
		Username: username,
		TaskID: TaskID,
		Action: ActionStart,
	}
	data, err := json.Marshal(request)
	if err != nil {
  	log.Println(err)
		return err
	}

	return client.write(data)
}