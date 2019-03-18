package geolang

import (
	"encoding/json"
	"github.com/golang/geo/s2"
)

func PointFromDegrees(degrees []float64) s2.Point {
	lat := degrees[1]
	lng := degrees[0]
	latLng := s2.LatLngFromDegrees(lat, lng)
	return s2.PointFromLatLng(latLng)
}

type geoJSONPolygon struct {
	Type        string        `json:"type"`
	Coordinates [][][]float64 `json:"coordinates"`
}

func (l *GeoJSONPolygon) UnmarshalJSON(data []byte) error {
	var geoJSONPolygon geoJSONPolygon
	if err := json.Unmarshal(data, &geoJSONPolygon); err != nil {
		return err
	}

	loops := make([]*s2.Loop, len(geoJSONPolygon.Coordinates))
	for i, loop := range geoJSONPolygon.Coordinates {
		points := make([]s2.Point, len(loop))
		for ii, degrees := range loop {
			points[ii] = PointFromDegrees(degrees)
		}
		loops[i] = s2.LoopFromPoints(points)
	}

	l.Polygon = s2.PolygonFromLoops(loops)

	return nil
}
