package main

import (
	"time"
)

func HeldKarpByCost(start *Vertex, vertices VertexSet, v Vertex, path *[]Vertex, cost *float64,
	currentPath []Vertex, currentCost float64, tickets *[]TrainTicket, currentTickets []TrainTicket) float64 {
	if vertices.size() == 1 && vertices.has(&v) {
		ticket := start.getTicketByPrice(v.stationID)
		currentTickets = append(currentTickets, *ticket)
		currentPath = append(currentPath, *start)

		currentCost += ticket.price

		if *cost >= currentCost {
			*path = currentPath
			*cost = currentCost
			*tickets = currentTickets
		}

		return ticket.price
	}

	vertices.remove(&v)
	otherVertices := vertices.getList()

	minCost := FakeHugeCost
	localCost := FakeHugeCost

	for _, currentVertex := range otherVertices {
		tempSet := *NewVertexSet()
		tempSet = *tempSet.union(&vertices)

		ticket := currentVertex.getTicketByPrice(v.stationID)
		currentAdjacentCost := ticket.price

		if currentAdjacentCost == FakeHugeCost {
			continue
		}

		currentPath = append(currentPath, *currentVertex)
		currentTickets = append(currentTickets, *ticket)
		currentCost += currentAdjacentCost

		currentHeldKarp := HeldKarpByCost(start, tempSet, *currentVertex, path, cost, currentPath, currentCost, tickets, currentTickets)

		if currentHeldKarp == FakeHugeCost {
			currentPath = currentPath[:len(currentPath)-1]
			currentTickets = currentTickets[:len(currentTickets)-1]
			currentCost -= currentAdjacentCost

			continue
		}

		localCost = currentHeldKarp + currentAdjacentCost

		if minCost == FakeHugeCost || minCost > localCost {
			minCost = localCost
		}

		currentPath = currentPath[:len(currentPath)-1]
		currentTickets = currentTickets[:len(currentTickets)-1]
		currentCost -= currentAdjacentCost
	}

	return minCost
}

func HeldKarpByDuration(start *Vertex, currentTime *time.Time, vertices VertexSet, v Vertex, path *[]Vertex, tickets *[]TrainTicket, globalVertices *VertexSet) time.Duration {
	if vertices.size() == 1 && vertices.has(&v) {
		_, duration := start.getTicketByDuration(v.stationID, *currentTime)

		return duration
	}

	vertices.remove(&v)
	otherVertices := vertices.getList()

	currentDuration := FakeTravelDuration
	minDuration := FakeTravelDuration

	minVertex := otherVertices[0]
	minTicket := NewFakeTicket()
	minVertexFound := false

	for _, currentVertex := range otherVertices {
		tempSet := *NewVertexSet()
		tempSet.addMany(otherVertices...)

		ticket, duration := currentVertex.getTicketByDuration(v.stationID, *currentTime)

		if duration == FakeTripDuration || duration == FakeTravelDuration || !globalVertices.has(currentVertex) {
			continue
		}

		currentHeldKarp := HeldKarpByDuration(start, currentTime, tempSet, *currentVertex, path, tickets, globalVertices)

		if currentHeldKarp == FakeTripDuration || currentHeldKarp == FakeTravelDuration {
			continue
		}

		currentDuration = currentHeldKarp + duration

		if minDuration == FakeTravelDuration || minDuration > currentDuration {
			minDuration = currentDuration
			minVertex = currentVertex
			minTicket = ticket
			minVertexFound = true
		}
	}

	if minVertexFound {
		*path = append(*path, *minVertex)
		*tickets = append(*tickets, *minTicket)
		*currentTime = minTicket.arrival
		globalVertices.remove(minVertex)
	}

	return minDuration
}
