package db

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"

	"github.com/ferueda/go-postgres-server-example/pokemons"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

var db *gorm.DB
var err error

func Init(dbURI string, dialect string) (*gorm.DB, error) {
	fmt.Printf("connecting to %s db\n", dialect)

	db, err = gorm.Open(dialect, dbURI)
	if err != nil {
		return nil, err
	}

	if err := runMigrations(db); err != nil {
		return nil, err
	}

	return db, nil
}

func runMigrations(db *gorm.DB) error {
	pokemonModel := pokemons.Pokemon{}

	db.DropTable(&pokemonModel)
	db.AutoMigrate(&pokemonModel)

	pokemons, err := getPokemonsFromFile("pokemons.json")
	if err != nil {
		return err
	}

	fmt.Printf("loading initial %v pokemons\n", len(pokemons))
	for _, p := range pokemons {
		db.Create(p)
	}

	return nil
}

func getPokemonsFromFile(fileName string) ([]*pokemons.Pokemon, error) {
	if fileName == "" {
		return nil, errors.New("must enter a valid file name")
	}

	file, err := os.Open(fileName)
	if err != nil {
		return nil, errors.New("unable to open pokemons json file")
	}

	var pokemons []*pokemons.Pokemon
	if err = json.NewDecoder(file).Decode(&pokemons); err != nil {
		return nil, errors.New("unable to unmarshal pokemons")
	}

	return pokemons, nil
}