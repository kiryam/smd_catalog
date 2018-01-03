package comdb

import (
	"net/http"
	"smd_catalog/comdb/storage"
	"strconv"
	"fmt"
	"encoding/json"
	"log"
)

type ComDB struct {
	storage storage.Storage
}


type PutAnswer struct {
	Id int `json:"id"`
}

func NewComDB() (*ComDB, error)  {
	s := storage.NewBoltDBStorage()

	err := s.Init("components")
	if err != nil {
		return nil, err
	}
	catalog := &ComDB{
		storage: s,
	}
	return catalog, nil
}

func (c *ComDB) Close() {
	if c.storage != nil {
		c.Close()
	}
}

func (c *ComDB) GetServeMux() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "PUT" {
			c.Put(w, r)
		} else if r.Method == "DELETE" {
			c.Delete(w, r)
		} else if r.Method == "POST" {
			c.Post(w, r)
		} else {
		}
	})
	return mux
}

func (c *ComDB) Post(w http.ResponseWriter, req *http.Request) {
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

	var item storage.ComItem
	item, err = c.storage.Get(i)

	if err != nil {
		http.Error(w, "Item not found", 404)
		return
	}

	catalog := req.PostFormValue("catalog_id")
	if catalog != "" {
		catalog_id, err := strconv.Atoi(catalog)
		if err != nil {
			http.Error(w, "`parent` failed to parse", 400)
			return
		}
		item.Catalog = catalog_id
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

func (c *ComDB) Put(w http.ResponseWriter, req *http.Request) {
	item := storage.ComItem{
		Name: req.PostFormValue("name"),
	}

	if item.Name == "" {
		http.Error(w, "Name field can't be empty", 400)
		return
	}

	id, err := c.storage.Add(item)

	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	answer := PutAnswer{
		Id: id,
	}

	jsoned, err := json.Marshal(answer)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsoned)
}

func (c *ComDB) Delete(w http.ResponseWriter, req *http.Request) {
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