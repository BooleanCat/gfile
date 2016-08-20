gfile
=====

gfile is an implementation of [gbytes](http://github.com/onsi/gomega)'
BufferProvider interface that allows you to make assertions on the contents of
files as they are updated using the [ginkgo](http://github.com/onsi/ginkgo)
testing framework.

Installation
------------

`go get github.com/BooleanCat/gfile`

Usage
-----

Below is an example usage of `gfile.Buffer` to check that a logfile has had a
specific message written to it.

```
...
Context("when writing to a log file", func() {
    var buffer *gfile.Buffer

    BeforeEach(func() {
        buffer = gfile.NewBuffer("/tmp/foo")
    })

    AfterEach(func() {
        buffer.Close()
    })

    It("can be read back", func() {
        WriteToLog("/tmp/foo")
        Eventually(buffer).Should(gbytes.Say("I'm a log message"))
    })
})
...
```
