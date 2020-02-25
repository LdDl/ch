package ch

import (
	"context"
	"fmt"
	"log"
	"math"
	"os"
	"strings"

	"github.com/paulmach/osm"
	"github.com/paulmach/osm/osmpbf"
)

const (
	earthRadius = 6370.986884258304
	pi180       = math.Pi / 180.0
	pi180Rev    = 180.0 / math.Pi
)

// geoPoint Representation of point on Earth
type geoPoint struct {
	Lat float64
	Lon float64
}

// String Pretty printing for geoPoint
func (gp geoPoint) String() string {
	return fmt.Sprintf("Lon: %f | Lat: %f", gp.Lon, gp.Lat)
}

// edgeComponent Representation of edge (vertex_from -> vertex_to)
type edgeComponent struct {
	from int64
	to   int64
}

// wayComponent First and last edges of osm.Way
type wayComponent struct {
	FirstEdge edgeComponent
	LastEdge  edgeComponent
}

// restrictionComponent Representation of member of restriction relation. Could be way or node.
type restrictionComponent struct {
	ID   int64
	Type string
}

// expandedEdge New edge built on top of two adjacent edges
type expandedEdge struct {
	ID   int64
	Cost float64
	Geom []geoPoint
}

// ExpandedGraph Representation of edge expanded graph
/*
	map[newSourceVertexID]map[newTargetVertexID]newExpandedEdge
*/
type ExpandedGraph map[int64]map[int64]expandedEdge

