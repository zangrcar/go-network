package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/jmCodeCraft/go-network/model"
	"github.com/nvkp/turtle"
)

//"github.com/jmCodeCraft/go-network/model"

func Create_graph(file_name string) *model.NewGraph {
	file, err := os.Open(file_name)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return nil
	}
	defer file.Close()

	// Declare a map to hold the data
	var all_atribute_map map[string]map[string]interface{}

	// Decode the JSON data into the map
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&all_atribute_map); err != nil {
		fmt.Println("Error decoding JSON:", err)
		return nil
	}

	g := model.NewGraph{
		Nodes: map[string]model.NewNode{},
		Edges: map[int]model.NewEdge{},
	}

	for key, value := range all_atribute_map {
		g.AddNode(model.NewNode{
			ID:         key,
			Attributes: value,
		})
	}

	for _, node := range g.Nodes {
		if node.Attributes["neighbors"] == nil {
			continue
		}
		strings := make([]string, len(node.Attributes["neighbors"].([]interface{})))
		for i, v := range node.Attributes["neighbors"].([]interface{}) {
			// Perform type assertion
			str, ok := v.(string)
			if !ok {
				fmt.Printf("Value at index %d is not a string\n", i)
				continue
			}
			strings[i] = str
		}

		for _, neighbor := range strings {
			g.AddEdge(model.NewEdge{
				First_node:  node,
				Second_node: g.Nodes[neighbor],
				Attributes:  map[string]interface{}{},
			})
		}
		delete(node.Attributes, "neighbors")
	}

	fmt.Println("Graph creation successful!")

	return &g

}

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

	// Create a map to store the parsed JSON
	var result map[string]interface{}

	// Read the response body
	content, err := io.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	// Parse the JSON into the map
	err = json.Unmarshal(content, &result)
	if err != nil {
		return nil
		//log.Printf("Failed to parse JSON. Content: %s\n", content)
		//log.Fatal("Error parsing JSON:", err)
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

func GetText(api_map map[string]interface{}) []float64 {
	inverted_text, ok := api_map["abstract_inverted_index"].(map[string]interface{})
	var text string
	if ok {
		text = reconstructText(inverted_text)
	}

	cmd := exec.Command("python", "text_embedding_model.py", text)
	output, err := cmd.Output()
	if err != nil {
		log.Fatal("Error getting text embedding: ", err)
	}

	outputStr := strings.TrimSpace(string(output))

	outputStr = strings.ReplaceAll(outputStr, "[", "")
	outputStr = strings.ReplaceAll(outputStr, "]", "")

	parts := strings.Split(outputStr, ", ")

	var result []float64
	for _, part := range parts {
		value, err := strconv.ParseFloat(part, 64)
		if err != nil {
			log.Fatal("Error transformint string to float64: ", err)
		}
		result = append(result, value)
	}
	return result
}

func GetNodeAttributes(api_map map[string]interface{}) map[string]interface{} {
	attr := map[string]interface{}{}
	attr["keywords"] = GetKeywords(api_map)
	attr["text_embedding"] = GetText(api_map)
	return attr
}

func Extract(file_name string) string {
	file, err := os.Open(file_name) //"citation_network_tiny.ttl"
	if err != nil {
		fmt.Println("Error opening file:", err)
		return ""
	}
	defer file.Close()

	var triples []Triple

	byteFile, err := io.ReadAll(file)
	if err != nil {
		fmt.Println("Error converting file:", err)
		return ""
	}
	err = turtle.Unmarshal(byteFile, &triples)
	if err != nil {
		fmt.Println("Error unmarshaling Turtle data:", err)
		return ""
	}

	allAttributes := make(map[string]map[string]interface{})
	i := 0
	for _, triple := range triples {
		fmt.Println(i)
		i++
		_, ok1 := allAttributes[triple.Subject]
		_, ok2 := allAttributes[triple.Object]
		if !ok1 {
			//	str := strings.Replace(triple.Subject, "https://semopenalex.org/work", "https://api.openalex.org/works", 1)
			//	api_map := OnPage(str)
			allAttributes[triple.Subject] = map[string]interface{}{} // GetNodeAttributes(api_map)
		}

		if !ok2 {
			//	str := strings.Replace(triple.Object, "https://semopenalex.org/work", "https://api.openalex.org/works", 1)
			//	api_map := OnPage(str)
			allAttributes[triple.Object] = map[string]interface{}{} //GetNodeAttributes(api_map)
		}
		_, exists := allAttributes[triple.Subject]["neighbors"]

		if !exists {
			var neighbors []string
			allAttributes[triple.Subject]["neighbors"] = neighbors
		}

		allAttributes[triple.Subject]["neighbors"] = append(allAttributes[triple.Subject]["neighbors"].([]string), triple.Object)
	}

	newfile_name := "citation_network_tiny_extracted_data_test.json"

	newfile, err := os.Create(newfile_name)
	if err != nil {
		fmt.Println("Error creating file:", err)
		return ""
	}
	defer newfile.Close()

	encoder := json.NewEncoder(newfile)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(allAttributes); err != nil {
		fmt.Println("Error encoding JSON:", err)
		return ""
	}

	fmt.Println("Data extraction succesful!")
	return newfile_name
}

func main() {
	start := time.Now()
	extracted_file_name := "citation_network_tiny_extracted_data_test.json"
	g := Create_graph(extracted_file_name)
	fmt.Println(len(g.Edges), len(g.Nodes))

	g.CombineLeaves()

	fmt.Println(len(g.Edges), len(g.Nodes))

	g.WriteToFile("graph_data.json")

	fmt.Println("Success")

	elapsed := time.Since(start)
	fmt.Println(elapsed)

}

//TODO: download spletnih strani posebej, skupaj z embeddingi in podobno
//potem posebej sestavljanje grafa
//izmeri čas za vsak korak posebej: branje podatkov, ustvarjanje grafa, redčenje grafa
//dokumentacija novih metod
//datasete v gitignore

//graph creation tiny: 200-300ms
//graph creatin & dilution tiny: 500-550ms
//graph to pyg format tiny: ~6500000ns
//shallow embedding of graph tiny: ~15000000ns
