package gfile_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestGfile(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Gfile Suite")
}
