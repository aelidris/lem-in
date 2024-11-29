package Finder

import (
	"fmt"
	//"strings"
)

func PrintResult(paths [][]string, paths2 [][]string, ants int) {
	type Path struct {
		stepsRemaining int
		rooms          []string
	}

	antPaths := make([][]string, ants)

	antPaths2 := make([][]string, ants)

	pathData := make([]Path, len(paths))

	pathData2 := make([]Path, len(paths2))

	///////////////////////////////////////////////////////////

	for i, path := range paths {
		pathData[i] = Path{stepsRemaining: len(path) - 1, rooms: path[1:]}
	}

	for i, path := range paths2 {
		pathData2[i] = Path{stepsRemaining: len(path) - 1, rooms: path[1:]}
	}

	////////////////://////////////////////////////////////////

	for i := 0; i <= ants-1; i++ {

		selectedPath := 0
		selectedPath2 := 0

		for index, path := range pathData {
			if pathData[selectedPath].stepsRemaining > path.stepsRemaining {
				selectedPath = index
			}
		}

		for index, path := range pathData2 {
			if pathData2[selectedPath2].stepsRemaining > path.stepsRemaining {
				selectedPath2 = index
			}
		}

		// Initialize wait time slots for current ant path.

		antPaths[i] = make([]string, pathData[selectedPath].stepsRemaining-len(pathData[selectedPath].rooms))

		antPaths2[i] = make([]string, pathData2[selectedPath2].stepsRemaining-len(pathData2[selectedPath2].rooms))

		for _, room := range pathData[selectedPath].rooms {
			antPaths[i] = append(antPaths[i], fmt.Sprintf("L%d-%s", i+1, room))
		}

		for _, room := range pathData2[selectedPath2].rooms {
			antPaths2[i] = append(antPaths2[i], fmt.Sprintf("L%d-%s", i+1, room))
		}

		pathData[selectedPath].stepsRemaining++

		pathData2[selectedPath2].stepsRemaining++
	}

	displayResult(antPaths, antPaths2)
}

// displayResult outputs the results of the path calculation.
func displayResult(result [][]string, result2 [][]string) {
	 
	flag := 0
	flag2 := 0

	for step := 0; step < len(result[len(result)-1]); step++ {
		for _, antSteps := range result {
			if step <= len(antSteps)-1 && antSteps[step] != "" {
				flag = len(antSteps)
			}
		}
	}

	for step := 0; step < len(result2[len(result2)-1]); step++ {
		for _, antSteps := range result2 {
			if step <= len(antSteps)-1 && antSteps[step] != "" {
				flag2 = len(antSteps)
			}
		}
	}

	if flag2 <= flag {

		
		for step := 0; step < len(result2[len(result2)-1]); step++ {

			res2 := []string{}

			for _, antSteps := range result2 {
				if step <= len(antSteps)-1 && antSteps[step] != "" {
					res2 = append(res2, antSteps[step])
				}
			}

			for _, v := range res2 {
				fmt.Print(v + " ")
			}

			fmt.Println()
		}


	} else {
		for step := 0; step < len(result[len(result)-1]); step++ {

			res := []string{}

			for _, antSteps := range result {
				if step <= len(antSteps)-1 && antSteps[step] != "" {
					res = append(res, antSteps[step])
				}
			}

			for _, v := range res {
				fmt.Print(v + " ")
			}

			fmt.Println()

		}
	}
}
