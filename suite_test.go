package geolang_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestGeolang(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Geolang Suite")
}
