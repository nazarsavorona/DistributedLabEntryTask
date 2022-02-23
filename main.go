package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

type TrainTicket struct {
	trainID   int
	from      int
	to        int
	price     float64
	departure time.Time
	arrival   time.Time
}

func (ticket TrainTicket) String() string {
	return fmt.Sprintf("{TrainID: %d, from: %d, to: %d, price: %f, departure: %s, arrival: %s}",
		ticket.trainID, ticket.from, ticket.to, ticket.price, ticket.departure.Format("15:04:05"), ticket.arrival.Format("15:04:05"))
}

func getTicketsString(tickets []TrainTicket) string {
	ticketsString := fmt.Sprint("[")
	for i, ticket := range tickets {
		if i == 0 {
			ticketsString += fmt.Sprint(ticket)
		} else {
			ticketsString += fmt.Sprint(";\n", ticket)
		}
	}
	ticketsString += fmt.Sprint("]")

	return ticketsString
}

func usage() {
	fmt.Fprintf(os.Stderr, "usage: %s [inputfile]\n", os.Args[0])
	flag.PrintDefaults()
	os.Exit(2)
}

func main() {
	if len(os.Args) != 2 {
		usage()
	}

	filepath := os.Args[1]

	file := openFile(filepath)

	defer file.Close()

	data := readCsv(file)

	tickets := createTicketList(data)

	fmt.Println(len(tickets))
	fmt.Println(getTicketsString(tickets))
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

func openFile(filepath string) *os.File {
	file, err := os.Open(filepath)

	if err != nil {
		log.Fatal(err)
	}

	return file
}

func createTicketList(data [][]string) []TrainTicket {
	var ticketList []TrainTicket

	for _, line := range data {
		var ticket TrainTicket

		for j, field := range line {
			field = strings.Trim(field, " ")
			field = strings.Trim(field, "{")
			field = strings.Trim(field, "}")

			switch j {
			case 0:
				ticket.trainID, _ = strconv.Atoi(field)
			case 1:
				ticket.from, _ = strconv.Atoi(field)
			case 2:
				ticket.to, _ = strconv.Atoi(field)
			case 3:
				ticket.price, _ = strconv.ParseFloat(field, 64)
			case 4:
				ticket.departure, _ = time.Parse("15:04:05", field)
			case 5:
				ticket.arrival, _ = time.Parse("15:04:05", field)

			}
		}

		ticketList = append(ticketList, ticket)
	}

	return ticketList
}
