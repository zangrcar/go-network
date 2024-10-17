package model

import (
	"fmt"
)

type NewGraph struct {
	Nodes map[string]NewNode
	Edges map[int]NewEdge
	Type  string
}

type NewNode struct {
	Attributes map[string]interface{}
	ID         string
}

type NewEdge struct {
	First_node  NewNode
	Second_node NewNode
	Attributes  map[string]interface{}
}

func BasicGraph() NewGraph {
	return NewGraph{
		Nodes: map[string]NewNode{},
		Edges: map[int]NewEdge{},
		Type:  "graph",
	}
}

func DiGraph() NewGraph {
	return NewGraph{
		Nodes: map[string]NewNode{},
		Edges: map[int]NewEdge{},
		Type:  "digraph",
	}
}

func MultiGraph() NewGraph {
	return NewGraph{
		Nodes: map[string]NewNode{},
		Edges: map[int]NewEdge{},
		Type:  "multigraph",
	}
}

func CompareNodes(n1, n2 NewNode) bool {
	if n1.ID != n2.ID /*|| len(n1.Attributes) != len(n2.Attributes)*/ {
		return false
	}

	/*for key, value1 := range n1.Attributes {
		value2, exists := n2.Attributes[key]
		if !exists || value1 != value2 {
			return false
		}
	}*/
	return true
}

func (g NewGraph) AddNode(node NewNode) NewGraph {
	for _, value := range g.Nodes {
		if CompareNodes(node, value) {
			for key, value := range node.Attributes {
				existingvalue, isfound := g.GetNode(node.ID).Attributes[key]
				if isfound {
					v, ok := existingvalue.(int)
					u, okk := value.(int)
					if ok && okk {
						g.GetNode(node.ID).Attributes[key] = v + u
					} else {
						g.GetNode(node.ID).Attributes[key] = value
					}
				} else {
					g.GetNode(node.ID).Attributes[key] = value
				}
			}
			return g
		}
	}
	g.Nodes[node.ID] = node
	return g
}

func (g NewGraph) AddNodesFrom(arr []NewNode) NewGraph {
	for _, node := range arr {
		g.AddNode(node)
	}
	return g
}

func (g NewGraph) HasNode(n NewNode) bool {
	for _, node := range g.Nodes {
		if CompareNodes(n, node) {
			return true
		}
	}
	return false
}

func (g NewGraph) GetNode(id string) *NewNode {
	for _, node := range g.Nodes {
		if node.ID == id {
			return &node
		}
	}
	return nil
}

func (g NewGraph) AddNodeAttribute(node_id string, att string, value interface{}) NewGraph {
	if g.GetNode(node_id) == nil {
		fmt.Println("this graph does not have node with id:", node_id)
		return g
	}
	g.GetNode((node_id)).Attributes[att] = value
	return g
}

func (g NewGraph) NumberOfNodes() int {
	return len(g.Nodes)
}

func (g NewGraph) Neighbors(n NewNode) []NewNode {
	neighbors := []NewNode{}

	for _, edge := range g.Edges {
		if CompareNodes(edge.First_node, n) {
			neighbors = append(neighbors, edge.First_node)
		} else if CompareNodes(edge.Second_node, n) {
			neighbors = append(neighbors, edge.Second_node)
		}
	}

	return neighbors

}

func (g NewGraph) RemoveNode(n NewNode) NewGraph {
	for i, node := range g.Nodes {
		if CompareNodes(node, n) {
			delete(g.Nodes, i)
		}
	}
	return g
}

func (g NewGraph) NodeDegree(n NewNode) int {
	counter := 0
	for _, edge := range g.Edges {
		if CompareNodes(n, edge.First_node) || CompareNodes(n, edge.Second_node) {
			counter++
		}
	}
	return counter
}

func (g NewGraph) CompareEdges(e1, e2 NewEdge) bool {
	if !((CompareNodes(e1.First_node, e2.First_node) && CompareNodes(e1.Second_node, e2.Second_node)) || (CompareNodes(e1.First_node, e2.Second_node) && CompareNodes(e1.Second_node, e2.First_node))) {
		return false
	}

	if g.Type != "multigraph" {
		return true
	}

	if len(e1.Attributes) != len(e2.Attributes) {
		return false
	}

	for key, value1 := range e1.Attributes {
		value2, exists := e2.Attributes[key]
		if !exists || value1 != value2 {
			return false
		}
	}
	return true
}

func (g NewGraph) HasEdge(e NewEdge) bool {
	for _, edge := range g.Edges {
		if g.CompareEdges(e, edge) {
			return true
		}
	}
	return false
}

func (g NewGraph) AddEdge(edge NewEdge) NewGraph {
	for _, value := range g.Edges {
		if g.CompareEdges(edge, value) {
			for key, value := range edge.Attributes {
				existingvalue, isfound := g.GetEdge(edge).Attributes[key]
				if isfound {
					v, ok := existingvalue.(int)
					u, okk := value.(int)
					if ok && okk {
						g.GetEdge(edge).Attributes[key] = v + u
					} else {
						g.GetEdge(edge).Attributes[key] = value
					}
				} else {
					g.GetEdge(edge).Attributes[key] = value
				}
			}
			return g
		}
	}
	if !g.HasNode(edge.First_node) {
		g.AddNode(edge.First_node)
	}
	if !g.HasNode(edge.Second_node) {
		g.AddNode(edge.Second_node)
	}
	g.Edges[len(g.Edges)] = edge
	return g
}

func (g NewGraph) AddEdgesFrom(arr []NewEdge) NewGraph {
	for _, edge := range arr {
		g.AddEdge(edge)
	}
	return g
}

func (g NewGraph) GetEdge(e NewEdge) *NewEdge {
	for _, edge := range g.Edges {
		if g.CompareEdges(e, edge) {
			return &edge
		}
	}
	return nil
}

func (g NewGraph) AddEdgeAttribute(e NewEdge, att string, value interface{}) NewGraph {
	if g.GetEdge(e) == nil {
		fmt.Println("this graph does not have this edge:", e)
		return g
	}
	g.GetEdge(e).Attributes[att] = value
	return g
}

func (g NewGraph) NumberOfEdges() int {
	return len(g.Edges)
}

func (g NewGraph) RemoveEdge(e NewEdge) NewGraph {
	for i, edge := range g.Edges {
		if g.CompareEdges(e, edge) {
			delete(g.Edges, i)
		}
	}
	return g
}
