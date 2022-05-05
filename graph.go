package main

import (
	"fmt"
	"math"
	"strings"
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

func (vertex *Vertex) getTicketByCondition(anotherID int, condition ConditionType) *TrainTicket {
	currentMinIndex := -1

	for i := 0; i < len(vertex.possibleTickets); i++ {
		if currentMinIndex == -1 && vertex.possibleTickets[i].to == anotherID {
			currentMinIndex = i
			continue
		}

		if vertex.possibleTickets[i].to == anotherID {
			switch condition {
			case ByCost:
				if vertex.possibleTickets[i].price < vertex.possibleTickets[currentMinIndex].price {
					currentMinIndex = i
				}
			case ByDuration:
				if vertex.possibleTickets[i].duration < vertex.possibleTickets[currentMinIndex].duration {
					currentMinIndex = i
				}
			}
		}
	}

	if currentMinIndex == -1 {
		return nil
	}

	return &vertex.possibleTickets[currentMinIndex]
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

func (graph *Graph) getDistanceMatrix(condition ConditionType) (map[int]int, [][]*TrainTicket) {
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
			matrix[indexMapping[vertex.stationID]][indexMapping[adjacentVertex.stationID]] = vertex.getTicketByCondition(adjacentVertex.stationID, condition)
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

func (graph *Graph) optimalRoutes(condition ConditionType) ([][]Vertex, [][]TrainTicket) {
	paths := make([][]Vertex, 0)
	tickets := make([][]TrainTicket, 0)

	currentPath := make([]Vertex, 0)
	currentTickets := make([]TrainTicket, 0)

	mapping, costs := graph.getDistanceMatrix(condition)

	for _, startVertex := range graph.vertices {
		set := NewVertexSet()
		set.AddMulti(graph.vertices...)

		fmt.Printf("\n\ns:%d\n", startVertex.stationID)

		currentPath = make([]Vertex, 0)
		currentTickets = make([]TrainTicket, 0)

		HeldKarp(startVertex, *set, *startVertex, mapping, costs,
			&currentPath, &currentTickets)

		paths = append(paths, currentPath)
		tickets = append(tickets, currentTickets)
	}

	return paths, tickets
}

func (graph *Graph) printCostsMatrix(conditionType ConditionType) {
	_, costs := graph.getDistanceMatrix(conditionType)

	fmt.Printf("\t")
	for _, v := range graph.vertices {
		fmt.Printf("%d\t", v.stationID)
	}
	fmt.Printf("\n")

	for i, row := range costs {
		fmt.Printf("%d\t", graph.vertices[i].stationID)
		for j := range row {
			if costs[i][j].price == math.MaxFloat64 {
				fmt.Printf("max\t")
			} else {
				fmt.Printf("%.2f\t", costs[i][j].price)

			}
		}
		fmt.Printf("\n")
	}
}

func HeldKarp(start *Vertex, vertices VertexSet, v Vertex, mapping map[int]int,
	costs [][]*TrainTicket, path *[]Vertex, tickets *[]TrainTicket) float64 {
	if vertices.Size() == 1 && vertices.Has(&v) {
		return costs[mapping[start.stationID]][mapping[v.stationID]].price
	}

	vertices.Remove(&v)
	otherVertices := vertices.GetList()

	currentValue := math.MaxFloat64
	min := math.MaxFloat64

	minVertex := otherVertices[0]
	minTicket := NewFakeTicket()
	minVertexFound := false

	for _, currentVertex := range otherVertices {
		tempSet := *NewVertexSet()
		tempSet = *tempSet.Union(&vertices)

		ticket := costs[mapping[currentVertex.stationID]][mapping[v.stationID]]
		currentAdjacentCost := ticket.price

		if currentAdjacentCost == math.MaxFloat64 {
			continue
		}

		currentHeldKarp := HeldKarp(start, tempSet, *currentVertex, mapping, costs, path, tickets)

		if currentHeldKarp == math.MaxFloat64 {
			continue
		}

		currentValue = currentHeldKarp + currentAdjacentCost

		if min == math.MaxFloat64 || min > currentValue {
			min = currentValue
			minVertex = currentVertex
			minTicket = ticket
			minVertexFound = true
		}
	}

	if minVertexFound {
		*path = append(*path, *minVertex)
		*tickets = append(*tickets, *minTicket)
	}

	return min
}
