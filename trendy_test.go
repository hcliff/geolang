package geolang_test

import (
	"github.com/hcliff/geolang"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("checking which bars are trendy", func() {

	var venues []geolang.Venue
	var kiezs []geolang.Neighborhood
	BeforeEach(func() {
		var err error
		venues, err = geolang.LoadVenues("fixtures/venues.json")
		Ω(err).ShouldNot(HaveOccurred())
		kiezs, err = geolang.LoadNeighborhoods("fixtures/super_trendy_neighborhoods.json")
		Ω(err).ShouldNot(HaveOccurred())
	})

	Measure("raycasting", func(b Benchmarker) {
		t := b.Time("runtime", func() {
			geolang.Trendy(venues, kiezs)
		})

		b.RecordValueWithPrecision("time nano", float64(t), "ns", 1)
	}, 500)

	Measure("s2", func(b Benchmarker) {
		t := b.Time("runtime", func() {
			geolang.FastTrendy(venues, kiezs)
		})

		b.RecordValueWithPrecision("time nano", float64(t), "ns", 1)
	}, 500)

	Context("with prebuilt coverings", func() {

		var venuesWithCell []geolang.VenueWithCell
		var kiezsMitCovering []geolang.NeighborhoodWithCovering
		BeforeEach(func() {
			venuesWithCell = geolang.NewVenuesWithCell(venues)
			kiezsMitCovering = geolang.NewNeighborhoodsWithCovering(kiezs)
		})

		Measure("prebuilt s2", func(b Benchmarker) {
			t := b.Time("runtime", func() {
				geolang.FastestTrendy(venuesWithCell, kiezsMitCovering)
			})

			b.RecordValueWithPrecision("time nano", float64(t), "ns", 1)
		}, 500)
	})

})
