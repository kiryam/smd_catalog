package catalog

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"smd_catalog/catalog/storage"
	"strconv"
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
		} else if r.Method == "DELETE" {
			c.Delete(w, r)
		} else if r.Method == "POST" {
			c.Post(w, r)
		} else {
			c.List(w, r)
		}
	})
	return mux
}


func (c *Catalog) Post(w http.ResponseWriter, req *http.Request) {
	id := req.PostFormValue("id")
	if id == "" {
		http.Error(w, "`id` can't be empty", 400)
		return
	}

	i, err := strconv.Atoi(id)
	if err != nil {
		http.Error(w, "`id` failed to parse", 400)
		return
	}

	var item storage.CatalogItem
	item, err = c.storage.Get(i)

	if err != nil {
		http.Error(w, "Item not found", 404)
		return
	}

	parent := req.PostFormValue("parent")
	if parent != "" {
		parent_id, err := strconv.Atoi(parent)
		if err != nil {
			http.Error(w, "`parent` failed to parse", 400)
			return
		}
		item.Parent = parent_id
	}

	name := req.PostFormValue("name")
	if name != "" {
		item.Name = name
	}

	var jsoned []byte
	jsoned, err = json.Marshal(item)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "Failed to serialize", 400)
		return
	}

	err = c.storage.Update(i, item)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "Failed to update", 400)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsoned)
}

func (c *Catalog) Put(w http.ResponseWriter, req *http.Request) {
	item := storage.CatalogItem{
		Name: req.PostFormValue("name"),
	}

	if item.Name == "" {
		http.Error(w, "Name field can't be empty", 400)
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
	id := req.FormValue("id")
	if id == "" {
		http.Error(w, "`id` can't be empty", 400)
		return
	}

	i, err := strconv.Atoi(id)
	if err != nil {
		http.Error(w, "`id` failed to parse", 400)
		return
	}

	err = c.storage.Delete(i)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "Failed to delete", 400)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprint(w, "{\"status\": \"ok\"}")
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

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonded)
}

func (c *Catalog) Error(w http.ResponseWriter, e error, code int) {
	log.Println(e.Error())

	http.Error(w, e.Error(), code)
}
