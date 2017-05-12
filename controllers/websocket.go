// websocket.go introduces websocket handling functionality
// 866
// All Rights Reserved

package controllers

import (
	"fmt"
	"math/rand"
	"net/http"
	"time"

	"union/messages"

	"github.com/astaxie/beego"
	"github.com/gorilla/websocket"
)

type MainController struct {
	beego.Controller
}

// WebSocketController handles WebSocket requests.
type WebSocketController struct {
	beego.Controller
}

func (c *MainController) Get() {
	c.Data["Website"] = "union.org"
	c.Data["Email"] = "comrazvictor@gmail.com"
	c.TplName = "index.tpl"
}

// Get method handles WebSocket requests for WebSocketController.
func (this *WebSocketController) Get() {
	// Make the websocket upgrader with compression
	u := websocket.Upgrader{ReadBufferSize: 1024, WriteBufferSize: 1024, EnableCompression: true}
	u.Error = func(w http.ResponseWriter, r *http.Request, status int, reason error) {
		// don't return errors to maintain backwards compatibility
	}
	u.CheckOrigin = func(r *http.Request) bool {
		// allow all connections by default
		return true
	}
	// Upgrade from http request to WebSocket.
	ws, err := u.Upgrade(this.Ctx.ResponseWriter, this.Ctx.Request, nil)
	defer ws.Close()
	if _, ok := err.(websocket.HandshakeError); ok {
		http.Error(this.Ctx.ResponseWriter, "Not a websocket handshake", 400)
		return
	} else if err != nil {
		beego.Error("Cannot setup WebSocket connection:", err)
		return
	}
	beego.Info(fmt.Sprintf("Websocket connection: %s", ws.RemoteAddr().String()))

	// Generating random data
	wsdata := messages.WSData{}
	wsdata.FillRandom(0, 40, 0, 15)
	send, err := wsdata.Jsonify()
	if err != nil {
		beego.BeeLogger.Error("wsdata.Jsonify error: %#v", err)
		return
	}
	// Send messages until everything ok.
	for ws.WriteMessage(websocket.TextMessage, send) == nil {
		time.Sleep(time.Millisecond * time.Duration(rand.Intn(3000)+900))
		wsdata.FillRandom(0, 40, 0, 15)
		send, err = wsdata.Jsonify()
		if err != nil {
			beego.BeeLogger.Error("wsdata.Jsonify error: %#v", err)
			return
		}
	}
	beego.BeeLogger.Info("Disconnected: %s", ws.RemoteAddr().String())
}
