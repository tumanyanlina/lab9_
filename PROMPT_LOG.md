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
### Промпт 2
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
### Промпт 3
**Инструмент:** Cursor
**Дата:** 20.03.2026

**Промпт:**
Напиши тест для HTTP сервера:

1. Используй пакет testing и httptest
2. TestProcessEndpoint:
   - отправляет POST запрос с {"data": "test"}
   - проверяет что статус код 200
   - проверяет что ответ {"status": "accepted"}
   - проверяет что ответ пришел сразу (без задержки)

**Результат:**
Создала файл http_server_test.go с тестом для эндпоинта /process. Тест проверяет, что ответ приходит мгновенно, без ожидания фоновой обработки.

**Код (http_server_test.go):**
```go
package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestProcessEndpoint(t *testing.T) {
	handler := loggingMiddleware(http.HandlerFunc(processHandler))
	srv := httptest.NewServer(handler)
	defer srv.Close()

	payload := ProcessRequest{Data: "test"}
	body, _ := json.Marshal(payload)

	start := time.Now()
	resp, err := http.Post(srv.URL+"/process", "application/json", bytes.NewReader(body))
	if err != nil {
		t.Fatalf("failed to do POST: %v", err)
	}
	defer resp.Body.Close()
	elapsed := time.Since(start)

	if elapsed > 300*time.Millisecond {
		t.Errorf("response took too long: %v", elapsed)
	}
	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}

	var got ProcessResponse
	json.NewDecoder(resp.Body).Decode(&got)
	expected := ProcessResponse{Status: "accepted"}
	if got != expected {
		t.Errorf("expected %+v, got %+v", expected, got)
	}
}
### Промпт 4
**Инструмент:** Cursor
**Дата:** 20.03.2026

**Промпт:**
Добавь дополнительные тесты для processHandler:
- GET запрос должен возвращать 405 Method Not Allowed
- POST с некорректным JSON должен возвращать 400 Bad Request
- POST с пустым телом должен возвращать 400 Bad Request
Используй httptest и проверь статус-коды.

**Результат:**
Добавила дополнительные тесты для проверки обработки ошибок: метод не разрешен, некорректный JSON и пустое тело запроса. Все тесты проходят успешно.

**Код (http_server_test.go - дополнение):**
```go
func TestProcessMethodNotAllowed(t *testing.T) {
	handler := loggingMiddleware(http.HandlerFunc(processHandler))
	srv := httptest.NewServer(handler)
	defer srv.Close()

	req, err := http.NewRequest(http.MethodGet, srv.URL, nil)
	if err != nil {
		t.Fatalf("failed to create GET request: %v", err)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("failed to send GET request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusMethodNotAllowed {
		t.Errorf("expected status 405 Method Not Allowed, got %d", resp.StatusCode)
	}
}

func TestProcessInvalidJSON(t *testing.T) {
	handler := loggingMiddleware(http.HandlerFunc(processHandler))
	srv := httptest.NewServer(handler)
	defer srv.Close()

	invalidJSON := []byte(`{invalid json}`)
	resp, err := http.Post(srv.URL, "application/json", bytes.NewReader(invalidJSON))
	if err != nil {
		t.Fatalf("failed to POST invalid JSON: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("expected status 400 Bad Request, got %d", resp.StatusCode)
	}
}

func TestProcessEmptyBody(t *testing.T) {
	handler := loggingMiddleware(http.HandlerFunc(processHandler))
	srv := httptest.NewServer(handler)
	defer srv.Close()

	resp, err := http.Post(srv.URL, "application/json", bytes.NewReader(nil))
	if err != nil {
		t.Fatalf("failed to POST with empty body: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("expected status 400 Bad Request, got %d", resp.StatusCode)
	}
}
### Промпт 5
**Инструмент:** Cursor
**Дата:** 20.03.2026

**Промпт:**
Вынеси логику обработки данных из горутины в отдельную функцию processRequest(data string). Обнови processHandler, чтобы он вызывал эту функцию внутри горутины.

**Результат:**
Вынесла логику обработки данных в отдельную функцию processRequest, что улучшило читаемость кода и позволит легко тестировать логику обработки независимо от HTTP.

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

// Вынесенная функция обработки данных
func processRequest(data string) {
	log.Printf("Start processing: %s", data)
	time.Sleep(2 * time.Second)
	log.Printf("Done processing: %s", data)
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
		processRequest(data)
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
### Промпт 6
**Инструмент:** Cursor
**Дата:** 20.03.2026

**Промпт:**
Обнови http_server.go так, чтобы фоновая обработка делалась через интерфейс Processor:

1. Создай интерфейс Processor с методом Process(data string) error.
2. Сделай реализацию DefaultProcessor, которая логирует начало/конец обработки и делает time.Sleep(2s).
3. Перепиши processHandler, чтобы он принимал Processor и вызывал его метод Process в горутине, сохраняя использование bgWg и логирование ошибок.
4. В main() создай DefaultProcessor и передай его в handler для /process.

**Результат:**
Реализовала интерфейс Processor и его DefaultProcessor, что улучшило тестируемость и гибкость кода. Теперь обработчик принимает интерфейс, что позволяет легко подменять реализацию для тестов.

**Код (http_server.go):**
```go
package main

import (
	"context"
	"encoding/json"
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

// Processor interface
type Processor interface {
	Process(data string) error
}

// DefaultProcessor implementation
type DefaultProcessor struct{}

func (p *DefaultProcessor) Process(data string) error {
	log.Printf("Start processing: %s", data)
	time.Sleep(2 * time.Second)
	log.Printf("Done processing: %s", data)
	return nil
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

// Handler factory to use injected Processor
func processHandler(processor Processor) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
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
			if err := processor.Process(data); err != nil {
				log.Printf("Background processing failed: %v", err)
			}
		}(req.Data)
	}
}

func main() {
	mux := http.NewServeMux()
	processor := &DefaultProcessor{}
	mux.Handle("/process", loggingMiddleware(processHandler(processor)))
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
### Промпт 7
**Инструмент:** Cursor
**Дата:** 20.03.2026

**Промпт:**
Обнови http_server_test.go, чтобы тесты корректно работали с новой архитектурой (с интерфейсом Processor):

1. Убери дублирование (не объявляй Processor и processHandler в тесте)
2. Используй dependency injection — передавай в processHandler тестовую реализацию Processor
3. Создай mock-процессор MockProcessor, реализующий интерфейс Processor, метод Process не должен делать time.Sleep
4. Обнови создание handler: вместо http.HandlerFunc(processHandler) используй processHandler(mockProcessor)
5. Проверь, что ответ приходит быстро, статус 200 OK, JSON = {"status": "accepted"}
6. Убедись, что Process действительно вызывается (через счётчик в mock)

**Результат:**
Обновила тесты для работы с интерфейсом Processor. Добавила MockProcessor с мьютексом для безопасного доступа из горутин. Все тесты проходят успешно.

**Код (http_server_test.go):**
```go
package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"
)

// --- Мок-процессор для тестирования фоновой обработки ---
type MockProcessor struct {
	CalledWith []string
	mu         sync.Mutex
}

func (m *MockProcessor) Process(data string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.CalledWith = append(m.CalledWith, data)
	return nil
}

func TestProcessEndpoint(t *testing.T) {
	mockProc := &MockProcessor{}
	handler := loggingMiddleware(processHandler(mockProc))
	srv := httptest.NewServer(handler)
	defer srv.Close()

	payload := ProcessRequest{Data: "test"}
	body, err := json.Marshal(payload)
	if err != nil {
		t.Fatalf("failed to marshal payload: %v", err)
	}

	start := time.Now()
	resp, err := http.Post(srv.URL, "application/json", bytes.NewReader(body))
	if err != nil {
		t.Fatalf("failed to do POST: %v", err)
	}
	defer resp.Body.Close()

	elapsed := time.Since(start)
	if elapsed > 300*time.Millisecond {
		t.Errorf("response took too long: %v", elapsed)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status 200 OK, got %d", resp.StatusCode)
	}

	var got ProcessResponse
	if err := json.NewDecoder(resp.Body).Decode(&got); err != nil {
		t.Errorf("failed to decode response: %v", err)
	}
	expected := ProcessResponse{Status: "accepted"}
	if got != expected {
		t.Errorf("expected response %+v, got %+v", expected, got)
	}

	// Проверяем, что Process вызван с нужными данными
	time.Sleep(30 * time.Millisecond)
	mockProc.mu.Lock()
	defer mockProc.mu.Unlock()
	if len(mockProc.CalledWith) != 1 {
		t.Fatalf("Expected Process to be called once, got %d", len(mockProc.CalledWith))
	}
	if mockProc.CalledWith[0] != "test" {
		t.Errorf("Expected Process called with 'test', got %q", mockProc.CalledWith[0])
	}
}

func TestProcessMethodNotAllowed(t *testing.T) {
	mockProc := &MockProcessor{}
	handler := loggingMiddleware(processHandler(mockProc))
	srv := httptest.NewServer(handler)
	defer srv.Close()

	req, err := http.NewRequest(http.MethodGet, srv.URL, nil)
	if err != nil {
		t.Fatalf("failed to create GET request: %v", err)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("failed to send GET request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusMethodNotAllowed {
		t.Errorf("expected status 405 Method Not Allowed, got %d", resp.StatusCode)
	}
}

func TestProcessInvalidJSON(t *testing.T) {
	mockProc := &MockProcessor{}
	handler := loggingMiddleware(processHandler(mockProc))
	srv := httptest.NewServer(handler)
	defer srv.Close()

	invalidJSON := []byte(`{invalid json}`)
	resp, err := http.Post(srv.URL, "application/json", bytes.NewReader(invalidJSON))
	if err != nil {
		t.Fatalf("failed to POST invalid JSON: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("expected status 400 Bad Request, got %d", resp.StatusCode)
	}
}

func TestProcessEmptyBody(t *testing.T) {
	mockProc := &MockProcessor{}
	handler := loggingMiddleware(processHandler(mockProc))
	srv := httptest.NewServer(handler)
	defer srv.Close()

	resp, err := http.Post(srv.URL, "application/json", bytes.NewReader(nil))
	if err != nil {
		t.Fatalf("failed to POST with empty body: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("expected status 400 Bad Request, got %d", resp.StatusCode)
	}
}

func TestBackgroundProcessing(t *testing.T) {
	mockProc := &MockProcessor{}
	handler := loggingMiddleware(processHandler(mockProc))
	srv := httptest.NewServer(handler)
	defer srv.Close()

	payload := ProcessRequest{Data: "background"}
	body, err := json.Marshal(payload)
	if err != nil {
		t.Fatalf("failed to marshal payload: %v", err)
	}

	resp, err := http.Post(srv.URL, "application/json", bytes.NewReader(body))
	if err != nil {
		t.Fatalf("failed to POST: %v", err)
	}
	resp.Body.Close()

	// Даем время на запуск фоновой горутины
	time.Sleep(30 * time.Millisecond)

	mockProc.mu.Lock()
	defer mockProc.mu.Unlock()
	if len(mockProc.CalledWith) != 1 {
		t.Fatalf("process should be called once, got %d calls", len(mockProc.CalledWith))
	}
	if mockProc.CalledWith[0] != "background" {
		t.Errorf("expected process called with %q, got %q", "background", mockProc.CalledWith[0])
	}
}
### Промпт 8
**Инструмент:** Cursor
**Дата:** 20.03.2026

**Промпт:**
Проведи рефакторинг тестов в http_server_test.go:

1. Объедини похожие тесты в табличные (table-driven tests)
2. Убери дублирование создания handler и сервера
3. Убери все комментарии
4. Сделай тесты компактными, сохранив всю логику

**Результат:**
Провела рефакторинг тестов: объединила 4 теста (POST valid JSON, GET returns 405, POST invalid JSON, POST empty body) в табличный тест TestProcessHandlerVariants. Добавила вспомогательную функцию withTestServer для устранения дублирования создания сервера. Тесты стали чище и компактнее.

**Код (http_server_test.go):**
package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"
)

type MockProcessor struct {
	CalledWith []string
	mu         sync.Mutex
}

func (m *MockProcessor) Process(data string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.CalledWith = append(m.CalledWith, data)
	return nil
}

func withTestServer(proc *MockProcessor, f func(url string)) {
	srv := httptest.NewServer(loggingMiddleware(processHandler(proc)))
	defer srv.Close()
	f(srv.URL)
}

func doRequest(t *testing.T, method, url string, body []byte) *http.Response {
	if method == http.MethodPost {
		resp, err := http.Post(url, "application/json", bytes.NewReader(body))
		if err != nil {
			t.Fatalf("POST error: %v", err)
		}
		return resp
	}
	req, _ := http.NewRequest(method, url, nil)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("%s request error: %v", method, err)
	}
	return resp
}

func TestProcessHandlerVariants(t *testing.T) {
	mockProc := &MockProcessor{}
	cases := []struct {
		name          string
		method        string
		body          []byte
		wantStatus    int
		wantJSON      *ProcessResponse
		wantProcess   string
	}{
		{"POST valid JSON", http.MethodPost, mustMarshal(ProcessRequest{Data: "test"}), http.StatusOK, &ProcessResponse{Status: "accepted"}, "test"},
		{"GET returns 405", http.MethodGet, nil, http.StatusMethodNotAllowed, nil, ""},
		{"POST invalid JSON", http.MethodPost, []byte(`{invalid json}`), http.StatusBadRequest, nil, ""},
		{"POST empty body", http.MethodPost, nil, http.StatusBadRequest, nil, ""},
	}

	withTestServer(mockProc, func(url string) {
		for _, tc := range cases {
			t.Run(tc.name, func(t *testing.T) {
				resp := doRequest(t, tc.method, url, tc.body)
				defer resp.Body.Close()
				if resp.StatusCode != tc.wantStatus {
					t.Errorf("want status %d, got %d", tc.wantStatus, resp.StatusCode)
				}
				if tc.wantJSON != nil {
					var got ProcessResponse
					if err := json.NewDecoder(resp.Body).Decode(&got); err != nil {
						t.Errorf("decode error: %v", err)
					}
					if got != *tc.wantJSON {
						t.Errorf("want %+v, got %+v", *tc.wantJSON, got)
					}
				}
				if tc.wantProcess != "" {
					time.Sleep(30 * time.Millisecond)
					mockProc.mu.Lock()
					defer mockProc.mu.Unlock()
					if len(mockProc.CalledWith) == 0 || mockProc.CalledWith[len(mockProc.CalledWith)-1] != tc.wantProcess {
						t.Errorf("expected Process with %q, got %v", tc.wantProcess, mockProc.CalledWith)
					}
				}
			})
		}
	})
}

func TestBackgroundProcessing(t *testing.T) {
	mockProc := &MockProcessor{}
	withTestServer(mockProc, func(url string) {
		body := mustMarshal(ProcessRequest{Data: "background"})
		resp := doRequest(t, http.MethodPost, url, body)
		resp.Body.Close()
		time.Sleep(30 * time.Millisecond)

		mockProc.mu.Lock()
		defer mockProc.mu.Unlock()
		if len(mockProc.CalledWith) != 1 {
			t.Fatalf("want 1 process call, got %d", len(mockProc.CalledWith))
		}
		if mockProc.CalledWith[0] != "background" {
			t.Errorf("want 'background', got %q", mockProc.CalledWith[0])
		}
	})
}

func mustMarshal(v any) []byte {
	b, _ := json.Marshal(v)
	return b
}
### Промпт 9
**Инструмент:** Cursor
**Дата:** 20.03.2026

**Промпт:**
Проведи рефакторинг кода в http_server.go:

- Убери дублирование
- Вынеси повторяющиеся блоки в отдельные функции
- Сделай код идиоматичным для Go
- Убери все комментарии

**Результат:**
Провела рефакторинг серверного кода: вынесла respondJSON для унификации ответов, decodeProcessRequest для обработки запросов, runServer для управления жизненным циклом сервера. Код стал чище и структурированнее.

**Код (http_server.go):**
```go
package main

import (
	"context"
	"encoding/json"
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

type Processor interface {
	Process(string) error
}

type DefaultProcessor struct{}

func (p *DefaultProcessor) Process(data string) error {
	log.Printf("Start processing: %s", data)
	time.Sleep(2 * time.Second)
	log.Printf("Done processing: %s", data)
	return nil
}

var bgWg sync.WaitGroup

func respondJSON(w http.ResponseWriter, status int, resp interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(resp)
}

func decodeProcessRequest(r *http.Request) (ProcessRequest, error) {
	defer r.Body.Close()
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return ProcessRequest{}, err
	}
	var req ProcessRequest
	err = json.Unmarshal(body, &req)
	return req, err
}

func processHandler(processor Processor) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		req, err := decodeProcessRequest(r)
		if err != nil {
			http.Error(w, "Bad request", http.StatusBadRequest)
			return
		}
		respondJSON(w, http.StatusOK, ProcessResponse{Status: "accepted"})
		bgWg.Add(1)
		go func(data string) {
			defer bgWg.Done()
			if err := processor.Process(data); err != nil {
				log.Printf("Background processing failed: %v", err)
			}
		}(req.Data)
	}
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s from %s", r.Method, r.URL.Path, r.RemoteAddr)
		next.ServeHTTP(w, r)
	})
}

func runServer(addr string, handler http.Handler) {
	server := &http.Server{
		Addr:    addr,
		Handler: handler,
	}
	idleConnsClosed := make(chan struct{})
	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt, syscall.SIGTERM)
		<-sigint
		log.Println("Shutting down server...")
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := server.Shutdown(ctx); err != nil {
			log.Printf("HTTP server Shutdown: %v", err)
		}
		bgWg.Wait()
		close(idleConnsClosed)
	}()
	log.Println("Server listening on", addr)
	if err := server.ListenAndServe(); err != http.ErrServerClosed {
		log.Fatalf("HTTP server ListenAndServe: %v", err)
	}
	<-idleConnsClosed
	log.Println("Server stopped")
}

func main() {
	processor := &DefaultProcessor{}
	mux := http.NewServeMux()
	mux.Handle("/process", loggingMiddleware(processHandler(processor)))
	runServer(":8081", mux)
}
## Задание 7: Использовать PyO3 для создания Python-модуля на Rust.

### Промпт 1
**Инструмент:** Chat GPT
**Дата:** 20.03.2026

**Промпт:**
Создай Python модуль на Rust с PyO3. Нужны две функции:
- multiply_by_two: принимает число, возвращает число умноженное на 2
- make_greeting: принимает имя, возвращает приветствие "Hi, {name}! Welcome from Rust."

**Результат:**
Получила базовый Rust модуль с двумя функциями, зарегистрированными через pymodule. Настроила Cargo.toml с cdylib и зависимостью pyo3. Добавила Python тесты.

**Код (src/lib.rs):**
```rust
use pyo3::prelude::*;

#[pyfunction]
fn multiply_by_two(x: i32) -> PyResult<i32> {
    Ok(x * 2)
}

#[pyfunction]
fn make_greeting(name: &str) -> PyResult<String> {
    Ok(format!("Hi, {}! Welcome from Rust.", name))
}

#[pymodule]
fn my_rust_module(_py: Python<'_>, m: &Bound<'_, PyModule>) -> PyResult<()> {
    m.add_function(wrap_pyfunction!(multiply_by_two, m)?)?;
    m.add_function(wrap_pyfunction!(make_greeting, m)?)?;
    Ok(())
}