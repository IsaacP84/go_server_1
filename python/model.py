import socket, time, pickle, struct
from dataclasses import dataclass



HOST = '127.0.0.1'  # The server's hostname or IP address
PORT = 27015        # The port used by the server

# @dataclass
# class JoinPacket:
#     name: str='default'

with socket.socket(socket.AF_INET, socket.SOCK_DGRAM) as s:
    # s.connect((HOST, PORT))
    
    # s.sendto(b"Hello, server!", (HOST, PORT))

    data = [b'default']
    msg = bytearray(b'JN  ')
    msg += struct.pack('16s', data[0])
    # msg += pickle.dumps(data)

    s.sendto(msg, (HOST, PORT))

    address, port = s.getsockname()
    print(f"UDP socket bound to address: {address}, port: {port}")
    
    while True:
        data, addr = s.recvfrom(512)
        if len(data) != 0:
            print(addr)
            print(f"Received from server: {data.decode()}")

