package graph

import (
	"errors"
	"fmt"
)

type Graph struct {
	edges    []*Edge
	vertices map[string]*Vertex
}

type GraphData struct {
	From  string     `json:"from,omitempty"`
	To    string     `json:"to,omitempty"`
	Infos []EdgeInfo `json:"infos"`
}

func (v *Graph) Init() {
	v.vertices = map[string]*Vertex{}
	v.edges = []*Edge{}
}

func (v *Graph) AddVertex(key string) *Vertex {
	vertex := &Vertex{}
	vertex.Init(key)

	v.vertices[vertex.GetKey()] = vertex

	return vertex
}

func (v *Graph) AddEdge(edge *Edge) error {

	if !v.HasVertex(edge.From.key) {
		return errors.New(fmt.Sprintf("The Vertex with key \"%s\" not exists on Graph", edge.From.key))
	}

	if !v.HasVertex(edge.To.key) {
		return errors.New(fmt.Sprintf("The Vertex with key \"%s\" not exists on Graph", edge.To.key))
	}

	v.edges = append(v.edges, edge)

	return nil
}

func (v *Graph) HasVertex(key string) bool {
	for _, value := range v.getVerticesKey() {
		if value == key {
			return true
		}
	}
	return false
}

func (v *Graph) getVerticesKey() []string {
	keys := make([]string, len(v.vertices))

	for key := range v.vertices {
		keys = append(keys, key)
	}

	return keys
}

func (v *Graph) ToJSON() []GraphData {
	data := make([]GraphData, len(v.edges))

	for index, edge := range v.edges {
		data[index] = GraphData{
			From:  edge.From.key,
			To:    edge.To.key,
			Infos: edge.Infos,
		}
	}

	return data
}
