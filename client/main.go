package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
)

func main() {
	conn, err := net.Dial("tcp", ":8080")
	if err != nil {
		panic(err)
	}
	defer func() {
		if err := conn.Close(); err != nil {
			log.Println("Error closing connection:", err)
		}
	}()
	fmt.Println("Please input your name below: ")
	name, err := bufio.NewReader(os.Stdin).ReadString('\n')
	name = strings.TrimSpace(strings.Trim(name, "\n"))

	fmt.Println("Hello ", name)

	// Send the message to the server
	//if _, err := conn.Write([]byte("Hello World")); err != nil {
	//	panic(err)
	//}

	// Receive the message back from the server
	go func() {
		for {
			buf := make([]byte, 1024)
			n, err := conn.Read(buf)
			s := struct{
				Name string
				Message string
			}{}
			_ = json.Unmarshal(buf[:n], &s)
			if err != nil {
				panic(err)
			}
			fmt.Printf("%s: %s\n", s.Name, s.Message)
		}
	}()
	log.Println("Please input your message below: ")
	for {
		message, err := bufio.NewReader(os.Stdin).ReadString('\n')
		if err != nil {
			log.Println(err)
			break
		}
		message = strings.TrimSpace(strings.Trim(message, "\n"))

		s := struct{
			Name string
			Message string
		}{
			Name: name,
			Message: message,
		}
		data, _ := json.Marshal(&s)
		_, err = conn.Write(data)

	}
}
