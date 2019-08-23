package ch

import (
	"encoding/csv"
	"os"
	"strconv"
)

// ExportToFile Exports graph to file of CSV-format
//
// Header of CSV-file:
// from_vertex_id;to_vertex_id;from_vertex_internal_id;to_vertex_internal_id;edge_weight;contract_id;contract_internal_id
// int;int;int;int;float64;-1 (no contract) else external ID of contracted vertex; -1 (no contract) else internal ID of contracted vertex
//
func (graph *Graph) ExportToFile(fname string) error {

	file, err := os.Create(fname)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()
	writer.Comma = ';'
	err = writer.Write([]string{"from_vertex_id", "to_vertex_id", "from_vertex_internal_id", "to_vertex_internal_id", "edge_weight", "contract_id", "contract_internal_id"})
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
			isContractExternal := -1
			isContractInternal := -1
			if v, ok := graph.contracts[fromVertexInternal][currentVertexInternal]; ok {
				isContractExternal = graph.Vertices[v].Label
				isContractInternal = v
			}
			err = writer.Write([]string{
				strconv.Itoa(fromVertexExternal),
				strconv.Itoa(currentVertexExternal),
				strconv.Itoa(fromVertexInternal),
				strconv.Itoa(currentVertexInternal),
				strconv.FormatFloat(cost, 'f', -1, 64),
				strconv.Itoa(isContractExternal),
				strconv.Itoa(isContractInternal),
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
			isContractExternal := -1
			isContractInternal := -1
			if v, ok := graph.contracts[currentVertexInternal][toVertexInternal]; ok {
				isContractExternal = graph.Vertices[v].Label
				isContractInternal = v
			}
			err = writer.Write([]string{
				strconv.Itoa(currentVertexExternal),
				strconv.Itoa(toVertexExternal),
				strconv.Itoa(currentVertexInternal),
				strconv.Itoa(toVertexInternal),
				strconv.FormatFloat(cost, 'f', -1, 64),
				strconv.Itoa(isContractExternal),
				strconv.Itoa(isContractInternal),
			})
			if err != nil {
				return err
			}
		}
	}

	return err
}
