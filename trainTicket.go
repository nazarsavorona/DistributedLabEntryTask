package main

import (
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"
)

const FakeHugeCost = math.MaxFloat64

const FakeTravelDuration = time.Hour * 24 * 30

const FakeTimeUnixNano = int64(-62169984000) // time.Date(0, 0, 0, 0, 0, 0, 0, time.UTC).Unix() representation
var FakeTime = time.Unix(FakeTimeUnixNano, 0)

type ConditionType int

const (
	ByCost ConditionType = iota
	ByDuration
)

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
		price:     FakeHugeCost,
		departure: FakeTime,
		arrival:   FakeTime,
		duration:  FakeTravelDuration,
	}
}

func (ticket TrainTicket) String() string {
	if ticket.trainID == -1 {
		return fmt.Sprint("Fake")
	}
	return fmt.Sprintf("{TrainID: %d, from: %d, to: %d, price: %.2f, departure: %s, arrival: %s, duration: %s}",
		ticket.trainID, ticket.from, ticket.to, ticket.price, ticket.departure.Format("15:04:05"), ticket.arrival.Format("15:04:05"), ticket.duration)
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

func (ticket *TrainTicket) getDuration(currentTime time.Time) time.Duration {
	currentDepartureTime := ticket.departure
	currentWaitingTime := time.Duration(0)

	if currentTime != FakeTime {
		for currentDepartureTime.Before(currentTime) {
			currentDepartureTime = currentDepartureTime.Add(time.Hour * 24)
		}

		currentWaitingTime = currentDepartureTime.Sub(currentTime)
	}

	return ticket.duration + currentWaitingTime
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
