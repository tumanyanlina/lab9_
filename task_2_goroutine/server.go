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

type Response struct {
    Sum      int   `json:"sum,omitempty"`
    Original []int `json:"original,omitempty"`
    Error    string `json:"error,omitempty"`
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

        if len(nums.Nums) == 0 {
            resp, _ := json.Marshal(Response{Error: "no numbers provided"})
            conn.Write(append(resp, '\n'))
            continue
        }

        tooLarge := false
        for _, n := range nums.Nums {
            if n > 1000 {
                tooLarge = true
                break
            }
        }

        if tooLarge {
            resp, _ := json.Marshal(Response{Error: "number too large"})
            conn.Write(append(resp, '\n'))
            continue
        }

        sum := 0
        for _, n := range nums.Nums {
            sum += n * n
        }

        resp, _ := json.Marshal(Response{
            Sum:      sum,
            Original: nums.Nums,
        })
        conn.Write(append(resp, '\n'))
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