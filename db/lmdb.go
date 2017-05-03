package db

import (
	"github.com/bmatsuo/lmdb-go/lmdb"
	"runtime"
)

// Modifier is a functor interface that changes the content.
// It is used by Modify function of DBHandler
type Modifier interface {
	Apply(input []byte) ([]byte, error)
}

// lmdbop is a basic lmdb operation
type lmdbop struct {
	op  lmdb.TxnOp
	res chan<- error
}

// LMDB is a thread-safe wrapper over lmdb environment.
type LMDB struct {
	env lmdb.Env
	db     lmdb.DBI
	worker chan lmdbop
	update func(lmdb.TxnOp) error
}

// writer is a background goroutine that accepts requests and updates the datavase
func (dbh *LMDB) writer() {
	runtime.LockOSThread()
	defer runtime.LockOSThread()

	for {
		select {
		case work, open := <-dbh.worker:
			if !open {
				return
			}
			work.res <- dbh.env.UpdateLocked(work.op)
		}
	}
}

// Write writes the content val at address key.
// If DBHandler is not initialized the function panics.
func (dbh *LMDB) Write(key, val []byte) error {
	return dbh.update(func(txn *lmdb.Txn) (err error) {
		return txn.Put(dbh.db, key, val, 0)
	})
}

// Read reads the value at address key.
// If DBHandler is not initialized the function panics.
func (dbh *LMDB) Read(key []byte) (v []byte, err error) {
	err = dbh.env.View(func(txn *lmdb.Txn) (err error) {
		v, err = txn.Get(dbh.db, key)
		return err
	})
	return
}

// Modify extracts content which corresponds to the key then it modifies it by means of functor m.
// If DBHandler is not initialized the function panics.
func (dbh *LMDB) Modify(key []byte, m Modifier) error {
	return dbh.update(func(txn *lmdb.Txn) (err error) {
		var v []byte
		// Read
		v, err = txn.Get(dbh.db, key)
		if err != nil {
			return
		}
		// Change the content
		v, err = m.Apply(v)
		if err != nil {
			return
		}
		// Write the update
		err = txn.Put(dbh.db, key, v, 0)
		return
	})
}

// Close finishes the work with an environment. Should be called when the work is finished.
func (dbh *LMDB) Close() {
	dbh.env.Close()
}

// MakeLMDBHandler returns LMDB object with opened database db at the specified path.
// If the db doesn't exit, the function creates it.
func MakeLMDBHandler(path, db string) (l *LMDB, err error) {
	var env lmdb.Env
	// Open the environment
	env, err = lmdb.NewEnv()
	if err != nil {
		return
	}
	// Set max db size
	err = env.SetMapSize(1024 * 1024 * 1024) // 1GB
	if err != nil {
		return
	}
	// Change this value if you want to have more dbs
	err = env.SetMaxDBs(1)
	if err != nil {
		return
	}
	// Open the db at the given path
	err = env.Open(path, 0, 0664)
	if err != nil {
		return
	}
	// Open or create the database
	var dbi lmdb.DBI
	err = env.Update(func(txn *lmdb.Txn) (err error) {
		dbi, err = txn.OpenDBI(db, 0)
		if err != nil {
			dbi, err = txn.OpenDBI(db, lmdb.Create)
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
	l = &LMDB{env, dbi, make(chan lmdbop), update}
	return
}