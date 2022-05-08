package main

import (
	"fmt"
	"strings"
	"time"
)

type ConditionType int

const (
	ByCost ConditionType = iota
	ByDuration
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

		//info += fmt.Sprintf("\"%d\" -> \"%d\" [ label = \"%s, train: %d\" ]\n",
		//	ticket.from, ticket.to, edgeInfo, ticket.trainID)
		info += fmt.Sprintf("\"%d\" -> \"%d\" [ label = \"%s\" ]\n",
			ticket.from, ticket.to, edgeInfo)
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
			currentTicket := vertex.possibleTickets[i]
			currentDepartureTime := currentTicket.departure
			currentWaitingTime := time.Duration(0)

			if currentTime != fakeTime {
				for currentDepartureTime.Before(currentTime) {
					currentDepartureTime = currentDepartureTime.Add(time.Hour * 24)
				}

				currentWaitingTime = currentDepartureTime.Sub(currentTime)
			}

			currentDuration = currentTicket.duration + currentWaitingTime

			if minIndex == -1 || currentDuration < minDuration {
				minIndex = i
				minDuration = currentDuration
			}
		}
	}

	if minIndex == -1 {
		return NewFakeTicket(), fakeTripDuration
	}

	return &vertex.possibleTickets[minIndex], minDuration
}

func (vertex *Vertex) getNeighbours() *VertexSet {
	neighbours := NewVertexSet()

	for _, another := range vertex.adjacent {
		if !neighbours.Has(another) {
			neighbours.Add(another)
		}
	}

	return neighbours
}

type Graph struct {
	vertices []*Vertex
}

func generateGraph(tickets []TrainTicket) *Graph {
	graph := Graph{vertices: []*Vertex{}}

	for _, ticket := range tickets {
		graph.addVertex(ticket.from)
		graph.addVertex(ticket.to)
	}

	for _, ticket := range tickets {
		graph.addEdge(ticket)
	}

	return &graph
}

func (graph *Graph) getVertex(stationID int) *Vertex {
	for _, vertex := range graph.vertices {
		if vertex.stationID == stationID {
			return vertex
		}
	}
	return nil
}

func (graph *Graph) addVertex(stationID int) {
	if graph.exists(stationID) {
		//fmt.Printf("Station \"%d\" already exists in a graph\n", stationID)
		return
	}

	graph.vertices = append(graph.vertices, NewVertex(stationID))
}

func (graph *Graph) exists(stationID int) bool {
	return graph.getVertex(stationID) != nil
}

func (graph *Graph) addEdge(ticket TrainTicket) {
	fromVertex := graph.getVertex(ticket.from)
	toVertex := graph.getVertex(ticket.to)

	fromVertex.adjacent = append(fromVertex.adjacent, toVertex)
	fromVertex.possibleTickets = append(fromVertex.possibleTickets, ticket)
}

func (graph *Graph) getDistanceMatrix() (map[int]int, [][]*TrainTicket) {
	vertexCount := len(graph.vertices)
	matrix := make([][]*TrainTicket, vertexCount)

	indexMapping := make(map[int]int)

	for i := 0; i < vertexCount; i++ {
		matrix[i] = make([]*TrainTicket, vertexCount)

		for j := 0; j < vertexCount; j++ {
			matrix[i][j] = NewFakeTicket()
		}

		indexMapping[graph.vertices[i].stationID] = i
	}

	for _, vertex := range graph.vertices {
		for _, adjacentVertex := range vertex.adjacent {
			matrix[indexMapping[vertex.stationID]][indexMapping[adjacentVertex.stationID]] = vertex.getTicketByPrice(adjacentVertex.stationID)
		}
	}

	return indexMapping, matrix
}

func (graph *Graph) getGraphvizInfo(name string, edgeType ConditionType) string {
	name = strings.ReplaceAll(name, " ", "_")

	graphString := fmt.Sprintf("digraph %s {\n", name)

	for _, vertex := range graph.vertices {
		graphString += vertex.getGraphvizInfo(edgeType)
	}

	graphString += "}"

	return graphString
}

func (graph *Graph) printCostsMatrix() {
	_, costs := graph.getDistanceMatrix()

	fmt.Printf("\t")
	for _, v := range graph.vertices {
		fmt.Printf("%d\t", v.stationID)
	}
	fmt.Printf("\n")

	for i, row := range costs {
		fmt.Printf("%d\t", graph.vertices[i].stationID)
		for j := range row {
			if costs[i][j].price == fakeHugeCost {
				fmt.Printf("max\t")
			} else {
				fmt.Printf("%.2f\t", costs[i][j].price)

			}
		}
		fmt.Printf("\n")
	}
}

