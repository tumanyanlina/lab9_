package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
)

type Numbers struct {
	Nums []int `json:"numbers"`
}

type Result struct {
	Sum int `json:"sum"`
}

func handleConnection(conn net.Conn) {
	defer conn.Close()
	reader := bufio.NewReader(conn)
	for {
		line, err := reader.ReadBytes('\n')
		if err != nil {
			return
		}
		var nums Numbers
		if err := json.Unmarshal(line, &nums); err != nil {
			conn.Write([]byte("{\"error\":\"invalid json\"}\n"))
			continue
		}
		sum := 0
		for _, n := range nums.Nums {
			sum += n * n
		}
		res, _ := json.Marshal(Result{Sum: sum})
		conn.Write(append(res, '\n'))
	}
}

func main() {
	ln, err := net.Listen("tcp", ":8080")
	if err != nil {
		fmt.Println("Error starting server:", err)
		return
	}
	defer ln.Close()
	fmt.Println("TCP server listening on :8080")
	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println("Failed to accept connection:", err)
			continue
		}
		go handleConnection(conn)
	}
}