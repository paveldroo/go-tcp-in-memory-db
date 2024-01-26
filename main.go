package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"strings"
)

func main() {
	li, err := net.Listen("tcp", ":8000")
	if err != nil {
		panic(err)
	}

	defer func() {
		err := li.Close()
		if err != nil {
			panic(err)
		}
	}()

	fmt.Println("Server was started!")

	for {
		conn, err := li.Accept()
		if err != nil {
			panic(err)
		}
		go handle(conn)
	}
}

func handle(conn net.Conn) {
	defer func() {
		err := conn.Close()
		if err != nil {
			panic(err)
		}
	}()
	io.WriteString(conn, "\nIN-MEMORY DATABASE\n\n"+
		"USE:\n"+
		"SET key value \n"+
		"GET key \n"+
		"DEL key \n\n"+
		"EXAMPLE:\n"+
		"SET fruit banana \n"+
		"GET fruit \n\n\n",
	)

	m := make(map[string]string)

	s := bufio.NewScanner(conn)
	for s.Scan() {
		ln := s.Text()
		fs := strings.Fields(ln)

		switch fs[0] {
		case "GET":
			k := fs[1]
			v := m[k]
			fmt.Fprintln(conn, v)
		case "SET":
			if len(fs) != 3 {
				fmt.Fprintln(conn, "Expected 3 elements in SET command: SET {key} {value")
				continue
			}
			k := fs[1]
			v := fs[2]
			m[k] = v
			fmt.Fprintln(conn, "New key and value were set!")
		case "DEL":
			k := fs[1]
			delete(m, k)
			fmt.Fprintln(conn, "Key and value were deleted!")
		default:
			fmt.Fprintln(conn, "Invalid command!")
		}

	}

}
