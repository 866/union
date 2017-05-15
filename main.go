// main.go is the server's entry point
// 866
// All Rights Reserved

package main

import (
	"encoding/json"

	"union/db"
	"union/messages"
	_ "union/routers"

	"github.com/astaxie/beego"
	"github.com/satori/go.uuid"
)

// initialize lmdb database
// this is a temporary function provided for testing
// it will be changed in the future
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
	// Initialize the database
	initLMDB()
	beego.Info("DB is initialized.")
	// Run the beego
	beego.Run()
}
