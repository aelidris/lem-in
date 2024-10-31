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
			if antNumbers <= 0 {
				log.Fatal("the number of ant should >= 1 !!")
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

func displayFinalPaths(paths [][]string, antNumbers int, lastRoom string) {
    fmt.Println(antNumbers)

    // Initialize positions of ants (-1 means not started yet)
    antPosition := make([]int, antNumbers)
    for i := 0; i < antNumbers; i++ {
        antPosition[i] = -1
    }

    // Assign initial paths to ants
    assignedPaths := make([]int, antNumbers)
    for i := 0; i < antNumbers; i++ {
        assignedPaths[i] = i % len(paths)
    }

    turn := 1
    completedAnts := 0

    // Iterate until all ants have reached the end room
    for completedAnts < antNumbers {
        fmt.Printf("Turn %d: ", turn)
        roomsOccupied := make(map[string]bool) // Track occupied rooms for this turn

        // Move ants in order
        for i := 0; i < antNumbers; i++ {
            pathIdx := assignedPaths[i]
            currentPosition := antPosition[i]

            // If the ant has reached the end room, skip it
            if currentPosition >= len(paths[pathIdx])-1 {
                continue
            }

            if currentPosition == -1 {
                // Ant has not started yet; try to move it out of the start room
                nextRoom := paths[pathIdx][1]
                if !roomsOccupied[nextRoom] {
                    fmt.Printf("L%d-%s ", i+1, nextRoom)
                    antPosition[i] = 1
                    roomsOccupied[nextRoom] = true
                }
            } else {
                // Ant is already on its path; try to move to the next room
                nextRoom := paths[pathIdx][currentPosition+1]
                if nextRoom == lastRoom || !roomsOccupied[nextRoom] {
                    // Move the ant to the next room if it is not occupied or if it’s the end room
                    fmt.Printf("L%d-%s ", i+1, nextRoom)
                    antPosition[i]++
                    if nextRoom != lastRoom {
                        roomsOccupied[nextRoom] = true
                    } else {
                        completedAnts++
                    }
                }
            }
        }

        fmt.Println()
        turn++
    }
}


// Function to filter paths with unique rooms, excluding rooms "1" and "0"
func filterUniquePaths(paths [][]string) [][]string {
	uniqueRooms := make(map[string]bool) // Track unique rooms
	var filteredPaths [][]string

	for _, path := range paths {
		hasUniqueRooms := true
		roomSet := make(map[string]bool)

		// Check each room in the path, skip room "1" and room "0"
		for _, room := range path {
			if room == "1" || room == "0" {
				continue
			}
			// Check if the room is already used in another path
			if uniqueRooms[room] || roomSet[room] {
				hasUniqueRooms = false
				break
			}
			roomSet[room] = true
		}

		// If all rooms in this path are unique, add the path to the result
		if hasUniqueRooms {
			filteredPaths = append(filteredPaths, path)
			// Mark these rooms as used
			for room := range roomSet {
				uniqueRooms[room] = true
			}
		}
	}

	return filteredPaths
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
	fmt.Println("Display the graph: ")
	fmt.Println(graph)


	fmt.Println()
	// Use BFS to find all shortest paths from start to end
	paths := BFSAllPaths(graph, start, end)
	if len(paths) == 0 {
		log.Fatal("In your data there is no path(s) from the start room to end room !!")
	}
	// Display the found paths
	fmt.Println("Found Paths:")
	fmt.Println(paths)
	fmt.Println()
	for _, path := range paths {
		fmt.Println(path)
	}
	fmt.Println("Paths after filtring: ")
	filtredPaths := filterUniquePaths(paths)
	fmt.Println(filtredPaths)
	fmt.Println()
	//Find the final paths to make all the ants in the end room
	displayFinalPaths(filtredPaths, antNumbers, end)
}
