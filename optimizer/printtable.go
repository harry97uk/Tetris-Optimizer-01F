package optimizer

import "fmt"

func PrintTable(arr [][]int) {
	for i := 0; i < len(arr); i++ {
		for j := 0; j < len(arr); j++ {
			if arr[i][j] == 0 {
				fmt.Print(string(rune(46)))
			}
			fmt.Print(string(rune(arr[i][j])))
		}
		fmt.Println()
	}
}
