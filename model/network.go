package model

import (
	"fmt"
	"reflect"
	"strings"
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

func (g NewGraph) IsEqual(other NewGraph) bool {

	if g.Type != other.Type || len(g.Nodes) != len(other.Nodes) || len(g.Edges) != len(other.Edges) {
		return false
	}

	for key := range g.Nodes {
		_, ok := other.Nodes[key]
		if !ok {
			return false
		}
	}

	for i := 0; i < len(g.Edges); i++ {
		if !g.CompareEdges(g.Edges[i], other.Edges[i]) {
			return false
		}
	}

	return true
}

func (g NewGraph) Combine(h NewGraph) NewGraph {
	if g.Type != h.Type {
		fmt.Println("Graphs must be the same type!")
		return NewGraph{}
	}

	k := NewGraph{
		Nodes: map[string]NewNode{},
		Edges: map[int]NewEdge{},
		Type:  g.Type,
	}

	for key, value := range g.Nodes {
		k.Nodes[key] = value
	}
	for key, value := range h.Nodes {
		k.Nodes[key] = value
	}
	for key, value := range g.Edges {
		k.Edges[key] = value
	}
	for key, value := range h.Edges {
		k.Edges[key] = value
	}
	return k
}

func (g NewGraph) ToString() string {
	var str strings.Builder
	str.WriteString(fmt.Sprintf("Type: %s\n", g.Type))
	str.WriteString("Nodes:\n")
	for key, value := range g.Nodes {
		str.WriteString(fmt.Sprintf("- ID: %s, attributes: %+v\n", key, value.Attributes))
	}
	str.WriteString("Edges:\n")
	for _, value := range g.Edges {
		str.WriteString(fmt.Sprintf("- node1_ID: %s, node2_ID: %s, attributes: %+v\n", value.First_node.ID, value.Second_node.ID, value.Attributes))
	}
	return str.String()
}

func (g NewGraph) NodeToString(node NewNode) string {
	var str strings.Builder
	str.WriteString("Nodes:\n")
	str.WriteString(fmt.Sprintf("- ID: %s, attributes: %+v\n", node.ID, node.Attributes))
	return str.String()
}

func (g NewGraph) NewDFS(startNode NewNode) NewGraph {
	_, exists := g.Nodes[startNode.ID]
	if !exists {
		return NewGraph{}
	}
	visited := make(map[string]bool)
	visitedGraph := NewGraph{
		Nodes: make(map[string]NewNode),
		Edges: make(map[int]NewEdge),
	}
	g.NewDfsUtil(startNode, visited, visitedGraph)
	return visitedGraph
}

func (g NewGraph) NewDfsUtil(node NewNode, visited map[string]bool, visitedGraph NewGraph) {
	visited[node.ID] = true
	visitedGraph.AddNode(node)

	for _, neighbor := range g.Neighbors(node) {
		if !visited[neighbor.ID] {
			e := g.GetEdgeByNodes(node, neighbor)
			visitedGraph.AddEdge(NewEdge{First_node: node, Second_node: neighbor, Attributes: e.Attributes})
			g.NewDfsUtil(neighbor, visited, visitedGraph)
		}
	}
}

func (g NewGraph) GetComponents() []NewGraph {
	visited := make(map[string]bool)
	for key := range g.Nodes {
		visited[key] = false
	}
	components := []NewGraph{}
	x := GetTrueString(visited)
	for x != "" {
		components = append(components, g.NewDFS(*g.GetNode(x)))
	}
	return components
}

func GetTrueString(x map[string]bool) string {
	for key, value := range x {
		if !value {
			return key
		}
	}
	return ""
}

func CompareNodes(n1, n2 NewNode) bool {
	return n1.ID == n2.ID
	/*if n1.ID != n2.ID || len(n1.Attributes) != len(n2.Attributes) {
		return false
	}

	for key, value1 := range n1.Attributes {
		value2, exists := n2.Attributes[key]
		if !exists || value1 != value2 {
			return false
		}
	}
	return true*/
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

func (g NewGraph) ContractNewNode(n NewNode) {
	if !g.HasNode(n) {
		fmt.Println("This graph does not have this node.")
		return
	}
	neighbors := g.Neighbors(n)
	for i := 0; i < len(neighbors); i++ {
		for j := i + 1; j < len(neighbors); j++ {
			e := NewEdge{
				First_node:  neighbors[i],
				Second_node: neighbors[j],
				Attributes:  map[string]interface{}{},
			}
			g.AddEdge(e)
		}
	}
	g.RemoveNode(n)
}

type CombineStrategy interface {
	CombineInt(n1, n2 interface{}) interface{}
	CombineFloat32(n1, n2 interface{}) interface{}
	CombineFloat64(n1, n2 interface{}) interface{}
	CombineString(n1, n2 interface{}) interface{}
}

type StrategyAvgNum struct{}
type StrategyArray struct{}
type StrategyRetainMax struct{}
type StrategyRetainMin struct{}

func (s StrategyRetainMin) CombineInt(n1, n2 int) int {
	if n1 <= n2 {
		return n1
	}
	return n2
}

func (s StrategyRetainMin) CombineFloat32(n1, n2 float32) float32 {
	if n1 <= n2 {
		return n1
	}
	return n2
}

func (s StrategyRetainMin) CombineFloat64(n1, n2 float64) float64 {
	if n1 <= n2 {
		return n1
	}
	return n2
}

func (s StrategyRetainMin) CombineString(n1, n2 string) string {
	if n1 <= n2 {
		return n1
	}
	return n2
}

func (s StrategyRetainMax) CombineInt(n1, n2 int) int {
	if n1 >= n2 {
		return n1
	}
	return n2
}

func (s StrategyRetainMax) CombineFloat32(n1, n2 float32) float32 {
	if n1 >= n2 {
		return n1
	}
	return n2
}

func (s StrategyRetainMax) CombineFloat64(n1, n2 float64) float64 {
	if n1 >= n2 {
		return n1
	}
	return n2
}

func (s StrategyRetainMax) CombineString(n1, n2 string) string {
	if n1 >= n2 {
		return n1
	}
	return n2
}

func (s StrategyAvgNum) CombineInt(n1, n2 int) int {
	return (n1 + n2) / 2
}

func (s StrategyAvgNum) CombineFloat32(n1, n2 float32) float32 {
	return (n1 + n2) / 2
}

func (s StrategyAvgNum) CombineFloat64(n1, n2 float64) float64 {
	return (n1 + n2) / 2
}

func (s StrategyArray) CombineInt(n1, n2 interface{}) []int {
	if !(reflect.TypeOf(n1).Kind() == reflect.Slice) {
		nv := n1.(int)
		n1 = []int{nv}
	}
	if !(reflect.TypeOf(n2).Kind() == reflect.Slice) {
		nv := n2.(int)
		n2 = []int{nv}
	}
	return append(n1.([]int), n2.([]int)...)
}

func (s StrategyArray) CombineFloat32(n1, n2 interface{}) []float32 {
	if !(reflect.TypeOf(n1).Kind() == reflect.Slice) {
		nv := n1.(float32)
		n1 = []float32{nv}
	}
	if !(reflect.TypeOf(n2).Kind() == reflect.Slice) {
		nv := n2.(float32)
		n2 = []float32{nv}
	}
	return append(n1.([]float32), n2.([]float32)...)
}

func (s StrategyArray) CombineFloat64(n1, n2 interface{}) []float64 {
	if !(reflect.TypeOf(n1).Kind() == reflect.Slice) {
		nv := n1.(float64)
		n1 = []float64{nv}
	}
	if !(reflect.TypeOf(n2).Kind() == reflect.Slice) {
		nv := n2.(float64)
		n2 = []float64{nv}
	}
	return append(n1.([]float64), n2.([]float64)...)
}

func (s StrategyArray) CombineString(n1, n2 interface{}) []string {
	if !(reflect.TypeOf(n1).Kind() == reflect.Slice) {
		nv := n1.(string)
		n1 = []string{nv}
	}
	if !(reflect.TypeOf(n2).Kind() == reflect.Slice) {
		nv := n2.(string)
		n2 = []string{nv}
	}
	return append(n1.([]string), n2.([]string)...)
}

func (g NewGraph) CombineNodes(n1, n2 NewNode, strat_num, strat_string CombineStrategy) NewNode {
	id := ""
	if n1.ID < n2.ID {
		id += n1.ID + "_" + n2.ID
	} else {
		id += n2.ID + "_" + n1.ID
	}

	att := map[string]interface{}{}
	for key, value := range n1.Attributes {
		att[key] = value
	}
	for key, newValue := range n2.Attributes {
		if existingValue, found := att[key]; found {
			switch v := existingValue.(type) {
			case int:
				if nv, ok := newValue.(int); ok {
					att[key] = strat_num.CombineInt(v, nv)
				}
			case float64:
				if nv, ok := newValue.(float64); ok {
					att[key] = strat_num.CombineFloat64(v, nv)
				}
			case string:
				if nv, ok := newValue.(string); ok {
					strat_string.CombineString(v, nv)
				}
			}
		} else {
			att[key] = newValue
		}
	}

	return NewNode{ID: id, Attributes: att}
}

func (g NewGraph) ContractNewEdge(e NewEdge, strategy_num CombineStrategy, strategy_string CombineStrategy) {
	if !g.HasEdge(e) {
		fmt.Println("This graph does not have this edge.")
		return
	}
	n := g.CombineNodes(e.First_node, e.Second_node, strategy_num, strategy_string)
	g.AddNode(n)

	first := g.Neighbors(e.First_node)
	sec := g.Neighbors(e.Second_node)

	for x := 0; x < len(first); x++ {
		g.AddEdge(NewEdge{
			First_node:  n,
			Second_node: first[x],
			Attributes:  map[string]interface{}{},
		})
	}

	for x := 0; x < len(sec); x++ {
		g.AddEdge(NewEdge{
			First_node:  n,
			Second_node: sec[x],
			Attributes:  map[string]interface{}{},
		})
	}

	g.RemoveNode(e.First_node)
	g.RemoveNode(e.Second_node)

}

func (g NewGraph) EdgeID(e NewEdge) string {
	if e.First_node.ID < e.Second_node.ID || g.Type == "digraph" {
		return e.First_node.ID + "-" + e.Second_node.ID
	} else {
		return e.Second_node.ID + "-" + e.First_node.ID
	}
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

func (g NewGraph) GetEdgeByNodes(n1, n2 NewNode) NewEdge {
	for _, edge := range g.Edges {
		if (CompareNodes(edge.First_node, n1) && CompareNodes(edge.Second_node, n2)) || (CompareNodes(edge.First_node, n2) && CompareNodes(edge.Second_node, n1)) {
			return edge
		}
	}
	fmt.Println("this graph does not have edge between these two nodes.")
	return NewEdge{First_node: n1, Second_node: n2, Attributes: map[string]interface{}{}}
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
