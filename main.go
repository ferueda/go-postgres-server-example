package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/ferueda/go-postgres-server-example/db"
	"github.com/ferueda/go-postgres-server-example/pokemons"
	"github.com/gorilla/mux"
)

var dbURI, port, addr string

func init() {
	fmt.Println("loading env vars")

	port = os.Getenv("PORT")
	addr = os.Getenv("HOST") + ":" + port

	host := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbName := os.Getenv("DB_NAME")
	dbPassword := os.Getenv("DB_PASSWORD")

	dbURI = fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable", host, dbUser, dbPassword, dbName, dbPort)
}

func main() {
	// DB setup
	db, err := db.Init(dbURI)
	if err != nil {
		log.Fatal("error loading db")
	}

	// Server setup
	var wait time.Duration
	flag.DurationVar(&wait, "graceful-timeout", time.Second*15, "the duration for which the server gracefully wait for existing connections to finish")
	flag.Parse()

	l := log.New(os.Stdout, "pokemons-api", log.LstdFlags)
	ps := pokemons.NewStore(db)
	ph := pokemons.NewHandler(ps, l)

	r := mux.NewRouter()
	r.HandleFunc("/pokemons", ph.GetAll).Methods("GET")
	r.HandleFunc("/pokemons/{id:[0-9]+}", ph.GetOne).Methods("GET")

	http.Handle("/", r)

	srv := &http.Server{
		Addr: addr,

		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      r,
	}

	go func() {
		fmt.Printf("initializing server on port %s\n", port)
		if err := srv.ListenAndServe(); err != nil {
			log.Println(err)
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	<-c

	ctx, cancel := context.WithTimeout(context.Background(), wait)
	defer cancel()
	srv.Shutdown(ctx)
	log.Println("shutting down")
	os.Exit(0)
}
