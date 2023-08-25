package main

import (
	"balance-service/request"
	"balance-service/response"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
	"topup-storage/repository"

	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	router := mux.NewRouter()
	router.HandleFunc("/api/saldo", checkPrice).Methods(http.MethodGet)

	server := &http.Server{
		Handler:      router,
		Addr:         fmt.Sprintf(":%s", os.Getenv("PORT")),
		WriteTimeout: 5 * time.Second,
		ReadTimeout:  5 * time.Second,
	}
	log.Fatal(server.ListenAndServe())
}

func checkPrice(w http.ResponseWriter, r *http.Request) {
	var requestPayload request.BalanceCheckRequest
	var responsePayload response.BalanceCheckResponse
	err := json.NewDecoder(r.Body).Decode(&requestPayload)
	if err != nil {
		logrus.Error(err)
		responsePayload.Error = true
		responsePayload.Message = "invalid request payload"
		payload, _ := json.Marshal(responsePayload)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		w.Write(payload)
		return
	}

	db, err := sqlx.Connect("postgres", os.Getenv("DB_DSN"))
	if err != nil {
		log.Fatalln(err)
	}

	repo := repository.NewTopupRepo(db)
	account, err := repo.GetAccountByNorek(requestPayload.Norek)
	switch err {
	default:
		responsePayload.Error = true
		responsePayload.Message = err.Error()
		responsePayload.Data = nil
		payload, _ := json.Marshal(responsePayload)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(payload)
		return
	case sql.ErrNoRows:
		responsePayload.Error = true
		responsePayload.Message = "account not found"
		responsePayload.Data = nil
		payload, _ := json.Marshal(responsePayload)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		w.Write(payload)
		return
	case nil:
	}

	responsePayload.Data = &response.Data{
		Norek: account.Norek,
		Saldo: account.Saldo,
	}
	responsePayload.Error = false
	payload, _ := json.Marshal(responsePayload)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(payload)
}
