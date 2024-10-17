package main

import (
	"fmt"

	"github.com/jmCodeCraft/go-network/model"
)

func main() {
	a := model.NewGraph{
		Nodes: map[string]model.NewNode{},
		Edges: map[int]model.NewEdge{},
	}
	node_1 := model.NewNode{
		ID: "node_1",
		Attributes: map[string]interface{}{
			"order": 1,
			"value": 15,
		},
	}

	node_2 := model.NewNode{
		ID: "node_2",
		Attributes: map[string]interface{}{
			"order": 2,
			"value": 12,
		},
	}

	a.AddEdge(model.NewEdge{
		First_node:  node_1,
		Second_node: node_2,
		Attributes:  map[string]interface{}{},
	})

	a.AddNodeAttribute("node_1", "stacionary_vision", 15)

	fmt.Println(a)
}
