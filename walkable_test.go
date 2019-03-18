package geolang_test

import (
	"github.com/hcliff/geolang"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("walkable tests", func() {

	var venues []geolang.Venue
	BeforeEach(func() {
		var err error
		venues, err = geolang.LoadVenues("fixtures/venues.json")
		Ω(err).ShouldNot(HaveOccurred())
	})

	Measure("calculating walkable venues", func(b Benchmarker) {
		t := b.Time("runtime", func() {
			geolang.Walkable(venues, geolang.Meter(50))
		})

		b.RecordValueWithPrecision("time nano", float64(t), "ns", 1)
	}, 500)
})
