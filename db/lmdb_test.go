package db

import (
	"math/rand"
	"os"
	"path"
	"strconv"
	"sync"
	"testing"
)

const dbPath = "./"

func TestLMDB(t *testing.T) {
	lmdb, err := MakeLMDBHandler(dbPath)
	written := []byte("123")
	if err != nil {
		t.Errorf("MakeLMDBHandler error: %v", err)
	}
	err = lmdb.Write(PROPOSALS, []byte("val"), written)
	if err != nil {
		t.Errorf("lmdb.Write error: %v", err)
	}
	var read []byte
	read, err = lmdb.Read(PROPOSALS, []byte("val"))
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
	lmdb, err := MakeLMDBHandler(dbPath)
	written := []byte("1")
	expected := []byte("202323")
	if err != nil {
		t.Errorf("MakeLMDBHandler error: %v", err)
	}
	// Write the data
	err = lmdb.Write(PROPOSALS, []byte("key"), written)
	if err != nil {
		t.Errorf("lmdb.Write error: %v", err)
	}

	// Update the value
	u := Update{expected}
	lmdb.Modify(PROPOSALS, []byte("key"), u)

	// Check the result
	var read []byte
	read, err = lmdb.Read(PROPOSALS, []byte("key"))
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
	lmdb, err := MakeLMDBHandler(dbPath)
	if err != nil {
		b.Errorf("MakeLMDBHandler error: %v", err)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = lmdb.Write(PROPOSALS, []byte(strconv.Itoa(i)), message)
	}
}

func BenchmarkWrite10kbytesEntries(b *testing.B) {
	message := make([]byte, 10*1024)
	for i := range message {
		message[i] = byte(i % 256)
	}
	lmdb, err := MakeLMDBHandler(dbPath)
	if err != nil {
		b.Errorf("MakeLMDBHandler error: %v", err)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = lmdb.Write(PROPOSALS, []byte(strconv.Itoa(i)), message)
	}
}

func BenchmarkRead10kbytesEntries(b *testing.B) {
	message := make([]byte, 10*1024)
	for i := range message {
		message[i] = byte(i % 256)
	}
	lmdb, err := MakeLMDBHandler(dbPath)
	if err != nil {
		b.Errorf("MakeLMDBHandler error: %v", err)
	}
	for i := 0; i < 50; i++ {
		_ = lmdb.Write(PROPOSALS, []byte{byte(i)}, message)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = lmdb.Read(PROPOSALS, []byte{byte(rand.Intn(50))})
	}
}

func BenchmarkRead100bytesEntries(b *testing.B) {
	message := make([]byte, 100)
	for i := range message {
		message[i] = byte(i)
	}
	lmdb, err := MakeLMDBHandler(dbPath)
	if err != nil {
		b.Errorf("MakeLMDBHandler error: %v", err)
	}
	for i := 0; i < 50; i++ {
		_ = lmdb.Write(PROPOSALS, []byte{byte(i)}, message)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = lmdb.Read(PROPOSALS, []byte{byte(rand.Intn(50))})
	}
}

func BenchmarkRead10kbytesEntriesWithParallelReaders(b *testing.B) {
	message := make([]byte, 10*1024)
	for i := range message {
		message[i] = byte(i % 256)
	}
	lmdb, err := MakeLMDBHandler(dbPath)
	if err != nil {
		b.Errorf("MakeLMDBHandler error: %v", err)
	}
	for i := 0; i < 50; i++ {
		_ = lmdb.Write(PROPOSALS, []byte{byte(i)}, message)
	}
	b.ResetTimer()
	wg := sync.WaitGroup{}
	for i := 0; i < b.N; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			lmdb.Read(PROPOSALS, []byte{byte(rand.Intn(50))})
		}()
	}
	wg.Wait()
}

func compareBytes(a, b []byte, t *testing.T) {
	if len(a) != len(b) {
		t.Errorf("Length of the written value(%d) doesn't equal to the length of read value(%d)",
			len(a), len(b))
	}
	for i := range a {
		if a[i] != b[i] {
			t.Errorf("Written slice doesn't match with read slice at position %d", i)
		}
	}
}

func TestMain(m *testing.M) {
	code := m.Run()
	// Remove lmdb files
	_ = os.Remove(path.Join(dbPath, "lock.mdb"))
	_ = os.Remove(path.Join(dbPath, "data.mdb"))
	os.Exit(code)
}
