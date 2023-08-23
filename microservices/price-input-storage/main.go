package main

import (
	"log"
	"os"
	"price-input-storage/broker"
	"price-input-storage/repository"

	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	db, err := sqlx.Connect("postgres", os.Getenv("DB_DSN"))
	if err != nil {
		log.Fatalln(err)
	}

	repo := repository.NewPriceInputRepo(db)
	messageBroker := broker.NewMessageBroker(repo)
	err = messageBroker.Consume("input-harga")
	if err != nil {
		log.Println(err)
	}
}
