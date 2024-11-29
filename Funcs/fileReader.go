package Finder

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

func GetData(dataFile string) (start, end string, rooms []string, links []string, antNumbers int) {
	data, err := os.ReadFile(dataFile)
	if err != nil {
		log.Fatal(err)
	}

	if len(data) == 0 {
		log.Fatal("ERROR: No data in this file!!")
	}

	checkData := strings.Split(string(data), "\n")
	canStart, canEnd := false, false
	for _, flag := range checkData {
		if flag == "##start" {
			canStart = true
		}
		if flag == "##end" {
			canEnd = true
		}
	}
	if !canStart || !canEnd {
		log.Fatal("ERROR: Check if the ##start command and ##end command exist")
	}
	affStart, affEnd := false, false
	findStart, findEnd := false, false

	for i := len(checkData) - 1; i >= 0; i-- {

		tmp := strings.TrimSpace(checkData[i])

		if tmp == "" {
			checkData = checkData[:i]
		} else {
			break
		}

	}

	for i := 0; i < len(checkData)-1; i++ {

		tmp := strings.TrimSpace(checkData[i])

		if tmp == "" {
			checkData = checkData[i+1:]

			i--
		} else {
			break
		}

	}

	for i, line := range checkData {

		if len(line) == 0 && i != 0 {
			log.Fatal("ERROR: invalid data format, empty line data!")
		}

		if strings.HasPrefix(line, "#") {
			if line == "##start" {
				if findStart {
					log.Fatal("ERROR: Issue with the start command (##start)")
				}
				affStart = true
				findStart = true
			} else if line == "##end" {
				if findEnd {
					log.Fatal("ERROR: Issue with the end command (##end)")
				}
				affEnd = true
				findEnd = true
			} else if len(line) > 1 && strings.HasPrefix(line, "##") && line != "##start" && line != "##end" {
				log.Fatal("ERROR: Only ##start and ##end are allowed commands")
			}
			continue
		}

		if affStart {
			if affEnd {
				log.Fatal("ERROR: invalid data format, can't find the start room")
			}
			start = strings.Fields(line)[0]
			affStart = false

		}
		// If this is the first non-comment line after ##end, capture the end room
		if affEnd {
			if affEnd {
				log.Fatal("ERROR: invalid data format, can't find the end room")
			}
			end = strings.Fields(line)[0]
			affEnd = false
		}

		if antNumbers == 0 {
			antNumbers, err = strconv.Atoi(line)
			if err != nil {
				log.Fatal("Error converting number of ants to int:", err)
			}
			if antNumbers <= 0 {
				log.Fatal("ERROR: invalid number of ants")
			}
			continue
		}
		parts := strings.Fields(line)
		if len(parts) == 3 {
			rooms = append(rooms, parts[0])
		}
		if strings.Contains(line, "-") {
			x := strings.Split(line, "-")
			if len(x) != 2 || x[0] == "" || x[1] == "" || x[0] == x[1] {
				log.Fatal("ERROR: Invalid data format")
			}

			validRoomName := func(name string) bool {
				for _, room := range rooms {
					if room == name {
						return true
					}
				}

				return false
			}

			if !validRoomName(x[0]) || !validRoomName(x[1]) {
				log.Fatal("ERROR: Room doesn't exist")
			}

			links = append(links, line)
		}
	}

	for i, v := range rooms {
		for j, vv := range rooms {
			if i != j && vv == v {
				log.Fatal("ERROR: Room duplicated room !!!")
			}
		}
	}

	for i, v := range links {
		for j, vv := range links {
			if i != j && vv == v {
				log.Fatal("ERROR: Room duplicated links !!!")
			}
		}
	}

	for _, line := range checkData {
		fmt.Println(line)
	}

	fmt.Println()
	return start, end, rooms, links, antNumbers
}
