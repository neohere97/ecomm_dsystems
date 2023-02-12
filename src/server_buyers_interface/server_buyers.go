package main

import (
	"encoding/json"
	"fmt"
	"net"
	"os"
	"time"
)

const (
	SERVER_HOST        = "localhost"
	CUST_DATABASE_PORT = "9999"
	PROD_DATABASE_PORT = "9995"
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
	Data    []byte `json:"data"`
}

type Buyer struct {
	Name        string `json:"name"`
	Password    string `json:"passwd"`
	BuyerID     int    `json:buyerId`
	ItemsBought int    `json:itemsBought`
}

type Product struct {
	Name         string   `json:"name"`
	ItemCategory int      `json:"itemCategory"`
	ItemId       int      `json:"itemId"`
	Condition    string   `json:"condition"`
	Keywords     []string `json:"keywords"`
	Price        int      `json:"price"`
	SellerId     int      `json:"sellerId"`
	Availability int      `json:"availability"`
}

var Buyers []Buyer

type Products struct {
	Products []Product `json:"products"`
}

var prod Products

func main() {
	getCustomerDatabases()
	getProductDatabase()
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

func getProductDatabase() {

	fmt.Println("Getting product database...")
	connection, err := net.Dial(SERVER_TYPE, SERVER_HOST+":"+PROD_DATABASE_PORT)

	if err != nil {
		panic(err)
	}

	var req Request

	req.ReqType = "getProducts"

	reqBytes, _ := json.Marshal(req)
	_, err = connection.Write(reqBytes)

	buffer := make([]byte, RECEIVE_BUFFER)
	mLen, err := connection.Read(buffer)

	if err != nil {
		fmt.Println("Error reading:", err.Error())
	}
	json.Unmarshal(buffer[:mLen], &prod.Products)

	// for i := 0; i < len(prod.Products); i++ {
	// 	fmt.Printf("%v \n", (prod.Products[i]))
	// }

	defer connection.Close()

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

	// for i := 0; i < len(Buyers); i++ {
	// 	fmt.Printf("%v \n", (Buyers[i]))
	// }

	defer connection.Close()

}

func addNewBuyer(newbuyer Buyer) {

	fmt.Println("Adding a new buyer")
	connection, err := net.Dial(SERVER_TYPE, SERVER_HOST+":"+CUST_DATABASE_PORT)

	if err != nil {
		panic(err)
	}

	var req Request
	req.ReqType = "addBuyer"
	req.Data, _ = json.Marshal(newbuyer)

	reqBytes, _ := json.Marshal(req)
	_, err = connection.Write(reqBytes)

	buffer := make([]byte, RECEIVE_BUFFER)
	mLen, err := connection.Read(buffer)

	if err != nil {
		fmt.Println("Error reading:", err.Error())
	}

	json.Unmarshal(buffer[:mLen], &Buyers)

	// for i := 0; i < len(Buyers); i++ {
	// 	fmt.Printf("%v \n", (Buyers[i]))
	// }

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

	var marshalledBytes []byte

	if req.ReqType == "addBuyer" {
		var newbuyer Buyer

		json.Unmarshal(req.Data, &newbuyer)

		now := time.Now()
		newbuyer.BuyerID = int(now.UnixMilli())
		newbuyer.ItemsBought = 0

		// fmt.Printf("%v \n", newbuyer)

		Buyers = append(Buyers, newbuyer)

		marshalledBytes, _ = json.Marshal(string("200"))

		addNewBuyer(newbuyer)

	}

	if req.ReqType == "loginBuyer" {
		var newbuyer Buyer

		json.Unmarshal(req.Data, &newbuyer)

		for i := 0; i < len(Buyers); i++ {
			if newbuyer.Name == Buyers[i].Name {
				if newbuyer.Password == Buyers[i].Password {
					marshalledBytes, _ = json.Marshal(string("200"))
				} else {
					marshalledBytes, _ = json.Marshal(string("500"))
				}
			}
		}

	}

	_, err = connection.Write(marshalledBytes)
	connection.Close()
}
