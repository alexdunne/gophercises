package db

import (
	"encoding/binary"
	"path/filepath"
	"time"

	"github.com/boltdb/bolt"
	"github.com/mitchellh/go-homedir"
)

var taskBucket = []byte("tasks")

type handler struct {
	dbPath string
}

type HandlerOption func(h *handler)

func WithPath(dbPath string) HandlerOption {
	return func(h *handler) {
		h.dbPath = dbPath
	}
}

func buildHandler(opts ...HandlerOption) handler {
	h := handler{
		dbPath: defaultDatabasePath(),
	}

	for _, opt := range opts {
		opt(&h)
	}

	return h
}

// ConnectDB creates a new connection to the database
func ConnectDB(opts ...HandlerOption) (*bolt.DB, error) {
	h := buildHandler(opts...)

	db, err := bolt.Open(h.dbPath, 0600, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		return nil, err
	}

	err = db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists(taskBucket)
		return err
	})
	if err != nil {
		return nil, err
	}

	return db, nil
}

func defaultDatabasePath() string {
	home, _ := homedir.Dir()
	return filepath.Join(home, "tasks.db")
}

type Store struct {
	db *bolt.DB
}

type Task struct {
	Key   int
	Value string
}

func NewStore(db *bolt.DB) *Store {
	return &Store{db: db}
}

func (s *Store) CreateTask(task string) (int, error) {
	var id int
	err := s.db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(taskBucket)

		id64, _ := bucket.NextSequence()
		id = int(id64)
		key := itob(id)

		bucket.Put(key, []byte(task))

		return nil
	})
	if err != nil {
		return -1, err
	}

	return id, nil
}

func (s *Store) FindAllTasks() ([]Task, error) {
	var tasks []Task

	err := s.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(taskBucket)

		cursor := b.Cursor()

		for k, v := cursor.First(); k != nil; k, v = cursor.Next() {
			tasks = append(tasks, Task{
				Key:   btoi(k),
				Value: string(v),
			})
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return tasks, nil
}

// itob converts an int to a byte slice
// Big/Little endian determines the order of the bytes
// Big endian means the most significant bits are stores first
// Little endian is the reverse
func itob(v int) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, uint64(v))
	return b
}

func btoi(b []byte) int {
	return int(binary.BigEndian.Uint64(b))
}
