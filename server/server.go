package main

import (
	"io"
	"log"
	"net/http"
	"sync"
	"time"
)

type TransactionKind int

const (
	REPLACE TransactionKind = iota
	GET
)

type Transaction struct {
	kind TransactionKind
	data []byte
}

type Snapshot struct {
	journal []string
	body    string
}

var body string
var journal = make([]string, 0)
var transactions = make(chan Transaction)
var transactionMutex = &sync.Mutex{}
var responses = make(chan string)
var snapshot Snapshot

func replaceTransaction(data []byte) {
	transactionMutex.Lock()
	body = string(data)
	responses <- ""
	journal = append(journal, "REPLACE REQUEST: "+body)
	transactionMutex.Unlock()
}

func getTransaction() {
	transactionMutex.Lock()
	responses <- body
	journal = append(journal, "GET REQUEST")
	transactionMutex.Unlock()
}

func replaceHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	data, err := io.ReadAll(r.Body)
	if err != nil {
		log.Fatal(err)
	}
	transactions <- Transaction{REPLACE, data}
	response := <-responses
	_, err = w.Write([]byte(response))
	if err != nil {
		log.Fatal(err)
	}
}

func getHandler(w http.ResponseWriter, _ *http.Request) {
	transactions <- Transaction{GET, nil}
	response := <-responses
	_, err := w.Write([]byte(response))
	if err != nil {
		log.Fatal(err)
	}
}

func saveSnapshots() {
	for {
		time.Sleep(time.Minute)
		transactionMutex.Lock()
		journalCopy := make([]string, len(journal))
		copy(journalCopy, journal)
		snapshot = Snapshot{journalCopy, body}
		log.Println("Snapshot saved:", snapshot)
		transactionMutex.Unlock()
	}
}

func processTransactions() {
	go saveSnapshots()
	for {
		transaction := <-transactions
		switch transaction.kind {
		case REPLACE:
			replaceTransaction(transaction.data)
		case GET:
			getTransaction()
		}
	}
}

func main() {
	http.HandleFunc("/replace", replaceHandler)
	http.HandleFunc("/get", getHandler)
	go processTransactions()
	log.Fatal(http.ListenAndServe(":8080", nil))
}
