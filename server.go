package pomodorogo

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
)

const defaultSock = "/tmp/pomodorogo.sock"

type Server struct {
	listener net.Listener
}

func handleConnection(conn net.Conn) {
  log.Println("handling connection")
	buff := make([]byte, 1024*5)
	n, err := conn.Read(buff)
	if err != nil {
  	log.Println("error reading from connection")
		log.Println(err)
		// do something?
		return
	}
	log.Println(n)
	fmt.Println(string(buff[:n]))
}

// open socket
func (server *Server) Start() {
	listener, err := net.Listen("unix", defaultSock)
	if err != nil {
		log.Println(err)
		// crap
	}
	server.listener = listener
	server.Listen()
}

// listen on socket
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

func (server *Server) Stop() {
	server.listener.Close()
}
