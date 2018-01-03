package storage

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"github.com/boltdb/bolt"
	"github.com/pkg/errors"
	"log"
	"time"
)

const BUCKET = "components"

type BoltStorage struct {
	db *bolt.DB
}

func NewBoltDBStorage() *BoltStorage {
	storage := &BoltStorage{}

	return storage

}

func (s *BoltStorage) Init(dbname string) error {
	if dbname == "" {
		return errors.New("DBName can't be empty")
	}
	dbname = fmt.Sprintf("%s.db", dbname)

	var err error
	config := &bolt.Options{Timeout: 1 * time.Second}
	s.db, err = bolt.Open(dbname, 0600, config)
	if err != nil {
		log.Fatal(err)
		return err
	}

	log.Println(fmt.Sprintf("BoltDB successfully open (%s)", dbname))

	return nil
}

func (s *BoltStorage) Close() {
	log.Println("BoltDB closed")
	s.db.Close()
}

func (s *BoltStorage) Add(item ComItem) (int, error) {
	var id uint64
	err := s.db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte(BUCKET))
		if err != nil {
			return err
		}

		id, _ = b.NextSequence()

		item.ID = int(id)

		buf, err := json.Marshal(item)
		if err != nil {
			return err
		}

		return b.Put(itob(item.ID), buf)
	})

	return item.ID, err
}

func (s *BoltStorage) Get(id int) (ComItem, error) {
	var item ComItem
	err := s.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(BUCKET))
		if b == nil {
			return errors.New("Bucket does not exists")
		}

		buf := b.Get(itob(id))
		if buf == nil {
			return errors.New("Item not found")
		}
		return json.Unmarshal(buf, &item)
	})
	return item, err
}

func (s *BoltStorage) Delete(id int) error {
	return s.db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte(BUCKET))
		if err != nil {
			return err
		}
		return b.Delete(itob(id))
	})
}

func (s *BoltStorage) Update(id int, item ComItem) error {
	return s.db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte(BUCKET))
		if err != nil {
			return err
		}

		buf := b.Get(itob(id))
		if buf == nil {
			return errors.New("Not found")
		}

		jsoned, err := json.Marshal(item)
		if err != nil {
			return errors.New(fmt.Sprintf("Failed to marshal id: %d", id))
		}

		return b.Put(itob(id), jsoned)
	})
}

func itob(v int) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, uint64(v))
	return b
}