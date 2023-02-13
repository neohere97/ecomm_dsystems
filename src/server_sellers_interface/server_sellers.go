//  /***************************************************************************
//   * Distributed Systems assignment 1
//   * Author: Chinmay Shalawadi
//   * Institution: University of Colorado Boulder
//   * Mail id: chsh1552@colorado.edu
//   ***************************************************************************/

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
	SERVER_PORT        = "9996"
)

var (
	server net.Listener
	err    error
)

const RECEIVE_BUFFER = 102400

type Request struct {
	ReqType string `json:"reqType"`
	Data    []byte `json:"data"`
}

type Seller struct {
	Name        string `json:"name"`
	Password    string `json:"passwd"`
	SellerId    int    `json:sellerId`
	ItemsSold   int    `json:itemsSold`
	FeedbackPos int    `json:feedbackPos`
	FeedbackNeg int    `json:feedbackNeg`
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

var Sellers []Seller

type Products struct {
	Products []Product `json:"products"`
}

var prod Products

func main() {

	getCustomerDatabases()
	getProductDatabase()
	fmt.Println("Seller Database Updated")
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

	fmt.Println("Getting Sellers database...")
	connection, err := net.Dial(SERVER_TYPE, SERVER_HOST+":"+CUST_DATABASE_PORT)

	if err != nil {
		panic(err)
	}
	var req Request

	req.ReqType = "getSellers"
	reqBytes, _ := json.Marshal(req)
	_, err = connection.Write(reqBytes)

	buffer := make([]byte, RECEIVE_BUFFER)
	mLen, err := connection.Read(buffer)

	if err != nil {
		fmt.Println("Error reading:", err.Error())
	}

	json.Unmarshal(buffer[:mLen], &Sellers)

	for i := 0; i < len(Sellers); i++ {
		fmt.Printf("%v \n", (Sellers[i]))
	}

	defer connection.Close()

}

func addNewSeller(newseller Seller) {

	fmt.Println("Adding a new Seller")
	connection, err := net.Dial(SERVER_TYPE, SERVER_HOST+":"+CUST_DATABASE_PORT)

	if err != nil {
		panic(err)
	}

	var req Request
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

	json.Unmarshal(buffer[:mLen], &Sellers)

	// for i := 0; i < len(Sellers); i++ {
	// 	fmt.Printf("%v \n", (Sellers[i]))
	// }

	defer connection.Close()

}

func addNewProduct(newproduct Product) {

	connection, err := net.Dial(SERVER_TYPE, SERVER_HOST+":"+PROD_DATABASE_PORT)

	if err != nil {
		panic(err)
	}

	var req Request
	req.ReqType = "addProduct"
	req.Data, _ = json.Marshal(newproduct)

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

	json.Unmarshal(buffer[:mLen], &prod)

	// for i := 0; i < len(prod.Products); i++ {
	// 	fmt.Printf("%v \n", (prod.Products[i]))
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

	if req.ReqType == "addSeller" {
		var new_seller Seller

		json.Unmarshal(req.Data, &new_seller)

		now := time.Now()
		new_seller.SellerId = int(now.UnixMilli())
		new_seller.ItemsSold = 0
		new_seller.FeedbackNeg = 0
		new_seller.FeedbackPos = 0

		Sellers = append(Sellers, new_seller)

		marshalledBytes, _ = json.Marshal(new_seller)

		addNewSeller(new_seller)

	}

	if req.ReqType == "loginSeller" {
		var newseller Seller

		json.Unmarshal(req.Data, &newseller)

		for i := 0; i < len(Sellers); i++ {
			if newseller.SellerId == Sellers[i].SellerId {
				if newseller.Password == Sellers[i].Password {
					marshalledBytes, _ = json.Marshal(string("200"))
				} else {
					marshalledBytes, _ = json.Marshal(string("500"))
				}
			}
		}

	}
	if req.ReqType == "logoutSeller" {
		var newseller Seller

		json.Unmarshal(req.Data, &newseller)

		for i := 0; i < len(Sellers); i++ {
			if newseller.SellerId == Sellers[i].SellerId {
				if newseller.Password == Sellers[i].Password {
					marshalledBytes, _ = json.Marshal(string("200"))
				} else {
					marshalledBytes, _ = json.Marshal(string("500"))
				}
			}
		}

	}

	if req.ReqType == "addProduct" {
		var newproduct Product

		json.Unmarshal(req.Data, &newproduct)
		now := time.Now()
		newproduct.ItemId = int(now.UnixMilli())

		prod.Products = append(prod.Products, newproduct)
		marshalledBytes, _ = json.Marshal(newproduct)

		go addNewProduct(newproduct)
	}

	if req.ReqType == "updateProduct" {
		var newproduct Product

		json.Unmarshal(req.Data, &newproduct)

		marshalledBytes, _ = json.Marshal(string("500"))

		for i := 0; i < len(prod.Products); i++ {
			if prod.Products[i].ItemId == newproduct.ItemId {
				prod.Products[i] = newproduct
				marshalledBytes, _ = json.Marshal(string("200"))
				go updateProduct(newproduct)
				break
			}
		}
	}

	if req.ReqType == "getSellerProducts" {
		var seller Seller
		var prods []Product

		json.Unmarshal(req.Data, &seller)

		for i := 0; i < len(prod.Products); i++ {
			if prod.Products[i].SellerId == seller.SellerId {
				prods = append(prods, prod.Products[i])
			}
		}
		marshalledBytes, _ = json.Marshal(prods)
	}

	_, err = connection.Write(marshalledBytes)
	connection.Close()
}

func updateProduct(newproduct Product) {

	connection, err := net.Dial(SERVER_TYPE, SERVER_HOST+":"+PROD_DATABASE_PORT)

	if err != nil {
		panic(err)
	}

	var req Request
	req.ReqType = "updateProduct"
	req.Data, _ = json.Marshal(newproduct)

	reqBytes, _ := json.Marshal(req)
	_, err = connection.Write(reqBytes)

	if err != nil {
		fmt.Println("Error writing:", err.Error())
	}
	defer connection.Close()

}
