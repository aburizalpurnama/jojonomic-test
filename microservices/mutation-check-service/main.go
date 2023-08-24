package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"mutation-check-service/model"
	"mutation-check-service/repository"
	"net/http"
	"os"
	"time"
	topupRepo "topup-storage/repository"

	"github.com/gorilla/mux"
	"github.com/jinzhu/copier"
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
	router.HandleFunc("/api/mutasi", checkMutation).Methods(http.MethodGet)

	server := &http.Server{
		Handler:      router,
		Addr:         fmt.Sprintf(":%s", os.Getenv("PORT")),
		WriteTimeout: 5 * time.Second,
		ReadTimeout:  5 * time.Second,
	}
	log.Fatal(server.ListenAndServe())
}

func checkMutation(w http.ResponseWriter, r *http.Request) {
	var requestPayload model.MutationCheckRequest
	var responsePayload model.MutationCheckResponse
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

	topupRepo := topupRepo.NewTopupRepo(db)
	_, err = topupRepo.GetAccountByNorek(requestPayload.Norek)
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

	var transactions []model.Transcation

	mutationRepo := repository.NewMutationCheckRepo(db)
	result, err := mutationRepo.SelectTransactionsByNorek(requestPayload.Norek, requestPayload.StartDate, requestPayload.EndDate)
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
	case nil:
		for _, v := range result {
			var t model.Transcation
			copier.Copy(&t, &v)
			transactions = append(transactions, t)
		}
	}

	responsePayload.Data = transactions
	responsePayload.Error = false
	payload, _ := json.Marshal(responsePayload)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(payload)
}
