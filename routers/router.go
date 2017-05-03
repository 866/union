package routers

import (
	"union/controllers"
	"github.com/astaxie/beego"
	"union/db"
	"github.com/satori/go.uuid"
	"fmt"
	"union/messages"
	"encoding/json"
)

// initialize lmdb database
func initLMDB() {
	lmdb, err := db.MakeLMDBHandler("/tmp", "mydb")
	if err != nil {
		panic(err)
	}
	id := uuid.NewV4()
	fmt.Println("Prop ID: ", id.String())
	prop := messages.ChatMessage{}
	prop.FillRandom()
	data, _ := json.Marshal(prop)
	lmdb.Write(id.Bytes(), data)
	db.DB = lmdb
}

func init() {
	// init the database
	initLMDB()

	beego.Router("/", &controllers.MainController{})
	// Info controllers
	beego.Router("/proposal", &controllers.ProposalController{})
	// WebSocket connection
	beego.Router("/ws", &controllers.WebSocketController{})
}
