package pomogo

import (
	"encoding/json"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
)

const defaultSock = "/tmp/pomodorogo.sock"

// Server ...
type Server struct {
	listener net.Listener
}

func handleConnection(conn net.Conn) {
	log.Println("handling connection")
	buff := make([]byte, 1024*10)
	n, err := conn.Read(buff)
	if err != nil {
		log.Println("error reading from connection")
		log.Println(err)
		// do something?
		return
	}
	req := Request{}
	err = json.Unmarshal(buff[0:n], &req)
	if err != nil {
		log.Println(err)
		return
	}
	handleRequest(req, conn)
	conn.Close()
}

func handleRequest(req Request, conn net.Conn) {
	switch req.RequestType {
	case RequestTypeStart:
		state.StartSession(req.TaskID)
		//

	case RequestTypeStop:
		state.StopSession()
		//
	case RequestTypeStatus:
		resp := state.GetStatusResponse()
		respRaw, err := json.Marshal(resp)
		if err != nil {
			// log
			return
		}
		_, err = conn.Write(respRaw)
		if err != nil {
			// log
		}
	}
}

// Start open socket
func (server *Server) Start() {
	listener, err := net.Listen("unix", defaultSock)
	if err != nil {
		log.Println(err)
		// crap
	}
	server.listener = listener
	server.Listen()
}

// Listen on socket
func (server *Server) Listen() {
	log.Println("listening")
	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc, os.Interrupt, syscall.SIGTERM)
	go func(ln net.Listener, c chan os.Signal) {
		sig := <-c
		log.Printf("Caught signal %s: shutting down", sig)
		ln.Close()
		os.Exit(0)
	}(server.listener, sigc)

	for {
		conn, err := server.listener.Accept()
		if err != nil {
			log.Println(err)
			// do something
		}

		go handleConnection(conn)
	}
}

// Stop ...
func (server *Server) Stop() {
	server.listener.Close()
}
