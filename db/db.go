// 866
// All Rights Reserved

// Package db works with different key-value storages
// Currently it supports the work with lmdb
package db

// DBHandler is an instrument for working with any key-value storage.
// It can write, read and update values.
type DBHandler interface {
	Read(string, []byte) ([]byte, error)
	Write(string, []byte, []byte) error
	Modify(string, []byte, Modifier) error
	Append(string, []byte, Modifier) error
	Close()
}

const (
	// CHAT names the db which stores chat messages.
	CHAT = "chat"
	// PRIVATE names the private db which stores password hashes.
	PRIVATE = "private"
	// PROPOSALS names the database which contains info about proposals.
	PROPOSALS = "proposals"
	// USERS names the users db which contains public info of each user.
	USERS = "users"
	// DYNAMIC names dynamic db which stores dynamic data of proposals.
	DYNAMIC = "dynamic"
)

var (
	// DB is a global variable for handling the database.
	DB DBHandler
	// LastCB is a key pointer to last chat bucket ID.
	LastCB []byte
	// DBList is a list of all available dbs.
	DBList []string
)

func init() {
	// Initialize db stuff
	LastCB = []byte{0}
	DBList = []string{PRIVATE, PROPOSALS, USERS, CHAT, DYNAMIC}
}
