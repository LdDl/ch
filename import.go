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
// 		weight - float64, Weight of an edge
// Header of CSV-file containing information about vertices:
// 		vertex_id - int64, ID of vertex
// 		order_pos - int, Position of vertex in hierarchies (evaluted by library)
// 		importance - int, Importance of vertex in graph (evaluted by library)
// Header of CSV-file containing information about shortcuts between vertices:
// 		from_vertex_id - int64, ID of source vertex
// 		to_vertex_id - int64, ID of target vertex
// 		weight - float64, Weight of an shortcut
// 		via_vertex_id - int64, ID of vertex through which the shortcut exists
func ImportFromFile(edgesFname, verticesFname, contractionsFname string) (*Graph, error) {
	// Read edges first
	file, err := os.Open(edgesFname)
	if err != nil {
		return nil, err
	}
	defer file.Close()
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

		weight, err := strconv.ParseFloat(record[2], 64)
		if err != nil {
			return nil, err
		}

		err = graph.CreateVertex(sourceExternal)
		if err != nil {
			return nil, errors.Wrap(err, fmt.Sprintf("Can't add source vertex with external_ID = '%d'", sourceExternal))
		}
		err = graph.CreateVertex(targetExternal)
		if err != nil {
			return nil, errors.Wrap(err, fmt.Sprintf("Can't add target vertex with external_ID = '%d'", targetExternal))
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
	defer fileVertices.Close()
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
		vertexOrderPos, err := strconv.ParseInt(record[1], 10, 64)
		if err != nil {
			return nil, err
		}
		vertexImportance, err := strconv.Atoi(record[2])
		if err != nil {
			return nil, err
		}

		vertexInternal, vertexFound := graph.FindVertex(vertexExternal)
		if !vertexFound {
			return nil, fmt.Errorf("Vertex with Label = %d is not found in graph", vertexExternal)
		}
		graph.Vertices[vertexInternal].SetOrderPos(vertexOrderPos)
		graph.Vertices[vertexInternal].SetImportance(vertexImportance)
	}

	// Read contractions
	fileShortcuts, err := os.Open(contractionsFname)
	if err != nil {
		return nil, err
	}
	defer fileShortcuts.Close()
	readerShortcuts := csv.NewReader(fileShortcuts)
	readerShortcuts.Comma = ';'
	// Skip header of CSV-file
	_, err = readerShortcuts.Read()
	if err != nil {
		return nil, err
	}
	// Read file line by line
	for {
		record, err := readerShortcuts.Read()
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

		weight, err := strconv.ParseFloat(record[2], 64)
		if err != nil {
			return nil, err
		}
		contractionExternal, err := strconv.ParseInt(record[3], 10, 64)
		if err != nil {
			return nil, err
		}

		err = graph.AddEdge(sourceExternal, targetExternal, weight)
		if err != nil {
			return nil, errors.Wrap(err, fmt.Sprintf("Can't add shortcut with source_internal_ID = '%d' and target_internal_ID = '%d'", sourceExternal, targetExternal))
		}

		err = graph.AddShortcut(sourceExternal, targetExternal, contractionExternal, weight)
		if err != nil {
			return nil, errors.Wrap(err, fmt.Sprintf("Can't add shortcut with source_internal_ID = '%d' and target_internal_ID = '%d' to internal map", sourceExternal, targetExternal))
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
