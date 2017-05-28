// tcpengine.go implements tcp client that communicates with trigger, websockets and provides gameplay
// 866
// All Rights Reserved

package messages

import (
	"github.com/gorilla/websocket"
	"net"
	"union/db"
	"github.com/satori/go.uuid"
)

// CHAN_BUFFER is a length of waiting queue for buffered channels.
const CHAN_BUFFER = 10

// IDKey is a type for ID. ID should be a unique identifier for each object.
type IDKey uuid.UUID

// Subscriber is a user that can connect to the websocket communication.
type Subscriber struct {
	ID   IDKey
	Conn *websocket.Conn // Only for WebSocket users; otherwise nil.
}

// Subscribers represents a dictionary of connected user to the server.
type Subscribers map[IDKey]*websocket.Conn

// TCPWSEngine is a gameplay engine that works with TCP trigger server.
// It provides chatting ability, interaction with the database, websocket connections handling.
type TCPWSEngine struct {
	clients Subscribers
	trigger net.Conn
	// Database handler.
	db db.DBHandler
	// Channel for new join users.
	subscribe chan Subscriber
	// Channel for exit users.
	unsubscribe chan IDKey
	// Send events here to publish them.
	recieve chan WSData
}

// Close finishes all open objects.
func (e *TCPWSEngine) Close() {
	for _, conn := range e.clients {
		conn.Close()
	}
	e.db.Close()
	e.trigger.Close()
}

// Init initializes engine object.
// dbh is as DBHandler object. trigger is a connection to TCP trigger server.
func (e *TCPWSEngine) Init(dbh db.DBHandler, trigger net.Conn) {
	e.clients = make(Subscribers)
	e.trigger = trigger
	e.db = dbh
	e.subscribe = make(chan Subscriber, CHAN_BUFFER)
	e.unsubscribe = make(chan IDKey, CHAN_BUFFER)
	go e.internalLoop()
}

// Leave unsubscribes user from the websocket communication
func (e *TCPWSEngine) Leave(user IDKey) {
	e.unsubscribe <- user
}

// Join registers the user in system.
func (e *TCPWSEngine) Join(user IDKey, conn *websocket.Conn) {
	e.subscribe <- Subscriber{user, conn}
}

// internalLoop accepts new websocket connections and handles basic operations with them
func (e *TCPWSEngine) internalLoop() {
	// TODO
}
