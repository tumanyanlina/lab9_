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
### Промпт 4
**Инструмент:** Cursor
**Дата:** 19.03.2026

**Промпт:**
Рефакторинг TCP сервера на Go:

1. Вынеси всю бизнес-логику в отдельную функцию:
   func processNumbers(nums []int) Response

2. Внутри функции:
   - если массив пуст → вернуть Response с ошибкой "no numbers provided"
   - если число > 1000 → вернуть Response с ошибкой "number too large"
   - иначе вернуть сумму квадратов и original массив

3. В handleConnection:
   - убрать дублирующую логику
   - просто вызывать processNumbers(nums.Nums)

**Результат:**
Получила чистый код с отдельной функцией processNumbers. Теперь:
- Бизнес-логика изолирована от TCP-обработки
- Убрано дублирование кода
- Функцию processNumbers можно легко тестировать без запуска сервера
- handleConnection стал проще и понятнее

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
### Промпт 5
**Инструмент:** Cursor
**Дата:** 19.03.2026

**Промпт:**
Напиши unit-тесты для функции:
func processNumbers(nums []int) Response

Требования:
1. Используй пакет testing
2. Создай файл server_test.go
3. Покрой случаи:
   - нормальный ввод [1,2,3]
   - пустой массив []
   - число больше 1000
4. Проверяй:
   - правильную сумму
   - наличие ошибки
   - original массив

**Результат:**
Получила файл server_test.go с тремя тестами. Все тесты проходят.

**Код (server_test.go):**
```go
package main

import (
	"reflect"
	"testing"
)

func TestProcessNumbers_NormalInput(t *testing.T) {
	input := []int{1, 2, 3}
	wantSum := 1*1 + 2*2 + 3*3 // = 14

	resp := processNumbers(input)
	if resp.Error != "" {
		t.Errorf("unexpected error: %v", resp.Error)
	}
	if resp.Sum != wantSum {
		t.Errorf("got sum %v, want %v", resp.Sum, wantSum)
	}
	if !reflect.DeepEqual(resp.Original, input) {
		t.Errorf("got original %v, want %v", resp.Original, input)
	}
}

func TestProcessNumbers_EmptyArray(t *testing.T) {
	input := []int{}
	resp := processNumbers(input)
	if resp.Error != "no numbers provided" {
		t.Errorf("got error %v, want 'no numbers provided'", resp.Error)
	}
	if resp.Sum != 0 {
		t.Errorf("got sum %v, want 0", resp.Sum)
	}
	if resp.Original != nil && len(resp.Original) > 0 {
		t.Errorf("got original %v, want nil or empty", resp.Original)
	}
}

func TestProcessNumbers_NumberTooLarge(t *testing.T) {
	input := []int{10, 2000, 2}
	resp := processNumbers(input)
	if resp.Error != "number too large" {
		t.Errorf("got error %v, want 'number too large'", resp.Error)
	}
	if resp.Sum != 0 {
		t.Errorf("got sum %v, want 0", resp.Sum)
	}
	// For errors, the instruction does not specify about original, but let's check it is not filled
	if resp.Original != nil && len(resp.Original) > 0 {
		t.Errorf("got original %v, want nil or empty", resp.Original)
	}
}
### Промпт 6
**Инструмент:** Cursor
**Дата:** 19.03.2026

**Промпт:**
Улучши код TCP сервера:
1. Не игнорируй ошибки json.Marshal (убери _)
2. Если marshal падает — выведи ошибку через fmt.Println
3. Добавь логирование:
   - вывод входящего запроса
   - вывод ответа сервера

**Результат:**
Исправила игнорирование ошибок marshal. Добавила логирование запросов и ответов.

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

		fmt.Printf("Incoming request: %s", line)

		var nums Numbers
		if err := json.Unmarshal(line, &nums); err != nil {
			resp, marshalErr := json.Marshal(Response{Error: "invalid json"})
			if marshalErr != nil {
				fmt.Println("json.Marshal error:", marshalErr)
				return
			}
			fmt.Printf("Server response: %s\n", resp)
			conn.Write(append(resp, '\n'))
			continue
		}

		resp := processNumbers(nums.Nums)
		output, marshalErr := json.Marshal(resp)
		if marshalErr != nil {
			fmt.Println("json.Marshal error:", marshalErr)
			return
		}
		fmt.Printf("Server response: %s\n", output)
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
### Промпт 7
**Инструмент:** Cursor
**Дата:** 19.03.2026

**Промпт:**
Добавь в TCP сервер:
1. Graceful shutdown (обработка Ctrl+C)
2. Таймауты на чтение (5 мин) и запись (30 сек)
3. Проверку ошибок conn.Write
4. recover от паник в handleConnection

**Результат:**
Добавила graceful shutdown с контекстом и WaitGroup, таймауты на операции с соединением, проверку ошибок записи и защиту от паник.

**Код (server.go):**
```go
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