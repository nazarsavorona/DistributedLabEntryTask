package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
)

func openFile(filepath string) *os.File {
	file, err := os.Open(filepath)

	if err != nil {
		log.Fatal(err)
	}

	return file
}

func readCsv(file *os.File) [][]string {
	csvReader := csv.NewReader(file)
	csvReader.Comma = ';'
	data, err := csvReader.ReadAll()

	if err != nil {
		log.Fatal(err)
	}

	return data
}

func writeToFile(filepath string, ticketsLists []TicketsWithAlternatives) {
	file, err := os.Create(filepath)
	if err != nil {
		log.Fatal(err)
	} else {
		for i, tickets := range ticketsLists {
			_, err := file.WriteString(fmt.Sprintf("%d:\n%s\n", i, tickets.String()))
			if err != nil {
				return
			}
		}
	}

	err = file.Close()
	if err != nil {
		return
	}
}
