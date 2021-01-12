package beegoadapter_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestBeego(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Beego Adapter Suite")
}
