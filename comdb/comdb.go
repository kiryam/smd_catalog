package comdb

import (
	"net/http"
	"github.com/kiryam/smd_catalog/comdb/storage"
	"strconv"
	"fmt"
	"encoding/json"
	"log"
	"github.com/skip2/go-qrcode"
	"encoding/base64"
)

type ComDB struct {
	storage storage.Storage
}

type PutAnswer struct {
	Id int `json:"id"`
}

type DetailAnswer struct {
	storage.ComItem
	QrCode string `json:"qr_code"`
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
		c.storage.Close()
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
		} else if r.Method == "GET" {
			c.Detail(w, r)
		}
	})
	return mux
}

func (c *ComDB) Put(w http.ResponseWriter, req *http.Request) {
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

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsoned)
}

func (c *ComDB) Post(w http.ResponseWriter, req *http.Request) {
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

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
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

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "{\"status\": \"ok\"}")
}

func (c *ComDB) Detail(w http.ResponseWriter, req *http.Request) {
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

	var item storage.ComItem
	item, err = c.storage.Get(i)
	if err != nil {
		http.Error(w, "Not found", 404)
		return
	}

	var png []byte
	png, err = qrcode.Encode(fmt.Sprintf("http://%s/components/?id=", req.Host, i), qrcode.Medium, 256)

	answer := DetailAnswer{
		item,
		fmt.Sprintf("data:image/png;base64,%s", base64.StdEncoding.EncodeToString([]byte(png))),
	}

	buf, err := json.Marshal(answer)
	if err != nil {
		log.Println(fmt.Sprintf("Failed to marhal id: `%d`", i))
		http.Error(w, "Internal server error", 500)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(buf)
}