import socket
import os

ports = []

for i in range(5):
    s = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
    s.bind(('', 0))
    addr = s.getsockname()
    ports.append(addr[1])
    s.close()

port_list = ""
for i in range(5):
    p = str(ports[i])
    os.system("./server_node.sh " + p + " &")
    port_list += (p + " ")

os.system("./client.sh " + port_list)
