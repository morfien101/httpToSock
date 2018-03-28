package socketClient

import (
	"context"
	"fmt"
	"math/rand"
	"net"
	"os"
	"testing"
	"time"
)

func randString(length int64) string {
	rand.Seed(time.Now().UnixNano())
	s := make([]byte, length)

	for i := 0; i < len(s); i++ {
		s[i] = byte(rand.Intn(26) + 65)
	}

	return string(s)
}

type server struct {
	srv          net.Listener
	socketFile   string
	sleepSeconds int
	errorChan    chan error
}

func (s *server) start() {
	go func() {
		for {
			conn, err := s.srv.Accept()
			if err != nil {
				s.errorChan <- fmt.Errorf("Failed to accept socket connection: %s", err)
				return
			}
			go s.handleConn(conn)
		}
	}()
}

func (s *server) handleConn(conn net.Conn) {
	defer conn.Close()
	buf := make([]byte, 5120)
	numberOfBytes, err := conn.Read(buf)
	if err != nil {
		s.errorChan <- fmt.Errorf("Failed to read on the server side of the socket. Error: %s", err)
		return
	}
	if s.sleepSeconds > 0 {
		time.Sleep(time.Second * time.Duration(s.sleepSeconds))
	}
	_, err = conn.Write(buf[0:numberOfBytes])
	if err != nil {
		s.errorChan <- fmt.Errorf("Failed to write on the server side of the socket. Error: %s", err)
		return
	}
}

func (s *server) stop() {
	s.srv.Close()
	os.Remove(s.socketFile)
}

func newSocketServer(sleepSeconds int, t *testing.T) (*server, error) {
	s := &server{
		socketFile: fmt.Sprintf("./%s.sock", randString(20)),
		errorChan:  make(chan error, 1),
	}
	l, err := net.Listen("unix", s.socketFile)
	if err != nil {
		return &server{}, err
	}
	s.srv = l
	return s, nil
}

func TestGoodRequest(t *testing.T) {
	srv, err := newSocketServer(0, t)
	if err != nil {
		t.Fatalf("Failed to create a socket. Error: %s", err)
	}
	srv.start()
	defer srv.stop()

	requestString := "GoodTesting"
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	b, err := Request(ctx, srv.socketFile, requestString)
	if err != nil {
		t.Fatalf("Request returned an error. Error: %s", err)
	}
	if requestString != string(b) {
		t.Errorf("Socket did not return expected result.\nGot: %s\nWant: %s\n", string(b), requestString)
	}
}

func TestTimoutRequest(t *testing.T) {
	srv, err := newSocketServer(1, t)
	if err != nil {
		t.Fatalf("Failed to create a socket. Error: %s", err)
	}
	defer srv.stop()

	requestString := "GoodTesting"
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*200)
	defer cancel()

	_, err = Request(ctx, srv.socketFile, requestString)
	if err == nil {
		t.Fatal("Request did not return an error on a timeout\n")
	}
}

func TestOversizeRequest(t *testing.T) {
	srv, err := newSocketServer(0, t)
	if err != nil {
		t.Fatalf("Failed to create a socket. Error: %s", err)
	}
	defer srv.stop()

	requestString := randString(5200)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
	defer cancel()

	_, err = Request(ctx, srv.socketFile, string(requestString))
	if err == nil {
		t.Fatal("Request did not return an error on a oversized return\n")
	}
}

func TestBadSocketFile(t *testing.T) {
	_, err := Request(nil, "", "No Socket File")
	if err == nil {
		t.Error("Sending an empyt string for the socketfile did not throw an error")
	}
}
