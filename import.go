package ch

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"sort"
	"strconv"
	"strings"

	"github.com/pkg/errors"
)

var (
	ErrNotEnoughColumns = fmt.Errorf("not enough columns")
	ErrColumnNotFound   = fmt.Errorf("column not found")
)

// CSVHeaderImportEdges is just an helper structure to evaluate CSV columns for edges file
type CSVHeaderImportEdges struct {
	SourceExternal int
	TargetExternal int
	Weight         int
}

// CSVHeaderImportVertices is just an helper structure to evaluate CSV columns for vertices file
type CSVHeaderImportVertices struct {
	ID         int
	OrderPos   int
	Importance int
}

// CSVHeaderImportShortcuts is just an helper structure to evaluate CSV columns for shortcuts file
type CSVHeaderImportShortcuts struct {
	SourceExternal int
	TargetExternal int
	ViaExternal    int
	Weight         int
}

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
	edgesHeader, err := reader.Read()
	if err != nil {
		return nil, err
	}
	edgesColumns, err := prepareEdgesColumns(edgesHeader)
	if err != nil {
		return nil, err
	}
	// Read file line by line
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		sourceExternal, err := strconv.ParseInt(record[edgesColumns.SourceExternal], 10, 64)
		if err != nil {
			return nil, err
		}
		targetExternal, err := strconv.ParseInt(record[edgesColumns.TargetExternal], 10, 64)
		if err != nil {
			return nil, err
		}

		weight, err := strconv.ParseFloat(record[edgesColumns.Weight], 64)
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
	verticesHeader, err := readerVertices.Read()
	if err != nil {
		return nil, err
	}
	verticesColumns, err := prepareVerticesColumns(verticesHeader)
	if err != nil {
		return nil, err
	}
	// Read file line by line
	for {
		record, err := readerVertices.Read()
		if err == io.EOF {
			break
		}

		vertexExternal, err := strconv.ParseInt(record[verticesColumns.ID], 10, 64)
		if err != nil {
			return nil, err
		}
		vertexOrderPos, err := strconv.ParseInt(record[verticesColumns.OrderPos], 10, 64)
		if err != nil {
			return nil, err
		}
		vertexImportance, err := strconv.Atoi(record[verticesColumns.Importance])
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
	// Process header of CSV-file
	shortcutsHeader, err := readerShortcuts.Read()
	if err != nil {
		return nil, err
	}
	shortcutsColumns, err := prepareShortcutsColumns(shortcutsHeader)
	if err != nil {
		return nil, err
	}
	// Read file line by line
	for {
		record, err := readerShortcuts.Read()
		if err == io.EOF {
			break
		}
		sourceExternal, err := strconv.ParseInt(record[shortcutsColumns.SourceExternal], 10, 64)
		if err != nil {
			return nil, err
		}
		targetExternal, err := strconv.ParseInt(record[shortcutsColumns.TargetExternal], 10, 64)
		if err != nil {
			return nil, err
		}

		weight, err := strconv.ParseFloat(record[shortcutsColumns.Weight], 64)
		if err != nil {
			return nil, err
		}
		contractionExternal, err := strconv.ParseInt(record[shortcutsColumns.ViaExternal], 10, 64)
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

	// Finalize import: build contractionOrder, set chPrepared, freeze graph
	graph.FinalizeImport()

	return &graph, nil
}

// FinalizeImport should be called after manually importing a pre-computed CH graph.
// It builds the contractionOrder from vertices' orderPos values, sets up recustomization
// support structures, marks CH as prepared, and freezes the graph.
//
// Use this when you have your own import logic (e.g., reading additional data like
// GeoJSON coordinates) instead of using ImportFromFile().
//
// Example:
//
//	graph := ch.NewGraph()
//	// ... your custom import logic: CreateVertex, AddEdge, SetOrderPos, SetImportance, AddShortcut ...
//	graph.FinalizeImport()
//	// Now graph is ready for queries and recustomization
func (graph *Graph) FinalizeImport() {
	graph.buildContractionOrder()
	graph.chPrepared = true
	graph.Freeze()
}

// buildContractionOrder reconstructs contractionOrder slice from vertices' orderPos values.
func (graph *Graph) buildContractionOrder() {
	n := len(graph.Vertices)
	graph.contractionOrder = make([]int64, n)

	// Create slice of (vertexNum, orderPos) pairs and sort by orderPos
	type vertexOrder struct {
		vertexNum int64
		orderPos  int64
	}
	vertices := make([]vertexOrder, n)
	for i := range graph.Vertices {
		vertices[i] = vertexOrder{
			vertexNum: graph.Vertices[i].vertexNum,
			orderPos:  graph.Vertices[i].orderPos,
		}
	}
	sort.Slice(vertices, func(i, j int) bool {
		return vertices[i].orderPos < vertices[j].orderPos
	})

	// Build contractionOrder in correct order
	for i, v := range vertices {
		graph.contractionOrder[i] = v.vertexNum
	}
}

func prepareEdgesColumns(edgesHeader []string) (CSVHeaderImportEdges, error) {
	ans := CSVHeaderImportEdges{
		SourceExternal: -1,
		TargetExternal: -1,
		Weight:         -1,
	}
	if len(edgesHeader) < 3 {
		return ans, errors.Wrapf(ErrNotEnoughColumns, "Minimum 3 columns are needed. Provided: %d", len(edgesHeader))
	}
	for i, header := range edgesHeader {
		switch strings.ToLower(header) {
		case "from_vertex_id":
			ans.SourceExternal = i
		case "to_vertex_id":
			ans.TargetExternal = i
		case "weight":
			ans.Weight = i
		default:
			// Nothing
		}
	}
	if ans.SourceExternal < 0 {
		return ans, errors.Wrap(ErrColumnNotFound, "Not found: 'from_vertex_id'")
	}
	if ans.TargetExternal < 0 {
		return ans, errors.Wrap(ErrColumnNotFound, "Not found: 'to_vertex_id'")
	}
	if ans.Weight < 0 {
		return ans, errors.Wrap(ErrColumnNotFound, "Not found: 'weight'")
	}
	return ans, nil
}

func prepareVerticesColumns(verticesHeader []string) (CSVHeaderImportVertices, error) {
	ans := CSVHeaderImportVertices{
		ID:         -1,
		OrderPos:   -1,
		Importance: -1,
	}
	if len(verticesHeader) < 3 {
		return ans, errors.Wrapf(ErrNotEnoughColumns, "Minimum 3 columns are needed. Provided: %d", len(verticesHeader))
	}
	for i, header := range verticesHeader {
		switch strings.ToLower(header) {
		case "vertex_id":
			ans.ID = i
		case "order_pos":
			ans.OrderPos = i
		case "importance":
			ans.Importance = i
		default:
			// Nothing
		}
	}
	if ans.ID < 0 {
		return ans, errors.Wrap(ErrColumnNotFound, "Not found: 'vertex_id'")
	}
	if ans.OrderPos < 0 {
		return ans, errors.Wrap(ErrColumnNotFound, "Not found: 'order_pos'")
	}
	if ans.Importance < 0 {
		return ans, errors.Wrap(ErrColumnNotFound, "Not found: 'importance'")
	}
	return ans, nil
}

func prepareShortcutsColumns(verticesHeader []string) (CSVHeaderImportShortcuts, error) {
	ans := CSVHeaderImportShortcuts{
		SourceExternal: -1,
		TargetExternal: -1,
		ViaExternal:    -1,
		Weight:         -1,
	}
	if len(verticesHeader) < 4 {
		return ans, errors.Wrapf(ErrNotEnoughColumns, "Minimum 4 columns are needed. Provided: %d", len(verticesHeader))
	}
	for i, header := range verticesHeader {
		switch strings.ToLower(header) {
		case "from_vertex_id":
			ans.SourceExternal = i
		case "to_vertex_id":
			ans.TargetExternal = i
		case "via_vertex_id":
			ans.ViaExternal = i
		case "weight":
			ans.Weight = i
		default:
			// Nothing
		}
	}
	if ans.SourceExternal < 0 {
		return ans, errors.Wrap(ErrColumnNotFound, "Not found: 'from_vertex_id'")
	}
	if ans.TargetExternal < 0 {
		return ans, errors.Wrap(ErrColumnNotFound, "Not found: 'to_vertex_id'")
	}
	if ans.ViaExternal < 0 {
		return ans, errors.Wrap(ErrColumnNotFound, "Not found: 'via_vertex_id'")
	}
	if ans.Weight < 0 {
		return ans, errors.Wrap(ErrColumnNotFound, "Not found: 'weight'")
	}
	return ans, nil
}
