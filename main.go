package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"sort"
	"strings"

	"github.com/jmCodeCraft/go-network/model"
	"github.com/nvkp/turtle"
)

//"github.com/jmCodeCraft/go-network/model"

type Triple struct {
	Subject   string `turtle:"subject"`
	Predicate string `turtle:"predicate"`
	Object    string `turtle:"object"`
}

func OnPage(link string) map[string]interface{} {
	res, err := http.Get(link)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()

	// Read the response body
	content, err := io.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	// Create a map to store the parsed JSON
	var result map[string]interface{}

	// Parse the JSON into the map
	err = json.Unmarshal(content, &result)
	if err != nil {
		log.Fatal("Error parsing JSON:", err)
	}

	return result
}

func reconstructText(wordPositions map[string]interface{}) string {
	// Slice to hold words in order of positions
	var words []string

	// Collect positions and words
	positionWordMap := make(map[float64]string)
	for word, positions := range wordPositions {
		positions := positions.([]interface{})
		for _, pos := range positions {
			x := pos.(float64)
			positionWordMap[x] = word
		}
	}

	// Sort positions and build the ordered text
	var positions []float64
	for pos := range positionWordMap {
		positions = append(positions, pos)
	}
	sort.Float64s(positions)

	for _, pos := range positions {
		words = append(words, positionWordMap[pos])
	}

	return strings.Join(words, " ")
}

func GetKeywords(api_map map[string]interface{}) []string {
	x := api_map["keywords"]
	var words []string

	if intSlice, ok := x.([]interface{}); ok {
		for _, value := range intSlice {
			if slice, ok := value.(map[string]interface{}); ok {
				str, ok := slice["display_name"].(string)
				if ok {
					words = append(words, str)
				}
			}
		}
	} else {
		fmt.Println("x is not a slice of ints")
	}
	return words
}

func GetText(api_map map[string]interface{}) string {
	inverted_text, ok := api_map["abstract_inverted_index"].(map[string]interface{})
	var text string
	if ok {
		text = reconstructText(inverted_text)
	}
	return text

	/*apiKey := "your-api-key-here"

	client := openai.NewClient(apiKey)

	resp, err := client.CreateEmbeddings(context.Background(), openai.EmbeddingRequest{
		Input: []string{text},
		Model: "text-embedding-ada-002", // You can change the model if needed
	})
	if err != nil {
		log.Fatalf("Failed to get embedding: %v", err)
	}

	embedding := resp.Data[0].Embedding
	fmt.Println("Text Embedding:", embedding)
	return embedding*/
}

func GetNodeAttributes(api_map map[string]interface{}) map[string]interface{} {
	attr := map[string]interface{}{}
	attr["keywords"] = GetKeywords(api_map)
	attr["text_embedding"] = GetText(api_map)
	return attr
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

	g := model.NewGraph{
		Nodes: map[string]model.NewNode{},
		Edges: map[int]model.NewEdge{},
		Type:  "graph",
	}

	i := 0

	for _, triple := range triples {
		value, ok := g.Nodes[triple.Subject]
		value2, ok2 := g.Nodes[triple.Object]
		if !ok {
			str := strings.Replace(triple.Subject, "https://semopenalex.org/work", "https://api.openalex.org/works", 1)
			api_map := OnPage(str)
			value = model.NewNode{
				ID:         triple.Subject,
				Attributes: GetNodeAttributes(api_map),
			}
			g.AddNode(value)
		}
		if !ok2 {
			str := strings.Replace(triple.Object, "https://semopenalex.org/work", "https://api.openalex.org/works", 1)
			api_map := OnPage(str)
			value2 = model.NewNode{
				ID:         triple.Object,
				Attributes: GetNodeAttributes(api_map),
			}
			g.AddNode(value2)
		}
		g.AddEdge(model.NewEdge{
			First_node:  value,
			Second_node: value2,
			Attributes:  map[string]interface{}{},
		})

		fmt.Println(i)
		i++
	}

	fmt.Println(g.NodeToString(g.Nodes[triples[0].Object]))

	/*
		str := strings.Replace("https://semopenalex.org/work/W2140953464", "https://semopenalex.org/work", "https://api.openalex.org/works", 1)
		fmt.Println(str)

		m := OnPage(str)

		fmt.Printf("%d\n", len(m))

		x := m["abstract_inverted_index"]

		fmt.Println(reflect.TypeOf(x))
		GetKeywords(m)
		/*if intSlice, ok := x.(map[string]interface{}); ok {
			for key, value := range intSlice {
				fmt.Println(key, value)
			}
		} else {
			fmt.Println("x is not a slice of ints")
		}*/
}
