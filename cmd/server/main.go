package main

import (
	"fmt"
	"log"
	"os"

	"github.com/ferueda/go-postgres-server-example/pkg/api"
	"github.com/ferueda/go-postgres-server-example/pkg/app"
	"github.com/ferueda/go-postgres-server-example/pkg/db"
	"github.com/ferueda/go-postgres-server-example/pkg/repository"
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
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "error starting up the server: %s\\n", err)
		os.Exit(1)
	}
}

func run() error {
	// DB setup
	db, err := db.Init(dbURI)
	if err != nil {
		log.Fatal("error loading db")
	}

	uStore := repository.NewUserStore(db)
	pStore := repository.NewPokemonStore(db)

	uService := api.NewUserService(uStore)
	pService := api.NewPokemonService(pStore)
	r := mux.NewRouter()

	server := app.NewServer(r, addr, uService, pService)
	err = server.Run()
	if err != nil {
		return err
	}

	return nil
}
