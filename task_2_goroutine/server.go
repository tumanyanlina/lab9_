
import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
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
	defer conn.Close()
	reader := bufio.NewReader(conn)

	for {
		line, err := reader.ReadBytes('\n')
		if err != nil {
			return
		}

		var nums Numbers
		if err := json.Unmarshal(line, &nums); err != nil {
			resp, _ := json.Marshal(Response{Error: "invalid json"})
			conn.Write(append(resp, '\n'))
			continue
		}

		resp := processNumbers(nums.Nums)
		output, _ := json.Marshal(resp)
		conn.Write(append(output, '\n'))
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