func (graph *Graph) optimalRoutes(condition ConditionType) ([][]Vertex, [][]TrainTicket) {
	paths := make([][]Vertex, 0)
	tickets := make([][]TrainTicket, 0)

	currentPath := make([]Vertex, 0)
	currentTickets := make([]TrainTicket, 0)

	for _, startVertex := range graph.vertices {
		set := NewVertexSet()
		globalSet := NewVertexSet()

		set.AddMulti(graph.vertices...)
		globalSet.AddMulti(graph.vertices...)

		currentPath = make([]Vertex, 0)
		currentTickets = make([]TrainTicket, 0)

		currentPath = append(currentPath, *startVertex)

		switch condition {
		case ByCost:
			HeldKarpByCost(startVertex, *set, *startVertex, &currentPath, &currentTickets, globalSet)
		case ByDuration:
			initialTime := fakeTime
			HeldKarpByDuration(startVertex, &initialTime, *set, *startVertex, &currentPath, &currentTickets, globalSet)
		}

		if len(currentPath) == len(graph.vertices) {
			currentTickets = append([]TrainTicket{*currentPath[0].getTicketByPrice(currentPath[1].stationID)}, currentTickets...)

			paths = append(paths, currentPath)
			tickets = append(tickets, currentTickets)
		}
	}

	return paths, tickets
}

func HeldKarpByCost(start *Vertex, vertices VertexSet, v Vertex, path *[]Vertex, tickets *[]TrainTicket, globalVertices *VertexSet) float64 {
	if vertices.Size() == 1 && vertices.Has(&v) {
		return start.getTicketByPrice(v.stationID).price
	}

	vertices.Remove(&v)
	otherVertices := vertices.GetList()

	currentCost := fakeHugeCost
	minCost := fakeHugeCost

	minVertex := otherVertices[0]
	minTicket := NewFakeTicket()
	minVertexFound := false

	for _, currentVertex := range otherVertices {
		tempSet := *NewVertexSet()
		tempSet = *tempSet.Union(&vertices)

		ticket := currentVertex.getTicketByPrice(v.stationID)
		currentAdjacentCost := ticket.price

		if currentAdjacentCost == fakeHugeCost || !globalVertices.Has(currentVertex) {
			continue
		}

		currentHeldKarp := HeldKarpByCost(start, tempSet, *currentVertex, path, tickets, globalVertices)

		if currentHeldKarp == fakeHugeCost {
			continue
		}

		currentCost = currentHeldKarp + currentAdjacentCost

		if minCost == fakeHugeCost || minCost > currentCost {
			minCost = currentCost
			minVertex = currentVertex
			minTicket = ticket
			minVertexFound = true
		}
	}

	if minVertexFound {
		*path = append(*path, *minVertex)
		*tickets = append(*tickets, *minTicket)

		globalVertices.Remove(minVertex)
	}

	return minCost
}

func HeldKarpByDuration(start *Vertex, currentTime *time.Time, vertices VertexSet, v Vertex, path *[]Vertex, tickets *[]TrainTicket, globalVertices *VertexSet) time.Duration {
	if vertices.Size() == 1 && vertices.Has(&v) {
		_, duration := start.getTicketByDuration(v.stationID, *currentTime)

		return duration
	}

	vertices.Remove(&v)
	otherVertices := vertices.GetList()

	currentDuration := fakeTravelDuration
	minDuration := fakeTravelDuration

	minVertex := otherVertices[0]
	minTicket := NewFakeTicket()
	minVertexFound := false

	for _, currentVertex := range otherVertices {
		tempSet := *NewVertexSet()
		tempSet.AddMulti(otherVertices...)

		ticket, duration := currentVertex.getTicketByDuration(v.stationID, *currentTime)

		if duration == fakeTripDuration || duration == fakeTravelDuration || !globalVertices.Has(currentVertex) {
			continue
		}

		currentHeldKarp := HeldKarpByDuration(start, currentTime, tempSet, *currentVertex, path, tickets, globalVertices)

		if currentHeldKarp == fakeTripDuration || currentHeldKarp == fakeTravelDuration {
			continue
		}

		currentDuration = currentHeldKarp + duration

		if minDuration == fakeTravelDuration || minDuration > currentDuration {
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
		globalVertices.Remove(minVertex)
	}

	return minDuration
}
