package catalog

import (
	"encoding/json"
	"github.com/pkg/errors"
	"net/http"
	"smd_catalog/catalog/storage"
	"log"
)

type Catalog struct {
	storage storage.Storage
}

type ListAnswer struct {
	Items []storage.CatalogItem `json:"items"`
}

type PutAnswer struct {
	Id int `json:"id"`
}

func NewCatalog() (*Catalog, error) {
	s := storage.NewBoltDBStorage()

	err := s.Init("catalog")
	if err != nil {
		return nil, err
	}
	catalog := &Catalog{
		storage: s,
	}
	return catalog, nil
}

func (c *Catalog) Close() {
	if c.storage != nil {
		c.Close()
	}
}

func (c *Catalog) GetServeMux() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "PUT" {
			c.Put(w, r)
		} else {
			c.List(w, r)
		}
	})
	return mux
}

func (c *Catalog) Put(w http.ResponseWriter, req *http.Request) {
	item := storage.CatalogItem{
		Name: req.PostFormValue("name"),
	}

	if item.Name == "" {
		c.Error(w, errors.New("Name field can't be empty"), 400)
		return
	}

	id, err := c.storage.Add(item)

	if err != nil {
		c.Error(w, err, 500)
		return
	}

	answer := PutAnswer{
		Id: id,
	}

	jsoned, err := json.Marshal(answer)
	if err != nil {
		c.Error(w, err, 500)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsoned)
}

func (c *Catalog) Delete(w http.ResponseWriter, req *http.Request) {

}

func (c *Catalog) Assign(w http.ResponseWriter, req *http.Request) {

}

func (c *Catalog) List(w http.ResponseWriter, req *http.Request) {
	w.WriteHeader(http.StatusOK)
	answer := ListAnswer{}
	answer.Items = append(answer.Items, storage.CatalogItem{Name: "item"})

	jsonded, err := json.Marshal(answer)
	if err != nil {
		c.Error(w, err, 500)
		return
	}
	w.Write(jsonded)

}

func (c *Catalog) Error(w http.ResponseWriter, e error, code int) {
	log.Println(e.Error())

	http.Error(w, e.Error(), code)
}
