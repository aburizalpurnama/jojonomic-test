package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"price-input-service/broker"
	"price-input-service/usecases"
	"price-input-service/usecases/request"
	"price-input-service/usecases/response"
	"time"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	router := mux.NewRouter()
	router.HandleFunc("/api/input-harga", inputPrice).Methods(http.MethodPost)

	server := &http.Server{
		Handler:      router,
		Addr:         fmt.Sprintf(":%s", os.Getenv("PORT")),
		WriteTimeout: 5 * time.Second,
		ReadTimeout:  5 * time.Second,
	}
	log.Fatal(server.ListenAndServe())
}

func inputPrice(w http.ResponseWriter, r *http.Request) {
	var requestPayload request.UpdatePriceRequest
	var responsePayload response.InputPriceResponse
	err := json.NewDecoder(r.Body).Decode(&requestPayload)
	if err != nil {
		responsePayload.Error = true
		responsePayload.Message = "invalid request payload"
		payload, _ := json.Marshal(responsePayload)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		w.Write(payload)
		return
	}

	usecase := usecases.NewPriceInputUsecase(broker.NewMessageBroker())
	responsePayload = usecase.InputPrice(requestPayload)
	payload, _ := json.Marshal(responsePayload)

	w.Header().Set("Content-Type", "application/json")
	w.Write(payload)
	if responsePayload.Error {
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		w.WriteHeader(http.StatusOK)
	}
}
