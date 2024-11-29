package Finder

import (
	"strings"
	// Data  "lemin/Structs"
)

func GraphMaker(Rooms, Links []string) map[string][]string {
	Graph := make(map[string][]string)

	for _, v := range Rooms {
		// room := Data.LinksAndAnts{
		// 	Links: []string{},

		// }

		Graph[v] = []string{}

		// initialize the graph with empty data (links - Ants)
	}

	for _, Link := range Links {

		v := strings.Split(Link, "-")

		room1 := v[0]
		room2 := v[1]

		tmp := Graph[room1]

		tmp = append(tmp, room2)

		Graph[room1] = tmp

		tmp2 := Graph[room2]

		tmp2 = append(tmp2, room1)

		Graph[room2] = tmp2

		// append bidirectionally the rooms and links , so (3-1) will be a room 3 with link 1 , and room 1 with 3 as a link !!!!!
	}

	
	return Graph
}
