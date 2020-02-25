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

func main() {
	tagStr := flag.String("tags", "motorway,primary,primary_link,road,secondary,secondary_link,residential,tertiary,tertiary_link,unclassified,trunk,trunk_link", "Set of needed tags (separated by commas)")
	osmFileName := flag.String("file", "my_graph.osm.pbf", "Filename of *.osm.pbf file (it has to be compressed)")
	out := flag.String("out", "my_graph.csv", "Filename of 'Comma-Separated Values' (CSV) formatted file")
	flag.Parse()

	tags := strings.Split(*tagStr, ",")
	cfg := ch.OsmConfiguration{
		EntityName: "highway",
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
	err = writer.Write([]string{"from_vertex_id", "to_vertex_id", "weight", "geom"})
	if err != nil {
		log.Fatalln(err)
	}

	for source, targets := range *edgeExpandedGraph {
		for target, expEdge := range targets {
			err = writer.Write([]string{fmt.Sprintf("%d", source), fmt.Sprintf("%d", target), fmt.Sprintf("%f", expEdge.Cost), ch.PrepareWKTLinestring(expEdge.Geom)})
			if err != nil {
				log.Fatalln(err)
			}
		}
	}
}
