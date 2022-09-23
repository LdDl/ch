package ch

type VertexAlternative struct {
	Label              int64
	AdditionalDistance float64
}

type vertexAlternativeInternal struct {
	vertexNum          int64
	additionalDistance float64
}

const vertexNotFound = -1

func (graph *Graph) vertexAlternativeToInternal(alternative VertexAlternative) vertexAlternativeInternal {
	vertexNum, ok := graph.mapping[alternative.Label]
	if !ok {
		vertexNum = vertexNotFound
	}
	return vertexAlternativeInternal{
		vertexNum:          vertexNum,
		additionalDistance: alternative.AdditionalDistance,
	}
}

func (graph *Graph) vertexAlternativesToInternal(alternatives []VertexAlternative) []vertexAlternativeInternal {
	result := make([]vertexAlternativeInternal, 0, len(alternatives))
	for _, alternative := range alternatives {
		result = append(result, graph.vertexAlternativeToInternal(alternative))
	}
	return result
}
