### Description

This code is a minimal and runnable example of establishing a TCP server connection in Go. The application listens for incoming connections on :3000 and prints out data received from the client. <br />
Additionally, it maintains a dynamic list of connected peers and displays their unique IDs along with their respective network addresses.

### Quick Start

Start server
```
go run main.go
```

Make a connection
```
telnet localhost 3000
```

Type your messages
```
Hello, world!
```

See it on server

```
...
...
Message from %conn.RemoteAddr()%: Hello, world!
```

<br />
References:

https://okanexe.medium.com/the-complete-guide-to-tcp-ip-connections-in-golang-1216dae27b5a