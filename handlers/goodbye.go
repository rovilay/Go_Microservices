package handlers

import (
	"log"
	"net/http"
)

func NewGoodbye(l *log.Logger) *Goodbye {
	return &Goodbye{l}
}

func (g *Goodbye) ServeHTTP(res http.ResponseWriter, r *http.Request) {
	res.Write([]byte("Byeee!\n"))
}
