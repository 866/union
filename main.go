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

	chatb := messages.ChatBucket{}
	chatb.FillRandom(15)
	data, _ = json.Marshal(chatb)
	lmdb.Write(db.CHAT, id.Bytes(), data)

	prev := id.String()
	id = uuid.NewV4()
	chatb.FillRandom(20)
	chatb.Previous = &prev
	data, _ = json.Marshal(chatb)
	lmdb.Write(db.CHAT, id.Bytes(), data)

	beego.Info("Chat Bucket ID: ", id.String())
	// Global database
	db.DB = lmdb
}


func main() {
	// init the database
	initLMDB()
	beego.Info("DB is initialized.")
	beego.Run()
}

