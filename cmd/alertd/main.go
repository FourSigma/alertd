package main

import (
	"fmt"
	"log"

	"github.com/FourSigma/alertd/internal/http"
)

func main() {
	port := "4040"
	fmt.Println("Starting server on port -- ", port)
	log.Fatal(http.NewAPI(port).Run())
}
