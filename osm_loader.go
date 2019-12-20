package ch

import (
	"context"
	"math"
	"os"

	"github.com/paulmach/osm"
	"github.com/paulmach/osm/osmpbf"
)

// ImportFromOSMFile Imports graph from file of PBF-format (in OSM terms)
/*
	File should have PBF (Protocolbuffer Binary Format) extension according to https://github.com/paulmach/osm
*/
func ImportFromOSMFile(fileName string, cfg *OsmConfiguration) (*Graph, error) {
	graph := Graph{}
	f, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	scanner := osmpbf.New(context.Background(), f, 4)
	defer scanner.Close()

	nodes := make(map[int64][2]float64)

	for scanner.Scan() {
		obj := scanner.Object()
		if obj.ObjectID().Type() == "node" {
			nodes[obj.ObjectID().Ref()] = [2]float64{obj.(*osm.Node).Lon, obj.(*osm.Node).Lat}
		}

		if obj.ObjectID().Type() == "way" {
			tagMap := obj.(*osm.Way).TagMap()
			if tag, ok := tagMap[cfg.TagName]; ok {
				if cfg.CheckTag(tag) {
					oneway := false
					if v, ok := tagMap["oneway"]; ok {
						if v == "yes" || v == "1" {
							oneway = true
						}
					}
					ns := obj.(*osm.Way).Nodes
					for i := 1; i < len(ns); i++ {

						source := int64(ns[i-1].ID)
						target := int64(ns[i].ID)

						a := Coord{Lon: nodes[source][0], Lat: nodes[source][1]}
						b := Coord{Lon: nodes[target][0], Lat: nodes[target][1]}
						cost := DistanceBetweenPoints(a, b) * 1000

						graph.CreateVertex(source)
						graph.CreateVertex(target)

						graph.AddEdge(source, target, cost)

						if oneway == false {
							graph.AddEdge(target, source, cost)
						}
					}
				}
			}
		}
	}

	scanErr := scanner.Err()
	if scanErr != nil {
		return nil, scanErr
	}

	return &graph, nil
}

const (
	EarthRadius = 6370.986884258304
	pi180       = math.Pi / 180
)

type Coord struct {
	Lat float64
	Lon float64
}

func degreesToRadians(d float64) float64 {
	return d * pi180
}

func DistanceBetweenPoints(p, q Coord) float64 {
	lat1 := degreesToRadians(p.Lat)
	lon1 := degreesToRadians(p.Lon)
	lat2 := degreesToRadians(q.Lat)
	lon2 := degreesToRadians(q.Lon)
	diffLat := lat2 - lat1
	diffLon := lon2 - lon1
	a := math.Pow(math.Sin(diffLat/2), 2) + math.Cos(lat1)*math.Cos(lat2)*math.Pow(math.Sin(diffLon/2), 2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
	ans := c * EarthRadius
	return ans
}
