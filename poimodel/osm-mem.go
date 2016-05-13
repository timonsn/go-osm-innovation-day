package poimodel


import "github.com/qedus/osmpbf"

type OSM struct {
	Nodes map[int64]*osmpbf.Node
	Ways  map[int64]*osmpbf.Way
	Relations map[int64]*osmpbf.Relation
}
