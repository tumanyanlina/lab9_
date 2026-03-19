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