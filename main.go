package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/morfien101/httpToSock/httpServer"
	"github.com/morfien101/httpToSock/pathFlag"
	"github.com/morfien101/httpToSock/socketClient"
)

var (
	paths              pathFlag.PF
	helpFlag           = flag.Bool("h", false, "Shows help menu.")
	versionFlag        = flag.Bool("v", false, "Shows the version.")
	tlsEnabled         = flag.Bool("tls", false, "Enabled TLS support.")
	cert               = flag.String("cert", "./cert.pem", "TLS certificate.")
	key                = flag.String("key", "./cert.key", "TLS private key.")
	httpListenAddress  = flag.String("l", "0.0.0.0", "IP address to listen on.")
	httpListenPort     = flag.Int("p", 8080, "TCP port for HTTP(S) Server")
	socketTimeout      = flag.Int("socket-timeout", 3, "How long should we wait before timing out on the socket read.")
	socketFileLocation = flag.String("f", "/var/run/httpToSock.sock", "The location of the socket file.")
)

func main() {
	flag.Var(&paths, "path", "The paths to proxy")
	flag.Parse()

	if *helpFlag {
		flag.PrintDefaults()
		return
	}

	if *versionFlag {
		fmt.Println(version)
		return
	}

	pathMap, err := paths.Split()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	httpConfig := &httpServer.ServerConfig{
		TLS:           *tlsEnabled,
		Cert:          *cert,
		Key:           *key,
		ListenAddress: fmt.Sprintf("%s:%d", *httpListenAddress, *httpListenPort),
		Routes:        generateRoutes(pathMap),
	}

	httpsrv := httpServer.NewServer(httpConfig)

	c := make(chan os.Signal, 2)
	// We till the signals package where to send the signals. aka the channel we just made.
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	errChan := make(chan error, 1)
	go func() {
		errChan <- httpsrv.Start()
	}()

	select {
	case <-c:
		fmt.Println("Got stop signal. Stopting the HTTP Server...")
		err := httpsrv.Stop(5)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	case err := <-errChan:
		fmt.Println(err)
		os.Exit(1)
	}
}

func generateRoutes(pathMap map[string]string) map[string]func() http.Handler {
	routeMap := make(map[string]func() http.Handler)
	for k, v := range pathMap {
		routeMap[k] = func() http.Handler {
			return http.HandlerFunc(
				func(w http.ResponseWriter, r *http.Request) {
					ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(*socketTimeout))
					defer cancel()

					b, err := socketClient.Request(ctx, *socketFileLocation, v)
					if err != nil {
						fmt.Println(err)
						w.WriteHeader(http.StatusInternalServerError)
						return
					}
					if len(b) == 0 {
						w.WriteHeader(http.StatusInternalServerError)
						msg := "no response from socket"
						fmt.Fprint(w, msg)
						fmt.Println(msg)
						return
					}

					w.Header().Set("content-type", "application/json")
					w.Write(b)
				},
			)
		}
	}
	return routeMap
}
