package ch

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
)

// ExportToFile Exports graph to file of CSV-format
// Header of main CSV-file:
// 		from_vertex_id - int64, ID of source vertex
// 		to_vertex_id - int64, ID of arget vertex
// 		f_internal - int64, Internal ID of source vertex
// 		t_internal - int64, Internal ID of target vertex
// 		weight - float64, Weight of an edge
// 		via_vertex_id - int64, ID of vertex through which the contraction exists (-1 if no contraction)
// 		v_internal - int64, Internal ID of vertex through which the contraction exists (-1 if no contraction)
func (graph *Graph) ExportToFile(fname string) error {

	file, err := os.Create(fname)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()
	writer.Comma = ';'
	err = writer.Write([]string{"from_vertex_id", "to_vertex_id", "f_internal", "t_internal", "weight", "via_vertex_id", "v_internal"})
	if err != nil {
		return err
	}

	for i := 0; i < len(graph.Vertices); i++ {
		currentVertexExternal := graph.Vertices[i].Label
		currentVertexInternal := graph.Vertices[i].vertexNum

		incomingNeighbors := graph.Vertices[i].inEdges
		incomingCosts := graph.Vertices[i].inECost
		for j := 0; j < len(incomingNeighbors); j++ {
			fromVertexExternal := graph.Vertices[incomingNeighbors[j]].Label
			fromVertexInternal := incomingNeighbors[j]
			cost := incomingCosts[j]
			isContractExternal := int64(-1)
			isContractInternal := int64(-1)
			if v, ok := graph.contracts[fromVertexInternal][currentVertexInternal]; ok {
				isContractExternal = graph.Vertices[v].Label
				isContractInternal = v
			}
			err = writer.Write([]string{
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
		}

		outcomingNeighbors := graph.Vertices[i].outEdges
		outcomingCosts := graph.Vertices[i].outECost
		for j := 0; j < len(outcomingNeighbors); j++ {
			toVertexExternal := graph.Vertices[outcomingNeighbors[j]].Label
			toVertexInternal := outcomingNeighbors[j]
			cost := outcomingCosts[j]
			isContractExternal := int64(-1)
			isContractInternal := int64(-1)
			if v, ok := graph.contracts[currentVertexInternal][toVertexInternal]; ok {
				isContractExternal = graph.Vertices[v].Label
				isContractInternal = v
			}
			err = writer.Write([]string{
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
		}
	}

	return err
}
