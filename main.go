package main

import (
	_ "union/routers"
	"github.com/astaxie/beego"
	"union/db"
	"union/messages"
	"github.com/satori/go.uuid"
	"encoding/json"
)


// initialize lmdb database
func initLMDB() {
	lmdb, err := db.MakeLMDBHandler("./")
	if err != nil {
		panic(err)
	}
	// Add random proposal to the database
	id := uuid.NewV4()
	beego.Info("Prop ID: ", id.String())
	prop := messages.Proposal{}
	prop.FillRandom()
	data, _ := json.Marshal(prop)
	lmdb.Write(db.PROPOSALS, id.Bytes(), data)
	// Add random chat message to the database
	id = uuid.NewV4()
	beego.Info("Chat Message ID: ", id.String())
	chat := messages.ChatMessage{}
	chat.FillRandom()
	data, _ = json.Marshal(chat)
	lmdb.Write(db.CHAT, id.Bytes(), data)
	// Global database
	db.DB = lmdb
}


func main() {
	// init the database
	initLMDB()
	beego.Info("DB is initialized.")

	beego.Run()
}

