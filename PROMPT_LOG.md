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
### Промпт 8
**Инструмент:** Cursor
**Дата:** 19.03.2026

**Промпт:**
Создай Python клиент client.py для TCP сервера на 127.0.0.1:8080 с функцией calculate_squares(numbers), которая:

1. Отправляет список чисел на сервер в формате JSON и получает ответ.
2. Перед подключением выводит "Connecting to server...".
3. После отправки данных выводит "Sent:" и список чисел.
4. После получения ответа выводит "Received:" и словарь с результатом.
5. Обрабатывает ошибки соединения, таймаутов и некорректного JSON, выводя их через print("Error:", e).
6. Добавь общий except Exception для всех неожиданных ошибок с выводом print("Unexpected error:", e).
7. Используй with socket.create_connection(...) и recv до символа '\n'.
8. Добавь пример использования в if __name__ == "__main__":, где вызывается calculate_squares([2,3,4]).

**Результат:**
Создала Python клиент с отладочным выводом. Клиент успешно подключается к серверу, отправляет данные и получает ответ.

**Код (client.py):**
```python
import socket
import json

def calculate_squares(numbers, host="127.0.0.1", port=8080, timeout=5):
    """
    Отправляет список чисел на сервер в формате JSON и возвращает ответ сервера (dict).
    Обрабатывает ошибки соединения, таймауты и некорректный JSON.
    """
    try:
        print("Connecting to server...")
        data = json.dumps({"numbers": numbers}).encode("utf-8") + b'\n'
        response_data = b""
        with socket.create_connection((host, port), timeout=timeout) as sock:
            sock.sendall(data)
            print("Sent:", numbers)
            while True:
                chunk = sock.recv(4096)
                if not chunk:
                    break
                response_data += chunk
                if b'\n' in response_data:
                    break
        line = response_data.split(b'\n', 1)[0].decode("utf-8")
        result = json.loads(line)
        print("Received:", result)
        return result
    except (socket.timeout, socket.error) as e:
        print("Error:", e)
        return None
    except json.JSONDecodeError as e:
        print("Error:", e)
        return None
    except Exception as e:
        print("Unexpected error:", e)
        return None

if __name__ == "__main__":
    test = [2, 3, 4]
    result = calculate_squares(test)
    print(result)
### Промпт 9
**Инструмент:** Cursor
**Дата:** 19.03.2026

**Промпт:**
Создай файл test_client.py для Python TCP клиента (calculate_squares из client.py) с использованием unittest.  
Покрой следующие сценарии:

1. Нормальный массив чисел, например [1,2,3]: проверяй, что sum считается правильно (сумма квадратов), original совпадает с входным массивом, error отсутствует.
2. Пустой массив []: проверяй, что возвращается error "no numbers provided", sum и original равны None.
3. Массив с числом >1000, например [10,2000,2]: проверяй, что возвращается error "number too large", sum и original равны None.
4. В каждом тесте сначала проверяй, что результат не None (чтобы сервер ответил).
5. Добавь стандартный блок if __name__ == "__main__": для запуска unittest.main().
**Результат:**
Создала test_client.py с тремя тестами. Все тесты успешно проходят.

**Код (test_client.py):**
```python
import unittest
from client import calculate_squares

class TestCalculateSquares(unittest.TestCase):
    def test_normal_input(self):
        numbers = [1, 2, 3]
        result = calculate_squares(numbers)
        self.assertIsNotNone(result, "Result should not be None")
        self.assertEqual(result.get('sum'), 1*1 + 2*2 + 3*3)
        self.assertEqual(result.get('original'), numbers)
        self.assertFalse('error' in result and result['error'], f"Unexpected error: {result.get('error')}")

    def test_empty_array(self):
        numbers = []
        result = calculate_squares(numbers)
        self.assertIsNotNone(result, "Result should not be None")
        self.assertEqual(result.get('error'), "no numbers provided")

    def test_number_too_large(self):
        numbers = [10, 2000, 2]
        result = calculate_squares(numbers)
        self.assertIsNotNone(result, "Result should not be None")
        self.assertEqual(result.get('error'), "number too large")

