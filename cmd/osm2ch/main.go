package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/LdDl/ch"
)

var (
	tagStr      = flag.String("tags", "motorway,primary,primary_link,road,secondary,secondary_link,residential,tertiary,tertiary_link,unclassified,trunk,trunk_link", "Set of needed tags (separated by commas)")
	osmFileName = flag.String("file", "my_graph.osm.pbf", "Filename of *.osm.pbf file (it has to be compressed)")
	out         = flag.String("out", "my_graph.csv", "Filename of 'Comma-Separated Values' (CSV) formatted file")
	geomFormat  = flag.String("geomf", "wkt", "Format of output geometry. Expected values: wkt / geojson")
	units       = flag.String("units", "km", "Units of output weights. Expected values: km for kilometers / m for meters")
)

func main() {

	flag.Parse()

	tags := strings.Split(*tagStr, ",")
	cfg := ch.OsmConfiguration{
		EntityName: "highway", // Currrently we do not support others
		Tags:       tags,
	}

	edgeExpandedGraph, err := ch.ImportFromOSMFile(*osmFileName, &cfg)
	if err != nil {
		log.Fatalln(err)
	}

	file, err := os.Create(*out)
	if err != nil {
		log.Fatalln(err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()
	writer.Comma = ';'

	err = writer.Write([]string{"from_vertex_id", "to_vertex_id", "one_way", "weight", "geom"})
	if err != nil {
		log.Fatalln(err)
	}

	for source, targets := range *edgeExpandedGraph {
		for target, expEdge := range targets {
			geomStr := ""
			if strings.ToLower(*geomFormat) == "geojson" {
				geomStr = ch.PrepareGeoJSONLinestring(expEdge.Geom)
			} else {
				geomStr = ch.PrepareWKTLinestring(expEdge.Geom)
			}
			cost := expEdge.Cost
			if strings.ToLower(*units) == "m" {
				cost *= 1000.0
			}
			err = writer.Write([]string{fmt.Sprintf("%d", source), fmt.Sprintf("%d", target), "FT", fmt.Sprintf("%f", cost), geomStr})
			if err != nil {
				log.Fatalln(err)
			}
		}
	}
}
