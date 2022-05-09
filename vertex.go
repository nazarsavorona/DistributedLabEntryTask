package main

import (
	"fmt"
	"time"
)

type Vertex struct {
	stationID       int
	adjacent        []*Vertex
	possibleTickets []TrainTicket
}

func NewVertex(stationID int) *Vertex {
	return &Vertex{
		stationID:       stationID,
		adjacent:        []*Vertex{},
		possibleTickets: []TrainTicket{},
	}
}

func (vertex *Vertex) getGraphvizInfo(edgeType ConditionType) string {
	info := fmt.Sprintf("\"%d\"\n", vertex.stationID)

	edgeInfo := ""
	for _, ticket := range vertex.possibleTickets {
		edgeInfo = ""

		switch edgeType {
		case ByCost:
			edgeInfo += fmt.Sprintf("%.2f", ticket.price)
		case ByDuration:
			edgeInfo += ticket.duration.String()
		}

		info += fmt.Sprintf("\"%d\" -> \"%d\" [ label = \"%s, train: %d\" ]\n",
			ticket.from, ticket.to, edgeInfo, ticket.trainID)
	}

	return info
}

func (vertex *Vertex) getTicketByPrice(anotherID int) *TrainTicket {
	minIndex := -1

	for i := 0; i < len(vertex.possibleTickets); i++ {
		if vertex.possibleTickets[i].to == anotherID {
			if minIndex == -1 {
				minIndex = i
				continue
			}

			if vertex.possibleTickets[i].price < vertex.possibleTickets[minIndex].price {
				minIndex = i
			}
		}
	}

	if minIndex == -1 {
		return NewFakeTicket()
	}

	return &vertex.possibleTickets[minIndex]
}

func (vertex *Vertex) getTicketByDuration(anotherID int, currentTime time.Time) (*TrainTicket, time.Duration) {
	minIndex := -1
	minDuration := time.Duration(0)
	currentDuration := time.Duration(0)

	for i := 0; i < len(vertex.possibleTickets); i++ {
		if vertex.possibleTickets[i].to == anotherID {
			currentDuration = vertex.possibleTickets[i].getDuration(currentTime)

			if minIndex == -1 || currentDuration < minDuration {
				minIndex = i
				minDuration = currentDuration
			}
		}
	}

	if minIndex == -1 {
		return NewFakeTicket(), FakeTravelDuration
	}

	return &vertex.possibleTickets[minIndex], minDuration
}

func (vertex Vertex) getTicketAlternativesByPrice(ticket TrainTicket) []TrainTicket {
	tickets := make([]TrainTicket, 0)

	for _, currentTicket := range vertex.possibleTickets {
		if ticket.to == currentTicket.to && ticket.price == currentTicket.price {
			tickets = append(tickets, currentTicket)
		}
	}

	return tickets
}

func (vertex Vertex) getTicketAlternativesByDuration(ticket TrainTicket, currentTime time.Time) []TrainTicket {
	tickets := make([]TrainTicket, 0)

	for _, currentTicket := range vertex.possibleTickets {
		if ticket.to == currentTicket.to && ticket.getDuration(currentTime) == currentTicket.getDuration(currentTime) {
			tickets = append(tickets, currentTicket)
		}
	}

	return tickets
}
