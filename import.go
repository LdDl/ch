package ch

import (
	"container/heap"
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strconv"

	"github.com/pkg/errors"
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
		sourceExternal, err := strconv.ParseInt(record[0], 10, 64)
		if err != nil {
			return nil, err
		}
		targetExternal, err := strconv.ParseInt(record[1], 10, 64)
		if err != nil {
			return nil, err
		}

		sourceInternal, err := strconv.ParseInt(record[2], 10, 64)
		if err != nil {
			return nil, err
		}
		targetInternal, err := strconv.ParseInt(record[3], 10, 64)
		if err != nil {
			return nil, err
		}

		weight, err := strconv.ParseFloat(record[4], 64)
		if err != nil {
			return nil, err
		}

		isContractInternal, err := strconv.ParseInt(record[6], 10, 64)
		if err != nil {
			return nil, err
		}

		err = graph.AddVertex(sourceExternal, sourceInternal)
		if err != nil {
			return nil, errors.Wrap(err, fmt.Sprintf("Can't add source vertex with external_ID = '%d' and internal_ID = '%d'", sourceExternal, sourceInternal))
		}
		err = graph.AddVertex(targetExternal, targetInternal)
		if err != nil {
			return nil, errors.Wrap(err, fmt.Sprintf("Can't add target vertex with external_ID = '%d' and internal_ID = '%d'", targetExternal, targetInternal))
		}

		err = graph.AddEdge(sourceExternal, targetExternal, weight)
		if err != nil {
			return nil, errors.Wrap(err, fmt.Sprintf("Can't add edge with source_internal_ID = '%d' and target_internal_ID = '%d'", sourceExternal, targetExternal))
		}
		if isContractInternal != -1 {
			if _, ok := graph.contracts[sourceInternal]; !ok {
				graph.contracts[sourceInternal] = make(map[int64]int64)
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
		sourceExternal, err := strconv.ParseInt(record[0], 10, 64)
		if err != nil {
			return err
		}
		viaExternal, err := strconv.ParseInt(record[1], 10, 64)
		if err != nil {
			return err
		}
		targetExternal, err := strconv.ParseInt(record[2], 10, 64)
		if err != nil {
			return err
		}

		err = g.AddTurnRestriction(sourceExternal, viaExternal, targetExternal)
		if err != nil {
			return errors.Wrap(err, fmt.Sprintf("Can't add restriction between source_external_ID = '%d' and target_external_ID = '%d' via via_external_id = '%d'", sourceExternal, targetExternal, viaExternal))
		}
	}
	return nil
}
