package ch

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// ExportToFile Exports graph to file of CSV-format
// Header of edges CSV-file:
// 		from_vertex_id - int64, ID of source vertex
// 		to_vertex_id - int64, ID of target vertex
// 		f_internal - int64, Internal ID of source vertex
// 		t_internal - int64, Internal ID of target vertex
// 		weight - float64, Weight of an edge
// Header of vertices CSV-file:
// 		vertex_id - int64, ID of vertex
// 		internal_id - int64, Internal ID of vertex
// 		order_pos - int, Position of vertex in hierarchies (evaluted by library)
// 		importance - int, Importance of vertex in graph (evaluted by library)
// Header of contractios CSV-file:
// 		from_vertex_id - int64, ID of source vertex
// 		to_vertex_id - int64, ID of arget vertex
// 		f_internal - int64, Internal ID of source vertex
// 		t_internal - int64, Internal ID of target vertex
// 		weight - float64, Weight of an edge
// 		via_vertex_id - int64, ID of vertex through which the contraction exists
// 		v_internal - int64, Internal ID of vertex through which the contraction exists
func (graph *Graph) ExportToFile(fname string) error {

	fnamePart := strings.Split(fname, ".csv") // to guarantee proper filename and its extension
	file, err := os.Create(fnamePart[0] + ".csv")
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()
	writer.Comma = ';'
	err = writer.Write([]string{"from_vertex_id", "to_vertex_id", "f_internal", "t_internal", "weight"})
	if err != nil {
		return err
	}

	fileVertices, err := os.Create(fnamePart[0] + "_vertices.csv")
	if err != nil {
		return err
	}
	defer fileVertices.Close()

	writerVertices := csv.NewWriter(fileVertices)
	defer writerVertices.Flush()
	writerVertices.Comma = ';'
	err = writerVertices.Write([]string{"vertex_id", "internal_id", "order_pos", "importance"})
	if err != nil {
		return err
	}

	fileContractions, err := os.Create(fnamePart[0] + "_contractions.csv")
	if err != nil {
		return err
	}
	defer fileContractions.Close()

	writerContractions := csv.NewWriter(fileContractions)
	defer writerContractions.Flush()
	writerContractions.Comma = ';'
	err = writerContractions.Write([]string{"from_vertex_id", "to_vertex_id", "f_internal", "t_internal", "weight", "via_vertex_id", "v_internal"})
	if err != nil {
		return err
	}

	for i := 0; i < len(graph.Vertices); i++ {
		currentVertexExternal := graph.Vertices[i].Label
		currentVertexInternal := graph.Vertices[i].vertexNum

		// Write reference information about vertex
		err = writerVertices.Write([]string{
			fmt.Sprintf("%d", currentVertexExternal),
			fmt.Sprintf("%d", currentVertexInternal),
			fmt.Sprintf("%d", graph.Vertices[i].orderPos),
			fmt.Sprintf("%d", graph.Vertices[i].importance),
		})
		if err != nil {
			return err
		}

		// Write reference information about "incoming" adjacent vertices
		incomingNeighbors := graph.Vertices[i].inEdges
		incomingCosts := graph.Vertices[i].inECost
		for j := 0; j < len(incomingNeighbors); j++ {
			fromVertexExternal := graph.Vertices[incomingNeighbors[j]].Label
			fromVertexInternal := incomingNeighbors[j]
			cost := incomingCosts[j]
			if v, ok := graph.contracts[fromVertexInternal][currentVertexInternal]; ok {
				isContractExternal := graph.Vertices[v].Label
				isContractInternal := v
				err = writerContractions.Write([]string{
					fmt.Sprintf("%d", fromVertexExternal),
					fmt.Sprintf("%d", currentVertexExternal),
					fmt.Sprintf("%d", fromVertexInternal),
					fmt.Sprintf("%d", currentVertexInternal),
					strconv.FormatFloat(cost, 'f', -1, 64),
					fmt.Sprintf("%d", isContractExternal),
					fmt.Sprintf("%d", isContractInternal),
				})
				if err != nil {
					return err
				}
			} else {
				err = writer.Write([]string{
					fmt.Sprintf("%d", fromVertexExternal),
					fmt.Sprintf("%d", currentVertexExternal),
					fmt.Sprintf("%d", fromVertexInternal),
					fmt.Sprintf("%d", currentVertexInternal),
					strconv.FormatFloat(cost, 'f', -1, 64),
				})
				if err != nil {
					return err
				}
			}
		}

		// Write reference information about "outcoming" adjacent vertices
		outcomingNeighbors := graph.Vertices[i].outEdges
		outcomingCosts := graph.Vertices[i].outECost
		for j := 0; j < len(outcomingNeighbors); j++ {
			toVertexExternal := graph.Vertices[outcomingNeighbors[j]].Label
			toVertexInternal := outcomingNeighbors[j]
			cost := outcomingCosts[j]
			if v, ok := graph.contracts[currentVertexInternal][toVertexInternal]; ok {
				isContractExternal := graph.Vertices[v].Label
				isContractInternal := v
				err = writerContractions.Write([]string{
					fmt.Sprintf("%d", currentVertexExternal),
					fmt.Sprintf("%d", toVertexExternal),
					fmt.Sprintf("%d", currentVertexInternal),
					fmt.Sprintf("%d", toVertexInternal),
					strconv.FormatFloat(cost, 'f', -1, 64),
					fmt.Sprintf("%d", isContractExternal),
					fmt.Sprintf("%d", isContractInternal),
				})
				if err != nil {
					return err
				}
			} else {
				err = writer.Write([]string{
					fmt.Sprintf("%d", currentVertexExternal),
					fmt.Sprintf("%d", toVertexExternal),
					fmt.Sprintf("%d", currentVertexInternal),
					fmt.Sprintf("%d", toVertexInternal),
					strconv.FormatFloat(cost, 'f', -1, 64),
				})
				if err != nil {
					return err
				}
			}
		}
	}

	return err
}
