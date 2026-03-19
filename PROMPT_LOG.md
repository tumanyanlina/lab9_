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
### Промпт 2
**Инструмент:** Cursor
**Дата:** 19.03.2026 11:00

**Промпт:**
"Улучши код TCP сервера:
1. Добавь структуру Response с полями Sum и Original
2. Если массив чисел пустой, возвращай {"error": "no numbers provided"}
3. Если число больше 1000, возвращай {"error": "number too large"}
4. Возвращай original массив в ответе"

**Результат:**
Сервер теперь возвращает original массив и обрабатывает ошибки. Но обнаружилась проблема: при числе >1000 сервер отправляет ошибку, но всё равно считает сумму.

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
        for _, n := range nums.Nums {
            if n > 1000 {
                resp, _ := json.Marshal(Response{Error: "number too large"})
                conn.Write(append(resp, '\n'))
                goto NEXT
            }
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
    NEXT:
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
### Промпт 3
**Инструмент:** Cursor
**Дата:** 19.03.2026

**Промпт:**
Улучши код TCP сервера на Go:

1. Убери использование goto.
2. Сделай чистую обработку ошибок:
   - если число > 1000 → сразу отправлять ошибку и переходить к следующему запросу
   - не закрывай соединение при ошибке — просто continue
3. Код должен быть читаемым и без дублирования
4. Сохрани всю функциональность:
   - JSON форматы {“numbers”:[1,2,3]} и {“sum”:14,“original”:[1,2,3]}
   - Обработка пустого массива → error
   - Обработка чисел >1000 → error
   - Обработка некорректного JSON → error

**Результат:**
Получила улучшенный код без goto. Теперь обработка ошибок чистая, соединение не закрывается при ошибках, код стал более читаемым.

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