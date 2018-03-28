# HTTP to SOCKET

## Why

Application don't always want to expose http endpoints. This makes it difficult to gather metrics from a healthcheck endpoint.
We want metrics, they are there for very good reasons. However putting in a HTTP endpoint can make a simple application quite tricky to handle.

Unix Sockets are generally easier to setup and to use. They are also faster as there is less overhead to deal with. HTTP requires TCP and lots of context switching.

So, all you need to do as a developer is listen on a unix socket. You tell httpToSock what to send when it gets a request on a path. It will return to the caller the bytes you send back down the socket.

## How

httpToSock has been written to be a side car container or background process. You pass it in the arguments that describe your path and socket requests. It then handles the TLS and HTTP for you. You will need to supply certificates if you are using TLS.

Flag | Default | Description
---|---|---
-h | false | Shows the help menu
-v | false | Shows the version
-l | "0.0.0.0" | IP address to listen on. "0.0.0.0" means all.
-p | 8080 | The port the HTTP server will listen on.
-tls | false | Enable TLS
-cert | "./cert.pem" | Location of the certificate file.
-key | "./cert.key" | Location of the certificate private key.
-socket-timeout | 3 | Seconds to wait till timing out while reading from the socket.
-f | "/var/run/httpToSock.sock" | Location of the socket that httpToSock should use to talk to your app.
-path | "path" | The paths and strings to send to the socket. See explaination below. Can be called multiple times.

### Path parameter

The path is the only tricky parameter. It is a mix of your http path and what we should send to your app down the socket. It is has 2 parts seperated by a ":" (Semicolon). Part 1 is the http path, part 2 is the string to send down the socket. The application will complain if there is more than 1 : in the path argument.

The easist way to see it is if we give an example.

```bash
./httpToSock -path "/healthz:GET /healthz" -path "/status:sendStatus"
```

We will get 2 http paths:

1. `/healthz` which will send `GET /healthz` down the socket.
1. `/status` which will send `sendStatus` down the socket.

## Example

```bash
./httpToSock \
    -f /tmp/testSocket.sock \
    -l "0.0.0.0" \
    -p 8080 \
    -tls \
    -cert "./cert.pem" \
    -key "./cert.key" \
    -path "/_status:GET /healthz"
```