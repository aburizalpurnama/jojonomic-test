package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"price-check-service/repository"
	"price-check-service/usecases"
	"time"

	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	router := mux.NewRouter()
	router.HandleFunc("/api/check-harga", checkPrice).Methods(http.MethodGet)

	server := &http.Server{
		Handler:      router,
		Addr:         fmt.Sprintf(":%s", os.Getenv("PORT")),
		WriteTimeout: 5 * time.Second,
		ReadTimeout:  5 * time.Second,
	}
	log.Fatal(server.ListenAndServe())
}

func checkPrice(w http.ResponseWriter, r *http.Request) {
	db, err := sqlx.Connect("postgres", os.Getenv("DB_DSN"))
	if err != nil {
		log.Fatalln(err)
	}

	repo := repository.NewPriceCheckRepo(db)
	usecase := usecases.NewPriceCheckUsecase(repo)
	responsePayload, err := usecase.CheckPrice()
	payload, _ := json.Marshal(responsePayload)

	w.Header().Set("Content-Type", "application/json")
	w.Write(payload)
	if responsePayload.Error {
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		w.WriteHeader(http.StatusOK)
	}
}
