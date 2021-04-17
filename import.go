package ch

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strconv"

	"github.com/pkg/errors"
)

// ImportFromFile Imports graph (prepared by ExportToFile(fname string) function) from file of CSV-format
// Header of CSV-file containing information about edges:
// 		from_vertex_id - int64, ID of source vertex
// 		to_vertex_id - int64, ID of arget vertex
// 		f_internal - int64, Internal ID of source vertex
// 		t_internal - int64, Internal ID of target vertex
// 		weight - float64, Weight of an edge
// Header of CSV-file containing information about vertices:
// 		vertex_id - int64, ID of vertex
// 		internal_id - int64, internal ID of target vertex
// 		order_pos - int, Position of vertex in hierarchies (evaluted by library)
// 		importance - int, Importance of vertex in graph (evaluted by library)
// Header of CSV-file containing information about contractios between vertices:
// 		from_vertex_id - int64, ID of source vertex
// 		to_vertex_id - int64, ID of arget vertex
// 		f_internal - int64, Internal ID of source vertex
// 		t_internal - int64, Internal ID of target vertex
// 		weight - float64, Weight of an edge
// 		via_vertex_id - int64, ID of vertex through which the contraction exists
// 		v_internal - int64, Internal ID of vertex through which the contraction exists
func ImportFromFile(edgesFname, verticesFname, contractionsFname string) (*Graph, error) {
	// Read edges first
	file, err := os.Open(edgesFname)
	if err != nil {
		return nil, err
	}
	reader := csv.NewReader(file)
	reader.Comma = ';'

	graph := Graph{}

	// Fill graph with edges informations
	// Skip header of CSV-file
	_, err = reader.Read()
	if err != nil {
		return nil, err
	}
	// Read file line by line
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
	}

	// Read vertices
	fileVertices, err := os.Open(verticesFname)
	if err != nil {
		return nil, err
	}
	readerVertices := csv.NewReader(fileVertices)
	readerVertices.Comma = ';'

	// Skip header of CSV-file
	_, err = readerVertices.Read()
	if err != nil {
		return nil, err
	}
	// Read file line by line
	for {
		record, err := readerVertices.Read()
		if err == io.EOF {
			break
		}

		vertexExternal, err := strconv.ParseInt(record[0], 10, 64)
		if err != nil {
			return nil, err
		}
		vertexInternal, err := strconv.ParseInt(record[1], 10, 64)
		if err != nil {
			return nil, err
		}
		vertexOrderPos, err := strconv.Atoi(record[2])
		if err != nil {
			return nil, err
		}
		vertexImportance, err := strconv.Atoi(record[3])
		if err != nil {
			return nil, err
		}

		if graph.Vertices[vertexInternal].Label != vertexExternal {
			return nil, fmt.Errorf("Vertex with Label = %d has wrong reference information. Incoming label info is '%d'", graph.Vertices[vertexInternal].Label, vertexExternal)
		}
		graph.Vertices[vertexInternal].SetOrderPos(vertexOrderPos)
		graph.Vertices[vertexInternal].SetImportance(vertexImportance)
	}

	// Read contractions
	fileContractions, err := os.Open(contractionsFname)
	if err != nil {
		return nil, err
	}
	readerContractions := csv.NewReader(fileContractions)
	readerContractions.Comma = ';'
	// Skip header of CSV-file
	_, err = readerContractions.Read()
	if err != nil {
		return nil, err
	}
	// Read file line by line
	for {
		record, err := readerContractions.Read()
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

		contractionInternal, err := strconv.ParseInt(record[6], 10, 64)
		if err != nil {
			return nil, err
		}

		err = graph.AddEdge(sourceExternal, targetExternal, weight)
		if err != nil {
			return nil, errors.Wrap(err, fmt.Sprintf("Can't add edge with source_internal_ID = '%d' and target_internal_ID = '%d'", sourceExternal, targetExternal))
		}
		if _, ok := graph.shortcuts[sourceInternal]; !ok {
			graph.shortcuts[sourceInternal] = make(map[int64]*ContractionPath)
			graph.shortcuts[sourceInternal][targetInternal] = &ContractionPath{
				ViaVertex: contractionInternal,
				Cost:      weight,
			}
		}
		graph.shortcuts[sourceInternal][targetInternal] = &ContractionPath{
			ViaVertex: contractionInternal,
			Cost:      weight,
		}
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
