// Messages package is used for converting and transmitting the information for server to client and vice versa.
package messages

import (
	"math/rand"
	"encoding/json"
	"github.com/satori/go.uuid"
)

// Server --> Client messages

// Proposal represents proposal information which can be send via websockets.
// Example of JSON representation of proposal's object:
// var proposal = {
//	id: “qwefwrgaweqwfeg” ,
//	author: {
//		name: “Paulo”,
//		id: “qweaer22drg124”,
//		rate: 2.72
//	}
//	type: 1, // 0 - buystop(buy before), 1 - buylimit(buy after), 2 - sellstop(sell after), 3 - selllimit(sell before)
//	price: 1.2345,
//	stoploss: 1.2327, // price - 0.0020 + spread(0.0002)
//	takeprofit: 1.2397, // price + 0.0050 + spread(0.0002)
//	score: 2.72,
//	deadline: 1493062222, // UNIX timestamp in sec
//	pendexp: 2400, // in sec
//	posexp: 7320, // in sec
// }
type Proposal struct {
	AuthorID    string `json:"authorid"`
	ID          string `json:"id"`
	Type        byte `json:"type"`
	Price       float32 `json:"price"`
	StopLoss    float32 `json:"stoploss"`
	TakeProfit  float32 `json:"takeprofit"`
	Score       float32 `json:"score"`
	Deadline    int64 `json:"deadline"`
	PendingExp  int64 `json:"pendexp"`
	PositionExp int64 `json:"posexp"`
}

// ProposalUpdate sends the update message for the proposal with given ID
type ProposalUpdate struct {
	ID    string `json:"id"`
	Score float32 `json:"score"`
}

// ProposalUpdate sends the update message for the proposal with given ID
type ProposalRemove struct {
	ID    string `json:"id"`
	Score float32 `json:"score"`
}

// FillRandom fills the proposal object with some random data.
// ID has 16 runes length.
// Author
// Type is within the range [0, 3]
// Pending expiration is within the [900, 10000] sec range
// Position expiration is within the [3600, 10000] sec range
func (p *Proposal) FillRandom() {
	p.AuthorID = randStringRunes(16)
	p.ID = string(uuid.NewV4().Bytes())
	p.Type = byte(rand.Intn(4) % 256)
	p.Price = rand.Float32()
	p.StopLoss = rand.Float32()
	p.TakeProfit = rand.Float32()
	p.Score = rand.Float32()
	p.Deadline = rand.Int63()
	p.PendingExp = rand.Int63n(10000 - 900 + 1) + 900
	p.PositionExp = rand.Int63n(10000 - 3600 + 1) + 3600
}

// ChatMessage represents chat message which can be send via websockets from server to client.
type ChatMessage struct {
	AuthorID string `json:"authorid"`
	Text     string `json:"text"`
}

// FillRandom fills the ChatMessage object with a random data.
// The sentence has length from 1 to 30
// AuthorID has 16 runes length
func (cm *ChatMessage) FillRandom() {
	cm.AuthorID = randStringRunes(16)
	n := rand.Intn(100) + 1
	cm.Text = randSentence(n)
}


// Message is a single websocket message
// Type value can be:
//	0 - all proposals
//	1 - all chat messages
//	2 - add a proposal
//	3 - add a message
//	4 - update a proposal
type Message struct {
	Type byte `json:"type"`
	Data interface{} `json:"data"`
}

// WSData stores multiple Messages. Can be Marshalled to json
type WSData struct {
	Data []Message
}

// FillRandom fills up the WSData object with random number of elements that lies in the range [min, max].
func (wsd *WSData) FillRandom(chatsmin, chatsmax, propsmin, propsmax int) {
	// Get the number of props/messages
	nchats := rand.Intn(chatsmax - chatsmin + 1) + chatsmin
	nprops := rand.Intn(propsmax - propsmin + 1) + propsmin
	chats := make([]ChatMessage, nchats)
	props := make([]Proposal, nprops)
	// Fill chats by the random data
	for i := range chats {
		chats[i].FillRandom()
	}
	// Fill props by the random data
	for i := range props {
		props[i].FillRandom()
	}
	// Make data slice and fill it with randomly generated data
	wsd.Data = make([]Message, 2)
	var place int
	if rand.Float32() < .5 {
		place = 1
	}
	wsd.Data[1 - place] = Message{1, chats}
	wsd.Data[place] = Message{0, props}
}

// Converts WSData to JSON object
func (wsd *WSData) Jsonify() ([]byte, error) {
	return json.Marshal(wsd)
}

// Possible runes are listed here.
var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ,.")

// Generates the random string of size n.
func randStringRunes(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

// Generates fake message of length n.
func randSentence(n int) string {
	w := make([]rune, n)
	for i := range w {
		if rand.Float64() < 0.2 {
			w[i] = rune(' ')
		} else {
			w[i] = letterRunes[rand.Intn(len(letterRunes))]
		}
	}
	return string(w)
}