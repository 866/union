package routers

import (
	"union/controllers"
	"github.com/astaxie/beego"
)

func init() {
	beego.Router("/", &controllers.MainController{})
	// Info controllers
	beego.Router("/proposal", &controllers.ProposalController{})
	beego.Router("/chat", &controllers.ProposalController{})
	// WebSocket connection
	beego.Router("/ws", &controllers.WebSocketController{})
}
