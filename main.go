// https://jspaint.app/#local:7836ad97ab865

// Add connections between rooms (constant distance assumed)
// graph.AddConnection("Room0", "Room4")
// graph.AddConnection("Room0", "Room6")
// graph.AddConnection("Room1", "Room3")
// graph.AddConnection("Room4", "Room3")
// graph.AddConnection("Room5", "Room2")
// graph.AddConnection("Room3", "Room5")
// graph.AddConnection("Room4", "Room2")
// graph.AddConnection("Room2", "Room1")
// graph.AddConnection("Room7", "Room6")
// graph.AddConnection("Room7", "Room2")
// graph.AddConnection("Room7", "Room4")
// graph.AddConnection("Room6", "Room5")

package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Room struct {
	name  string
	links []string
}

type Graph struct {
	rooms     map[string]*Room
	startRoom string
	endRoom   string
}

func NewGraph() *Graph {
	return &Graph{
		rooms: make(map[string]*Room),
	}
}

func (g *Graph) AddRoom(name string) {
	if _, exists := g.rooms[name]; !exists {
		g.rooms[name] = &Room{name: name, links: []string{}}
	}
}

func (g *Graph) AddConnection(room1, room2 string) {
	g.AddRoom(room1)
	g.AddRoom(room2)
	g.rooms[room1].links = append(g.rooms[room1].links, room2)
	g.rooms[room2].links = append(g.rooms[room2].links, room1)
}

func parseFile(fileName string) (*Graph, int, error) {
	file, err := os.Open(fileName)
	if err != nil {
		return nil, 0, fmt.Errorf("ERROR: could not open file")
	}
	defer file.Close()

	var numAnts int
	graph := NewGraph()
	var startRoom, endRoom string

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()

		if strings.HasPrefix(line, "#") && !strings.HasPrefix(line, "##") {
			continue
		}

		if strings.HasPrefix(line, "##start") {
			scanner.Scan()
			startRoom = parseRoom(scanner.Text(), graph)
			graph.startRoom = startRoom
			continue
		}

		if strings.HasPrefix(line, "##end") {
			scanner.Scan()
			endRoom = parseRoom(scanner.Text(), graph)
			graph.endRoom = endRoom
			continue
		}

		if numAnts == 0 {
			numAnts, err = strconv.Atoi(line)
			if err != nil {
				return nil, 0, fmt.Errorf("ERROR: invalid number of ants")
			}
			continue
		}

		if strings.Contains(line, " ") {
			parseRoom(line, graph)
		}

		if strings.Contains(line, "-") {
			parts := strings.Split(line, "-")
			if len(parts) == 2 {
				graph.AddConnection(parts[0], parts[1])
			}
		}
	}

	if startRoom == "" || endRoom == "" {
		return nil, 0, fmt.Errorf("ERROR: invalid data format, no start or end room found")
	}

	return graph, numAnts, nil
}

func parseRoom(line string, graph *Graph) string {
	parts := strings.Split(line, " ")
	if len(parts) != 3 {
		return ""
	}
	roomName := parts[0]
	graph.AddRoom(roomName)
	return roomName
}

// Find multiple shortest paths using BFS
func findAllShortestPaths(graph *Graph, start, end string) [][]string {
	var paths [][]string
	queue := [][]string{{start}}

	visited := map[string]bool{start: true}

	for len(queue) > 0 {
		path := queue[0]
		queue = queue[1:]

		lastRoom := path[len(path)-1]

		if lastRoom == end {
			paths = append(paths, path)
			continue
		}

		for _, neighbor := range graph.rooms[lastRoom].links {
			if !visited[neighbor] || neighbor == end {
				newPath := append([]string{}, path...)
				newPath = append(newPath, neighbor)
				queue = append(queue, newPath)
				visited[neighbor] = true
			}
		}
	}

	return paths
}

// Simulate the movement of ants
func simulateAntMovement(paths [][]string, numAnts int) {
	antPositions := make(map[int]int) // Ant number to path position
	antsFinished := 0
	step := 0
	numPaths := len(paths)

	// Assign ants to paths in a round-robin manner
	for antsFinished < numAnts {
		step++
		moves := []string{}

		for ant := 1; ant <= numAnts; ant++ {
			pathIdx := (ant - 1) % numPaths
			path := paths[pathIdx]

			// If the ant has already reached the end, skip it
			if antPositions[ant] >= len(path)-1 {
				continue
			}

			// Move ant to the next room if the room is not occupied
			if antPositions[ant] < len(path)-1 {
				antPositions[ant]++
				move := fmt.Sprintf("L%d-%s", ant, path[antPositions[ant]])
				moves = append(moves, move)

				// If the ant reached the end, increase finished count
				if antPositions[ant] == len(path)-1 {
					antsFinished++
				}
			}
		}

		// Only print moves if any ants moved this step
		if len(moves) > 0 {
			fmt.Println(strings.Join(moves, " "))
		}
	}
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("ERROR: no input file specified")
		return
	}

	filename := os.Args[1]
	graph, numAnts, err := parseFile(filename)
	if err != nil {
		fmt.Println(err)
		return
	}

	// Find all shortest paths from start to end
	paths := findAllShortestPaths(graph, graph.startRoom, graph.endRoom)
	if len(paths) == 0 {
		fmt.Println("ERROR: no path from start to end")
		return
	}

	// Simulate and display the movement of ants
	simulateAntMovement(paths, numAnts)
}
