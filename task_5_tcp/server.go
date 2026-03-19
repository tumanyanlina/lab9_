package main

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

type Numbers struct {
	Nums []int `json:"numbers"`
}

type Response struct {
	Sum      int   `json:"sum,omitempty"`
	Original []int `json:"original,omitempty"`
	Error    string `json:"error,omitempty"`
}

func processNumbers(nums []int) Response {
	if len(nums) == 0 {
		return Response{Error: "no numbers provided"}
	}

	for _, n := range nums {
		if n > 1000 {
			return Response{Error: "number too large"}
		}
	}

	sum := 0
	for _, n := range nums {
		sum += n * n
	}

	return Response{
		Sum:      sum,
		Original: nums,
	}
}

func handleConnection(conn net.Conn) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("Recovered from panic in handleConnection: %v\n", r)
		}
		conn.Close()
	}()
	reader := bufio.NewReader(conn)

	for {
		// Set read deadline for 5 minutes from now
		if err := conn.SetReadDeadline(time.Now().Add(5 * time.Minute)); err != nil {
			fmt.Println("Failed to set read deadline:", err)
			return
		}
		line, err := reader.ReadBytes('\n')
		if err != nil {
			return
		}

		fmt.Printf("Incoming request: %s", line)

		var nums Numbers
		if err := json.Unmarshal(line, &nums); err != nil {
			resp, marshalErr := json.Marshal(Response{Error: "invalid json"})
			if marshalErr != nil {
				fmt.Println("json.Marshal error:", marshalErr)
				return
			}
			fmt.Printf("Server response: %s\n", resp)
			// Set write deadline for 30 seconds from now
			conn.SetWriteDeadline(time.Now().Add(30 * time.Second))
			if _, err := conn.Write(append(resp, '\n')); err != nil {
				fmt.Println("conn.Write error:", err)
				return
			}
			continue
		}

		resp := processNumbers(nums.Nums)
		output, marshalErr := json.Marshal(resp)
		if marshalErr != nil {
			fmt.Println("json.Marshal error:", marshalErr)
			return
		}
		fmt.Printf("Server response: %s\n", output)
		conn.SetWriteDeadline(time.Now().Add(30 * time.Second))
		if _, err := conn.Write(append(output, '\n')); err != nil {
			fmt.Println("conn.Write error:", err)
			return
		}
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

	// Graceful shutdown vars
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	var wg sync.WaitGroup

	// Channel for accepted connections
	connections := make(chan net.Conn)

	// Goroutine to accept connections
	go func() {
		for {
			conn, err := ln.Accept()
			if err != nil {
				select {
				case <-ctx.Done():
					// Listener closed, exit accept loop
					close(connections)
					return
				default:
					fmt.Println("Failed to accept connection:", err)
					continue
				}
			}
			select {
			case <-ctx.Done():
				conn.Close()
			case connections <- conn:
			}
		}
	}()

acceptLoop:
	for {
		select {
		case <-ctx.Done():
			fmt.Println("Shutting down gracefully...")
			ln.Close()
			// Wait a moment for ongoing connections to wrap up
			break acceptLoop
		case conn, ok := <-connections:
			if !ok {
				break acceptLoop
			}
			wg.Add(1)
			go func(c net.Conn) {
				defer wg.Done()
				handleConnection(c)
			}(conn)
		}
	}

	wg.Wait()
	fmt.Println("Server exited.")
}