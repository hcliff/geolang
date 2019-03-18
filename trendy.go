package geolang

import (
	"encoding/json"
	"io/ioutil"
	"os"

	"github.com/golang/geo/s2"
)

type GeoJSONPolygon struct {
	*s2.Polygon
}

type Neighborhood struct {
	Name  string
	Shape GeoJSONPolygon
}

func LoadNeighborhoods(path string) ([]Neighborhood, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	b, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, err
	}

	kiezs := []Neighborhood{}
	if err := json.Unmarshal(b, &kiezs); err != nil {
		return nil, err
	}
	return kiezs, nil
}

func Trendy(venues []Venue, kiezs []Neighborhood) []Venue {
	output := make([]Venue, 0, len(venues))
	for _, venue := range venues {
		// internally stored as radians, normally communicated as degrees
		latLng := s2.LatLngFromDegrees(venue.Location.Lat, venue.Location.Lng)
		// convert the latlng to an s2 "point" (a d3 vector on the unit sphere)
		point := s2.PointFromLatLng(latLng)

		for _, kiez := range kiezs {
			if kiez.Shape.ContainsPoint(point) {
				output = append(output, venue)
				break
			}
		}
	}
	return output
}

type NeighborhoodWithCovering struct {
	Neighborhood
	Covering s2.CellUnion
}

func NewNeighborhoodsWithCovering(kiezs []Neighborhood) []NeighborhoodWithCovering {
	rc := &s2.RegionCoverer{MaxLevel: 30, MaxCells: 20}
	output := make([]NeighborhoodWithCovering, len(kiezs))
	for i, kiez := range kiezs {
		covering := rc.FastCovering(kiez.Shape.Polygon)
		output[i] = NeighborhoodWithCovering{
			kiez,
			covering,
		}
	}
	return output
}

type VenueWithCell struct {
	Venue
	Cell s2.Cell
}

func NewVenuesWithCell(venues []Venue) []VenueWithCell {
	output := make([]VenueWithCell, len(venues))
	for i, venue := range venues {
		// internally stored as radians, normally communicated as degrees
		latLng := s2.LatLngFromDegrees(venue.Location.Lat, venue.Location.Lng)
		// convert the latlng to an s2 "point" (a d3 vector on the unit sphere)
		cell := s2.CellFromLatLng(latLng)
		output[i] = VenueWithCell{
			venue,
			cell,
		}
	}
	return output
}

func FastTrendy(venues []Venue, kiezs []Neighborhood) []Venue {
	kiezsMitCovering := NewNeighborhoodsWithCovering(kiezs)
	venuesWithCell := NewVenuesWithCell(venues)

	output := make([]Venue, 0, len(venues))
	for _, venue := range venuesWithCell {
		for _, kiez := range kiezsMitCovering {
			if kiez.Covering.ContainsCell(venue.Cell) {
				output = append(output, venue.Venue)
				break
			}
		}
	}

	return output
}

func FastestTrendy(venues []VenueWithCell, kiezs []NeighborhoodWithCovering) []Venue {
	output := make([]Venue, 0, len(venues))
	for _, venue := range venues {
		for _, kiez := range kiezs {
			if kiez.Covering.ContainsCell(venue.Cell) {
				output = append(output, venue.Venue)
				break
			}
		}
	}

	return output
}
