package main

import (
	"encoding/json"
	"io"
	"log"
	"net"
)
var m = make(map[string]net.Conn)

func main() {
	log.SetFlags(log.Lshortfile)
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		panic(err)
	}

	for {
		conn, _ := listener.Accept()
		if len(m) > 3 {
			log.Printf("too many connections")
			conn.Close()
			continue
		}
		m[conn.RemoteAddr().String()] = conn
		log.Println("New connection from", conn.RemoteAddr().String())
		welcome, _ := json.Marshal(&struct {
			Name string
			Message string
		}{"Server", "Welcome to the server"})
		conn.Write(welcome)

		go func(conn net.Conn) {
			for {
				buf := make([]byte, 1024)
				n, err := conn.Read(buf)
				if err != nil {
					if err == io.EOF {
						log.Printf("%s disconnected\n", conn.RemoteAddr().String())
						break
					}
					log.Println(err)
					break
				}
				for i, v := range m {
					if v.RemoteAddr().String() != conn.RemoteAddr().String() {
						_, err := m[i].Write(buf[:n])
						if err != nil {
							log.Println(err)
							break
						}
					}
				}
			}
		}(conn)

	}
}
