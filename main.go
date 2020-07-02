package main

import (
	"go_micorservices/handlers"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/mux"
	"golang.org/x/net/context"
)

func main() {
	l := log.New(os.Stdout, "product-api", log.LstdFlags)
	// hh := handlers.NewHello(l)
	// gh := handlers.NewGoodbye(l)
	ps := handlers.NewProducts(l)

	// create new serveMux
	sm := mux.NewRouter()

	getRouter := sm.Methods("GET").Subrouter()
	getRouter.HandleFunc("/", ps.GetProducts)

	postRouter := sm.Methods("POST").Subrouter()
	postRouter.HandleFunc("/", ps.AddProduct)
	postRouter.Use(ps.MiddlewareProductValidation)

	putRouter := sm.Methods("PUT").Subrouter()
	putRouter.HandleFunc("/{id:[0-9]+}", ps.UpdateProduct)
	putRouter.Use(ps.MiddlewareProductValidation)

	// create custom server to better configure your sever to your needs
	s := &http.Server{
		Addr:         ":9090",
		Handler:      sm,
		IdleTimeout:  120 * time.Second,
		ReadTimeout:  1 * time.Second,
		WriteTimeout: 1 * time.Second,
	}

	s.ListenAndServe()

	go func() {
		err := s.ListenAndServe()

		if err != nil {
			l.Fatal(err)
		}
	}()

	// check https://www.youtube.com/watch?v=hodOppKJm5Y&t=9s to understand more
	sigChan := make(chan os.Signal)
	signal.Notify(sigChan, os.Interrupt)
	signal.Notify(sigChan, os.Kill)

	sig := <-sigChan

	l.Println("Recieved terminate, graceful shutdown", sig)

	tc, _ := context.WithTimeout(context.Background(), 30*time.Second)
	s.Shutdown(tc)
}
