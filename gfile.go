package gfile

import (
	"io"
	"io/ioutil"
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

	if err := buffer.initialise(path); err != nil {
		return nil, err
	}
	go buffer.start()

	return buffer, nil
}

func (buffer *Buffer) initialise(path string) error {
	err := buffer.openFile(path)
	if err != nil {
		return err
	}

	buffer.stopChan = make(chan bool)
	buffer.buffer = gbytes.NewBuffer()

	return nil
}

//Buffer satisfies the interface gbytes.BufferProvider
func (buffer *Buffer) Buffer() *gbytes.Buffer {
	return buffer.buffer
}

//Close stops the buffer from scanning the target file
func (buffer *Buffer) Close() error {
	if atomic.CompareAndSwapInt32(&buffer.closed, 0, 1) {
		buffer.stopChan <- true
	}

	return buffer.closeFile()
}

func (buffer *Buffer) openFile(path string) (err error) {
	buffer.file, err = os.OpenFile(path, os.O_RDONLY, 0)
	return
}

func (buffer *Buffer) closeFile() (err error) {
	err = buffer.file.Close()
	if err != nil && err.Error() == "invalid argument" {
		err = nil
	}
	return
}

func (buffer *Buffer) readFile() (bytesBuffer []byte) {
	bytesBuffer, err := ioutil.ReadAll(buffer.file)
	if err != nil && err != io.EOF {
		panic(err.Error())
	}
	return
}

func (buffer *Buffer) start() {
	defer close(buffer.stopChan)

	for {
		select {
		case <-time.After(time.Millisecond * 50):
			buffer.update()
		case <-buffer.stopChan:
			return
		}
	}
}

func (buffer *Buffer) update() {
	buffer.write(buffer.readFile())
}

func (buffer *Buffer) write(bytesBuffer []byte) {
	if _, err := buffer.buffer.Write(bytesBuffer); err != nil {
		panic(err.Error())
	}
}
