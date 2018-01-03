package smd_catalog

import (
	"encoding/json"
	"fmt"
	"net/http"
	"smd_catalog/catalog"
)

type ApiServer struct {
	port    int
	catalog *catalog.Catalog
}

type smdCatalogAnswer struct {
	Message string `json:"Message"`
	Error   string `json:"Error"`
}

func NewApiServer(port int, catalog *catalog.Catalog) *ApiServer {
	api := &ApiServer{port: port, catalog: catalog}
	return api
}

func (s *ApiServer) Start() error {
	mux := http.NewServeMux()
	mux.Handle("/catalog/",  http.StripPrefix("/catalog", s.catalog.GetServeMux()))

	mux.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		answer := smdCatalogAnswer{Message: "SmdCatalog"}

		data, err := json.Marshal(answer)
		if err != nil {
			fmt.Errorf("Failed to marshal answer")
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, "Internal server error")
		}

		w.WriteHeader(http.StatusOK)
		w.Write(data)
	})

	err := http.ListenAndServe(fmt.Sprintf(":%d", s.port), mux)

	if err != nil {
		return err
	}

	return nil
}