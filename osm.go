package main

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"io"
	"log"
	"os"
	"runtime"

	"github.com/boltdb/bolt"
	"github.com/qedus/osmpbf"
	"github.com/timonsn/go-osm-innovation-day/paint2d"
	"github.com/timonsn/go-osm-innovation-day/poimodel"
)

type Store struct {
	db *bolt.DB
}

func (s *Store) Open(filename string) error {
	// Open bolt
	var err error
	s.db, err = bolt.Open(filename, 0600, nil)
	if err != nil {
		log.Fatal(err)
	}

	return s.db.Batch(func(tx *bolt.Tx) error {
		tx.CreateBucket([]byte("Node"))
		tx.CreateBucket([]byte("Way"))
		tx.CreateBucket([]byte("Relation"))
		return nil
	})

}

func (s *Store) Close() {
	s.db.Close()
}

func (s *Store) CreateNode(obj *osmpbf.Node) error {
	var blob bytes.Buffer
	enc := gob.NewEncoder(&blob)
	err := enc.Encode(obj)
	if err != nil {
		log.Fatal(err)
	}
	return s.db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte("Node"))
		err = bucket.Put([]byte(fmt.Sprintf("%d", obj.ID)), blob.Bytes())
		//log.Printf("Stored node %d\n", obj.ID)
		return err
	})
}

func (s *Store) CreateWay(obj *osmpbf.Way) error {
	var blob bytes.Buffer
	enc := gob.NewEncoder(&blob)
	err := enc.Encode(obj)
	if err != nil {
		log.Fatal(err)
	}
	return s.db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte("Way"))
		err = bucket.Put([]byte(fmt.Sprintf("%d", obj.ID)), blob.Bytes())
		//log.Printf("Stored way %d\n", obj.ID)
		return err
	})
}

func (s *Store) CreateRelation(obj *osmpbf.Relation) error {
	var blob bytes.Buffer
	enc := gob.NewEncoder(&blob)
	err := enc.Encode(obj)
	if err != nil {
		log.Fatal(err)
	}
	return s.db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte("Relation"))
		err = bucket.Put([]byte(fmt.Sprintf("%d", obj.ID)), blob.Bytes())
		//log.Printf("Stored relation %d\n", obj.ID)
		return err
	})
}

func loadOSM(store *Store, filename string) poimodel.OSM {

	o := poimodel.OSM{}
	o.Nodes = make(map[int64]*osmpbf.Node)
	o.Ways = make(map[int64]*osmpbf.Way)
	o.Relations = make(map[int64]*osmpbf.Relation)

	f, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	d := osmpbf.NewDecoder(f)
	err = d.Start(runtime.GOMAXPROCS(-1)) // use several goroutines for faster decoding
	if err != nil {
		log.Fatal(err)
	}
	for i := 0; i < 100000; i++ {
		if v, err := d.Decode(); err == io.EOF {
			break
		} else if err != nil {
			log.Fatal(err)
		} else {

			switch v := v.(type) {
			case *osmpbf.Node:
				err = store.CreateNode(v)
				if err != nil {
					log.Fatal(err)
				}
				case *osmpbf.Way:
				err = store.CreateWay(v)
				if err != nil {
					log.Fatal(err)
				}
			case *osmpbf.Relation:
				err = store.CreateRelation(v)
				if err != nil {
					log.Fatal(err)
				}
			default:
				log.Fatalf("unknown type %T\n", v)
			}
		}
	}
	return o
}

func main() {
	var store Store
	err := store.Open("geo.db")
	if err != nil {
		log.Fatal(err)
	}
	defer store.Close()

	mp := &poimodel.PoiCollection{}
	fmt.Println("Loading")
	osmDump := loadOSM(&store, "netherlands-latest.osm.pbf")
	fmt.Println("OSM loaded")
	poimodel.ExtractSupermarket(osmDump)
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
