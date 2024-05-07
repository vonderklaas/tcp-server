package main

import (
	"encoding/hex"
	"fmt"
	"log"
	"math/rand"
	"net"
)

type Message struct {
	from    string
	payload []byte
}

type Server struct {
	listenAddr     string
	listener       net.Listener
	quitChannel    chan struct{}
	messageChannel chan Message
	peersMap       map[string]string
}

func NewServer(listenAddr string) *Server {
	return &Server{
		listenAddr:     listenAddr,
		quitChannel:    make(chan struct{}),
		messageChannel: make(chan Message, 10),
		peersMap:       make(map[string]string),
	}
}

func (server *Server) Start() error {

	// Listen for incoming connections
	listener, err := net.Listen("tcp", server.listenAddr)
	if err != nil {
		printError("listenAddr", err)
		return err
	}

	// Cleanup after function return
	defer listener.Close()

	server.listener = listener

	// Start accepting connections
	go server.acceptLoop()

	fmt.Println("Server is listening at: http://localhost:3000")

	// Cleanup (after quit channel is closed)
	// Why? Because if we stop server, other people still can read from it,
	// now we close channel and notify users so they can also gracefuly handle the exit
	<-server.quitChannel

	return nil
}

func (server *Server) acceptLoop() {
	for {
		// Accept incoming connections if someone dials our server
		conn, err := server.listener.Accept()
		if err != nil {
			printError("acceptLoop", err)
			// We want to continue even after an error
			continue
		}

		fmt.Println("New connection:", conn.RemoteAddr())

		// Add new connection to our peers map
		peerID := generateUniqueID()
		server.peersMap[peerID] = conn.RemoteAddr().String()

		// Each time we have a connection, we spin up a goroutine,
		// so we can have millions of connection without any issues
		go server.readLoop(conn)
	}
}

func (server *Server) readLoop(conn net.Conn) {
	// Cleanup after function return
	defer conn.Close()

	// Create a buffer to read client data
	buffer := make([]byte, 2048)

	// Spin up a loop to read bytes from new connections
	for {
		n, err := conn.Read(buffer)
		if err != nil {
			printError("readLoop", err)
			continue
		}

		// Message is going to be amount we read from buffer, then we
		// convert byte slice into string each time we receieve the message.

		// We gonna write it into the channel and print it from there,
		// instead of writing plain bytes to this channel
		server.messageChannel <- Message{
			from:    conn.RemoteAddr().String(),
			payload: buffer[:n],
		}

		// Implementation of 'ping-pong' (writing back to peer)
		conn.Write([]byte("Thank you for your message!\n"))

		// Display a map of all connected peers
		server.showAllPeers()
	}
}

func (server *Server) showAllPeers() {
	fmt.Println("Peers connected:")
	for peerID, conn := range server.peersMap {
		fmt.Printf("ID: %s, Address: %s\n", peerID, conn)
	}
}

func main() {
	// Start the server
	server := NewServer(":3000")

	// A new goroutine to loop in the messageChannel and print messages
	go func() {
		for msg := range server.messageChannel {
			fmt.Printf("Message from (%s):%s\n", msg.from, string(msg.payload))
		}
	}()

	log.Fatal(server.Start())
}

// Helpers
func printError(variant string, err error) {
	fmt.Printf("Error in %s: %v\n", variant, err)
}

func generateUniqueID() string {
	randomBytes := make([]byte, 8)
	_, err := rand.Read(randomBytes)
	if err != nil {
		printError("generateUniqueID", err)
		return ""
	}
	uniqueID := hex.EncodeToString(randomBytes)
	return uniqueID
}
