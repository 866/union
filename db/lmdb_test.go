package db

import (
	"testing"
	"strconv"
	"math/rand"
)

func TestLMDB(t *testing.T) {
	lmdb, err := MakeLMDBHandler("/tmp", "lmdbtest")
	written := []byte("123")
	if err != nil {
		t.Errorf("MakeLMDBHandler error: %v", err)
	}
	err = lmdb.Write([]byte("val"), written)
	if err != nil {
		t.Errorf("lmdb.Write error: %v", err)
	}
	var read []byte
	read, err = lmdb.Read([]byte("val"))
	if err != nil {
		t.Errorf("lmdb.Write error: %v", err)
	}
	compareBytes(written, read, t)
	lmdb.Close()
}

// Update is a fake Modifier interface for testing.
type Update struct {
	Expected []byte
}

// Apply is a faked Modifier method for testing. It changes the data to Expected
func (u Update) Apply(input []byte) ([]byte, error) {
	return u.Expected, nil
}

func TestUpdate(t *testing.T) {
	lmdb, err := MakeLMDBHandler("/tmp", "lmdbtest")
	written := []byte("1")
	expected := []byte("202323")
	if err != nil {
		t.Errorf("MakeLMDBHandler error: %v", err)
	}
	// Write the data
	err = lmdb.Write([]byte("key"), written)
	if err != nil {
		t.Errorf("lmdb.Write error: %v", err)
	}

	// Update the value
	u := Update{expected}
	lmdb.Modify([]byte("key"), u)

	// Check the result
	var read []byte
	read, err = lmdb.Read([]byte("key"))
	if err != nil {
		t.Errorf("lmdb.Write error: %v", err)
	}
	compareBytes(read, expected, t)
	lmdb.Close()
}

func BenchmarkWrite100bytesEntries(b *testing.B) {
	message := make([]byte, 100)
	for i := range message {
		message[i] = byte(i)
	}
	lmdb, err := MakeLMDBHandler("/tmp", "lmdbtest")
	if err != nil {
		b.Errorf("MakeLMDBHandler error: %v", err)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = lmdb.Write([]byte(strconv.Itoa(i)), message)
	}
}

func BenchmarkRead100bytesEntries(b *testing.B) {
	message := make([]byte, 100)
	for i := range message {
		message[i] = byte(i)
	}
	lmdb, err := MakeLMDBHandler("/tmp", "lmdbtest")
	if err != nil {
		b.Errorf("MakeLMDBHandler error: %v", err)
	}
	for i := 0; i < 50; i++ {
		_ = lmdb.Write([]byte{byte(i)}, message)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = lmdb.Read([]byte{byte(rand.Intn(50))})
	}
}


func compareBytes(a, b []byte, t *testing.T) {
	if len(a) != len(b) {
		t.Errorf("Length of the written value(%d) doesn't equal to the length of read value(%d)",
			len(a), len(b))
	}
	for i := range(a) {
		if a[i] != b[i] {
			t.Errorf("Written slice doesn't match with read slice at position %d", i)
		}
	}
}