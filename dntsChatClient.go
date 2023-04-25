package main

import (
	"bufio"
	"bytes"
	"fmt"
	"net"
	"os"
	"time"
)

func main() {
	conn, err := net.Dial("tcp", "serverIP:4753") // replace server ip
	if err != nil {
		fmt.Println("Error connecting:", err)
		os.Exit(1)
	}

	reader := bufio.NewReader(os.Stdin)
	fmt.Println("###########################")
	fmt.Println("### wlecome to dntsChat ###")
	fmt.Println("### function: /history  ###")
	fmt.Println("###########################\n")
	fmt.Print("Please enter your name: ")
	name, _ := reader.ReadString('\n')
	name = name[:len(name)-1]

	fmt.Fprintf(conn, name+"\n")

	go func() {
		for {
			//response, err := bufio.NewReader(conn).ReadString('\n')
			response, err := readString(conn)
			if err != nil {
				fmt.Println("Error receiving response:", err)
				os.Exit(1)
			}
			fmt.Println(response)
		}
	}()

	for {
		//fmt.Print("Enter message: ")
		message, _ := reader.ReadString('\n')
		now := time.Now()
		formatted := now.Format("15时04分")
		fmt.Printf("\033[1A\r[%s][%s]:%s", formatted, name, message)
		fmt.Fprintf(conn, message)
	}
}

func readString(conn net.Conn) (string, error) {
	var buffer bytes.Buffer
	for {
		b := make([]byte, 1)
		_, err := conn.Read(b)
		if err != nil {
			return "", err
		}
		if b[0] == '\n' {
			break
		}
		buffer.Write(b)
	}
	return buffer.String(), nil
}
