package main

import (
	"bytes"
	"context"
	"flag"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/websocket"
)

var (
	network   *string
	payload   *string
	checkType *string
	host      *string
)

func init() {
	network = flag.String("network", "bsc", "Network")
	payload = flag.String("payload", "", "Payload of the testing")
	checkType = flag.String("checkType", "ws", "check type http or ws")
	host = flag.String("host", "127.0.0.1", "Target host")
}

func wsCheck() {
	port := "8546"
	if *network == "solana" {
		port = "8900"
	}
	ws, _, err := websocket.DefaultDialer.Dial("ws://"+*host+":"+port, nil)
	if err != nil {
		log.Fatalf("%v", err)
		os.Exit(1)
	}
	defer ws.Close()
	if err := ws.WriteMessage(websocket.TextMessage, []byte(*payload)); err != nil {
		log.Fatalf("%v", err)
		os.Exit(1)
	}
	os.Exit(0)
}

func httpCheck() {
	port := "8545"
	if *network == "solana" {
		port = "8899"
	}
	timeout := time.Second * 5
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, "http://"+*host+":"+port, bytes.NewBuffer([]byte(*payload)))
	req.Header.Set("Content-Type", "application/json")
	if err != nil {
		log.Fatalf("%v", err)
		os.Exit(1)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatalf("%v", err)
		os.Exit(1)
	}
	defer resp.Body.Close()
	_, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("%v", err)
		os.Exit(1)
	}
	os.Exit(0)
}

func main() {
	flag.Parse()
	if *payload == "" {
		log.Fatalf("%v", "please specify a payload")
		os.Exit(1)
	}
	if *checkType == "http" {
		httpCheck()
	} else if *checkType == "ws" {
		wsCheck()
	} else {
		log.Fatalf("%v", "check type invalid, allowed http or ws")
		os.Exit(1)
	}
}
