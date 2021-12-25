package main

import (
	"bufio"
	"fmt"
	"net"
)

func main() {
	conn, err := net.Dial("tcp", "hacktheoxidation.xyz:80")
	if err != nil {
		fmt.Println("Could not connect.")
	}

	fmt.Fprintf(conn, "GET / HTTP/1.0\r\n\r\n")
	status, err := bufio.NewReader(conn).ReadString('\n')

	fmt.Printf("Status: %s\n", status)
}
