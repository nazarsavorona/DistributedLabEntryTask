package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"log"
	"math"
	"os"
	"strconv"
	"strings"
	"time"
)

const fakeHugeCost = math.MaxFloat64

const fakeTripDuration = time.Hour * 24 * 3
const fakeTravelDuration = fakeTripDuration * 10

const fakeTimeUnixNano = int64(-6829751778871345152) // time.Date(0, 0, 0, 0, 0, 0, 0, time.UTC) UnixNano representation
var fakeTime = time.Unix(0, fakeTimeUnixNano)

type TrainTicket struct {
	trainID   int
	from      int
	to        int
	price     float64
	departure time.Time
	arrival   time.Time
	duration  time.Duration
}

type Tickets []TrainTicket
type TicketsWithAlternatives []Tickets

func NewFakeTicket() *TrainTicket {
	return &TrainTicket{
		trainID:   -1,
		from:      0,
		to:        0,
		price:     fakeHugeCost,
		departure: fakeTime,
		arrival:   fakeTime,
		duration:  fakeTripDuration,
	}
}

func (ticket TrainTicket) String() string {
	if ticket.trainID == -1 {
		return fmt.Sprint("Fake")
	}
	return fmt.Sprintf("{TrainID: %d, from: %d, to: %d, price: %.2f, departure: %s, arrival: %s, duration: %s}",
		ticket.trainID, ticket.from, ticket.to, ticket.price, ticket.departure.Format("15:04:05"), ticket.arrival.Format("15:04:05"), ticket.duration)
}

func (ticket *TrainTicket) getDuration(currentTime time.Time) time.Duration {
	currentDepartureTime := ticket.departure
	currentWaitingTime := time.Duration(0)

	if currentTime != fakeTime {
		for currentDepartureTime.Before(currentTime) {
			currentDepartureTime = currentDepartureTime.Add(time.Hour * 24)
		}

		currentWaitingTime = currentDepartureTime.Sub(currentTime)
	}

	return ticket.duration + currentWaitingTime
}

func (tickets TicketsWithAlternatives) String() string {
	ticketsString := fmt.Sprint("[")
	for i, ticketAlternatives := range tickets {
		if i == 0 {
			ticketsString += fmt.Sprint(ticketAlternatives.String())
		} else {
			ticketsString += fmt.Sprint(";\n", ticketAlternatives.String())
		}
	}
	ticketsString += fmt.Sprint("]")

	return ticketsString
}

func (tickets Tickets) String() string {
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
				if ticket.departure.After(ticket.arrival) {
					ticket.arrival = ticket.arrival.Add(time.Hour * 24)
				}

				ticket.duration = ticket.arrival.Sub(ticket.departure)
			}
		}

		ticketList = append(ticketList, ticket)
	}

	return ticketList
}