if __name__ == "__main__":
    unittest.main()

## Задание 2: Добавить горутину для фоновой обработки запросов.

### Промпт 1
**Инструмент:** Cursor
**Дата:** 20.03.2026

**Промпт:**
Напиши HTTP сервер на Go с эндпоинтом POST /process.

Требования:
1. Сервер слушает порт 8081
2. Эндпоинт принимает JSON {"data": "текст"}
3. Сразу отвечает {"status": "accepted"}
4. Запускает горутину, которая:
   - ждет 2 секунды (time.Sleep)
   - выводит в консоль: "Processed: <текст>"
5. Добавь логирование запросов

**Результат:**
Создала HTTP сервер с горутиной для фоновой обработки. Сервер принимает POST запросы, сразу отвечает и обрабатывает данные в фоне.

**Код (http_server.go):**
```go
package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

type ProcessRequest struct {
	Data string `json:"data"`
}

type ProcessResponse struct {
	Status string `json:"status"`
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s from %s", r.Method, r.URL.Path, r.RemoteAddr)
		next.ServeHTTP(w, r)
	})
}

func processHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	var req ProcessRequest
	if err := json.Unmarshal(body, &req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// Respond immediately
	resp := ProcessResponse{Status: "accepted"}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)

	// Process in background
	go func(data string) {
		time.Sleep(2 * time.Second)
		fmt.Printf("Processed: %s\n", data)
	}(req.Data)
}

func main() {
	mux := http.NewServeMux()
	mux.Handle("/process", loggingMiddleware(http.HandlerFunc(processHandler)))
	log.Println("Server listening on :8081")
	log.Fatal(http.ListenAndServe(":8081", mux))
}
### Промпт 2.2
**Инструмент:** Cursor
**Дата:** 20.03.2026

**Промпт:**
Добавь в сервер graceful shutdown:

При нажатии Ctrl+C сервер должен:
   - перестать принимать новые запросы
   - дать текущим запросам 5 секунд на завершение
   - дождаться завершения фоновых горутин
   - вывести "Server stopped" и выйти

**Результат:**
Добавила graceful shutdown с использованием контекста и WaitGroup для отслеживания фоновых горутин. Сервер корректно завершает работу, дожидаясь обработки текущих запросов.

**Код (http_server.go):**
```go
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

type ProcessRequest struct {
	Data string `json:"data"`
}

type ProcessResponse struct {
	Status string `json:"status"`
}

var (
	bgWg sync.WaitGroup
)

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s from %s", r.Method, r.URL.Path, r.RemoteAddr)
		next.ServeHTTP(w, r)
	})
}

func processHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	var req ProcessRequest
	if err := json.Unmarshal(body, &req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// Respond immediately
	resp := ProcessResponse{Status: "accepted"}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)

	// Process in background and track completion
	bgWg.Add(1)
	go func(data string) {
		defer bgWg.Done()
		time.Sleep(2 * time.Second)
		fmt.Printf("Processed: %s\n", data)
	}(req.Data)
}

func main() {
	mux := http.NewServeMux()
	mux.Handle("/process", loggingMiddleware(http.HandlerFunc(processHandler)))
	server := &http.Server{
		Addr:    ":8081",
		Handler: mux,
	}

	// Signal handling for graceful shutdown
	idleConnsClosed := make(chan struct{})
	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt, syscall.SIGTERM)
		<-sigint

		log.Println("Shutting down server...")

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		// Stop accepting new connections, allow up to 5s for existing requests
		if err := server.Shutdown(ctx); err != nil {
			log.Printf("HTTP server Shutdown: %v", err)
		}

		// Wait for background goroutines to finish
		bgWg.Wait()

		close(idleConnsClosed)
	}()

	log.Println("Server listening on :8081")
	if err := server.ListenAndServe(); err != http.ErrServerClosed {
		log.Fatalf("HTTP server ListenAndServe: %v", err)
	}

	<-idleConnsClosed
	log.Println("Server stopped")
}
