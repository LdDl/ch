package ch

import (
	"context"
	"fmt"
	"log"
	"math"
	"os"
	"time"

	"github.com/paulmach/osm"
	"github.com/paulmach/osm/osmpbf"
)

// restrictionComponent Representation of member of restriction relation. Could be way or node.
type restrictionComponent struct {
	ID   int64
	Type string
}

type expandedEdge struct {
	ID   int64
	Cost float64
}

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
	vertices := make(map[int64]bool)
	newEdges := make(map[int64]map[int64]expandedEdge)
	newEdgeID := int64(1)

	restrictions := make(map[restrictionComponent]map[restrictionComponent]restrictionComponent)

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

						vertices[source] = true
						vertices[target] = true

						if _, ok := newEdges[source]; !ok {
							newEdges[source] = make(map[int64]expandedEdge)
						}
						newEdges[source][target] = expandedEdge{
							ID:   newEdgeID,
							Cost: cost,
						}
						newEdgeID++

						if source == 5107446176 {
							fmt.Println("first", target)
						}
						// graph.CreateVertex(source)
						// graph.CreateVertex(target)
						// graph.AddEdge(source, target, cost)
						if oneway == false {
							if _, ok := newEdges[target]; !ok {
								newEdges[target] = make(map[int64]expandedEdge)
							}
							newEdges[target][source] = expandedEdge{
								ID:   newEdgeID,
								Cost: cost,
							}
							newEdgeID++
							// graph.AddEdge(target, source, cost)
						}
					}
				}
			}
		}

		// Handling restrictions
		if obj.ObjectID().Type() == "relation" {
			relation := obj.(*osm.Relation)
			tagMap := relation.TagMap()
			if tag, ok := tagMap["restriction"]; ok {
				_ = tag
				members := relation.Members
				if len(members) != 3 {
					fmt.Printf("Restriction does not contain 3 members, relation ID: %d", relation.ID)
				}
				// log.Println(tag, members)
				firstMember := restrictionComponent{-1, ""}
				secondMember := restrictionComponent{-1, ""}
				thirdMember := restrictionComponent{-1, ""}

				switch members[0].Role {
				case "from":
					firstMember = restrictionComponent{members[0].Ref, string(members[0].Type)}
					break
				case "via":
					thirdMember = restrictionComponent{members[0].Ref, string(members[0].Type)}
					break
				case "to":
					secondMember = restrictionComponent{members[0].Ref, string(members[0].Type)}
					break
				default:
					fmt.Printf("Something went wrong for first member of relation with ID: %d", relation.ID)
					break
				}

				switch members[1].Role {
				case "from":
					firstMember = restrictionComponent{members[1].Ref, string(members[1].Type)}
					break
				case "via":
					thirdMember = restrictionComponent{members[1].Ref, string(members[1].Type)}
					break
				case "to":
					secondMember = restrictionComponent{members[1].Ref, string(members[1].Type)}
					break
				default:
					fmt.Printf("Something went wrong for second member of relation with ID: %d", relation.ID)
					break
				}

				switch members[2].Role {
				case "from":
					firstMember = restrictionComponent{members[2].Ref, string(members[2].Type)}
					break
				case "via":
					thirdMember = restrictionComponent{members[2].Ref, string(members[2].Type)}
					break
				case "to":
					secondMember = restrictionComponent{members[2].Ref, string(members[2].Type)}
					break
				default:
					fmt.Printf("Something went wrong for third member of relation with ID: %d", relation.ID)
					break
				}

				if _, ok := restrictions[firstMember]; !ok {
					restrictions[firstMember] = make(map[restrictionComponent]restrictionComponent)
				}
				if _, ok := restrictions[firstMember][secondMember]; !ok {
					restrictions[firstMember][secondMember] = thirdMember
				}
			}
		}
	}

	expandedGraph := make(map[int64]map[int64]float64)
	for source := range newEdges {
		for target := range newEdges[source] {
			sourceExpandVertex := newEdges[source][target]
			sourceCost := sourceExpandVertex.Cost
			if targetAsSource, ok := newEdges[target]; ok {
				for subTarget := range targetAsSource {
					targetExpandVertex := newEdges[target][subTarget]
					targetCost := targetExpandVertex.Cost
					if _, ok := expandedGraph[sourceExpandVertex.ID]; !ok {
						expandedGraph[sourceExpandVertex.ID] = make(map[int64]float64)
					}
					// if target == 5107446176 {
					// 	fmt.Println("second", source, target, subTarget, sourceExpandVertex.ID, targetExpandVertex.ID, sourceCost+targetCost)
					// }
					expandedGraph[sourceExpandVertex.ID][targetExpandVertex.ID] = (sourceCost + targetCost) / 2.0
				}
			}
		}
	}

	source := int64(2574862283)
	target := int64(5107446176)
	subTarget := int64(96489258)

	edge := newEdges[source][target]
	subEdge := newEdges[target][subTarget]
	// toTarget1 := graph.Vertices[15770].inEdges
	// toTargetCost1 := graph.Vertices[15770].inECost[0]

	fmt.Println(edge, subEdge, expandedGraph[edge.ID][subEdge.ID])

	for source := range expandedGraph {
		for target := range expandedGraph[source] {
			cost := expandedGraph[source][target]
			graph.CreateVertex(source)
			graph.CreateVertex(target)
			graph.AddEdge(source, target, cost)
		}
	}

	log.Println("Number of edges:", len(newEdges))
	log.Println("Number of vertices:", len(vertices))
	log.Println("Number of new edges:", len(expandedGraph))
	log.Println("Number of new vertices:", len(graph.Vertices))
	log.Println("Number of restrictions:", len(restrictions))

	st := time.Now()
	graph.PrepareContracts()
	log.Println("Elapsed to prepare contracts:", time.Since(st))
	u := int64(17756)
	v := int64(17757)
	ans, path := graph.ShortestPath(u, v)
	log.Println(ans, path)

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
