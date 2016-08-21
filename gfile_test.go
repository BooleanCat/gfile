package gfile_test

import (
	"io/ioutil"
	"os"
	"time"

	"github.com/BooleanCat/gfile"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gbytes"
)

var Say = gbytes.Say

var _ = Describe("gfile", func() {
	var (
		err    error
		file   *os.File
		buffer *gfile.Buffer
	)

	Describe("#NewBuffer", func() {
		BeforeEach(func() {
			var fileErr error
			file, fileErr = ioutil.TempFile("", "gfile")
			Expect(fileErr).NotTo(HaveOccurred())
			file.Close()
		})

		AfterEach(func() {
			bufErr := buffer.Close()
			Expect(bufErr).NotTo(HaveOccurred())
			fileErr := os.Remove(file.Name())
			Expect(fileErr).NotTo(HaveOccurred())
		})

		JustBeforeEach(func() {
			buffer, err = gfile.NewBuffer(file.Name())
		})

		Context("when the target file exists", func() {
			It("does not return an error", func() {
				Expect(err).NotTo(HaveOccurred())
			})
		})

		Context("when the target file does not exist", func() {
			var badBuffer *gfile.Buffer

			JustBeforeEach(func() {
				badBuffer, err = gfile.NewBuffer("/foo/bar/baz")
			})

			It("returns an error", func() {
				Expect(err).To(HaveOccurred())
			})

			It("returns nil", func() {
				Expect(badBuffer).To(BeNil())
			})
		})
	})

	Describe("Buffer", func() {
		BeforeEach(func() {
			var fileErr error
			file, fileErr = ioutil.TempFile("", "gfile")
			Expect(fileErr).NotTo(HaveOccurred())
		})

		AfterEach(func() {
			buffer.Close()
			file.Close()
			os.Remove(file.Name())
		})

		JustBeforeEach(func() {
			buffer, err = gfile.NewBuffer(file.Name())
		})

		It("does not return an error", func() {
			Expect(err).NotTo(HaveOccurred())
		})

		Context("when the target file has static contents", func() {
			BeforeEach(func() {
				_, err := file.WriteString("This is a line of text")
				Expect(err).NotTo(HaveOccurred())
			})

			It("reads the contents of the file", func() {
				Eventually(buffer).Should(Say("This is a line of text"))
			})
		})

		Context("when the target file has multi-line static contents", func() {
			var lines = "This is a line of text\nand this is another"

			BeforeEach(func() {
				_, err := file.WriteString(lines)
				Expect(err).NotTo(HaveOccurred())
			})

			It("reads the contents of the file", func() {
				Eventually(buffer).Should(Say(lines))
			})

			It("doesn't re-read the earlier file contents", func() {
				Eventually(buffer).Should(Say("This is a line of text"))
				Expect(buffer).NotTo(Say("This is a line of text"))
			})

			It("reads contents in sequence", func() {
				Eventually(buffer).Should(Say("This is a line of text"))
				Eventually(buffer).Should(Say("\nand this is another"))
			})
		})

		Context("when the target file is continuously written to", func() {
			var longLivedBuffer *gfile.Buffer

			BeforeEach(func() {
				var err error

				longLivedBuffer, err = gfile.NewBuffer(file.Name())
				Expect(err).NotTo(HaveOccurred())

				_, err = file.WriteString("An initial line\n")
				Expect(err).NotTo(HaveOccurred())
				Eventually(longLivedBuffer).Should(Say("An initial line\n"))

				_, err = file.WriteString("And then another line")
				Expect(err).NotTo(HaveOccurred())
			})

			It("reads the new contents", func() {
				Eventually(longLivedBuffer).Should(Say("And then another line"))
			})
		})

		Context("when the target file is written to asynchronously", func() {
			var syncChan chan bool

			BeforeEach(func() {
				syncChan = make(chan bool)

				go func() {
					time.Sleep(time.Millisecond * 250)
					_, fileErr := file.WriteString("I came from a go func!")
					Expect(fileErr).NotTo(HaveOccurred())
					syncChan <- true
				}()
			})

			AfterEach(func() {
				<-syncChan
				close(syncChan)
			})

			It("reads the new contents", func() {
				Eventually(buffer).Should(Say("I came from a go func!"))
			})
		})

		Describe("#Close", func() {
			var closeErr error

			JustBeforeEach(func() {
				closeErr = buffer.Close()
			})

			It("does not return an error", func() {
				Expect(closeErr).NotTo(HaveOccurred())
			})

			Context("when closed repeatedly", func() {
				JustBeforeEach(func() {
					Expect(closeErr).NotTo(HaveOccurred())
					closeErr = buffer.Close()
				})

				It("does not return an error", func() {
					Expect(closeErr).NotTo(HaveOccurred())
				})
			})
		})

		Describe("#Buffer", func() {
			It("is a gbytes.Buffer", func() {
				Expect(buffer.Buffer()).To(BeAssignableToTypeOf(gbytes.NewBuffer()))
			})
		})
	})
})
