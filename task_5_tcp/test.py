import socket
import json

client = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
client.connect(('localhost', 8080))

data = {"numbers": [1, 2, 3]}
client.send((json.dumps(data) + '\n').encode())

response = client.recv(1024).decode()
print("Ответ:", response)

client.close()