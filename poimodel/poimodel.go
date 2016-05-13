package poimodel

import "math"

type Poi struct {
	Name     string
	Location Location
}
type Location struct {
	Lat float64
	Lon float64
}
type PoiCollection struct {
	MinLat, MinLon, MaxLat, MaxLon float64
	pois                           []Poi
}

func (c *PoiCollection) Add(p Poi) {
	if c.pois == nil {
		c.MinLat = p.Location.Lat
		c.MaxLat = p.Location.Lat
		c.MinLon = p.Location.Lon
		c.MaxLon = p.Location.Lon
	}

	c.MinLat = math.Min(c.MinLat, p.Location.Lat)
	c.MinLon = math.Min(c.MinLon, p.Location.Lon)
	c.MaxLat = math.Max(c.MaxLat, p.Location.Lat)
	c.MaxLon = math.Max(c.MaxLon, p.Location.Lon)

	c.pois = append(c.pois, p)
}

func (c *PoiCollection) ForEach(f func(poi Poi)) {
	for i := range c.pois {
		f(c.pois[i])
	}
}
