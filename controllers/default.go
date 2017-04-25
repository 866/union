package controllers

import (
	"github.com/gorilla/websocket"
	"github.com/astaxie/beego"
	"net/http"
	"time"
	"math/rand"
	"fmt"
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
	if _, ok := err.(websocket.HandshakeError); ok {
		http.Error(this.Ctx.ResponseWriter, "Not a websocket handshake", 400)
		return
	} else if err != nil {
		beego.Error("Cannot setup WebSocket connection:", err)
		return
	}
	beego.Info(fmt.Sprintf("Websocket connection: %#v", ws.RemoteAddr().String()))

	var earnings float64
	send := "0 $"

	// Send messages until everything ok.
	for ws.WriteMessage(websocket.TextMessage, []byte(send)) == nil {
		add := rand.Float64() * 0.1
		send = fmt.Sprintf("%s Assets: %0.2f $", time.Now().Format("15:04:05.000"), earnings)
		earnings += add
		time.Sleep(time.Millisecond * 50)
	}

	beego.BeeLogger.Error("Disconnected")
}