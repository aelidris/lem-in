package main

import (
	"container/list"
	"fmt"
	"sync"
	"time"
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
	g.rooms[room1] = append(g.rooms[room1], room2)
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

// Goroutine that simulates an ant moving along a path
func antMovement(antID int, path []string, roomsAvailability map[string]chan struct{}, wg *sync.WaitGroup) {
	defer wg.Done()

	for i, room := range path {
		// Check if the room is available
		if i > 0 { // The first room is where the ant starts, so it doesn't need to check
			fmt.Printf("Ant %d waiting to enter %s...\n", antID, room)
			roomsAvailability[room] <- struct{}{} // Block until the room becomes available
		}

		fmt.Printf("Ant %d is in %s\n", antID, room)
		time.Sleep(1 * time.Second) // Simulate some delay between rooms

		// If it's the destination, no need to release the room
		if i == len(path)-1 {
			fmt.Printf("Ant %d has reached the destination: %s\n", antID, room)
		} else {
			// Release the previous room so the next ant can move in
			prevRoom := path[i]
			if i > 0 && prevRoom != "Room1" { // Release only if it's not the starting room
				fmt.Printf("Ant %d leaving %s\n", antID, prevRoom)
				<-roomsAvailability[prevRoom] // Leave the current room
			}
		}
	}
}

func main() {
	// Create a graph of rooms
	graph := NewGraph()

	// Add connections between rooms (constant distance assumed)
	graph.AddConnection("Room0", "Room4")
	graph.AddConnection("Room0", "Room6")
	graph.AddConnection("Room1", "Room3")
	graph.AddConnection("Room4", "Room3")
	graph.AddConnection("Room5", "Room2")
	graph.AddConnection("Room3", "Room5")
	graph.AddConnection("Room4", "Room2")
	graph.AddConnection("Room2", "Room1")
	graph.AddConnection("Room7", "Room6")
	graph.AddConnection("Room7", "Room2")
	graph.AddConnection("Room7", "Room4")
	graph.AddConnection("Room6", "Room5")

	// Specify the number of ants
	var numAnts int
	fmt.Print("Enter the number of ants: ")
	fmt.Scan(&numAnts)

	// Find the shortest path for an ant from Room1 to Room0
	path := graph.ShortestPath("Room1", "Room0")

	// Print the result
	if path != nil {
		fmt.Printf("Shortest path: %v\n", path)
	} else {
		fmt.Println("No path found")
		return
	}

	// Create room availability channels (to simulate room occupancy)
	roomsAvailability := make(map[string]chan struct{})
	for room := range graph.rooms {
		if room == "Room1" || room == "Room0" {
			// Room1 (starting room) and Room0 (ending room) can be occupied by all ants at once
			roomsAvailability[room] = make(chan struct{}, numAnts)
		} else {
			// Other rooms can only be occupied by one ant at a time
			roomsAvailability[room] = make(chan struct{}, 1)
		}
	}

	// Use goroutines to simulate each ant moving concurrently
	var wg sync.WaitGroup
	wg.Add(numAnts)

	// Start the ants moving concurrently using goroutines
	for i := 1; i <= numAnts; i++ {
		go antMovement(i, path, roomsAvailability, &wg)
	}

	// Wait for all ants to finish moving
	wg.Wait()

	fmt.Println("All ants have reached the destination.")
}

// package main

// import (
// 	"container/list"
// 	"fmt"
// 	"sync"
// 	"time"
// )

// // Define a graph to represent the rooms and their connections
// type Graph struct {
// 	rooms map[string][]string // Each room is connected to a list of other rooms
// }

// // Create a new graph
// func NewGraph() *Graph {
// 	return &Graph{rooms: make(map[string][]string)}
// }

// // Add a connection (edge) between two rooms (nodes)
// func (g *Graph) AddConnection(room1, room2 string) {
// 	g.rooms[room1] = append(g.rooms[room1], room2)
// 	g.rooms[room2] = append(g.rooms[room2], room1)
// }

// // Find the shortest path between startRoom and endRoom using BFS
// func (g *Graph) ShortestPath(startRoom, endRoom string) []string {
// 	// Check if the rooms exist
// 	if _, ok := g.rooms[startRoom]; !ok {
// 		return nil
// 	}
// 	if _, ok := g.rooms[endRoom]; !ok {
// 		return nil
// 	}

// 	// BFS queue and visited set
// 	queue := list.New()
// 	queue.PushBack([]string{startRoom})
// 	visited := make(map[string]bool)
// 	visited[startRoom] = true

// 	// Perform BFS to find the shortest path
// 	for queue.Len() > 0 {
// 		path := queue.Remove(queue.Front()).([]string)
// 		currentRoom := path[len(path)-1]

// 		// If we reached the target room, return the path
// 		if currentRoom == endRoom {
// 			return path
// 		}

// 		// Visit neighbors of the current room
// 		for _, neighbor := range g.rooms[currentRoom] {
// 			if !visited[neighbor] {
// 				visited[neighbor] = true
// 				newPath := append(path, neighbor)
// 				queue.PushBack(newPath)
// 			}
// 		}
// 	}
// 	return nil
// }

// // Goroutine that simulates an ant moving along a path
// func antMovement(antID int, path []string, wg *sync.WaitGroup) {
// 	defer wg.Done()

// 	for i, room := range path {
// 		fmt.Printf("Ant %d is in %s\n", antID, room)
// 		time.Sleep(1 * time.Second) // Simulate some delay between rooms
// 		if i == len(path)-1 {
// 			fmt.Printf("Ant %d has reached the destination: %s\n", antID, room)
// 		}
// 	}
// }

// func main() {
// 	// Create a graph of rooms
// 	graph := NewGraph()

// 	// Add connections between rooms (constant distance assumed)
// 	graph.AddConnection("Room0", "Room4")
// 	graph.AddConnection("Room0", "Room6")
// 	graph.AddConnection("Room1", "Room3")
// 	graph.AddConnection("Room4", "Room3")
// 	graph.AddConnection("Room5", "Room2")
// 	graph.AddConnection("Room3", "Room5")
// 	graph.AddConnection("Room4", "Room2")
// 	graph.AddConnection("Room2", "Room1")
// 	graph.AddConnection("Room7", "Room6")
// 	graph.AddConnection("Room7", "Room2")
// 	graph.AddConnection("Room7", "Room4")
// 	graph.AddConnection("Room6", "Room5")

// 	// Specify the number of ants
// 	var numAnts int
// 	fmt.Print("Enter the number of ants: ")
// 	fmt.Scan(&numAnts)

// 	// Find the shortest path for an ant from Room1 to Room4
// 	path := graph.ShortestPath("Room1", "Room4")

// 	// Print the result
// 	if path != nil {
// 		fmt.Printf("Shortest path: %v\n", path)
// 	} else {
// 		fmt.Println("No path found")
// 		return
// 	}

// 	// Use goroutines to simulate each ant moving concurrently
// 	var wg sync.WaitGroup
// 	wg.Add(numAnts)

// 	// Start the ants moving concurrently using goroutines
// 	for i := 1; i <= numAnts; i++ {
// 		go antMovement(i, path, &wg)
// 	}

// 	// Wait for all ants to finish moving
// 	wg.Wait()

// 	fmt.Println("All ants have reached the destination.")
// }
