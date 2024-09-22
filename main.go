package main

import (
	"container/list"
	"fmt"
)

// Define a graph to represent the rooms and their connections
type Graph struct {
	rooms map[string][]string // Each room is connected to a list of other rooms
}

// Create a new graph
func NewGraph() *Graph {
	return &Graph{rooms: make(map[string][]string)}
}

// Add a connection (edge) between two rooms (nodes)
func (g *Graph) AddConnection(room1, room2 string) {
	g.rooms[room1] = append(g.rooms[room2], room2)
	g.rooms[room2] = append(g.rooms[room2], room1)
}

// Find the shortest path between startRoom and endRoom using BFS
func (g *Graph) ShortestPath(startRoom, endRoom string) []string {
	// Check if the rooms exist
	if _, ok := g.rooms[startRoom]; !ok {
		return nil
	}
	if _, ok := g.rooms[endRoom]; !ok {
		return nil
	}

	// BFS queue and visited set
	queue := list.New()
	queue.PushBack([]string{startRoom})
	visited := make(map[string]bool)
	visited[startRoom] = true

	// Perform BFS to find the shortest path
	for queue.Len() > 0 {
		path := queue.Remove(queue.Front()).([]string)
		currentRoom := path[len(path)-1]

		// If we reached the target room, return the path
		if currentRoom == endRoom {
			return path
		}

		// Visit neighbors of the current room
		for _, neighbor := range g.rooms[currentRoom] {
			if !visited[neighbor] {
				visited[neighbor] = true
				newPath := append(path, neighbor)
				queue.PushBack(newPath)
			}
		}
	}
	return nil
}

func main() {
	// Create a graph of rooms
	graph := NewGraph()

	// Add connections between rooms (constant distance assumed)
	graph.AddConnection("Room1", "Room2")
	graph.AddConnection("Room2", "Room3")
	graph.AddConnection("Room3", "Room4")
	graph.AddConnection("Room1", "Room5")
	graph.AddConnection("Room5", "Room4")

	// Find the shortest path for an ant from Room1 to Room4
	path := graph.ShortestPath("Room1", "Room4")

	// Print the result
	if path != nil {
		fmt.Printf("Shortest path: %v\n", path)
	} else {
		fmt.Println("No path found")
	}
}
