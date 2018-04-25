package main

import (
	"fmt"

	log "github.com/Sirupsen/logrus"

	"github.com/FourSigma/alertd/internal/http"
)

func main() {
	l := log.New()
	port := "4040"
	fmt.Println("Starting server on port -- ", port)
	log.Fatal(http.NewAPI(port, l).Run())
}
