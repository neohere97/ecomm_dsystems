//  /***************************************************************************
//   * Distributed Systems assignment 1
//   * Author: Chinmay Shalawadi
//   * Institution: University of Colorado Boulder
//   * Mail id: chsh1552@colorado.edu
//   ***************************************************************************/

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
	Name         string    `json:"name"`
	ItemCategory int       `json:"itemCategory"`
	ItemId       int       `json:"itemId"`
	Condition    string    `json:"condition"`
	Keywords     [5]string `json:"keywords"`
	Price        int       `json:"price"`
	SellerId     int       `json:"sellerId"`
	Availability int       `json:"availability"`
}

var logs bool = true

func perfTestingResponse(sessionnum int) {
	avg := 0
	for i := 0; i < 10; i++ {

		now := time.Now()
		epoch := now.UnixMilli()

		var newseller Seller

		var nameAttachment string
		nameAttachment = strconv.Itoa(int(epoch))
		newseller.SellerId = 0
		newseller.Name = "SK Holdings" + nameAttachment[7:]
		newseller.Password = "qwerty"
		newseller.ItemsSold = 0
		newseller.FeedbackNeg = 0
		newseller.FeedbackPos = 0

		if addNewSeller(&newseller) && logs {
			// fmt.Printf("SessionId:%v , User Registered \n", sessionNum)
		} else {
			// fmt.Printf("SessionId:%v , User Registration Failed \n", sessionNum)
		}
		end := time.Now()
		endepoch := end.UnixMilli()

		avg = avg + int(endepoch-epoch)

	}
	fmt.Printf("\n***************************************\n")
	fmt.Printf("Instance-%v Average Response Time is %vms\n", sessionnum, avg/10)
	fmt.Printf("***************************************\n")

}

func perfTestingThroughput(sessionnum int) {

	now := time.Now()
	epoch := now.UnixMilli()
	for i := 0; i < 1000; i++ {
		var newseller Seller

		var nameAttachment string
		nameAttachment = strconv.Itoa(int(epoch))
		newseller.SellerId = 0
		newseller.Name = "SK Holdings" + nameAttachment[7:]
		newseller.Password = "qwerty"
		newseller.ItemsSold = 0
		newseller.FeedbackNeg = 0
		newseller.FeedbackPos = 0

		if addNewSeller(&newseller) && logs {
			// fmt.Printf("SessionId:%v , User Registered \n", sessionNum)
		} else {
			// fmt.Printf("SessionId:%v , User Registration Failed \n", sessionNum)
		}

	}
	end := time.Now()
	endepoch := end.UnixMilli()
	fmt.Printf("\n***************************************\n")
	fmt.Printf("Instance-%v Time taken for 1000 operations is %vms\n", sessionnum, endepoch-epoch)
	fmt.Printf("***************************************\n")

}

// --------------------------main--------------------------------------------------
func main() {
	for {
		fmt.Println("\n1. Run Average Response Time Test")
		fmt.Println("2. Run Average Server Throughput Test")
		fmt.Println("3. Run Session Test")

		fmt.Printf("\n\nPlease select an option:\n")
		var i int
		fmt.Scanf("%d", &i)

		fmt.Printf("\n\nPlease give number of users (1-100):")
		var j int
		fmt.Scanf("%d", &j)

		if i == 1 {

			fmt.Printf("\n\nPress Enter to start Test")
			input := bufio.NewScanner(os.Stdin)
			input.Scan()
			fmt.Printf("\n\nTest Running....")
			for k := 0; k < j; k++ {
				go perfTestingResponse(k)
			}

		} else if i == 2 {
			fmt.Printf("\n\nPress Enter to start Test")
			input := bufio.NewScanner(os.Stdin)
			input.Scan()
			fmt.Printf("\n\nTest Running....")
			for k := 0; k < j; k++ {
				go perfTestingThroughput(k)
			}

		} else if i == 3 {
			fmt.Printf("\n\nPress Enter to start Test")
			input := bufio.NewScanner(os.Stdin)
			input.Scan()
			fmt.Printf("\n\nTest Running....")
			for k := 0; k < j; k++ {
				go session(k)
			}
		}
	}
}

// --------------------------addNewSeller--------------------------------------------------

func addNewSeller(newseller *Seller) bool {

	// fmt.Println("Adding a new Buyer")
	connection, err := net.Dial(SERVER_TYPE, SERVER_HOST+":"+SERVER_PORT)

	if err != nil {
		panic(err)
	}

	var req Request
	req.ReqType = "addSeller"
	req.Data, _ = json.Marshal(*newseller)

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

	json.Unmarshal(buffer[:mLen], newseller)

	return true
}

// --------------------------login--------------------------------------------------

func login(credentials Seller) bool {
	connection, err := net.Dial(SERVER_TYPE, SERVER_HOST+":"+SERVER_PORT)

	if err != nil {
		panic(err)
	}

	var req Request
	req.ReqType = "loginSeller"
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

	defer connection.Close()

	if string(buffer[:mLen]) == `"200"` {
		return true
	} else {
		return false
	}
}

