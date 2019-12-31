package ch

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"os"
	"time"

	geojson "github.com/paulmach/go.geojson"
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
	geom []Coord
}

type wayComponent struct {
	FirstEdge edgeComponent
	LastEdge  edgeComponent
}

type edgeComponent struct {
	from int64
	to   int64
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
	allWays := make(map[int64]*wayComponent)
	restrictions := make(map[string]map[restrictionComponent]map[restrictionComponent]restrictionComponent)
	possibleRestrictionCombos := make(map[string]map[string]bool)

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
					way := obj.(*osm.Way)
					allWays[int64(way.ID)] = &wayComponent{}

					for i := 1; i < len(ns); i++ {
						source := int64(ns[i-1].ID)
						target := int64(ns[i].ID)

						if (i - 1) == 0 {
							allWays[int64(way.ID)].FirstEdge = edgeComponent{
								from: source,
								to:   target,
							}
						}

						if i == len(ns)-1 {
							allWays[int64(way.ID)].LastEdge = edgeComponent{
								from: source,
								to:   target,
							}
						}

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
							geom: []Coord{a, b},
						}

						newEdgeID++

						// if source == 5107446176 {
						// 	fmt.Println("first", target)
						// }
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
								geom: []Coord{b, a},
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

				if _, ok := possibleRestrictionCombos[tag]; !ok {
					possibleRestrictionCombos[tag] = make(map[string]bool)
				}
				possibleRestrictionCombos[tag][fmt.Sprintf("%s;%s;%s", firstMember.Type, secondMember.Type, thirdMember.Type)] = true

				if _, ok := restrictions[tag]; !ok {
					restrictions[tag] = make(map[restrictionComponent]map[restrictionComponent]restrictionComponent)
				}

				if _, ok := restrictions[tag][firstMember]; !ok {
					restrictions[tag][firstMember] = make(map[restrictionComponent]restrictionComponent)
				}
				if _, ok := restrictions[tag][firstMember][secondMember]; !ok {
					restrictions[tag][firstMember][secondMember] = thirdMember
				}
			}
		}
	}

	for i, v := range possibleRestrictionCombos {
		fmt.Println(i)
		for j := range v {
			fmt.Println("\t", j)
		}
	}

	expandedGraph := make(map[int64]map[int64]expandedEdge)
	for source := range newEdges {
		for target := range newEdges[source] {
			sourceExpandVertex := newEdges[source][target]
			sourceCost := sourceExpandVertex.Cost

			sourceMiddlePoint := MiddlePoint(sourceExpandVertex.geom[0], sourceExpandVertex.geom[1])
			if targetAsSource, ok := newEdges[target]; ok {
				for subTarget := range targetAsSource {
					targetExpandVertex := newEdges[target][subTarget]
					targetCost := targetExpandVertex.Cost
					targetMiddlePoint := MiddlePoint(targetExpandVertex.geom[0], targetExpandVertex.geom[1])
					if _, ok := expandedGraph[sourceExpandVertex.ID]; !ok {
						expandedGraph[sourceExpandVertex.ID] = make(map[int64]expandedEdge)
					}
					// if target == 5107446176 {
					// 	fmt.Println("second", source, target, subTarget, sourceExpandVertex.ID, targetExpandVertex.ID, sourceCost+targetCost)
					// 	fmt.Println("From edge", sourceExpandVertex.geom)
					// 	fmt.Println("To edge", targetExpandVertex.geom)
					// 	fmt.Println("New edge", sourceMiddlePoint, sourceExpandVertex.geom[1], targetMiddlePoint)
					// }

					expandedGraph[sourceExpandVertex.ID][targetExpandVertex.ID] = expandedEdge{
						Cost: (sourceCost + targetCost) / 2.0,
						geom: []Coord{sourceMiddlePoint, sourceExpandVertex.geom[1], targetMiddlePoint},
					}

				}
			}
		}
	}

	for i, k := range restrictions {
		if i != "no_left_turn" {
			continue
		}
		for j, v := range k {
			if j.Type != "way" {
				continue
			}
			from, ok := allWays[j.ID]
			if !ok {
				continue
			}
			for n := range v {
				if n.Type != "way" {
					continue
				}
				if v[n].Type != "node" {
					continue
				}

				to, ok := allWays[n.ID]
				if !ok {
					continue
				}

				fromRestVertex := from.LastEdge.from
				toRestVertex := to.FirstEdge.to
				viaRestVertex := v[n].ID

				if j.ID == 23178249 {
					fromExp, toExp := newEdges[fromRestVertex][viaRestVertex].ID, newEdges[viaRestVertex][toRestVertex].ID
					fmt.Println("gotcha", fromExp, toExp, expandedGraph[fromExp][toExp])
					delete(expandedGraph[fromExp], toExp)
				}
			}
		}
	}

	source := int64(2574862283)
	target := int64(5107446176)
	subTarget := int64(96489258)

	edge := newEdges[source][target]
	subEdge := newEdges[target][subTarget]
	_, _ = edge, subEdge
	// toTarget1 := graph.Vertices[15770].inEdges
	// toTargetCost1 := graph.Vertices[15770].inECost[0]

	// fmt.Println(edge, subEdge, expandedGraph[edge.ID][subEdge.ID])

	for source := range expandedGraph {
		for target := range expandedGraph[source] {
			cost := expandedGraph[source][target]
			graph.CreateVertex(source)
			graph.CreateVertex(target)
			graph.AddEdge(source, target, cost.Cost)
		}
	}

	fc := geojson.NewFeatureCollection()
	for i, k := range expandedGraph {
		for j := range k {
			// fmt.Println(i, j, k[j])
			ls := make([][]float64, len(k[j].geom))
			for v := range k[j].geom {
				ls[v] = []float64{k[j].geom[v].Lon, k[j].geom[v].Lat}
			}
			f := geojson.NewLineStringFeature(ls)
			f.SetProperty("from", i)
			f.SetProperty("to", j)
			fc.AddFeature(f)
		}
	}
	bytesFC, _ := fc.MarshalJSON()
	_ = ioutil.WriteFile("expanded_graph.json", bytesFC, 0644)

	log.Println("Number of edges:", len(newEdges))
	log.Println("Number of vertices:", len(vertices))
	log.Println("Number of new edges:", len(expandedGraph))
	log.Println("Number of new vertices:", len(graph.Vertices))

	st := time.Now()
	graph.PrepareContracts()
	log.Println("Elapsed to prepare contracts:", time.Since(st))
	u := int64(17752)
	v := int64(19811)

	u = int64(2274)
	v = int64(17709)

	ans, path := graph.ShortestPath(u, v)
	_, _ = ans, path
	// log.Println("answer", ans, path)

	fcAnswer := geojson.NewFeatureCollection()
	for i := 1; i < len(path); i++ {
		from := path[i-1]
		to := path[i]
		edge := expandedGraph[from][to]

		ls := make([][]float64, len(edge.geom))
		for v := range edge.geom {
			ls[v] = []float64{edge.geom[v].Lon, edge.geom[v].Lat}
		}
		f := geojson.NewLineStringFeature(ls)
		f.SetProperty("from", from)
		f.SetProperty("to", to)
		fcAnswer.AddFeature(f)

		// fmt.Println("path edge", edge)
	}
	bytesAnswer, _ := fcAnswer.MarshalJSON()
	_ = ioutil.WriteFile("answer_path.json", bytesAnswer, 0644)

	scanErr := scanner.Err()
	if scanErr != nil {
		return nil, scanErr
	}

	return &graph, nil
}

const (
	EarthRadius = 6370.986884258304
	pi180       = math.Pi / 180
	pi180_rev   = 180 / math.Pi
)

type Coord struct {
	Lat float64
	Lon float64
}

func degreesToRadians(d float64) float64 {
	return d * pi180
}

func radiansTodegrees(d float64) float64 {
	return d * pi180_rev
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

func MiddlePoint(p, q Coord) Coord {
	lat1 := degreesToRadians(p.Lat)
	lon1 := degreesToRadians(p.Lon)
	lat2 := degreesToRadians(q.Lat)
	lon2 := degreesToRadians(q.Lon)

	Bx := math.Cos(lat2) * math.Cos(lon2-lon1)
	By := math.Cos(lat2) * math.Sin(lon2-lon1)

	latMid := math.Atan2(math.Sin(lat1)+math.Sin(lat2), math.Sqrt((math.Cos(lat1)+Bx)*(math.Cos(lat1)+Bx)+By*By))
	lonMid := lon1 + math.Atan2(By, math.Cos(lat1)+Bx)
	return Coord{Lat: radiansTodegrees(latMid), Lon: radiansTodegrees(lonMid)}
}
