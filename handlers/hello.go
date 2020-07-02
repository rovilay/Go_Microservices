package handlers

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

func NewHello(l *log.Logger) *Hello {
	return &Hello{l}
}

func (h *Hello) ServeHTTP(res http.ResponseWriter, r *http.Request) {
	log.Println("Hello World!")
	d, err := ioutil.ReadAll(r.Body)

	if err != nil {
		http.Error(res, "OOps", http.StatusBadRequest)
		return
	}

	// log.Printf("Data %s\n", d)
	fmt.Fprintf(res, "Hello %s\n", d)
}
