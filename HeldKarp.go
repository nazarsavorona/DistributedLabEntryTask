package main

import "time"

func HeldKarpByCost(start *Vertex, vertices VertexSet, v Vertex, path *[]Vertex, tickets *[]TrainTicket, globalVertices *VertexSet) float64 {
	if vertices.size() == 1 && vertices.has(&v) {
		return start.getTicketByPrice(v.stationID).price
	}

	vertices.remove(&v)
	otherVertices := vertices.getList()

	currentCost := FakeHugeCost
	minCost := FakeHugeCost

	minVertex := otherVertices[0]
	minTicket := NewFakeTicket()
	minVertexFound := false

	for _, currentVertex := range otherVertices {
		tempSet := *NewVertexSet()
		tempSet = *tempSet.union(&vertices)

		ticket := currentVertex.getTicketByPrice(v.stationID)
		currentAdjacentCost := ticket.price

		if currentAdjacentCost == FakeHugeCost || !globalVertices.has(currentVertex) {
			continue
		}

		currentHeldKarp := HeldKarpByCost(start, tempSet, *currentVertex, path, tickets, globalVertices)

		if currentHeldKarp == FakeHugeCost {
			continue
		}

		currentCost = currentHeldKarp + currentAdjacentCost

		if minCost == FakeHugeCost || minCost > currentCost {
			minCost = currentCost
			minVertex = currentVertex
			minTicket = ticket
			minVertexFound = true
		}
	}

	if minVertexFound {
		*path = append(*path, *minVertex)
		*tickets = append(*tickets, *minTicket)

		globalVertices.remove(minVertex)
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
