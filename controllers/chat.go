package controllers

import (
	"github.com/astaxie/beego"
	"union/messages"
	"union/db"
	"github.com/satori/go.uuid"
)

// ProposalController handles Proposal requests.
type ChatController struct {
	beego.Controller
}

// Get method handles Chat requests for ChatController.
func (this *ChatController) Get() {
	// Read url param
	idbytes, err := chatUUIDstring(this.GetString("id"))
	if err != nil {
		messages.SendError(this.Ctx.WriteString, err)
		return
	}
	// Read the underlying data
	var data []byte
	data, err = db.DB.Read(idbytes)
	if err != nil {
		messages.SendError(this.Ctx.WriteString, err)
		return
	}
	this.Ctx.WriteString(string(data))
	return
}

// Converts idstr into the uuid byte sequence.
// If idstr is empty it fetches the last bucket
func chatUUIDstring(idstr string) (idbytes []byte, err error) {
	if idstr == "" {
		idbytes, err = db.DB.Read(db.LastCB)
	} else {
		var id uuid.UUID
		id, err = uuid.FromString(idstr)
		idbytes = id.Bytes()
		return
	}
	return
}