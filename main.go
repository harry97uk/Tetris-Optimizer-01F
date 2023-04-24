package main

import (
	"fmt"
	"tetris-optimizer/optimizer"
	"time"
)

func main() {
	fmt.Println("Program start: ", time.Now())
	defer timer("Main")()
	pieces := optimizer.FileReader("sample.txt")
	ft := optimizer.FindSmallestSquare(pieces)
	optimizer.PrintTable(ft)
	fmt.Println("Smallest square size is: ", len(ft))
}

func timer(name string) func() {
	start := time.Now()
	return func() {
		fmt.Printf("%s took %v\n", name, time.Since(start))
	}
}
