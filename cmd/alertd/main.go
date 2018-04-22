package main

import (
	"log"

	"github.com/FourSigma/alertd/internal/http"
)

func main() {
	log.Fatal(http.NewAPI("8080").Run())
}
