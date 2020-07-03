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
		ErrorLog:     l,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	s.ListenAndServe()

	go func() {
		l.Println("Starting server on port n9090")

		err := s.ListenAndServe()
		if err != nil {
			l.Printf("Error starting server: %s\n", err)
			os.Exit(1)
		}
	}()

	// check https://www.youtube.com/watch?v=hodOppKJm5Y&t=9s to understand more
	sigChan := make(chan os.Signal)
	signal.Notify(sigChan, os.Interrupt)
	signal.Notify(sigChan, os.Kill)

	sig := <-sigChan
	log.Println("Recieved terminate, graceful shutdown", sig)

	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	s.Shutdown(ctx)
}
