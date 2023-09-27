package main

import (
	"io"
	"log"
	"net/http"
	"strings"
)

func printResponseBody(resp *http.Response) {
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Response:", string(body))
}

func sendReplace(newBody string) {
	resp, err := http.Post("http://127.0.0.1:8080/replace", "text/html", strings.NewReader(newBody))
	if err != nil {
		log.Fatal(err)
	}
	printResponseBody(resp)
}

func sendGet() {
	resp, err := http.Get("http://127.0.0.1:8080/get")
	if err != nil {
		log.Fatal(err)
	}
	printResponseBody(resp)
}

func main() {
	sendGet()
	sendReplace("lol")
	sendGet()
	sendReplace("kek")
	sendGet()
	sendGet()
}
