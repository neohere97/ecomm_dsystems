package main

import (
	"fmt"
	"net"
)

const (
	SERVER_HOST = "localhost"
	SERVER_PORT = "9999"
	SERVER_TYPE = "tcp"
)

func main() {
	fmt.Println("Hello World, This one is going to be a server and a client whew")

	connection, err := net.Dial(SERVER_TYPE, SERVER_HOST+":"+SERVER_PORT)

	if err != nil {
		panic(err)
	}

	_, err = connection.Write([]byte("Hello Server! Greetings."))

	buffer := make([]byte, 1024)
	mLen, err := connection.Read(buffer)

	if err != nil {
		fmt.Println("Error reading:", err.Error())
	}
	fmt.Printf("Received %v Bytes \n", mLen)
	fmt.Println("Received: ", string(buffer[:mLen]))
	defer connection.Close()
}
