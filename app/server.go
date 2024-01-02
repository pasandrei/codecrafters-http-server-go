package main

import (
	"fmt"
	"net"
	"os"
	"strings"
)

func main() {
	// You can use print statements as follows for debugging, they'll be visible when running tests.
	fmt.Println("Logs from your program will appear here!")

	// Uncomment this block to pass the first stage

	l, err := net.Listen("tcp", "0.0.0.0:4221")
	if err != nil {
		fmt.Println("Failed to bind to port 4221")
		os.Exit(1)
	}

	defer l.Close()

	connection, err := l.Accept()
	if err != nil {
		fmt.Println("Error accepting connection: ", err.Error())
		os.Exit(1)
	}

	// Get the path from the request
	buffer := make([]byte, 1024)
	_, _ = connection.Read(buffer)

	first_line := strings.Split(string(buffer), "\r\n")[0]
	path := strings.Split(first_line, " ")[1]

	fmt.Printf("Path: %s\n", path)
	if path == "/" {
		handleConnectionWrite(connection, "HTTP/1.1 200 OK\r\n\r\n")
	} else if strings.HasPrefix(path, "/echo/") {
		echoedString := strings.TrimPrefix(path, "/echo/")
		response := fmt.Sprintf("HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\nContent-Length: %d\r\n\r\n%s", len(echoedString), echoedString)

		handleConnectionWrite(connection, response)
	} else {
		handleConnectionWrite(connection, "HTTP/1.1 404 Not Found\r\n\r\n")
	}
}

func handleConnectionWrite(connection net.Conn, response string) {
	_, err := connection.Write([]byte(response))
	if err != nil {
		fmt.Println("Error writing to connection: ", err.Error())
		os.Exit(1)
	}
}
