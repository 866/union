// tcpengine.go implements tcp client that communicates with trigger, websockets and provides gameplay
// 866
// All Rights Reserved

package messages

import (
	"github.com/gorilla/websocket"
	"net"
	"union/db"
)

// TCPWSEngine is a gameplay engine that works with TCP trigger server.
// It provides chatting ability, interaction with the database, websocket connections handling.
type TCPWSEngine struct {
	ws      []*websocket.Conn
	trigger net.Conn
	db      db.DBHandler
}

// Close finishes all open objects.
func (e *TCPWSEngine) Close() {
	for _, conn := range e.ws {
		conn.Close()
	}
	e.db.Close()
	e.trigger.Close()
}

// internalLoop accepts new websocket connections and handles basic operations with them
func (e *TCPWSEngine) internalLoop() {
	// TODO
}
