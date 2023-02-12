package main

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"io/ioutil"
	"net"
	"os"
)

type Customers struct {
	Sellers []Seller `json:sellers`
	Buyers  []Buyer  `json:"buyers"`
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

const (
	SERVER_HOST = "localhost"
	SERVER_PORT = "9999"
	SERVER_TYPE = "tcp"
)

const RECEIVE_BUFFER = 102400

var cust Customers

var (
	server net.Listener
	err    error
)

type Request struct {
	ReqType string `json:"reqType"`
	Data    []byte `json:"data"`
}

func main() {
	jsonFile, err := os.Open("customers.json")

	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Successfully Opened the JSON File")

	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)

	json.Unmarshal(byteValue, &cust)

	fmt.Println("Json loaded into memory\n\n")

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

func update_filedb() {
	marshalledBytes, _ := json.Marshal(cust)

	os.WriteFile("customers.json", marshalledBytes, fs.ModeAppend)
}

func processClient(connection net.Conn) {

	var req Request

	buffer := make([]byte, RECEIVE_BUFFER)

	mLen, err := connection.Read(buffer)

	if err != nil {
		fmt.Println("Error Reading:", err.Error())
	}

	json.Unmarshal(buffer[:mLen], &req)

	fmt.Printf("Request Type is : %v \n", req.ReqType)

	var marshalledBytes []byte

	if req.ReqType == "getSellers" {
		marshalledBytes, _ = json.Marshal(cust.Sellers)
	}

	if req.ReqType == "getBuyers" {
		marshalledBytes, _ = json.Marshal(cust.Buyers)
	}

	if req.ReqType == "addBuyer" {
		var new_buyer Buyer

		json.Unmarshal(req.Data, &new_buyer)

		fmt.Printf("%v \n", new_buyer)

		cust.Buyers = append(cust.Buyers, new_buyer)

		marshalledBytes, _ = json.Marshal(cust.Buyers)

		update_filedb()
	}

	if req.ReqType == "addSeller" {
		var new_seller Seller

		json.Unmarshal(req.Data, &new_seller)

		fmt.Printf("%v \n", new_seller)

		cust.Sellers = append(cust.Sellers, new_seller)

		marshalledBytes, _ = json.Marshal(cust.Sellers)

		update_filedb()
	}

	_, err = connection.Write(marshalledBytes)
	connection.Close()
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

// func addCustomer() {

// }

// func getAllCustomers(){

// }
