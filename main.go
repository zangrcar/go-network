package main

import (
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/PuerkitoBio/goquery"
	"github.com/jmCodeCraft/go-network/model"
	"github.com/nvkp/turtle"
)

//"github.com/jmCodeCraft/go-network/model"

type Triple struct {
	Subject   string `turtle:"subject"`
	Predicate string `turtle:"predicate"`
	Object    string `turtle:"object"`
}

func GetKeywordsAndAbstractEmbedding(url string) (string, []string, error) {
	res, err := http.Get(url)
	if err != nil {
		return "", nil, err
	}
	defer res.Body.Close()

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return "", nil, err
	}

	fmt.Println("aaa \n", res.Body, doc)
	return "", nil, nil
}

func main() {

	file, err := os.Open("citation_network_tiny.ttl")
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	var triples []Triple

	byteFile, err := io.ReadAll(file)
	if err != nil {
		fmt.Println("Error converting file:", err)
		return
	}
	err = turtle.Unmarshal(byteFile, &triples)
	if err != nil {
		fmt.Println("Error unmarshaling Turtle data:", err)
		return
	}

	a, _, _ := GetKeywordsAndAbstractEmbedding(triples[0].Subject)

	fmt.Println(a)

	g := model.NewGraph{
		Nodes: map[string]model.NewNode{},
		Edges: map[int]model.NewEdge{},
		Type:  "graph",
	}

	for _, triple := range triples {
		value, ok := g.Nodes[triple.Subject]
		value2, ok2 := g.Nodes[triple.Object]
		if !ok {
			value = model.NewNode{
				ID:         triple.Subject,
				Attributes: map[string]interface{}{},
			}
			g.AddNode(value)
		}
		if !ok2 {
			value2 = model.NewNode{
				ID:         triple.Object,
				Attributes: map[string]interface{}{},
			}
			g.AddNode(value2)
		}
		g.AddEdge(model.NewEdge{
			First_node:  value,
			Second_node: value2,
			Attributes:  map[string]interface{}{},
		})
	}

	fmt.Printf("%d, %d\n", len(g.Nodes), len(g.Edges))

	/*a := model.NewGraph{
		Nodes: map[string]model.NewNode{},
		Edges: map[int]model.NewEdge{},
		Type:  "graph",
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

	fmt.Println(a.ToString())*/
}
