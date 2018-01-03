package storage

type CatalogItem struct {
	ID int `json:"ID"`
	Name string `json:"name"`
	Parent int `json:"parent_id"`
}

type Storage interface {
	Init(dbname string) error
	Close()
	Add(item CatalogItem) (int, error)
	Get(id int) (CatalogItem, error)
	Delete(id int) error
	Update(id int, item CatalogItem) error
}