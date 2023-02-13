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
	"io/fs"
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

var prod Products

var (
	server net.Listener
	err    error
)

type Request struct {
	ReqType string `json:"reqType"`
	Data    []byte `json:"data"`
}

func main() {
	jsonFile, err := os.Open("products.json")

	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Successfully Opened the JSON File")

	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)

	json.Unmarshal(byteValue, &prod)

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
		marshalledBytes, _ = json.Marshal(prod.Products)
	}

	if req.ReqType == "addProduct" {
		var newproduct Product

		json.Unmarshal(req.Data, &newproduct)

		fmt.Printf("%v \n", newproduct)

		prod.Products = append(prod.Products, newproduct)

		marshalledBytes, _ = json.Marshal(prod.Products)

		update_filedb()
	}

	if req.ReqType == "updateProduct" {
		var newproduct Product

		json.Unmarshal(req.Data, &newproduct)

		marshalledBytes, _ = json.Marshal(string("500"))

		for i := 0; i < len(prod.Products); i++ {
			if prod.Products[i].ItemId == newproduct.ItemId {
				prod.Products[i] = newproduct
				marshalledBytes, _ = json.Marshal(string("200"))
				fmt.Println("product updated")
				break
			}
		}
		update_filedb()
	}

	_, err = connection.Write(marshalledBytes)
	connection.Close()
}

func update_filedb() {
	marshalledBytes, _ := json.Marshal(prod)

	os.WriteFile("products.json", marshalledBytes, fs.ModeAppend)
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
