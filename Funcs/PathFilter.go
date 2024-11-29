package Finder



func FilterUniquePaths(paths [][]string) ([][]string, [][]string) {

	
	result := [][]string{}
   
	result2 := [][]string{}

	type roomData struct {
		visited bool

	}


	roomSet := make(map[string]roomData) // Track rooms already used in selected paths
	
    roomSet2 := make(map[string]roomData)

	for _, path := range paths {

		exist := false

		for i := 1; i < len(path)-1; i++ {
			if roomSet2[path[i]].visited {

				exist = true

				break

			}
		}

		if !exist {

			result2 = append(result2, path)

			for i := 1; i < len(path)-1; i++ {
				roomSet2[path[i]] = roomData{
					visited: true,
					 
				}
			}

		}

	}



	for i := 0; i < len(paths)-1; {

		if len(paths[i]) > len(paths[i+1]) {

			tmp := paths[i]
			paths[i] = paths[i+1] /////       SORTING paths
			paths[i+1] = tmp

			if i > 0 {
				i--
			}

		} else {
			i++
		}
	}

	for _, path := range paths {

		exist := false

		for i := 1; i < len(path)-1; i++ {
			if roomSet[path[i]].visited {

				exist = true

				break

			}
		}

		if !exist {

			result = append(result, path)

			for i := 1; i < len(path)-1; i++ {
				roomSet[path[i]] = roomData{
					visited: true,
					 
				}
			}

		}

	}



	return result , result2
}
