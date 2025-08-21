import socket, time

HOST = '127.0.0.1'  # The server's hostname or IP address
PORT = 27015        # The port used by the server

with socket.socket(socket.AF_INET, socket.SOCK_DGRAM) as s:
    # s.connect((HOST, PORT))
    s.sendto(b"Hello, server!", (HOST, PORT))
    while True:
        data, addr = s.recvfrom(512)
        if len(data) != 0:
            print(addr)
            print(f"Received from server: {data.decode()}")

