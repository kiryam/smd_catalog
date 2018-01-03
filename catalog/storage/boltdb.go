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

const BUCKET = "catalog"

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

func (s *BoltStorage) Add(item CatalogItem) (int, error) {
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

func (s *BoltStorage) Delete(id int) {

}

func (s *BoltStorage) SetParent(id int, parent int) {

}

func itob(v int) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, uint64(v))
	return b
}