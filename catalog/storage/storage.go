package storage

type CatalogItem struct {
	ID int `json:"ID"`
	Name string `json:"name"`
}

type Storage interface {
	Init(dbname string) error
	Close()
	Add(item CatalogItem) (int, error)
	Delete(id int)
	SetParent(id int, parent int)
}
