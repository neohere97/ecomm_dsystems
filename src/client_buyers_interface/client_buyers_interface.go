package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"strconv"
	"time"
)

const (
	SERVER_HOST = "localhost"
	SERVER_TYPE = "tcp"
	SERVER_PORT = "9998"
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

type Buyer struct {
	Name        string `json:"name"`
	Password    string `json:"passwd"`
	BuyerID     int    `json:buyerId`
	ItemsBought int    `json:itemsBought`
}

type Request struct {
	ReqType string `json:"reqType"`
	Data    []byte `json:"data"`
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

var cart []Product
var logs bool = true

func main() {
	for {
		input := bufio.NewScanner(os.Stdin)
		input.Scan()
		session(1)
	}
}

func addNewBuyer(newbuyer Buyer) bool {

	// fmt.Println("Adding a new Buyer")
	connection, err := net.Dial(SERVER_TYPE, SERVER_HOST+":"+SERVER_PORT)

	if err != nil {
		panic(err)
	}

	var req Request
	req.ReqType = "addBuyer"
	req.Data, _ = json.Marshal(newbuyer)

	reqBytes, _ := json.Marshal(req)
	_, err = connection.Write(reqBytes)

	if err != nil {
		fmt.Println("Error writing:", err.Error())
	}

	buffer := make([]byte, 1024)
	mLen, err := connection.Read(buffer)

	if err != nil {
		fmt.Println("Error reading:", err.Error())
	}

	defer connection.Close()

	// fmt.Println(string(buffer[:mLen]))
	if string(buffer[:mLen]) == `"200"` {
		return true
	} else {
		return false
	}

}

func login(credentials Buyer) bool {
	// fmt.Println("Adding a new Buyer")
	connection, err := net.Dial(SERVER_TYPE, SERVER_HOST+":"+SERVER_PORT)

	if err != nil {
		panic(err)
	}

	var req Request
	req.ReqType = "loginBuyer"
	req.Data, _ = json.Marshal(credentials)

	reqBytes, _ := json.Marshal(req)
	_, err = connection.Write(reqBytes)

	if err != nil {
		fmt.Println("Error writing:", err.Error())
	}

	buffer := make([]byte, 1024)
	mLen, err := connection.Read(buffer)

	if err != nil {
		fmt.Println("Error reading:", err.Error())
	}

	// fmt.Println(string(buffer[:mLen]))
	defer connection.Close()

	if string(buffer[:mLen]) == `"200"` {
		return true
	} else {
		return false
	}
}

func session(sessionNum int) {
	now := time.Now()
	epoch := now.UnixMilli()

	//CreateNewUser
	var newbuyer Buyer
	var nameAttachment string
	newbuyer.BuyerID = 0
	nameAttachment = strconv.Itoa(int(epoch) / 100)
	newbuyer.Name = "Mona" + nameAttachment[5:]
	newbuyer.Password = "12345"
	if addNewBuyer(newbuyer) && logs {
		fmt.Printf("SessionId:%v , User Registered \n", sessionNum)
	} else {
		fmt.Printf("SessionId:%v , User Registration Failed \n", sessionNum)
	}

	//Login
	if login(newbuyer) && logs {
		fmt.Printf("SessionId:%v , Login Successful\n", sessionNum)
	} else {
		fmt.Printf("SessionId:%v , Login Failed\n", sessionNum)
	}

	//Display Products
	//Add Item to Shopping cart
	//Display Shopping cart
	//Remove Item from Shopping Cart
	//Display Shopping cart
	//Logout

	end := time.Now()
	endepoch := end.UnixMilli()

	fmt.Printf("time in session %vms \n", endepoch-epoch)
}
