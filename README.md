### Description

This code is a minimal and runnable example of establishing a TCP server connection in Go. <br/>
The application listens for incoming connections on :3000 and prints out data received from the client. <br />
Additionally, it maintains a dynamic list of connected peers and displays their unique IDs along with their respective network addresses.

### Quick Start

Start server
```console
go run main.go
```

Spawn as much peers as you want
```console
telnet localhost 3000
```

Send your bytes to the server
```console
Hello, world!
```

See your messages
```console
...
...
Message from %conn.RemoteAddr()%: Hello, world!
```

Receive "ping-pong"
```console
Thank you for your message!
```

### References

https://okanexe.medium.com/the-complete-guide-to-tcp-ip-connections-in-golang-1216dae27b5a
