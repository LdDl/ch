package main

import (
	"flag"
	"log"
	"strings"

	"github.com/LdDl/ch"
)

func main() {
	tagStr := flag.String("tags", "motorway,primary,primary_link,road,secondary,secondary_link,residential,tertiary,tertiary_link,unclassified,trunk,trunk_link", "Set of needed tags (separated by commas)")
	osmFileName := flag.String("file", "my_graph.osm.pbf", "Filename of *.osm.pbf file (it has to be compressed)")
	flag.Parse()

	tags := strings.Split(*tagStr, ",")
	cfg := ch.OsmConfiguration{
		TagName: "highway",
		Tags:    tags,
	}
	log.Println(cfg, *osmFileName)

	graph, err := ch.ImportFromOSMFile("data/moscow_center_reduced.osm.pbf", &cfg)
	if err != nil {
		log.Fatalln(err)
	}

	_ = graph
}
