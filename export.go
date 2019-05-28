package ch

import (
	"container/heap"
	"encoding/csv"
	"io"
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
	err = writer.Write([]string{"from_vertex_id", "to_vertex_id", "from_vertex_internal_id", "to_vertex_internal_id", "edge_weight", "is_contract"})
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

// ImportFromFile Imports graph from file of CSV-format
//
// Header of CSV-file:
// from_vertex_id;to_vertex_id;from_vertex_internal_id;to_vertex_internal_id;edge_weight;contract_id;contract_internal_id
// int;int;int;int;float64;-1 (no contract) else external ID of contracted vertex; -1 (no contract) else internal ID of contracted vertex
//
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

		// log.Println(sourceExternal, targetExternal, sourceInternal, targetInternal, weight)

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
