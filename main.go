package main

import (
	"bufio"
	"fmt"
	"net"
	_ "net"
	"os"
	"strings"
)

// from https://jameshfisher.com/2017/04/18/golang-tcp-server/

func main() {
	fmt.Println("chat client")

	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter username: ")
	username, _ := reader.ReadString('\n')
	username = strings.TrimSpace(username)
	fmt.Printf("your username is: '%s' \n", username)

	fmt.Print("Enter server: ")
	server, _ := reader.ReadString('\n')
	server = strings.TrimSpace(server)
	fmt.Printf("target server is: '%s' \n", server)

	// connect to server
	conn, err := net.Dial("tcp", server)
	if err != nil {
		panic(err)
	}

	// fire off a loop that always reads
	go func() {
		buf := make([]byte, 1024)
		for {
			fmt.Println("waiting to read")
			nbyte, err := conn.Read(buf)
			fmt.Println("not waiting. read", nbyte, "bytes")
			fmt.Println("from", conn.RemoteAddr())
			if err != nil {
				// deadConns <- conn
				fmt.Println("lost connection")
				break
			} else {
				fragment := make([]byte, nbyte)
				copy(fragment, buf[:nbyte])
				// publishes <- fragment

				// string from bytes
				st := string(fragment[:nbyte])
				fmt.Println("recd msg", st)
			}
		}
	}()

	for {
		fmt.Print("say: ")
		msg, _ := reader.ReadString('\n')
		msg = strings.TrimSpace(msg)
		fmt.Println(">", string(len(msg))+" "+msg)

		//for {
		sz, err := fmt.Fprintf(conn, msg)
		fmt.Println("sz", sz, "err", err)
		status, err := bufio.NewReader(conn).ReadString('\n')
		fmt.Println(status)
		//}
	}

	/* this is from server
	newConns := make(chan net.Conn, 128)
	deadConns := make(chan net.Conn, 128)
	publishes := make(chan []byte, 128)
	conns := make(map[net.Conn]bool)
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		panic(err)
	}
	go func() {
		for {
			conn, err := listener.Accept()
			if err != nil {
				panic(err)
			}
			newConns <- conn
		}
	}()
	for {
		select {
		case conn := <-newConns:
			conns[conn] = true
			go func() {
				buf := make([]byte, 1024)
				for {
					nbyte, err := conn.Read(buf)
					if err != nil {
						deadConns <- conn
						break
					} else {
						fragment := make([]byte, nbyte)
						copy(fragment, buf[:nbyte])
						publishes <- fragment
					}
				}
			}()
		case deadConn := <-deadConns:
			_ = deadConn.Close()
			delete(conns, deadConn)
		case publish := <-publishes:
			for conn, _ := range conns {
				go func(conn net.Conn) {
					totalWritten := 0
					for totalWritten < len(publish) {
						writtenThisCall, err := conn.Write(publish[totalWritten:])
						if err != nil {
							deadConns <- conn
							break
						}
						totalWritten += writtenThisCall
					}
				}(conn)
			}
		}
	}
	listener.Close()
	*/
}
