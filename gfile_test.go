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

var _ = Describe("gfile", func() {
	Describe("#NewBuffer", func() {
		Context("when the target file exists", func() {
			var (
				file   *os.File
				buffer *gfile.Buffer
				err    error
			)

			BeforeEach(func() {
				var fileErr error
				file, fileErr = ioutil.TempFile("/tmp", "gfile")
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
					Eventually(buffer).Should(gbytes.Say("This is a line of text"))
				})
			})

			Context("when the target file has multi-line static contents", func() {
				var lines = "This is a line of text\nand this is another"

				BeforeEach(func() {
					_, err := file.WriteString(lines)
					Expect(err).NotTo(HaveOccurred())
				})

				It("reads the contents of the file", func() {
					Eventually(buffer).Should(gbytes.Say(lines))
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
					Eventually(longLivedBuffer).Should(gbytes.Say("An initial line\n"))

					_, err = file.WriteString("And then another line")
					Expect(err).NotTo(HaveOccurred())
				})

				It("reads the new contents", func() {
					Eventually(longLivedBuffer).Should(gbytes.Say("And then another line"))
				})
			})

			Context("when the target file is written to asynchonously", func() {
				BeforeEach(func() {
					go func() {
						time.Sleep(time.Millisecond * 250)
						_, fileErr := file.WriteString("I came from a go func!")
						Expect(fileErr).NotTo(HaveOccurred())
					}()
				})

				It("reads the new contents", func() {
					Eventually(buffer).Should(gbytes.Say("I came from a go func!"))
				})
			})
		})

		Context("when the target file does not exist", func() {
			var err error

			JustBeforeEach(func() {
				_, err = gfile.NewBuffer("/foo/bar/baz")
			})

			It("returns an error", func() {
				Expect(err).To(HaveOccurred())
			})
		})
	})
})
