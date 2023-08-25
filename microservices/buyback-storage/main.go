package main

import (
	"buyback-storage/broker"
	"log"
	"os"
	"topup-storage/repository"

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

	repo := repository.NewTopupRepo(db)
	messageBroker := broker.NewMessageBroker(repo)
	err = messageBroker.Consume("buyback")
	if err != nil {
		log.Println(err)
	}
}
