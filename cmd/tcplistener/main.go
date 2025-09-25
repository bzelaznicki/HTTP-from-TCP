package main

import (
	"fmt"
	"github.com/bzelaznicki/HTTP-from-TCP/internal/request"
	"log"
	"net"
)

const port = ":42069"

func main() {

	l, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatal(err)
	}
	defer l.Close()

	fmt.Println("Listening on", port)
	for {
		c, err := l.Accept()

		if err != nil {
			log.Fatal(err)
		}

		fmt.Println("Connection accepted")

		data, err := request.RequestFromReader(c)

		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("Request line:")
		fmt.Println("- Method:", data.RequestLine.Method)
		fmt.Println("- Target:", data.RequestLine.RequestTarget)
		fmt.Println("- Version:", data.RequestLine.HttpVersion)

		fmt.Printf("Connection to %s closed.\n", c.RemoteAddr())
	}
}
