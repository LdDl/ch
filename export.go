package ch

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/pkg/errors"
)

// ExportToFile Exports graph to file of CSV-format
// Header of edges CSV-file:
// 		from_vertex_id - int64, ID of source vertex
// 		to_vertex_id - int64, ID of target vertex
// 		weight - float64, Weight of an edge
// Header of vertices CSV-file:
// 		vertex_id - int64, ID of vertex
// 		order_pos - int, Position of vertex in hierarchies (evaluted by library)
// 		importance - int, Importance of vertex in graph (evaluted by library)
// Header of shortcuts CSV-file:
// 		from_vertex_id - int64, ID of source vertex
// 		to_vertex_id - int64, ID of target vertex
// 		weight - float64, Weight of an shortcut
// 		via_vertex_id - int64, ID of vertex through which the shortcut exists
func (graph *Graph) ExportToFile(fname string) error {

	fnamePart := strings.Split(fname, ".csv") // to guarantee proper filename and its extension

	err := graph.ExportEdgesToFile(fnamePart[0] + ".csv")
	if err != nil {
		return errors.Wrap(err, "Can't export edges")
	}

	// Write reference information about vertices
	err = graph.ExportVerticesToFile(fnamePart[0] + "_vertices.csv")
	if err != nil {
		return errors.Wrap(err, "Can't export shortcuts")
	}

	// Write reference information about contractions
	err = graph.ExportShortcutsToFile(fnamePart[0] + "_shortcuts.csv")
	if err != nil {
		return errors.Wrap(err, "Can't export shortcuts")
	}

	return nil
}

// ExportVerticesToFile Exports edges information to CSV-file with header:
// 	from_vertex_id - int64, ID of source vertex
// 	to_vertex_id - int64, ID of target vertex
// 	weight - float64, Weight of an edge
func (graph *Graph) ExportEdgesToFile(fname string) error {
	file, err := os.Create(fname)
	if err != nil {
		return errors.Wrap(err, "Can't create edges file")
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()
	writer.Comma = ';'
	err = writer.Write([]string{"from_vertex_id", "to_vertex_id", "weight"})
	if err != nil {
		return errors.Wrap(err, "Can't write header to edges file")
	}

	for i := 0; i < len(graph.Vertices); i++ {
		currentVertexExternal := graph.Vertices[i].Label
		currentVertexInternal := graph.Vertices[i].vertexNum
		// Write reference information about "outcoming" adjacent vertices
		// Why don't write info about "incoming" adjacent vertices also? Because all edges will be covered due the loop iteration mechanism
		outcomingNeighbors := graph.Vertices[i].outIncidentEdges
		for j := 0; j < len(outcomingNeighbors); j++ {
			toVertexExternal := graph.Vertices[outcomingNeighbors[j].vertexID].Label
			toVertexInternal := outcomingNeighbors[j].vertexID
			cost := outcomingNeighbors[j].cost
			if _, ok := graph.shortcuts[currentVertexInternal][toVertexInternal]; !ok {
				err = writer.Write([]string{
					fmt.Sprintf("%d", currentVertexExternal),
					fmt.Sprintf("%d", toVertexExternal),
					strconv.FormatFloat(cost, 'f', -1, 64),
				})
				if err != nil {
					return errors.Wrap(err, "Can't write edge information")
				}
			}
		}
	}

	return nil
}

// ExportVerticesToFile Exports vertices information to CSV-file with header:
// 	vertex_id - int64, ID of vertex
// 	order_pos - int, Position of vertex in hierarchies (evaluted by library)
// 	importance - int, Importance of vertex in graph (evaluted by library)
func (graph *Graph) ExportVerticesToFile(fname string) error {
	fileVertices, err := os.Create(fname)
	if err != nil {
		return errors.Wrap(err, "Can't create vertices file")
	}
	defer fileVertices.Close()
	writerVertices := csv.NewWriter(fileVertices)
	defer writerVertices.Flush()
	writerVertices.Comma = ';'
	err = writerVertices.Write([]string{"vertex_id", "order_pos", "importance"})
	if err != nil {
		return errors.Wrap(err, "Can't write header to vertices file")
	}
	for i := 0; i < len(graph.Vertices); i++ {
		currentVertexExternal := graph.Vertices[i].Label
		err = writerVertices.Write([]string{
			fmt.Sprintf("%d", currentVertexExternal),
			fmt.Sprintf("%d", graph.Vertices[i].orderPos),
			fmt.Sprintf("%d", graph.Vertices[i].importance),
		})
		if err != nil {
			return errors.Wrap(err, "Can't write vertex information")
		}
	}
	return nil
}

// ExportShortcutsToFile Exports shortcuts information to CSV-file with header:
// 	from_vertex_id - int64, ID of source vertex
// 	to_vertex_id - int64, ID of target vertex
// 	weight - float64, Weight of an shortcut
// 	via_vertex_id - int64, ID of vertex through which the shortcut exists
func (graph *Graph) ExportShortcutsToFile(fname string) error {
	fileContractions, err := os.Create(fname)
	if err != nil {
		return errors.Wrap(err, "Can't create shortcuts file")
	}
	defer fileContractions.Close()
	writerContractions := csv.NewWriter(fileContractions)
	defer writerContractions.Flush()
	writerContractions.Comma = ';'
	err = writerContractions.Write([]string{"from_vertex_id", "to_vertex_id", "weight", "via_vertex_id"})
	if err != nil {
		return errors.Wrap(err, "Can't write header to shortucts file")
	}
	for sourceVertexInternal, to := range graph.shortcuts {
		sourceVertexExternal := graph.Vertices[sourceVertexInternal].Label
		for targetVertexInternal, viaNodeInternal := range to {
			targetVertexExternal := graph.Vertices[targetVertexInternal].Label
			viaNodeExternal := graph.Vertices[viaNodeInternal.ViaVertex].Label
			err = writerContractions.Write([]string{
				fmt.Sprintf("%d", sourceVertexExternal),
				fmt.Sprintf("%d", targetVertexExternal),
				strconv.FormatFloat(viaNodeInternal.Cost, 'f', -1, 64),
				fmt.Sprintf("%d", viaNodeExternal),
			})
			if err != nil {
				return errors.Wrap(err, "Can't write shortcut information")
			}
		}
	}
	return nil
}
