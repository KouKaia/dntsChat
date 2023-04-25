package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
	"time"
)

type Client struct {
	conn net.Conn
	name string
}

var clients []Client

func main() {
	ln, err := net.Listen("tcp", "serverIP:4753")
	if err != nil {
		fmt.Println("Error listening:", err)
		return
	}
	defer ln.Close()

	fmt.Println("Server started, listening on port 4753...")

	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println("Error accepting:", err)
			continue
		}

		client := Client{
			conn: conn,
			name: "",
		}
		clients = append(clients, client)

		go handleClient(client)
	}
}

func handleClient(client Client) {
	defer client.conn.Close()

	reader := bufio.NewReader(client.conn)

	for {
		message, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Error reading message:", err)
			break
		}

		if client.name == "" {
			client.name = message[:len(message)-1]
			fmt.Println("New client:", client.name)
			continue
		}
		if message == "/history\n" {
			fmt.Print("histry request")
			logStr := readLog()
			fmt.Print(logStr)
			fmt.Fprintf(client.conn, "===histry===\n"+logStr+"\n===histry===\n")

		} else {
			now := time.Now()
			formatted := now.Format("15时04分")
			writeToLog("[" + formatted + "][" + client.name + "]:" + message)
			for _, c := range clients {
				if c.conn != client.conn {
					fmt.Fprintf(c.conn, "["+formatted+"]["+client.name+"]:"+message)
				}
			}
		}
	}

	fmt.Println("Client disconnected:", client.name)

	for i, c := range clients {
		if c.conn == client.conn {
			clients = append(clients[:i], clients[i+1:]...)
			break
		}
	}
}

func writeToLog(str string) {
	// 打开文件，如果不存在则创建文件，如果文件已存在则将数据追加到文件末尾
	file, err := os.OpenFile("log.txt", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()

	// 写入字符串到文件
	_, err = fmt.Fprint(file, str)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("log write")
}

func readLog() string {
	// open file
	file, err := os.Open("log.txt")
	if err != nil {
		fmt.Println(err)
		return ""
	}
	defer file.Close()

	// count all line
	scanner := bufio.NewScanner(file)
	lineCount := 0
	for scanner.Scan() {
		lineCount++
	}

	// read last 20 line
	startLine := lineCount - 20
	if startLine < 0 {
		startLine = 0
	}
	endLine := lineCount - 1

	// reopen log
	file, err = os.Open("log.txt")
	if err != nil {
		fmt.Println(err)
		return ""
	}
	defer file.Close()

	// read last 20 line
	scanner = bufio.NewScanner(file)
	var lines []string
	lineNum := 0
	for scanner.Scan() {
		if lineNum >= startLine {
			lines = append(lines, scanner.Text())
		}
		lineNum++
		if lineNum > endLine {
			break
		}
	}
	result := strings.Join(lines, "\n")
	// return 20 line with ¥n
	return result
}
