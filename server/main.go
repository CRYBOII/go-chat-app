package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
)

type client chan string

var chMsg = make(chan string)
var enClient = make(chan client)
var leClient = make(chan client)

func main() {

	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatal(err)
	}

	go boardcaster()

	for {

		conn, err := listener.Accept()

		if err != nil {
			log.Println(err)
		}

		go handleConnection(conn)

	}

}

func boardcaster() {

	clients := make(map[client]bool)

	for {

		select {
		case m := <-chMsg:
			for k := range clients {
				k <- m
			}

		case cli := <-enClient:
			clients[cli] = true

		case cli := <-leClient:
			clients[cli] = false
		}

	}

}

func handleConnection(con net.Conn) {
	fmt.Println("Client conencted : ", con.RemoteAddr())
	chOutMsg := make(chan string)

	enClient <- chOutMsg

	clientName := con.RemoteAddr()

	go clientWriter(chOutMsg, con)

	scanner := bufio.NewScanner(con)
	fmt.Println(scanner.Scan())

	for scanner.Scan() {
		fmt.Println(scanner.Scan())

		msg := scanner.Text()

		chMsg <- fmt.Sprintf("%s => %s \n", clientName, msg)

	}

	leClient <- chOutMsg

	fmt.Printf("%s  : disconect \n", con.RemoteAddr())
}

func clientWriter(ch <-chan string, conn net.Conn) {

	for msg := range ch {

		fmt.Fprint(conn, msg)

	}

}
