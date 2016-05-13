package poimodel

func ExtractSupermarket (o OSM) PoiCollection {
	poi := PoiCollection{}

	for _, node := range o.Nodes {
		if node.Tags["shop"] == "supermarket" {
			poi.Add(Poi{"AH", Location{node.Lat, node.Lon}})
		}

		if node.Tags["source"] == "BAG" {
			poi.Add(Poi{"BAG", Location{node.Lat, node.Lon}})
		}
	}

	for _, way := range o.Ways {
		if way.Tags["shop"] == "supermarket" {
			//zoek node
		}
	}
	return poi
}
