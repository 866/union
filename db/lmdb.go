// lmdb.go simplifies the work with lmdb
// 866
// All Rights Reserved

package db

import (
	"runtime"

	"github.com/bmatsuo/lmdb-go/lmdb"
	"github.com/satori/go.uuid"
)


const (
	// FULLERROR appears when the bucket of values is full
	FULLERROR = "FULL"
	// MAXCHATBUCKETSIZE is a max possible size of bucket
	MAXCHATBUCKETSIZE = 30
)

// Modifier is a functor interface that changes the content.
// It is used by Modify function of DBHandler
type Modifier interface {
	Apply([]byte) ([]byte, error)
}

// lmdbop is a basic lmdb operation
type lmdbop struct {
	op  lmdb.TxnOp
	res chan<- error
}

// LMDB is a thread-safe wrapper over lmdb environment.
type LMDB struct {
	env    *lmdb.Env
	dbs    map[string]lmdb.DBI
	worker chan *lmdbop
	update func(lmdb.TxnOp) error
}

// writer is a background goroutine that accepts requests and updates the datavase
func (dbh *LMDB) writer() {
	runtime.LockOSThread()
	defer runtime.LockOSThread()
	// Endless loop which receives tasks
	for {
		select {
		case work, open := <-dbh.worker:
			// Check for closed channel
			if !open {
				return
			}
			work.res <- dbh.env.UpdateLocked(work.op)
		}
	}
}

// Write writes the content val at address key.
// If DBHandler is not initialized the function panics.
func (dbh *LMDB) Write(db string, key, val []byte) error {
	return dbh.update(func(txn *lmdb.Txn) (err error) {
		return txn.Put(dbh.dbs[db], key, val, 0)
	})
}

// Read reads the value at address key.
// If DBHandler is not initialized the function panics.
func (dbh *LMDB) Read(db string, key []byte) (v []byte, err error) {
	err = dbh.env.View(func(txn *lmdb.Txn) (err error) {
		v, err = txn.Get(dbh.dbs[db], key)
		return err
	})
	return
}

// Modify extracts content which corresponds to the key then it modifies it by means of functor m.
// If DBHandler is not initialized the function panics.
func (dbh *LMDB) Modify(db string, key []byte, m Modifier) error {
	return dbh.update(func(txn *lmdb.Txn) (err error) {
		var v []byte
		// Read
		v, err = txn.Get(dbh.dbs[db], key)
		if err != nil {
			return
		}
		// Change the content
		v, err = m.Apply(v)
		if err != nil {
			return
		}
		// Write the update
		err = txn.Put(dbh.dbs[db], key, v, 0)
		return
	})
}

// Append extracts content which corresponds to the key then it modifies it by means of functor m.
// If functor returns FULLERROR it creates a new key-value cell and write the data there.
func (dbh *LMDB) Append(db string, pointer []byte, m Modifier) (error) {
	return dbh.update(func(txn *lmdb.Txn) (err error) {
		var v []byte
		// Read the tail key
		key, err := dbh.Read(CHAT, pointer)
		if err != nil {
			return err
		}
		// Read
		v, err = txn.Get(dbh.dbs[db], key)
		if err != nil {
			return
		}
		// Change the content
		v, err = m.Apply(v)
		if err == nil {
			err = txn.Put(dbh.dbs[db], key, v, 0)
			return
		}
		// The list is full or other error has been occurred
		if err.Error() == FULLERROR {
			// Generate the key
			// TODO: change to generate key
			key = uuid.NewV4().Bytes()
			// Create new entry
			v, err = m.Apply(nil)
			if err != nil {
				return
			}
			err = txn.Put(dbh.dbs[db], key, v, 0)
			if err != nil {
				return
			}
			// Create a new pointer to the last bucket
			err = txn.Put(dbh.dbs[db], pointer, key, 0)
		}
		// Return the error
		return
	})
}


// Close finishes the work with an environment. Should be called when the work is finished.
func (dbh *LMDB) Close() {
	close(dbh.worker)
	dbh.env.Close()
}

// MakeLMDBHandler returns LMDB object with opened database db at the specified path.
// If the db doesn't exit, the function creates it.
func MakeLMDBHandler(path string) (l *LMDB, err error) {
	var env *lmdb.Env
	// Open the environment
	env, err = lmdb.NewEnv()
	if err != nil {
		return
	}
	// Set max db size
	err = env.SetMapSize(100 * 1024 * 1024) // 100MB
	if err != nil {
		return
	}
	// Change this value if you want to have more dbs
	err = env.SetMaxDBs(10)
	if err != nil {
		return
	}
	// Open the db at the given path
	err = env.Open(path, 0, 0664)
	if err != nil {
		return
	}
	// Open or create the databases
	dbs := make(map[string]lmdb.DBI)
	err = env.Update(func(txn *lmdb.Txn) (err error) {
		for _, db := range DBList {
			dbs[db], err = txn.OpenDBI(db, 0)
			if err != nil {
				dbs[db], err = txn.OpenDBI(db, lmdb.Create)
				if err != nil {
					return
				}
			}
		}
		return
	})
	if err != nil {
		return
	}
	// Define update function
	update := func(op lmdb.TxnOp) error {
		res := make(chan error)
		l.worker <- &lmdbop{op, res}
		return <-res
	}
	// Create the lmdb object
	l = &LMDB{env, dbs, make(chan *lmdbop), update}
	go l.writer()
	return
}
