package gfile

import (
	"io"
	"os"
	"sync/atomic"
	"time"

	"github.com/onsi/gomega/gbytes"
)

//Buffer ...
type Buffer struct {
	buffer   *gbytes.Buffer
	stopChan chan bool
	file     *os.File
	closed   int32
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
	if atomic.CompareAndSwapInt32(&buffer.closed, 0, 1) {
		buffer.stopChan <- true
	}

	err = buffer.file.Close()
	if err != nil && err.Error() == "invalid argument" {
		err = nil
	}
	return
}

func (buffer *Buffer) start() {
	defer close(buffer.stopChan)

	for {
		bytesBuffer := make([]byte, 10000)
		select {
		case <-time.After(time.Millisecond * 50):
			read, err := buffer.file.Read(bytesBuffer)
			if err != nil && err != io.EOF {
				panic(err.Error())
			}
			if read > 0 {
				if _, err = buffer.buffer.Write(bytesBuffer); err != nil {
					panic(err.Error())
				}
			}
		case <-buffer.stopChan:
			return
		}
	}
}
