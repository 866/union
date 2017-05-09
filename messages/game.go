package messages

import (
	"github.com/satori/go.uuid"
	"math/rand"
)

type Event struct {
	// Time when the event occured
	Time int64 `json:"time"`
	// Value when the event occured
	Value float32 `json:"value"`
	// State after event
	State byte `json:"state"`
}

// FillRandom fills event with random numbers
func (e *Event) FillRandom() {
	e.Time = rand.Int63()
	e.Value = rand.Float32()
	e.State = byte(rand.Int() % 256)
}

// Engine represents game engine. It interacts with key-value database and trigger server.
// It adds and update proposals dynamically. Enigne contains all check mechanisms in order to identify
// wrong updates.
type Engine interface {
	AddProposal(p Proposal) (uuid.UUID, error)
	VoteProposal(propID, userID uuid.UUID) (error)
	UpgradeProposal(uuid.UUID, Event)
}