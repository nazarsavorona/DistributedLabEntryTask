package main

import (
	"flag"
	"fmt"
	"os"
)

func usage() {
	_, err := fmt.Fprintf(os.Stderr, "usage: %s [input file] [output file]\n", os.Args[0])
	if err != nil {
		return
	}
	flag.PrintDefaults()
	os.Exit(2)
}

func main() {
	if len(os.Args) != 3 {
		usage()
	}

	file := openFile(os.Args[1])

	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
		}
	}(file)

	data := readCsv(file)
	tickets := createTicketList(data)

	graph := generateGraph(tickets)

	//fmt.Print(graph.getGraphvizInfo("g", ByDuration))
	//fmt.Println(graph.getGraphvizInfo("g", ByCost))
	//graph.printCostsMatrix(ByCost)

	ticketsLists := graph.optimalRoutes(ByDuration)
	//ticketsLists := graph.optimalRoutes(ByCost)

	writeToFile(os.Args[2], ticketsLists)
}
