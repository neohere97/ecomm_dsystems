package main

import (
	"encoding/json"
	"fmt"
	"net"
)

const (
	SERVER_HOST        = "localhost"
	CUST_DATABASE_PORT = "9999"
	PROD_DATABASE_PORT = "9997"
	SERVER_TYPE        = "tcp"
	SERVER_PORT        = "9998"
)

const RECEIVE_BUFFER = 102400

type Request struct {
	ReqType string `json:"reqType"`
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

func main() {

	getCustomerDatabases()
	fmt.Println("Seller Database Updated")
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
