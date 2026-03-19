import socket
import json

def test_server():
    client = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
    client.connect(('localhost', 8080))
    
    # Тест 1: обычный запрос
    data = {"numbers": [1, 2, 3]}
    client.send((json.dumps(data) + '\n').encode())
    response = client.recv(1024).decode()
    print("Обычный запрос:", response)
    
    # Тест 2: пустой массив
    data = {"numbers": []}
    client.send((json.dumps(data) + '\n').encode())
    response = client.recv(1024).decode()
    print("Пустой массив:", response)
    
    # Тест 3: число >1000
    data = {"numbers": [1001, 2]}
    client.send((json.dumps(data) + '\n').encode())
    response = client.recv(1024).decode()
    print("Число >1000:", response)
    
    client.close()

if __name__ == "__main__":
    test_server()