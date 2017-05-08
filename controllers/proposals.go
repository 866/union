package controllers

import (
	"github.com/astaxie/beego"
	"github.com/satori/go.uuid"

	"union/db"
	"union/messages"
	"github.com/pkg/errors"
)

// ProposalController handles Proposal requests.
type ProposalController struct {
	beego.Controller
}

// Get method handles Proposal requests for ProposalController.
func (this *ProposalController) Get() {
	// Check the id for availability
	idstr := this.GetString("id")
	if idstr == "" {
		messages.SendError(this.Ctx.WriteString, errors.New("id parameter is missing"))
		return
	}
	// Convert id string into the uuid
	id, err := uuid.FromString(idstr)
	if err != nil {
		messages.SendError(this.Ctx.WriteString, err)
		return
	}
	// Read the data
	var data []byte
	data, err = db.DB.Read(db.PROPOSALS, id.Bytes())
	if err != nil {
		messages.SendError(this.Ctx.WriteString, err)
		return
	}
	// Write the response
	this.Ctx.WriteString(string(data))
	return
}