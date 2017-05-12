// proposals.go introduces proposal functionality
// 866
// All Rights Reserved

// Package controllers contains general handlers of the web server.
package controllers

import (
	"union/db"
	"union/messages"

	"github.com/astaxie/beego"
	"github.com/pkg/errors"
	"github.com/satori/go.uuid"
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
