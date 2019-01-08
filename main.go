package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/nikitsenka/bank-go/bank"
	"log"
	"net/http"
	"strconv"
	"strings"
)

var handler *bank.BankHandler

// our main function
func main() {
	dbService := bank.NewDbService()
	defer dbService.Close()

	handler = bank.NewHandler(dbService)

	router := mux.NewRouter()
	router.HandleFunc("/", HomeHandler)
	router.HandleFunc("/client/new/{deposit}", NewClientHandler).Methods("POST")
	router.HandleFunc("/transaction", NewTransactionHandler).Methods("POST")
	router.HandleFunc("/client/{id}/balance", BalanceHandler).Methods("GET")
	log.Fatal(http.ListenAndServe(":8080", router))
}

func HomeHandler(writer http.ResponseWriter, request *http.Request) {
	writer.Write([]byte(`{"status":"Ok"}`))
}

func NewTransactionHandler(writer http.ResponseWriter, request *http.Request) {
	decoder := json.NewDecoder(request.Body)
	var t bank.Transaction
	err := decoder.Decode(&t)
	if err != nil {
		errorResponse(writer, "Error deserializing request", err, http.StatusInternalServerError)
		return
	}
	new_transaction, err := handler.NewTransaction(t.From_client_id, t.To_client_id, t.Amount)
	if err != nil {
		errorResponse(writer, "Error in creating transaction", err, http.StatusInternalServerError)
		return
	}
	json.NewEncoder(writer).Encode(new_transaction)
}

func NewClientHandler(writer http.ResponseWriter, request *http.Request) {
	params := mux.Vars(request)
	s := params["deposit"]
	i, _ := strconv.Atoi(s)
	client, err := handler.NewClient(i)
	if err != nil {
		errorResponse(writer, "Error in creating client", err, http.StatusInternalServerError)
		return
	}
	json.NewEncoder(writer).Encode(client)
}

func BalanceHandler(writer http.ResponseWriter, request *http.Request) {
	params := mux.Vars(request)
	s := params["id"]
	i, _ := strconv.Atoi(s)
	response, err := handler.CheckBalance(i)
	if err != nil {
		errorResponse(writer, "Error in getting balance", err, http.StatusInternalServerError)
		return
	}

	json.NewEncoder(writer).Encode(response)
}


func errorResponse(w http.ResponseWriter, message string, err error, code int) {
	msg := fmt.Sprintf("{\"status\":\"Error\", \"message\":\"%v: %v\"}", message,
		strings.Replace(err.Error(), `"`, "'", -1))
	http.Error(w, msg, code)
}
