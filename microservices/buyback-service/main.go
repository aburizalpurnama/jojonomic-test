package main

import (
	"buyback-service/usecases"
	"buyback-service/usecases/request"
	"buyback-service/usecases/response"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	priceCheckRepo "price-check-service/repository"
	"price-input-service/broker"
	"time"
	topupRepo "topup-storage/repository"

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
	router.HandleFunc("/api/buyback", buyback).Methods(http.MethodPost)

	server := &http.Server{
		Handler:      router,
		Addr:         fmt.Sprintf(":%s", os.Getenv("PORT")),
		WriteTimeout: 5 * time.Second,
		ReadTimeout:  5 * time.Second,
	}
	log.Fatal(server.ListenAndServe())
}

func buyback(w http.ResponseWriter, r *http.Request) {
	var (
		reqPayload request.BuybackRequest
		resPayload response.BuybackResponse
	)

	err := json.NewDecoder(r.Body).Decode(&reqPayload)
	if err != nil {
		logrus.Error(err)
		resPayload.Error = true
		resPayload.Message = "invalid request payload"
		payload, _ := json.Marshal(resPayload)

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
	priceCheckRepo := priceCheckRepo.NewPriceCheckRepo(db)
	usecase := usecases.NewBuybackUsecase(broker.NewMessageBroker(), topupRepo, priceCheckRepo)
	resPayload, _ = usecase.Buyback(reqPayload)
	payload, _ := json.Marshal(resPayload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(resPayload.ResponseCode)
	w.Write(payload)
}
