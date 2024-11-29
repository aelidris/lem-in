package Finder

import (

//myStruct  "lemin/Structs"

)



func FindPaths(graph map[string][]string, start, end string) [][]string {


	var currentPath []string

	var dfs func(current string)

	visited := make(map[string]bool)

	var allPaths [][]string

	i := 0


	dfs = func(current string) {

		visited[current] = true

		currentPath = append(currentPath, current)




		if current == end {

			if len(allPaths) == i {

				allPaths = append(allPaths, []string{})
			}

			allPaths[i] = append(allPaths[i], currentPath...)

			i++



		} else {
	
			for _, Link := range graph[current] {

				if !visited[Link] {

					dfs(Link)

				}
			}
		}

	
		visited[current] = false
		currentPath = currentPath[:len(currentPath)-1]
	}












	dfs(start)


	return allPaths


}