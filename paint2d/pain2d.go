package paint2d

import (
	"image"
	"github.com/llgcode/draw2d/draw2dimg"
	"image/color"
	"github.com/timonsn/go-osm-innovation-day/poimodel"
)

func MapX(MinLon, MaxLon, lon float64, width int) float64 {
	return ((180 + lon) - (180 + MinLon)) / (MaxLon - MinLon) * float64(width)
}
func MapY(MinLat, MaxLat, lat float64, height int) float64 {
	return ((90 - lat) - (90 - MaxLat)) / (MaxLat - MinLat) * float64(height)
}
func Paint2d(c *poimodel.PoiCollection) {
	dim := 1000
	dest := image.NewRGBA(image.Rect(0, 0, dim, dim))
	gc := draw2dimg.NewGraphicContext(dest)

	// Set some properties
	gc.SetFillColor(color.RGBA{0, 255, 0, 255})
	gc.SetStrokeColor(color.RGBA{80, 80, 80, 255})
	gc.SetLineWidth(2)

	c.ForEach(func(poi poimodel.Poi) {
		switch poi.Name {
		case "BAG":
			Draw(gc, MapX(c.MinLon, c.MaxLon, poi.Location.Lon, dim),
				MapY(c.MinLat, c.MaxLat, poi.Location.Lat, dim))
		}

	})

	gc.SetFillColor(color.RGBA{0, 255, 0, 255})
	gc.SetStrokeColor(color.RGBA{255, 0, 0	, 255})
	gc.SetLineWidth(5)
	c.ForEach(func(poi poimodel.Poi) {
		switch poi.Name {
		case "AH":
			Draw(gc, MapX(c.MinLon, c.MaxLon, poi.Location.Lon, dim),
				MapY(c.MinLat, c.MaxLat, poi.Location.Lat, dim))
		}

	})

	// Save to file
	draw2dimg.SaveToPngFile("hello.png", dest)
}

func Draw(gc *draw2dimg.GraphicContext, x0, y0 float64) {
	// Draw a line
	gc.MoveTo(x0-0.5, y0-0.5)
	gc.LineTo(x0+1, y0+1)
	gc.Stroke()
}