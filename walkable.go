package geolang

import (
	"io/ioutil"
	"math"
	"os"

	"encoding/json"

	"github.com/golang/geo/s1"
	"github.com/golang/geo/s2"
)

type location struct {
	Lng float64 `json:"lng"`
	Lat float64 `json:"lat"`
}

type Venue struct {
	Name     string   `json:"name"`
	Location location `json:"location"`
}

type apiResponse struct {
	Response struct {
		Venues []Venue `json:"venues"`
	} `json:"response"`
}

type venueWithMeta struct {
	Venue
	cap s2.Cap
}

func LoadVenues(path string) ([]Venue, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	b, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, err
	}
	a := apiResponse{}
	if err := json.Unmarshal(b, &a); err != nil {
		return nil, err
	}
	return a.Response.Venues, nil
}

type Meter float64

const EarthCircumference Meter = 40075 * 1000
const radiansPerMeter = (2 * math.Pi) / s1.Angle(EarthCircumference)

func Walkable(venues []Venue, radius Meter) map[string][]string {
	venuesWithMeta := make([]venueWithMeta, len(venues))
	for i, venue := range venues {
		// internally stored as radians, normally communicated as degrees
		latLng := s2.LatLngFromDegrees(venue.Location.Lat, venue.Location.Lng)
		// convert the latlng to an s2 "point" (a r3 vector on the unit sphere)
		point := s2.PointFromLatLng(latLng)

		// not. a. disk.
		cap := s2.CapFromCenterAngle(point, radiansPerMeter*s1.Angle(radius))

		venuesWithMeta[i] = venueWithMeta{
			Venue: venue,
			cap:   cap,
		}
	}

	sets := make(map[string][]string, len(venuesWithMeta))
	for _, venue := range venuesWithMeta {
		intersections := []string{}
		for _, otherVenue := range venuesWithMeta {
			if venue.Name == otherVenue.Name {
				continue
			}
			if venue.cap.Intersects(otherVenue.cap) {
				intersections = append(intersections, otherVenue.Name)
			}
		}
		sets[venue.Name] = intersections
	}

	return sets
}
