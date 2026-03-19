# Prompt Log
## Лабораторная работа №9
**Студент:** Туманян Лина
**Вариант:** 7
**Дата:** 19.03.2026

---

## Задание 5: TCP сервер на Go + Python клиент

### Промпт 1
**Инструмент:** Cursor
**Дата:** 19.03.2026 10:00

**Промпт:**
"Напиши простой TCP сервер на Go, который принимает JSON с числами и возвращает сумму квадратов. Используй goroutine для обработки клиентов."

**Результат:**
Получила базовый TCP сервер. Сервер слушает порт 8080, принимает JSON вида `{"numbers":[1,2,3]}`, возвращает `{"sum":14}`. Проверила с помощью Python скрипта.

**Код (server.go):**
```go
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