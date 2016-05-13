package main

import (
	"fmt"
	"github.com/qedus/osmpbf"
	"image"
	"image/color"
	"image/png"
	"io"
	"log"
	"os"
	"runtime"
	"github.com/timonsn/go-osm-innovation-day/paint2d"
	"github.com/timonsn/go-osm-innovation-day/poimodel"
)

func loadOSM(c *poimodel.PoiCollection) {
	f, err := os.Open("netherlands-latest.osm.pbf")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	d := osmpbf.NewDecoder(f)
	err = d.Start(runtime.GOMAXPROCS(-1)) // use several goroutines for faster decoding
	if err != nil {
		log.Fatal(err)
	}

	var nc, wc, rc uint64
	for {
		if v, err := d.Decode(); err == io.EOF {
			break
		} else if err != nil {
			log.Fatal(err)
		} else {
			switch v := v.(type) {
			case *osmpbf.Node:
				if handleNode(v) {
					if v.Tags["shop"] == "supermarket" {
						c.Add(poimodel.Poi{"AH",poimodel.Location{v.Lat, v.Lon}})
					} else {
						c.Add(poimodel.Poi{"BAG",poimodel.Location{v.Lat, v.Lon}})
					}

				}
				nc++
			case *osmpbf.Way:
				// Process Way v.
				wc++
			case *osmpbf.Relation:
				rc++
			default:
				log.Fatalf("unknown type %T\n", v)
			}
		}
	}
}

func handleNode(node *osmpbf.Node) bool {
	return node.Tags["source"] == "BAG"
}


func main() {
	mp := &poimodel.PoiCollection{}
	fmt.Println("Every supermarket!");
	loadOSM(mp)
	paint2d.Paint2d(mp)
}

func MapX(MinLon, MaxLon, lon float64, width int) int {
	return int(((180 + lon) - (180 + MinLon)) / (MaxLon - MinLon) * float64(width))
}
func MapY(MinLat, MaxLat, lat float64, height int) int {
	return int(((90 - lat) - (90 - MaxLat)) / (MaxLat - MinLat) * float64(height))
}

func paint(c *poimodel.PoiCollection) {
	f, err := os.OpenFile("x.png", os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Panic(err)
	}
	m := image.NewRGBA(image.Rectangle{image.Point{0, 0}, image.Point{1600, 1600}})
	rcolor := color.RGBA{uint8(255), uint8(0), uint8(0), 255}
	bcolor := color.RGBA{uint8(0), uint8(0), uint8(255), 255}
	c.ForEach(func(poi poimodel.Poi) {
		switch poi.Name {
		case "AH":
			m.Set(MapX(c.MinLon, c.MaxLon, poi.Location.Lon, 1600),
				MapY(c.MinLat, c.MaxLat, poi.Location.Lat, 1600),
				bcolor)
		case "BAG":
			m.Set(MapX(c.MinLon, c.MaxLon, poi.Location.Lon, 1600),
				MapY(c.MinLat, c.MaxLat, poi.Location.Lat, 1600),
				rcolor)
		}

	})

	if err = png.Encode(f, m); err != nil {
		log.Panic(err)
	}
}
