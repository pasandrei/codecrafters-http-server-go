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

	for {
		connection, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
			continue
		}

		go handleConnection(connection)
	}
}

func handleConnection(connection net.Conn) {
	defer connection.Close()

	// Get the path from the request
	buffer := make([]byte, 1024)
	_, _ = connection.Read(buffer)

	lines := strings.Split(string(buffer), "\r\n")

	startLine := lines[0]
	parts := strings.Split(startLine, " ")
	method, path := parts[0], parts[1]

	body := ""
	for i, line := range lines[1:] {
		if line == "" {
			body = strings.Join(lines[i+2:], "\r\n")
			break
		}
	}

	fmt.Printf("Path: %s\n", path)
	if path == "/" {
		handleConnectionWrite(connection, "HTTP/1.1 200 OK\r\n\r\n")
	} else if strings.HasPrefix(path, "/echo/") {
		echoedString := strings.TrimPrefix(path, "/echo/")
		response := fmt.Sprintf("HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\nContent-Length: %d\r\n\r\n%s", len(echoedString), echoedString)

		handleConnectionWrite(connection, response)
	} else if strings.HasPrefix(path, "/user-agent") {
		lines := strings.Split(string(buffer), "\r\n")
		for _, line := range lines {
			if strings.HasPrefix(line, "User-Agent: ") {
				userAgent := strings.TrimPrefix(line, "User-Agent: ")
				response := fmt.Sprintf("HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\nContent-Length: %d\r\n\r\n%s", len(userAgent), userAgent)

				handleConnectionWrite(connection, response)
				break
			}
		}
	} else if method == "POST" && strings.HasPrefix(path, "/files/") {
		fileName := strings.TrimPrefix(path, "/files/")
		file, err := os.Create(fmt.Sprintf("%s%s", os.Args[2], fileName))

		if err != nil {
			handleConnectionWrite(connection, "HTTP/1.1 500 Internal Server Error\r\n\r\n")
		} else {
			file.Write([]byte(body))
			handleConnectionWrite(connection, "HTTP/1.1 201 Created\r\n\r\n")
		}
	} else if strings.HasPrefix(path, "/files/") {
		fileName := strings.TrimPrefix(path, "/files/")
		file, err := os.Open(fmt.Sprintf("%s%s", os.Args[2], fileName))

		if err != nil {
			handleConnectionWrite(connection, "HTTP/1.1 404 Not Found\r\n\r\n")
		} else {
			fileInfo, _ := file.Stat()
			fileSize := fileInfo.Size()
			fileBuffer := make([]byte, fileSize)
			_, _ = file.Read(fileBuffer)
			fmt.Printf("File: %s\n", string(fileBuffer))
			response := fmt.Sprintf("HTTP/1.1 200 OK\r\nContent-Type: application/octet-stream\r\nContent-Length: %d\r\n\r\n%s", fileSize, string(fileBuffer))

			handleConnectionWrite(connection, response)
		}
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
