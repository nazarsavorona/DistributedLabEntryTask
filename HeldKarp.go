package main

import (
	"time"
)

func HeldKarpByCost(start *Vertex, vertices VertexSet, v Vertex, path *[]Vertex, cost *float64,
	currentPath []Vertex, currentCost float64, tickets *Tickets, currentTickets Tickets) float64 {
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
		tempSet := *NewVertexSet().union(&vertices)

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

func HeldKarpByDuration(start *Vertex, currentTime time.Time, vertices VertexSet, v Vertex, paths *[][]Vertex, optimalDuration *time.Duration,
	currentPath []Vertex, currentDuration time.Duration, tickets *[]Tickets, currentTickets Tickets) time.Duration {
	if vertices.size() == 1 && vertices.has(&v) {
		ticket, duration := start.getTicketByDuration(v.stationID, currentTime)
		currentTickets = append(currentTickets, *ticket)
		currentPath = append(currentPath, *start)
		currentTime = ticket.arrival

		currentDuration += duration

		if *optimalDuration > currentDuration {
			*tickets = make([]Tickets, 0)
			*paths = make([][]Vertex, 0)

			*optimalDuration = currentDuration
		}

		if *optimalDuration >= currentDuration {
			*tickets = append(*tickets, currentTickets)
			*paths = append(*paths, currentPath)
		}

		return duration
	}

	vertices.remove(&v)

	otherVertices := vertices.getList()

	localDuration := FakeTravelDuration
	minDuration := FakeTravelDuration

	for _, currentVertex := range otherVertices {
		previousTime := currentTime
		tempSet := *NewVertexSet().union(&vertices)

		ticket, duration := currentVertex.getTicketByDuration(v.stationID, currentTime)

		if duration == FakeTravelDuration {
			continue
		}

		currentPath = append(currentPath, *currentVertex)
		currentTickets = append(currentTickets, *ticket)
		currentDuration += duration
		currentTime = ticket.arrival

		currentHeldKarp := HeldKarpByDuration(start, currentTime, tempSet, *currentVertex, paths, optimalDuration,
			currentPath, currentDuration, tickets, currentTickets)

		if currentHeldKarp == FakeTravelDuration {
			currentPath = currentPath[:len(currentPath)-1]
			currentTickets = currentTickets[:len(currentTickets)-1]
			currentDuration -= duration
			currentTime = previousTime

			continue
		}

		localDuration = currentHeldKarp + duration

		if minDuration == FakeTravelDuration || minDuration > localDuration {
			minDuration = localDuration
		}

		currentPath = currentPath[:len(currentPath)-1]
		currentTickets = currentTickets[:len(currentTickets)-1]
		currentDuration -= duration
		currentTime = previousTime
	}

	return minDuration
}
