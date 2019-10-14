package ch

import (
	"container/heap"
	"encoding/csv"
	"io"
	"os"
	"strconv"
)

// ImportFromFile Imports graph from file of CSV-format
// Header of main CSV-file:
// 		from_vertex_id - int64, ID of source vertex
// 		to_vertex_id - int64, ID of arget vertex
// 		f_internal - int64, Internal ID of source vertex
// 		t_internal - int64, Internal ID of target vertex
// 		weight - float64, Weight of an edge
// 		via_vertex_id - int64, ID of vertex through which the contraction exists (-1 if no contraction)
// 		v_internal - int64, Internal ID of vertex through which the contraction exists (-1 if no contraction)
func ImportFromFile(fname string) (*Graph, error) {
	file, err := os.Open(fname)
	if err != nil {
		return nil, err
	}
	reader := csv.NewReader(file)
	reader.Comma = ';'

	graph := Graph{}

	// skip header
	_, err = reader.Read()
	if err != nil {
		return nil, err
	}
	// read lines
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		sourceExternal, err := strconv.Atoi(record[0])
		if err != nil {
			return nil, err
		}
		targetExternal, err := strconv.Atoi(record[1])
		if err != nil {
			return nil, err
		}

		sourceInternal, err := strconv.Atoi(record[2])
		if err != nil {
			return nil, err
		}
		targetInternal, err := strconv.Atoi(record[3])
		if err != nil {
			return nil, err
		}

		weight, err := strconv.ParseFloat(record[4], 64)
		if err != nil {
			return nil, err
		}

		// isContractExternal, err := strconv.Atoi(record[5])
		// if err != nil {
		// 	return nil, err
		// }
		isContractInternal, err := strconv.Atoi(record[6])
		if err != nil {
			return nil, err
		}

		graph.AddVertex(sourceExternal, sourceInternal)
		graph.AddVertex(targetExternal, targetInternal)

		graph.AddEdge(sourceExternal, targetExternal, weight)

		if isContractInternal != -1 {
			if _, ok := graph.contracts[sourceInternal]; !ok {
				graph.contracts[sourceInternal] = make(map[int]int)
				graph.contracts[sourceInternal][targetInternal] = isContractInternal
			}
			graph.contracts[sourceInternal][targetInternal] = isContractInternal
		}
	}

	// Need to calculate order pos for every vertex to make work relaxEdgesBiForward() and relaxEdgesBiBackward() functions in bidirectional_ch.go
	graph.computeImportance()
	var extractNum int
	for graph.pqImportance.Len() != 0 {
		vertex := heap.Pop(graph.pqImportance).(*Vertex)
		vertex.computeImportance()
		if graph.pqImportance.Len() != 0 && vertex.importance > graph.pqImportance.Peek().(*Vertex).importance {
			graph.pqImportance.Push(vertex)
			continue
		}
		graph.Vertices[vertex.vertexNum].orderPos = extractNum
		extractNum = extractNum + 1
	}
	return &graph, nil
}

// ImportRestrictionsFromFile Imports turn restrictions from file of CSV-format into graph
//
// Header of CSV-file:
// from_vertex_id;via_vertex_id;to_vertex_id;
// int;int;int
//
func (g *Graph) ImportRestrictionsFromFile(fname string) error {
	file, err := os.Open(fname)
	if err != nil {
		return err
	}
	reader := csv.NewReader(file)
	reader.Comma = ';'

	// skip header
	_, err = reader.Read()
	if err != nil {
		return err
	}
	// read lines
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		sourceExternal, err := strconv.Atoi(record[0])
		if err != nil {
			return err
		}
		viaExternal, err := strconv.Atoi(record[1])
		if err != nil {
			return err
		}
		targetExternal, err := strconv.Atoi(record[2])
		if err != nil {
			return err
		}

		g.AddTurnRestriction(sourceExternal, viaExternal, targetExternal)
	}
	return nil
}