// --------------------------logout--------------------------------------------------

func logout(credentials Seller) bool {
	// fmt.Println("Adding a new Buyer")
	connection, err := net.Dial(SERVER_TYPE, SERVER_HOST+":"+SERVER_PORT)

	if err != nil {
		panic(err)
	}

	var req Request
	req.ReqType = "logoutSeller"
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

// --------------------------getAllProducts--------------------------------------------------

func getAllProducts(prods *[]Product) {
	connection, err := net.Dial(SERVER_TYPE, SERVER_HOST+":"+SERVER_PORT)

	if err != nil {
		panic(err)
	}

	var req Request
	req.ReqType = "getProducts"

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

	defer connection.Close()
	json.Unmarshal(buffer[:mLen], prods)
}

func addNewProduct(newproduct *Product) bool {
	connection, err := net.Dial(SERVER_TYPE, SERVER_HOST+":"+SERVER_PORT)

	if err != nil {
		panic(err)
	}

	var req Request
	req.ReqType = "addProduct"
	req.Data, _ = json.Marshal(*newproduct)

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

	json.Unmarshal(buffer[:mLen], &newproduct)

	return true

}

// --------------------------session--------------------------------------------------

func session(sessionNum int) {
	var sellerprods []Product
	// var cart []Product

	now := time.Now()
	epoch := now.UnixMilli()

	//Create New Seller
	var newseller Seller

	var nameAttachment string
	nameAttachment = strconv.Itoa(int(epoch))
	newseller.SellerId = 0
	newseller.Name = "SK Holdings" + nameAttachment[7:]
	newseller.Password = "qwerty"
	newseller.ItemsSold = 0
	newseller.FeedbackNeg = 0
	newseller.FeedbackPos = 0

	if addNewSeller(&newseller) && logs {
		fmt.Printf("SessionId:%v , Seller Registered \n", sessionNum)
	} else {
		fmt.Printf("SessionId:%v , Seller Registration Failed \n", sessionNum)
	}

	//Login New Seller
	if login(newseller) && logs {
		fmt.Printf("SessionId:%v , Login Successful\n", sessionNum)
	} else {
		fmt.Printf("SessionId:%v , Login Failed\n", sessionNum)
	}

	//Put One item for sale
	var newproduct Product
	newproduct.Name = "iPhone Sad " + nameAttachment[9:]
	newproduct.Condition = "New"
	newproduct.ItemCategory = 33
	newproduct.Availability = 23
	newproduct.Keywords[0] = "sad"
	newproduct.Keywords[1] = "iphone"
	newproduct.Keywords[2] = "apple"
	newproduct.SellerId = newseller.SellerId
	newproduct.Price = 87

	if addNewProduct(&newproduct) && logs {
		fmt.Printf("SessionId:%v , Product Added \n", sessionNum)
	} else {
		fmt.Printf("SessionId:%v , Product Additon Failed \n", sessionNum)
	}

	//Update Sale Price
	newproduct.Price = 8700
	newproduct.Availability -= 1
	if updateProduct(newproduct) && logs {
		fmt.Printf("SessionId:%v , Product Updated \n", sessionNum)
	} else {
		fmt.Printf("SessionId:%v , Product Update Failed \n", sessionNum)
	}

	//Display items on sale
	getSellerProducts(newseller, &sellerprods)

	fmt.Printf("Below are the products on sale by the seller\n")
	for i := 0; i < len(sellerprods); i++ {
		fmt.Printf("%+v \n", sellerprods[i])
	}

	//Logout Seller
	if logout(newseller) && logs {
		fmt.Printf("SessionId:%v , Logout Successful\n", sessionNum)
	} else {
		fmt.Printf("SessionId:%v , Logout Failed\n", sessionNum)
	}

	end := time.Now()
	endepoch := end.UnixMilli()

	fmt.Printf("time in session#%03v %vms \n", sessionNum, endepoch-epoch)
}

func updateProduct(newproduct Product) bool {
	connection, err := net.Dial(SERVER_TYPE, SERVER_HOST+":"+SERVER_PORT)

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

	buffer := make([]byte, 1024)
	mLen, err := connection.Read(buffer)

	if err != nil {
		fmt.Println("Error reading:", err.Error())
	}

	defer connection.Close()

	if string(buffer[:mLen]) == `"200"` {
		return true
	} else {
		return false
	}
}

func getSellerProducts(seller Seller, newproduct *[]Product) {
	connection, err := net.Dial(SERVER_TYPE, SERVER_HOST+":"+SERVER_PORT)

	if err != nil {
		panic(err)
	}

	var req Request
	req.ReqType = "getSellerProducts"
	req.Data, _ = json.Marshal(seller)

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

	defer connection.Close()

	json.Unmarshal(buffer[:mLen], newproduct)
}

// --------------------------End--------------------------------------------------
