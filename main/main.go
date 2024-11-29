package main

//"github.com/jmCodeCraft/go-network/model"

func main() {

	extraction_file_name := "citation_network_tiny.ttl"

	_ = Extract(extraction_file_name)
	/*
		k := Create_graph(atribute_map_name)
		g := *k

		fmt.Println(len(g.Nodes), len(g.Edges))

		g.CombineLeaves()

		err := g.WriteToFile("graph_data.json")
		if err != nil {
			fmt.Println("Error:", err)
		}
	*/
}

//TODO: download spletnih strani posebej, skupaj z embeddingi in podobno
//potem posebej sestavljanje grafa
//izmeri čas za vsak korak posebej: branje podatkov, ustvarjanje grafa, redčenje grafa
//dokumentacija novih metod
//datasete v gitignore
