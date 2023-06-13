package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

type StreamData struct {
	Hash string        `json:"hash"`
	Logs []LogParams   `json:"logs,omitempty"`
	Txs  []Transaction `json:"txs"`
}

type LogParams struct {
	Address string   `json:"address"`
	Topics  []string `json:"topics"`
	Data    string   `json:"data"`
}

type Transaction struct {
	To               string `json:"to,omitempty"`
	FunctionSelector string `json:"functionSelector,omitempty"`
	CallData         string `json:"callData,omitempty"`
}

func main() {

	go listenToStream()

	// Keep the main goroutine running
	for {
		time.Sleep(time.Second)
	}
}

func listenToStream() {
	url := "https://mev-share.flashbots.net"

	// Create an HTTP GET request
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatal("Error creating request:", err)
	}

	// Set headers to indicate event stream subscription
	req.Header.Set("Accept", "text/event-stream")
	req.Header.Set("Cache-Control", "no-cache")

	// Send the request
	client := http.DefaultClient
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal("Error sending request:", err)
	}
	defer resp.Body.Close()

	// Check the response status
	if resp.StatusCode != http.StatusOK {
		log.Fatalf("Unexpected response status: %s", resp.Status)
	}

	// Read the event stream
	reader := bufio.NewReader(resp.Body)
	for {

		res, err := reader.ReadBytes('\n')
		if err != nil {
			log.Fatal("Error reading event stream: ", err)
		}

		if string(res) == ":ping\n"  {
			fmt.Println(":ping")
			continue
		}

		if  string(res) == "\n" {
			continue
		}

		res = bytes.TrimPrefix(res, []byte("data: "))
		res = bytes.TrimSuffix(res, []byte("\n"))

		var data StreamData

		err = json.Unmarshal(res, &data)
		if err != nil {
			fmt.Println("error decoding", err)
		}
		// Process the event
		processData(data)
		fmt.Println("-------------------------------------------")
	}
}

func processData(data StreamData) {
	// Process the received data here
	fmt.Println("Received hash:", data.Hash)
	
	for _, log := range data.Logs {
		fmt.Println("Logs address: ", log.Address)
		fmt.Println("Logs data: ", log.Data)
		for i, topic := range log.Topics {
			fmt.Println("Topic ", i, topic)
		}

	}

	for _, tx := range data.Txs {
		fmt.Println("Transaction details:")
		fmt.Println("Call data:", tx.CallData)
		fmt.Println("Function selector:", tx.FunctionSelector)
		fmt.Println("To address:", tx.To)
	}
}
