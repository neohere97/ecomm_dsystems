package main

import (
	"encoding/json"
	"fmt"
	"net"
	"os"
)

const (
	SERVER_HOST        = "localhost"
	CUST_DATABASE_PORT = "9999"
	PROD_DATABASE_PORT = "9997"
	SERVER_TYPE        = "tcp"
	SERVER_PORT        = "9998"
)

const RECEIVE_BUFFER = 102400

var (
	server net.Listener
	err    error
)

type Request struct {
	ReqType string `json:"reqType"`
}

type Buyer struct {
	Name        string `json:"name"`
	Password    string `json:"passwd"`
	BuyerID     int    `json:buyerId`
	ItemsBought int    `json:itemsBought`
}

type Seller struct {
	Name        string `json:"name"`
	Password    string `json:"passwd"`
	SellerId    int    `json:sellerId`
	ItemsSold   int    `json:itemsSold`
	FeedbackPos int    `json:feedbackPos`
	FeedbackNeg int    `json:feedbackNeg`
}

var Sellers []Seller
var Buyers []Buyer

func main() {
	getCustomerDatabases()
	fmt.Println("Buyer Database Updated")

	//Setting up server part of the server
	setUpListener()
	defer server.Close()

	for {
		connection, err := server.Accept()
		if err != nil {
			fmt.Println("Error Accepting: ", err.Error())
			os.Exit(1)
		}
		fmt.Println("Client Connected...")
		go processClient(connection)
	}
}

func getCustomerDatabases() {

	fmt.Println("Getting Buyers database...")
	connection, err := net.Dial(SERVER_TYPE, SERVER_HOST+":"+CUST_DATABASE_PORT)

	if err != nil {
		panic(err)
	}

	var req Request

	req.ReqType = "getBuyers"
	reqBytes, _ := json.Marshal(req)
	_, err = connection.Write(reqBytes)

	buffer := make([]byte, RECEIVE_BUFFER)
	mLen, err := connection.Read(buffer)

	if err != nil {
		fmt.Println("Error reading:", err.Error())
	}

	json.Unmarshal(buffer[:mLen], &Buyers)

	for i := 0; i < len(Buyers); i++ {
		fmt.Printf("%v \n", (Buyers[i]))
	}

	defer connection.Close()

}

func setUpListener() {
	fmt.Println("Starting Server...\n")

	server, err = net.Listen(SERVER_TYPE, SERVER_HOST+":"+SERVER_PORT)

	if err != nil {
		fmt.Println("Error Listening:", err.Error())
		os.Exit(1)
	}
	fmt.Println("Listening on " + SERVER_HOST + ":" + SERVER_PORT + "\n")
	fmt.Println("Waiting for client... \n")
}

func processClient(connection net.Conn) {

	var req Request

	buffer := make([]byte, RECEIVE_BUFFER)

	mLen, err := connection.Read(buffer)

	if err != nil {
		fmt.Println("Error Reading:", err.Error())
	}

	json.Unmarshal(buffer[:mLen], &req)

	fmt.Println("Request Type is :", req.ReqType)

	// var marshalledBytes []byte

	// if req.ReqType == "getSellers" {
	// 	marshalledBytes, _ = json.Marshal(cust.Sellers)
	// }

	// if req.ReqType == "getBuyers" {
	// 	marshalledBytes, _ = json.Marshal(cust.Buyers)
	// }

	// _, err = connection.Write(marshalledBytes)
	connection.Close()
}
