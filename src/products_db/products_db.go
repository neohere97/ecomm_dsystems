package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"os"
)

type Products struct {
	Products []Product `json:products`
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

const (
	SERVER_HOST = "localhost"
	SERVER_PORT = "9995"
	SERVER_TYPE = "tcp"
)

const RECEIVE_BUFFER = 102400

var products Products

var (
	server net.Listener
	err    error
)

type Request struct {
	ReqType string `json:"reqType"`
}

func main() {
	jsonFile, err := os.Open("products.json")

	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Successfully Opened the JSON File")

	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)

	json.Unmarshal(byteValue, &products)

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

	// var new_cus Customer

	// new_cus.Email = "neo2@neohere.in"
	// new_cus.Name = "neo2"
	// new_cus.Password = "qtpie3.14"

	// cust.Customers = append(cust.Customers, new_cus)

	// for i := 0; i < len(cust.Customers); i++ {
	// 	fmt.Println(cust.Customers[i])
	// }

	// os.WriteFile("customers.json", marshalledBytes, fs.ModeAppend)
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

	if req.ReqType == "getProducts" {
		marshalledBytes, _ = json.Marshal(products)
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
