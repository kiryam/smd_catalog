package main

import (
	"fmt"
	"smd_catalog"
	"smd_catalog/catalog"
	"smd_catalog/comdb"
	"os"
	"log"
)

func main() {
	c, err := catalog.NewCatalog()

	if err != nil {
		log.Fatal(fmt.Sprintf("Failed to initialize catalog (%s)", err.Error()))
	}
	defer c.Close()

	comdb, err := comdb.NewComDB()
	if err != nil {
		log.Fatal(fmt.Sprintf("Failed to initialize components db (%s)", err.Error()))
	}
	defer comdb.Close()

	server := smd_catalog.NewApiServer(8080, c, comdb)
	err = server.Start()

	if err != nil {
		log.Fatal(err.Error())
	}

	os.Exit(3)
}
