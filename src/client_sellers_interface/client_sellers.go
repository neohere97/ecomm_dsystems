package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	"os"
)

const (
	SERVER_HOST = "localhost"
	SERVER_TYPE = "tcp"
	SERVER_PORT = "9996"
)

const RECEIVE_BUFFER = 102400

type Seller struct {
	Name        string `json:"name"`
	Password    string `json:"passwd"`
	SellerId    int    `json:sellerId`
	ItemsSold   int    `json:itemsSold`
	FeedbackPos int    `json:feedbackPos`
	FeedbackNeg int    `json:feedbackNeg`
}

type Request struct {
	ReqType string `json:"reqType"`
	Data    []byte `json:"data"`
}

func main() {
	for {
		input := bufio.NewScanner(os.Stdin)
		input.Scan()
		addNewSeller()
	}
}

func addNewSeller() {

	fmt.Println("Adding a new Seller")
	connection, err := net.Dial(SERVER_TYPE, SERVER_HOST+":"+SERVER_PORT)

	if err != nil {
		panic(err)
	}

	var req Request
	var newseller Seller

	newseller.SellerId = 9
	newseller.ItemsSold = 69
	newseller.Name = "M3gan"
	newseller.Password = "santa"
	newseller.FeedbackPos = 9
	newseller.FeedbackNeg = 990

	req.ReqType = "addSeller"
	req.Data, _ = json.Marshal(newseller)

	reqBytes, _ := json.Marshal(req)
	_, err = connection.Write(reqBytes)

	if err != nil {
		fmt.Println("Error writing:", err.Error())
	}

	buffer := make([]byte, RECEIVE_BUFFER)
	mLen, err := connection.Read(buffer)

	if err != nil {
		fmt.Println("Error reading:", err.Error())
	}

	fmt.Printf(string(buffer[:mLen]))
	defer connection.Close()

}
