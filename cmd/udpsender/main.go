package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
)

func main() {
	addr, err := net.ResolveUDPAddr("udp", "127.0.0.1:42069")

	if err != nil {
		log.Fatal(err)
	}

	d, err := net.DialUDP("udp", nil, addr)

	if err != nil {
		log.Fatal(err)
	}
	defer d.Close()
	r := bufio.NewReader(os.Stdin)

	for {
		fmt.Print(">")
		input, err := r.ReadString('\n')

		if err != nil {
			log.Fatal(err)
		}

		_, err = d.Write([]byte(input))

		if err != nil {
			log.Fatal(err)
		}
	}

}
