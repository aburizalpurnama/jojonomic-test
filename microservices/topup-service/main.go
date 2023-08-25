package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"price-check-service/repository"
	"price-input-service/broker"
	"strconv"
	"time"
	"topup-service/usecases"
	"topup-service/usecases/request"
	"topup-service/usecases/response"

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
	router.HandleFunc("/api/topup", topup).Methods(http.MethodPost)

	server := &http.Server{
		Handler:      router,
		Addr:         fmt.Sprintf(":%s", os.Getenv("PORT")),
		WriteTimeout: 5 * time.Second,
		ReadTimeout:  5 * time.Second,
	}
	log.Fatal(server.ListenAndServe())
}

func topup(w http.ResponseWriter, r *http.Request) {
	var requestPayload request.TopupRequest
	var responsePayload response.TopupResponse
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

	gram, err := strconv.ParseFloat(requestPayload.Gram, 32)
	if err != nil {
		logrus.Error(err)
		responsePayload.Error = true
		responsePayload.Message = "invalid gram input"
		payload, _ := json.Marshal(responsePayload)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		w.Write(payload)
		return
	}

	count := usecases.CountDecimalPlaces(float32(gram))
	if count > 3 {
		logrus.Info("invalid gram, max 3 decimal places")
		responsePayload.Error = true
		responsePayload.Message = "invalid gram, max 3 decimal places"
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

	repo := repository.NewPriceCheckRepo(db)

	price, err := repo.GetLatestPrice()
	if err != nil {
		responsePayload.Error = true
		responsePayload.Message = err.Error()
		payload, _ := json.Marshal(responsePayload)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(payload)
		return
	}

	harga, err := strconv.Atoi(requestPayload.Harga)
	if err != nil {
		logrus.Error(err)
		responsePayload.Error = true
		responsePayload.Message = "invalid harga input"
		payload, _ := json.Marshal(responsePayload)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		w.Write(payload)
		return
	}

	if harga != int(price.Topup) {
		msg := "invalid harga, not equal to latest topup price"
		logrus.Info(msg)
		responsePayload.Error = true
		responsePayload.Message = msg
		payload, _ := json.Marshal(responsePayload)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		w.Write(payload)
		return
	}

	strGram := fmt.Sprintf("%.3f", gram)
	newGram, _ := strconv.ParseFloat(strGram, 32)

	data := usecases.Topup{
		Gram:         float32(newGram),
		HargaTopup:   price.Topup,
		HargaBuyback: price.Buyback,
		Norek:        requestPayload.Norek,
	}

	usecase := usecases.NewTopupUsecase(broker.NewMessageBroker())
	responsePayload = usecase.Topup(data)
	payload, _ := json.Marshal(responsePayload)

	w.Header().Set("Content-Type", "application/json")
	if responsePayload.Error {
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		w.WriteHeader(http.StatusOK)
	}
	w.Write(payload)
}