// ImportFromOSMFile Imports graph from file of PBF-format (in OSM terms)
/*
	File should have PBF (Protocolbuffer Binary Format) extension according to https://github.com/paulmach/osm
*/
func ImportFromOSMFile(fileName string, cfg *OsmConfiguration) (*ExpandedGraph, error) {
	// graph := Graph{}
	f, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	scanner := osmpbf.New(context.Background(), f, 4)
	defer scanner.Close()

	nodes := make(map[int64]geoPoint)
	vertices := make(map[int64]bool)
	newEdges := make(ExpandedGraph)
	newEdgeID := int64(1)
	allWays := make(map[int64]*wayComponent)
	restrictions := make(map[string]map[restrictionComponent]map[restrictionComponent]restrictionComponent)
	possibleRestrictionCombos := make(map[string]map[string]bool)

	skipped := 0
	unsupportedRestrictionRoles := 0
	for scanner.Scan() {
		obj := scanner.Object()
		if obj.ObjectID().Type() == "node" {
			nodes[obj.ObjectID().Ref()] = geoPoint{Lon: obj.(*osm.Node).Lon, Lat: obj.(*osm.Node).Lat}
		}
		if obj.ObjectID().Type() == "way" {
			tagMap := obj.(*osm.Way).TagMap()
			if tag, ok := tagMap[cfg.EntityName]; ok {
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
						a := nodes[source]
						b := nodes[target]
						cost := greatCircleDistance(a, b) // kilometers
						vertices[source] = true
						vertices[target] = true
						if _, ok := newEdges[source]; !ok {
							newEdges[source] = make(map[int64]expandedEdge)
						}
						newEdges[source][target] = expandedEdge{
							ID:   newEdgeID,
							Cost: cost,
							Geom: []geoPoint{a, b},
						}
						newEdgeID++
						if oneway == false {
							if _, ok := newEdges[target]; !ok {
								newEdges[target] = make(map[int64]expandedEdge)
							}
							newEdges[target][source] = expandedEdge{
								ID:   newEdgeID,
								Cost: cost,
								Geom: []geoPoint{b, a},
							}
							newEdgeID++
						}
					}
				}
			}
		}

		// Collect restrictions
		if obj.ObjectID().Type() == "relation" {
			relation := obj.(*osm.Relation)
			tagMap := relation.TagMap()
			if tag, ok := tagMap["restriction"]; ok {

				members := relation.Members
				if len(members) != 3 {
					skipped++
					// fmt.Printf("Restriction does not contain 3 members, relation ID: %d. Skip it\n", relation.ID)
					continue
				}
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
					unsupportedRestrictionRoles++
					// fmt.Printf("Something went wrong for first member of relation with ID: %d\n", relation.ID)
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
					unsupportedRestrictionRoles++
					// fmt.Printf("Something went wrong for second member of relation with ID: %d\n", relation.ID)
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
					unsupportedRestrictionRoles++
					// fmt.Printf("Something went wrong for third member of relation with ID: %d\n", relation.ID)
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

	expandedGraph := make(ExpandedGraph)
	for source := range newEdges {

		for target := range newEdges[source] {

			sourceExpandVertex := newEdges[source][target]
			sourceCost := sourceExpandVertex.Cost
			sourceMiddlePoint := middlePoint(sourceExpandVertex.Geom[0], sourceExpandVertex.Geom[1])
			if targetAsSource, ok := newEdges[target]; ok {
				for subTarget := range targetAsSource {
					targetExpandVertex := newEdges[target][subTarget]
					targetCost := targetExpandVertex.Cost
					targetMiddlePoint := middlePoint(targetExpandVertex.Geom[0], targetExpandVertex.Geom[1])

					// Handle bidirectional edges
					if sourceExpandVertex.Geom[0] == targetExpandVertex.Geom[1] && sourceExpandVertex.Geom[1] == targetExpandVertex.Geom[0] {
						continue
					}

					if _, ok := expandedGraph[sourceExpandVertex.ID]; !ok {
						expandedGraph[sourceExpandVertex.ID] = make(map[int64]expandedEdge)
					}
					expandedGraph[sourceExpandVertex.ID][targetExpandVertex.ID] = expandedEdge{
						Cost: (sourceCost + targetCost) / 2.0,
						Geom: []geoPoint{sourceMiddlePoint, sourceExpandVertex.Geom[1], targetMiddlePoint},
					}
				}
			}
		}
	}

	immposibleRestrictions := 0
	// Handling restrictions of "only" type
	for i, k := range restrictions {
		switch i {
		case "only_left_turn", "only_right_turn", "only_straight_on":
			// handle only way(from)-way(to)-node(via)
			for j, v := range k {
				if j.Type != "way" { // way(from)
					continue
				}
				from, ok := allWays[j.ID]
				if !ok {
					continue
				}
				for n := range v {
					if n.Type != "way" { // way(to)
						continue
					}
					if v[n].Type != "node" { // node(via)
						continue
					}

					to, ok := allWays[n.ID]
					if !ok {
						continue
					}

					rvertexVia := v[n].ID

					var rvertexFrom, rvertexTo int64

					switch rvertexVia {
					case from.LastEdge.to:
						rvertexFrom = from.LastEdge.from
						break
					case from.LastEdge.from:
						rvertexFrom = from.LastEdge.to
						break
					case from.FirstEdge.from:
						rvertexFrom = from.FirstEdge.to
						break
					case from.FirstEdge.to:
						rvertexFrom = from.FirstEdge.from
						break
					default:
						immposibleRestrictions++
						break
					}

					switch rvertexVia {
					case to.FirstEdge.to:
						rvertexTo = to.FirstEdge.from
						break
					case to.FirstEdge.from:
						rvertexTo = to.FirstEdge.to
						break
					case to.LastEdge.to:
						rvertexTo = to.LastEdge.from
						break
					case to.LastEdge.from:
						rvertexTo = to.LastEdge.to
						break
					default:
						immposibleRestrictions++
						break
					}

					fromExp, toExp := newEdges[rvertexFrom][rvertexVia].ID, newEdges[rvertexVia][rvertexTo].ID

					saveExde := expandedGraph[fromExp][toExp]
					if _, ok := expandedGraph[fromExp]; ok {
						delete(expandedGraph, fromExp)
						expandedGraph[fromExp] = make(map[int64]expandedEdge)
						expandedGraph[fromExp][toExp] = saveExde
					}
				}
			}
			break
		default:
			// @todo: need to think about U-turns: "no_u_turn"
			break
		}

	}

	// Handling restrictions of "no" type
	for i, k := range restrictions {
		switch i {
		case "no_left_turn", "no_right_turn", "no_straight_on":
			// handle only way(from)-way(to)-node(via)
			for j, v := range k {
				if j.Type != "way" { // way(from)
					continue
				}
				from, ok := allWays[j.ID]
				if !ok {
					continue
				}
				for n := range v {
					if n.Type != "way" { // way(to)
						continue
					}
					if v[n].Type != "node" { // node(via)
						continue
					}

					to, ok := allWays[n.ID]
					if !ok {
						continue
					}

					rvertexFrom := from.LastEdge.from
					rvertexTo := to.FirstEdge.to
					rvertexVia := v[n].ID

					fromExp, toExp := newEdges[rvertexFrom][rvertexVia].ID, newEdges[rvertexVia][rvertexTo].ID
					if _, ok := expandedGraph[fromExp]; ok {
						delete(expandedGraph[fromExp], toExp)
					}
				}
			}
			break
		default:
			// @todo: need to think about U-turns: "no_u_turn"
			break
		}

	}

	log.Printf("Skipped restrictions (which have not exactly 3 members): %d\n", skipped)
	log.Printf("Not properly handeled restrictions: %d\n", immposibleRestrictions)
	log.Printf("Number of unknow restriction roles (only 'from', 'to' and 'via' supported): %d\n", unsupportedRestrictionRoles)

	return &expandedGraph, nil
}

// degreesToRadians deg = r * pi / 180
func degreesToRadians(d float64) float64 {
	return d * pi180
}

// radiansTodegrees r = deg  * 180 / pi
func radiansTodegrees(d float64) float64 {
	return d * pi180Rev
}

// greatCircleDistance Returns distance between two geo-points (kilometers)
func greatCircleDistance(p, q geoPoint) float64 {
	lat1 := degreesToRadians(p.Lat)
	lon1 := degreesToRadians(p.Lon)
	lat2 := degreesToRadians(q.Lat)
	lon2 := degreesToRadians(q.Lon)
	diffLat := lat2 - lat1
	diffLon := lon2 - lon1
	a := math.Pow(math.Sin(diffLat/2), 2) + math.Cos(lat1)*math.Cos(lat2)*math.Pow(math.Sin(diffLon/2), 2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
	ans := c * earthRadius
	return ans
}

func middlePoint(p, q geoPoint) geoPoint {
	lat1 := degreesToRadians(p.Lat)
	lon1 := degreesToRadians(p.Lon)
	lat2 := degreesToRadians(q.Lat)
	lon2 := degreesToRadians(q.Lon)

	Bx := math.Cos(lat2) * math.Cos(lon2-lon1)
	By := math.Cos(lat2) * math.Sin(lon2-lon1)

	latMid := math.Atan2(math.Sin(lat1)+math.Sin(lat2), math.Sqrt((math.Cos(lat1)+Bx)*(math.Cos(lat1)+Bx)+By*By))
	lonMid := lon1 + math.Atan2(By, math.Cos(lat1)+Bx)
	return geoPoint{Lat: radiansTodegrees(latMid), Lon: radiansTodegrees(lonMid)}
}

// PrepareWKTLinestring Creates WKT LineString from set of points
func PrepareWKTLinestring(pts []geoPoint) string {
	ptsStr := make([]string, len(pts))
	for i := range pts {
		ptsStr[i] = fmt.Sprintf("%f %f", pts[i].Lon, pts[i].Lat)
	}
	return fmt.Sprintf("LINESTRING(%s)", strings.Join(ptsStr, ","))
}
