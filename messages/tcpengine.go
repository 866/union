// tcpengine.go implements tcp client that communicates with trigger, websockets and provides gameplay
// 866
// All Rights Reserved

package messages

import (
	"github.com/gorilla/websocket"
	"net"
	"union/db"
	"github.com/satori/go.uuid"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/mitchellh/mapstructure"
	"time"
)

// CHAN_BUFFER is a length of waiting queue for buffered channels.
const CHAN_BUFFER = 10

// IDKey is a type for ID. ID should be a unique identifier for each object.
type IDKey uuid.UUID

// IDKeyFromBytes converts bytes slice to IDKey.
func IDKeyFromBytes(bytes []byte) (IDKey, error) {
	uid, err := uuid.FromBytes(bytes)
	return IDKey(uid), err
}

// String method for printing.
func (id IDKey) String() string {
	return uuid.UUID(id).String()
}

// Subscriber is a user that can connect to the websocket communication.
type Subscriber struct {
	ID   IDKey
	Conn *websocket.Conn // Only for WebSocket users; otherwise nil.
}

// WSEvent is an websocket data income from user ID.
type WSEvent struct {
	ID   IDKey
	Data WSData
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
	join chan Subscriber
	// Channel for exit users.
	leave chan IDKey
	// Send events here to publish them.
	receive chan WSEvent
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
	e.join = make(chan Subscriber, CHAN_BUFFER)
	e.leave = make(chan IDKey, CHAN_BUFFER)
	e.receive = make(chan WSEvent, CHAN_BUFFER)
	go e.internalLoop()
}

// Leave unsubscribes user from the websocket communication
func (e *TCPWSEngine) Leave(user IDKey) {
	e.leave <- user
}

// Join registers the user in system.
func (e *TCPWSEngine) Join(user IDKey, conn *websocket.Conn) {
	e.join <- Subscriber{user, conn}
}

// internalLoop accepts new websocket connections and handles basic operations with them
func (e *TCPWSEngine) internalLoop() {
	for {
		select {
		case sub := <-e.join:
			err := e.subscribe(sub)
			if err != nil {
				beego.Error("Error while connecting: ", err)
			} else {
				beego.Info("Client %s has been connected.", sub.ID.String())
			}
		case event := <-e.receive:
			err := e.handleMessage(event)
			if err != nil {
				beego.Error(err)

			}
		case id := <-e.leave:
			err := e.unsubscribe(id)
			if err != nil {
				beego.Error("Error while disconneting: ", err)
			} else {
				beego.Info("Client %s has been disconnected.", id)
			}
		}
	}
}

// Handles incoming WSData
func (e *TCPWSEngine) handleMessage(event WSEvent) (err error) {
	for _, msg := range event.Data.Messages {
		switch msg.Type {
		// TODO: Add all possible events
		case CADDCHAT:
			chatmsg := &ChatMessage{}
			err = mapstructure.Decode(msg.Data, chatmsg)
			if err != nil {
				return
			}
			idstr := event.ID.String()
			chatmsg.AuthorID = &idstr
			time := time.Now().Unix()
			chatmsg.Time = &time
			err = e.spreadChatMsg(event.ID, chatmsg)
			if err != nil {
				beego.Error("%s: %s", event.ID.String(), err.Error())
			}
		default:
			err = fmt.Errorf("Event type %d is not known", int(msg.Type))
		}
	}
	return
}

// Broadcasts chat message and writes it into the database.
// Runs paralelly.
func (e *TCPWSEngine) spreadChatMsg(id IDKey, cm *ChatMessage) error {
	// Add the message into the database
	lastKey, err := e.db.Read(db.CHAT, db.LastCB)
	if err != nil {
		return err
	}
	var lastid IDKey
	lastid, err = IDKeyFromBytes(lastKey)
	if err != nil {
		return err
	}
	mod := &AddChatMessage{lastid, cm}
	err = e.db.Append(db.CHAT, db.LastCB, mod)
	if err != nil {
		return err
	}
	// Convert to the json byte sequence
	var data []byte
	// Broadcast to all users
	for id, conn := range e.clients {
		// Immediately send event to WebSocket users.
		if conn != nil {
			if conn.WriteMessage(websocket.TextMessage, data) != nil {
				// User disconnected.
				e.leave <- id
			}
		}
	}
	return nil
}

// Unsubscibe user from the connection list.
func (e *TCPWSEngine) unsubscribe(id IDKey) (err error) {
	if _, contains := e.clients[id]; contains {
		delete(e.clients, id)
	} else {
		err = fmt.Errorf("Client %s doesn't exist in real-time list of users", id.String())
	}
	return
}

// Subscribe the new client and add new connection to the map.
func (e *TCPWSEngine) subscribe(s Subscriber) (err error) {
	if _, contains := e.clients[s.ID]; !contains {
		e.clients[s.ID] = s.Conn
	} else {
		err = fmt.Errorf("Client %s already exists in real-time list of users", s.ID.String())
	}
	return
}

