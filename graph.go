package main

import (
	"fmt"
	"strings"
	"time"
)

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
			if costs[i][j].price == FakeHugeCost {
				fmt.Printf("max\t")
			} else {
				fmt.Printf("%.2f\t", costs[i][j].price)

			}
		}
		fmt.Printf("\n")
	}
}

func (graph *Graph) optimalRoutes(condition ConditionType) []TicketsWithAlternatives {
	switch condition {
	case ByCost:
		return graph.optimalRoutesByCost()
	case ByDuration:
		return graph.optimalRoutesByDuration()
	}

	return nil
}

func (graph *Graph) optimalRoutesByDuration() []TicketsWithAlternatives {
	paths := make([][]Vertex, 0)
	tickets := make([][]TrainTicket, 0)

	currentPaths := make([][]Vertex, 0)
	currentTickets := make([]Tickets, 0)

	minDuration := FakeTravelDuration
	currentDuration := FakeTravelDuration

	for _, startVertex := range graph.vertices {
		set := NewVertexSet()
		set.addMany(graph.vertices...)

		currentTickets = make([]Tickets, 0)
		currentPaths = make([][]Vertex, 0)

		initialDuration := FakeTravelDuration

		currentDuration = HeldKarpByDuration(startVertex, FakeTime, *set, *startVertex, &currentPaths,
			&initialDuration, make([]Vertex, 0), time.Duration(0), &currentTickets, make([]TrainTicket, 0))

		for k, path := range currentPaths {
			if len(path) == len(graph.vertices) {

				for i, j := 0, len(path)-1; i < j; i, j = i+1, j-1 {
					path[i], path[j] = path[j], path[i]
				}

				for i, j := 0, len(currentTickets[k])-1; i < j; i, j = i+1, j-1 {
					currentTickets[k][i], currentTickets[k][j] = currentTickets[k][j], currentTickets[k][i]
				}

				if currentDuration < minDuration {
					paths = make([][]Vertex, 0)
					tickets = make([][]TrainTicket, 0)

					minDuration = currentDuration
				}

				if currentDuration <= minDuration {
					paths = append(paths, path)
					tickets = append(tickets, currentTickets[k])
				}
			}
		}
	}

	ticketsWithAlternatives := make([]TicketsWithAlternatives, 0)

	for i, path := range paths {
		currentTickets := make(TicketsWithAlternatives, 0)
		currentTime := FakeTime

		for j, vertex := range path {
			if j > 0 {
				currentTime = tickets[i][j-1].arrival
			}

			currentTickets = append(currentTickets, vertex.getTicketAlternativesByDuration(tickets[i][j], currentTime))
		}

		ticketsWithAlternatives = append(ticketsWithAlternatives, currentTickets)
	}

	return ticketsWithAlternatives
}

func (graph *Graph) optimalRoutesByCost() []TicketsWithAlternatives {
	paths := make([][]Vertex, 0)
	tickets := make([][]TrainTicket, 0)

	currentPath := make([]Vertex, 0)
	currentTickets := make([]TrainTicket, 0)

	minCost := FakeHugeCost

	currentCost := FakeHugeCost

	for _, startVertex := range graph.vertices {
		set := NewVertexSet()
		set.addMany(graph.vertices...)

		currentPath = make([]Vertex, 0)
		currentTickets = make(Tickets, 0)
		initialCost := FakeHugeCost

		currentCost = HeldKarpByCost(startVertex, *set, *startVertex, &currentPath, &initialCost,
			make([]Vertex, 0), 0, (*Tickets)(&currentTickets), make([]TrainTicket, 0))

		if len(currentPath) == len(graph.vertices) {
			for i, j := 0, len(currentPath)-1; i < j; i, j = i+1, j-1 {
				currentPath[i], currentPath[j] = currentPath[j], currentPath[i]
			}

			for i, j := 0, len(currentTickets)-1; i < j; i, j = i+1, j-1 {
				currentTickets[i], currentTickets[j] = currentTickets[j], currentTickets[i]
			}

			if currentCost < minCost {
				paths = make([][]Vertex, 0)
				tickets = make([][]TrainTicket, 0)

				minCost = currentCost
			}

			if currentCost <= minCost {
				paths = append(paths, currentPath)
				tickets = append(tickets, currentTickets)
			}
		}
	}

	ticketsWithAlternatives := make([]TicketsWithAlternatives, 0)
	for i, path := range paths {
		currentTickets := make(TicketsWithAlternatives, 0)

		for j, vertex := range path {
			currentTickets = append(currentTickets, vertex.getTicketAlternativesByPrice(tickets[i][j]))
		}

		ticketsWithAlternatives = append(ticketsWithAlternatives, currentTickets)
	}

	return ticketsWithAlternatives
}
