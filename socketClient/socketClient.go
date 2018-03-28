package socketClient

import (
	"context"
	"fmt"
	"io"
	"net"
)

// Request will send a request to the socket and return the bytes sent back.
// It will close the socket at the end.
// It takes a timeout value which will be used to wait for the output of the
// socket. This is a read timeout.
func Request(ctx context.Context, socketFile, requestText string) ([]byte, error) {
	c, err := net.Dial("unix", socketFile)
	if err != nil {
		return []byte{}, err
	}
	defer c.Close()
	_, err = c.Write([]byte(requestText))

	type readData struct {
		b   []byte
		err error
	}
	readChan := make(chan *readData)
	go func() {
		buf := make([]byte, 5120)
		nbytes, err := c.Read(buf)
		// We limit the size of the requests here to 5120 bytes. Anything more is most likely
		// an error.
		if err == nil && nbytes == len(buf) {
			err = fmt.Errorf("Buffer at maximum capacity of %d, likely overstaturation", len(buf))
		}
		if err == io.EOF {
			// We get EOF on a closed socket. This is expected behaviour.
			// So reset the error.
			err = nil
		}

		readChan <- &readData{b: buf[0:nbytes], err: err}
	}()

	select {
	case <-ctx.Done():
		return []byte{}, fmt.Errorf("timed out waiting for respose from socket")
	case out := <-readChan:
		return out.b, out.err
	}
}
