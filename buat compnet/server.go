package main

import (
	"fmt"
	"io"
	"net"
)

func main() {
	listener, err := net.Listen("tcp", "localhost:1234")
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("bound to %q\n", listener.Addr())
	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			return
		}

		go func(c net.Conn) {
			defer func() {
				c.Close()
			}()

			fmt.Println("koneksi diterima")

			buffer := make([]byte, 1024)
			for {
				n, err := c.Read(buffer)
				if err != nil {
					if err != io.EOF {
						fmt.Println(err)
					}
					return
				}
				fmt.Println("diterima: %q", buffer[:n])
			}

		}(conn)
	}
}
