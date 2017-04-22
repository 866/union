package routers

import (
	"union/controllers"
	"github.com/astaxie/beego"
)

func init() {
	beego.Router("/", &controllers.MainController{})
	// WebSocket.
	beego.Router("/ws", &controllers.WebSocketController{})
}
