package main

import (
	"log"
	"os"

	myFuncs "lemin/Funcs"
)

func main() {
	if len(os.Args) != 2 {
		log.Fatal("ERROR: Invalid Argument. Usage: ")
	}
	Start, End, Rooms, Links, antNumbers := myFuncs.GetData(os.Args[1])

	myGraph := myFuncs.GraphMaker(Rooms, Links)

	myPaths := myFuncs.FindPaths(myGraph, Start, End)

	if len(myPaths) == 0 {
		log.Fatal("ERROR: No path exists from start to end")
	}

	FilteredPaths, FilteredPaths2 := myFuncs.FilterUniquePaths(myPaths)

	myFuncs.PrintResult(FilteredPaths, FilteredPaths2, antNumbers)
}
