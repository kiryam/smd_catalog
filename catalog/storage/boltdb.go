package storage

import (
	"fmt"
	"github.com/boltdb/bolt"
	"github.com/pkg/errors"
	"log"
)

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

	var err error
	s.db, err = bolt.Open(fmt.Sprintf("%s.db", dbname), 0600, nil)
	if err != nil {
		log.Fatal(err)
		return err
	}

	return nil
}

func (s *BoltStorage) Close() {
	s.db.Close()
}

func (s *BoltStorage) Add(item CatalogItem) (int, error) {
	return 1, nil
}

func (s *BoltStorage) Delete(id int) {

}

func (s *BoltStorage) SetParent(id int, parent int) {

}
