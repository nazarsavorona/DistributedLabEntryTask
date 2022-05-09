package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
)

func usage() {
	_, err := fmt.Fprintf(os.Stderr, "usage: %s <input file> <output file> [{cost} / time]\n", os.Args[0])
	if err != nil {
		return
	}
	flag.PrintDefaults()
	os.Exit(2)
}

func main() {
	if len(os.Args) < 3 || len(os.Args) > 4 ||
		(len(os.Args) == 4 && !(strings.ToLower(os.Args[3]) == "cost" || strings.ToLower(os.Args[3]) == "time")) {
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
	ticketsLists := make([]TicketsWithAlternatives, 0)

	if len(os.Args) == 4 && strings.ToLower(os.Args[3]) == "time" {
		ticketsLists = graph.optimalRoutes(ByDuration)
	} else {
		ticketsLists = graph.optimalRoutes(ByCost)
	}

	//println(graph.getGraphvizInfo("f", ByCost))

	writeToFile(os.Args[2], ticketsLists)
	fmt.Printf("File sha256 hash: %s\n", getFileSha256(os.Args[2]))
}
