package db

import (
	"encoding/binary"
	"encoding/json"
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
	ID          int
	Description string
	Completed   bool
}

func NewStore(db *bolt.DB) *Store {
	return &Store{db: db}
}

func (s *Store) FindAllActiveTasks() ([]Task, error) {
	var tasks []Task

	err := forEachTask(s.db, func(k []byte, t Task) {
		if !t.Completed {
			tasks = append(tasks, Task{
				ID:          btoi(k),
				Description: t.Description,
				Completed:   t.Completed,
			})
		}
	})

	if err != nil {
		return nil, err
	}

	return tasks, nil
}

func (s *Store) FindAllCompletedTasks() ([]Task, error) {
	var tasks []Task

	err := forEachTask(s.db, func(k []byte, t Task) {
		if t.Completed {
			tasks = append(tasks, Task{
				ID:          btoi(k),
				Description: t.Description,
				Completed:   t.Completed,
			})
		}
	})

	if err != nil {
		return nil, err
	}

	return tasks, nil
}

func forEachTask(db *bolt.DB, fn func(k []byte, t Task)) error {
	return db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(taskBucket)

		cursor := b.Cursor()

		for k, v := cursor.First(); k != nil; k, v = cursor.Next() {
			var task Task
			err := json.Unmarshal(v, &task)
			if err != nil {
				return err
			}

			fn(k, task)
		}

		return nil
	})
}

func (s *Store) CreateTask(t *Task) (*Task, error) {
	err := s.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(taskBucket)

		id, _ := b.NextSequence()
		t.ID = int(id)

		buf, err := json.Marshal(t)
		if err != nil {
			return err
		}

		return b.Put(itob(t.ID), buf)
	})
	if err != nil {
		return nil, err
	}

	return t, nil
}

func (s *Store) MarkTaskAsComplete(id int) error {
	return s.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(taskBucket)

		key := itob(id)

		var task Task
		t := b.Get(key)
		err := json.Unmarshal(t, &task)
		if err != nil {
			return err
		}

		task.Completed = true

		buf, err := json.Marshal(task)
		if err != nil {
			return err
		}

		return b.Put(key, buf)
	})
}

func (s *Store) DeleteTask(key int) error {
	return s.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(taskBucket)

		return b.Delete(itob(key))
	})

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
