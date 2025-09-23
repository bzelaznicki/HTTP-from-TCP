package main

import (
	"fmt"
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

		ch := getLinesChannel(c)

		for line := range ch {
			fmt.Printf("%s\n", line)
		}
		fmt.Printf("Connection to %s closed.\n", c.RemoteAddr())
	}
}
