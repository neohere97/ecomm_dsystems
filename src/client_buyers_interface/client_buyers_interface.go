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

var logs bool = true
var sessionLatencies []int

func perfTestingResponse(sessionnum int) {
	avg := 0
	for i := 0; i < 10; i++ {

		now := time.Now()
		epoch := now.UnixMilli()

		//CreateNewUser
		var newbuyer Buyer
		newbuyer.BuyerID = 0
		newbuyer.Name = "Mona"
		newbuyer.Password = "12345"

		if addNewBuyer(newbuyer) && logs {
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
		//CreateNewUser
		var newbuyer Buyer
		newbuyer.BuyerID = 0
		newbuyer.Name = "Mona"
		newbuyer.Password = "12345"

		if addNewBuyer(newbuyer) && logs {
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

// --------------------------addNewBuyer--------------------------------------------------

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

// --------------------------login--------------------------------------------------

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

// --------------------------logout--------------------------------------------------

func logout(credentials Buyer) bool {
	// fmt.Println("Adding a new Buyer")
	connection, err := net.Dial(SERVER_TYPE, SERVER_HOST+":"+SERVER_PORT)

	if err != nil {
		panic(err)
	}

	var req Request
	req.ReqType = "logoutBuyer"
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

// --------------------------session--------------------------------------------------

func session(sessionNum int) {
	var prods []Product
	var cart []Product

	now := time.Now()
	epoch := now.UnixMilli()

	//CreateNewUser
	var newbuyer Buyer
	var nameAttachment string
	newbuyer.BuyerID = 0
	nameAttachment = strconv.Itoa(int(epoch))
	newbuyer.Name = "Mona" + nameAttachment[7:] + strconv.Itoa(sessionNum)
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
	getAllProducts(&prods)
	fmt.Println("Here are the product listings...")
	for i := 0; i < len(prods); i++ {
		fmt.Printf("Product %v, %+v\n", i+1, prods[i])
	}
	//Add Item to Shopping cart
	cart = append(cart, prods[1])

	//Display Shopping cart
	fmt.Println("1 Product added to cart")
	for i := 0; i < len(cart); i++ {
		fmt.Printf("Cart Item %v, %+v\n", i+1, cart[i])
	}

	//Remove Item from Shopping Cart
	cart = cart[:len(cart)-1]
	//Display Shopping cart
	fmt.Println("Product Removed from Shopping Cart. Current Cart..")
	for i := 0; i < len(cart); i++ {
		fmt.Printf("Cart Item %v, %+v\n", i+1, cart[i])
	}

	//Logout
	if logout(newbuyer) && logs {
		fmt.Printf("SessionId:%v , Logout Successful\n", sessionNum)
	} else {
		fmt.Printf("SessionId:%v , Logout Failed\n", sessionNum)
	}

	end := time.Now()
	endepoch := end.UnixMilli()

	fmt.Printf("time in session#%v %vms \n", sessionNum, endepoch-epoch)
}

// --------------------------End--------------------------------------------------
