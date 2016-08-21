package gfile

import (
	"io"
	"os"
	"time"

	"github.com/onsi/gomega/gbytes"
)

//Buffer ...
type Buffer struct {
	buffer   *gbytes.Buffer
	stopChan chan bool
	file     *os.File
	closed   bool
}

//NewBuffer returns a *gbytes.Buffer over the file at `path`
func NewBuffer(path string) (*Buffer, error) {
	buffer := new(Buffer)

	var err error
	buffer.file, err = os.OpenFile(path, os.O_RDONLY, 0)
	if err != nil {
		return nil, err
	}

	buffer.stopChan = make(chan bool)
	buffer.buffer = gbytes.NewBuffer()
	go buffer.start()

	return buffer, nil
}

//Buffer satisfies the interface gbytes.BufferProvider
func (buffer *Buffer) Buffer() *gbytes.Buffer {
	return buffer.buffer
}

//Close stops the buffer from scanning the target file
func (buffer *Buffer) Close() (err error) {
	if !buffer.closed {
		buffer.stopChan <- true
	}
	err = buffer.file.Close()
	if err != nil {
		if err.Error() == "invalid argument" {
			err = nil
		}
	}
	return
}

func (buffer *Buffer) start() {
	defer func() {
		close(buffer.stopChan)
		buffer.closed = true
	}()
	var index int64

	for {
		bytesBuffer := make([]byte, 10000)
		select {
		case <-time.After(time.Millisecond * 50):
			read, err := buffer.file.ReadAt(bytesBuffer, index)
			if err != nil && err != io.EOF {
				return
			}
			if read > 0 {
				index = index + int64(read)
				buffer.buffer.Write(bytesBuffer)
			}
		case <-buffer.stopChan:
			return
		}
	}
}
