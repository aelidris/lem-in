package main

import (
	"container/list"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

// Extend the function to return parsed links
func GetData(dataFile string) (start, end string, rooms []string, links [][2]string, antNumbers int) {
	data, err := os.ReadFile(dataFile)
	if err != nil {
		log.Fatal(err)
	}
	if len(data) == 0 {
		log.Fatal("There is no data in the file!!")
	}

	checkData := strings.Split(string(data), "\n")
	affStart, affEnd := false, false

	// Process the data in a single loop
	for _, line := range checkData {
		// Skip comments
		if strings.HasPrefix(line, "#") {
			if line == "##start" {
				affStart = true
			} else if line == "##end" {
				affEnd = true
			}
			continue
		}

		// If this is the first non-comment line after ##start, capture the start room
		if affStart {
			start = strings.Fields(line)[0]
			affStart = false
		}

		// If this is the first non-comment line after ##end, capture the end room
		if affEnd {
			end = strings.Fields(line)[0]
			affEnd = false
		}

		// Capture number of ants (first non-comment line that is a single number)
		if antNumbers == 0 {
			antNumbers, err = strconv.Atoi(line)
			if err != nil {
				log.Fatal("Error converting number of ants to int:", err)
			}
			continue
		}

		// Capture room names (lines with three parts, representing room definition)
		parts := strings.Fields(line)
		if len(parts) == 3 {
			rooms = append(rooms, parts[0])
		}

		// Capture links (lines in the format "name1-name2")
		if strings.Contains(line, "-") {
			roomsLink := strings.Split(line, "-")
			if len(roomsLink) == 2 {
				links = append(links, [2]string{roomsLink[0], roomsLink[1]})
			}
		}
	}

	return start, end, rooms, links, antNumbers
}


// BuildGraph takes the list of rooms and links, and constructs an adjacency list
func BuildGraph(rooms []string, links [][2]string) map[string][]string {
	graph := make(map[string][]string)

	// Initialize each room in the graph
	for _, room := range rooms {
		graph[room] = []string{}
	}

	// Add edges for each link
	for _, link := range links {
		room1, room2 := link[0], link[1]
		graph[room1] = append(graph[room1], room2)
		graph[room2] = append(graph[room2], room1) // because the graph is undirected
	}

	return graph
}

// BFS function to find all shortest paths using a variation of BFS
func BFSAllPaths(graph map[string][]string, start, end string) [][]string {
	var paths [][]string
	queue := list.New()
	queue.PushBack([]string{start})

	for queue.Len() > 0 {
		path := queue.Remove(queue.Front()).([]string)
		room := path[len(path)-1]

		if room == end {
			paths = append(paths, path)
		}

		for _, neighbor := range graph[room] {
			// Skip if the room is already in the path (to avoid cycles)
			if !contains(path, neighbor) {
				newPath := append([]string{}, path...) // Copy the current path
				newPath = append(newPath, neighbor)
				queue.PushBack(newPath)
			}
		}
	}

	return paths
}

// Helper function to check if a path contains a room
func contains(path []string, room string) bool {
	for _, r := range path {
		if r == room {
			return true
		}
	}
	return false
}

// SimulateAnts function to simulate the movement of ants along multiple paths
func SimulateAnts(paths [][]string, antNumbers int, end string, graph map[string][]string) {
	ants := make([]int, antNumbers)      // Array to track which path each ant is on
	positions := make([]int, antNumbers) // Array to track each ant's current position on their path

	// Room occupancy tracker for each round (key is room name, value is if occupied or not)
	occupied := make(map[string]bool)

	// Assign each ant to a path
	for i := 0; i < antNumbers; i++ {
		ants[i] = i % len(paths) // Cycle through the available paths
	}

	round := 1
	for !allAntsAtEnd(positions, paths) {
		fmt.Printf("Round %d:\n", round)

		// Reset room occupancy for the round
		occupied = make(map[string]bool)

		for i := 0; i < antNumbers; i++ {
			// Current position and next position for this ant
			nextPos := positions[i] + 1

			// Current room the ant is in
			currentRoom := paths[ants[i]][positions[i]]

			// Check if there is a direct link to the end room
			if hasDirectPath(graph, currentRoom, end) {
				positions[i]++       // Move the ant directly to the end room
				occupied[end] = true // Mark end room as occupied
				fmt.Printf("L%d-%s ", i+1, end)
				continue
			}

			// Check if the ant can move to the next room
			if nextPos < len(paths[ants[i]]) {
				nextRoom := paths[ants[i]][nextPos]

				// If the next room is the end room, move the ant there
				if nextRoom == end {
					positions[i]++
					occupied[nextRoom] = true // Mark the end room as occupied
					fmt.Printf("L%d-%s ", i+1, nextRoom)
					continue
				}

				// If the next room is not occupied, move the ant there
				if !occupied[nextRoom] {
					positions[i]++
					occupied[nextRoom] = true // Mark the room as occupied
					fmt.Printf("L%d-%s ", i+1, nextRoom)
				}
			}
		}
		fmt.Println()
		round++
	}
}

// Function to check if there is a direct path between two rooms in the graph
func hasDirectPath(graph map[string][]string, roomA, roomB string) bool {
	for _, neighbor := range graph[roomA] {
		if neighbor == roomB {
			return true
		}
	}
	return false
}

// Check if all ants have reached the end room
func allAntsAtEnd(positions []int, paths [][]string) bool {
	for i, pos := range positions {
		if pos < len(paths[i%len(paths)])-1 { // If any ant is not at the end
			return false
		}
	}
	return true
}

func main() {
	if len(os.Args) != 2 {
		log.Fatal("Usage: go run main.go <datafile>")
	}

	// Call GetData and retrieve the start, end, rooms, links, and number of ants
	start, end, rooms, links, antNumbers := GetData(os.Args[1])
	fmt.Println("start room: ", start)
	fmt.Println("end room: ", end)
	fmt.Println("rooms: ", rooms)
	fmt.Println("links between rooms: ", links)
	fmt.Println("number of ants: ", antNumbers)

	fmt.Println()

	// Build the graph (adjacency list) from rooms and links
	graph := BuildGraph(rooms, links)

	// Use BFS to find all shortest paths from start to end
	paths := BFSAllPaths(graph, start, end)

	// Display the found paths
	fmt.Println("Found Paths:")
	for _, path := range paths {
		fmt.Println(path)
	}

	// Simulate ant movement along the found paths
	fmt.Println("\nAnt Movement Simulation:")
	SimulateAnts(paths, antNumbers, end, graph)
}
