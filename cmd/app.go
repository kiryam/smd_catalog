package main

import (
	"fmt"
	"smd_catalog"
	"smd_catalog/catalog"
	"content_api/common/log"
)

func main() {
	c, err := catalog.NewCatalog()

	if err != nil {
		log.Fatal(fmt.Sprintf("Failed to initialize catalog (%s)", err.Error()))
	}
	defer c.Close()

	server := smd_catalog.NewApiServer(8080, c)
	err = server.Start()

	if err != nil {
		log.Fatal(err.Error())
	}
}
