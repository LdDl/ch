package ch

import (
	"context"
	"math"
	"os"

	"github.com/paulmach/osm"
	"github.com/paulmach/osm/osmpbf"
)

// Config - tags osm
type Config struct {
	Name string
	Tags []string
}

// LoadOsmGraph - loading osm graph
func LoadOsmGraph(fileName string, cfg Config) (*Graph, error) {
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
			if tag, ok := obj.(*osm.Way).TagMap()[cfg.Name]; ok {
				for _, t := range cfg.Tags {
					if t == tag {
						oneway := false
						if obj.(*osm.Way).TagMap()["oneway"] == "yes" {
							oneway = true
						}
						ns := obj.(*osm.Way).Nodes
						for i := 0; i < len(ns)-1; i++ {
							sourse := int64(ns[i].ID)
							target := int64(ns[i+1].ID)
							a := Coord{Lon: nodes[sourse][0], Lat: nodes[sourse][1]}
							b := Coord{Lon: nodes[target][0], Lat: nodes[target][1]}
							cost := DistanceBetweenPoints(a, b)

							graph.AddVertex(int(sourse), int(sourse))
							graph.AddVertex(int(target), int(target))

							graph.AddEdge(int(sourse), int(target), cost)

							if oneway == false {
								graph.AddEdge(int(target), int(sourse), cost)
							}
						}
						break
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
	earthRaidusKm = 6371 // radius of the earth in kilometers.
)

// Coord represents a geographic coordinate.
type Coord struct {
	Lat float64
	Lon float64
}

// degreesToRadians converts from degrees to radians.
func degreesToRadians(d float64) float64 {
	return d * math.Pi / 180
}

// DistanceBetweenPoints calculates the shortest path between two coordinates on the surface
// of the Earth. This function returns two units of measure, the first is the
// distance in miles, the second is the distance in kilometers.
func DistanceBetweenPoints(p, q Coord) (km float64) {
	lat1 := degreesToRadians(p.Lat)
	lon1 := degreesToRadians(p.Lon)
	lat2 := degreesToRadians(q.Lat)
	lon2 := degreesToRadians(q.Lon)
	diffLat := lat2 - lat1
	diffLon := lon2 - lon1
	a := math.Pow(math.Sin(diffLat/2), 2) + math.Cos(lat1)*math.Cos(lat2)*
		math.Pow(math.Sin(diffLon/2), 2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
	km = c * earthRaidusKm
	return km
}
