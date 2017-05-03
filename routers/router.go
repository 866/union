package routers

import (
	"union/controllers"
	"github.com/astaxie/beego"
	"union/db"
)

// initialize lmdb database
func initLMDB() db.DBHandler {
	lmdb, err := db.MakeLMDBHandler("/tmp", "mydb")
	if err != nil {
		panic(err)
	}
	return lmdb
}

func init() {
	// init the database
	dbh := initLMDB()

	beego.Router("/", &controllers.MainController{})
	// Info controllers
	beego.Router("/proposal", &controllers.ProposalController{beego.Controller{}, dbh})
	// WebSocket connection
	beego.Router("/ws", &controllers.WebSocketController{})
}
