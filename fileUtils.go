package main

import (
	"crypto/sha256"
	"encoding/csv"
	"fmt"
	"io"
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

func getFileSha256(filepath string) string {
	f, err := os.Open(filepath)
	if err != nil {
		log.Fatal(err)
	}
	defer func(f *os.File) {
		err := f.Close()
		if err != nil {

		}
	}(f)

	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		log.Fatal(err)
	}

	return fmt.Sprintf("%x", h.Sum(nil))
}
