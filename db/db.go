package db

// DBHandler is an instrument for working with any key-value storage.
// It can write, read and update values.
type DBHandler interface {
	Read([]byte) ([]byte, error)
	Write([]byte, []byte) error
	Modify([]byte, Modifier) error
	Close()
}

// DB is a global variable for handling the database
var DB DBHandler