package storage

type ComItem struct {
	ID int `json:"ID"`
	Name string `json:"name"`
	Catalog int `json:"catalog_id"`
}

type Storage interface {
	Init(dbname string) error
	Close()
	Add(item ComItem) (int, error)
	Get(id int) (ComItem, error)
	Delete(id int) error
	Update(id int, item ComItem) error
}