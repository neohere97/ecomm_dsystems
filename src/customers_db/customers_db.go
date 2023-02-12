package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"os"
)

type Customers struct {
	Customers []Customer `json:"customers"`
}

type Customer struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"passwd"`
}

const (
	SERVER_HOST = "localhost"
	SERVER_PORT = "9999"
	SERVER_TYPE = "tcp"
)

var cust Customers

func main() {
	jsonFile, err := os.Open("customers.json")

	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Successfully Opened the JSON File")

	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)

	json.Unmarshal(byteValue, &cust)

	for i := 0; i < len(cust.Customers); i++ {
		fmt.Println(cust.Customers[i])
	}

	fmt.Println("Json loaded into memory\n\n")
	fmt.Println("Starting Server...\n")

	server, err := net.Listen(SERVER_TYPE, SERVER_HOST+":"+SERVER_PORT)

	if err != nil {
		fmt.Println("Error Listening:", err.Error())
		os.Exit(1)
	}

	defer server.Close()

	fmt.Println("Listening on " + SERVER_HOST + ":" + SERVER_PORT + "\n")
	fmt.Println("Waiting for client... \n")

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

	// marshalledBytes, _ := json.Marshal(cust)

	// os.WriteFile("customers.json", marshalledBytes, fs.ModeAppend)
}

func processClient(connection net.Conn) {
	buffer := make([]byte, 1024)
	mLen, err := connection.Read(buffer)

	if err != nil {
		fmt.Println("Error Reading:", err.Error())
	}

	fmt.Println("Received: ", string(buffer[:mLen]))

	marshalledBytes, _ := json.Marshal(cust)

	_, err = connection.Write(marshalledBytes)
	connection.Close()
}

// func addCustomer() {

// }

// func getAllCustomers(){

// }
