package optimizer

import (
	"math"
	"sort"
)

type SolvedSquare struct {
	Values [][]int
}

type Position struct {
	row          int
	col          int
	adjacent     int
	centre       int
	availability int
}

func FindSmallestSquare(t []Tetrominoe) [][]int {
	sn := math.Sqrt(float64(len(t) * 4))
	snr := int(math.Round(sn) + 0.5)
	FinalSquare := SolvedSquare{}
	TempSquares := []SolvedSquare{}

	c := make(chan SolvedSquare, 1)

	for snr <= len(t)*4 {

		orderset := []int{}
		for i := 0; i < snr-1; i++ {
			orderset = append(orderset, i)
		}
		customOrders := GetDifferentVarieties(orderset)

		for n := 0; n < len(customOrders); n++ {
			TempSquares = append(TempSquares, SolvedSquare{})
			TempSquares[n].Values = make([][]int, snr)
			for j := 0; j < snr; j++ {
				TempSquares[n].Values[j] = make([]int, snr)
			}
			if n != 0 {
				go func(n int) {
					solved, finalarr := SolveTetrominoesCustomOrder(TempSquares[n].Values, snr, t, len(t)*4, customOrders[n])
					if solved {
						TempSquares[n].Values = finalarr
						c <- TempSquares[n]
					}
				}(n)
			}
		}
		go func() {
			solved, finalarr := SolveTetrominoes(TempSquares[0].Values, snr, t, len(t)*4)
			if solved {
				TempSquares[0].Values = finalarr
				c <- TempSquares[0]
			}
		}()
		FinalSquare = <-c
		if FinalSquare.Values != nil {
			break
		}
		snr++
	}
	return FinalSquare.Values
}

func SolveTetrominoes(arr [][]int, length int, t []Tetrominoe, total int) (bool, [][]int) {
	row := -1
	column := -1
	usedt := make([]Tetrominoe, 0)

	if IsValidFinalSquare(t, arr, length, total) {
		return true, arr
	}

	for i := 0; i < length; i++ {
		for j := 0; j < length; j++ {
			if (arr)[i][j] == 0 && CheckAdjacentSquares(arr, i, j) {
				row = i
				column = j
				for number := 0; number < len(t); number++ {
					samet := TetrominoeHasBeenUsed(usedt, t[number])
					if !samet && isCorrect(arr, row, column, t[number]) {
						SetSquareData(&arr, t[number], row, column, int(t[number].letter))
						tempt := CreateCopyOfTetrominoes(t, number)
						solved, returnarr := SolveTetrominoes(arr, length, tempt, total)
						if solved {
							return true, returnarr
						} else {
							SetSquareData(&arr, t[number], row, column, 0)
							usedt = append(usedt, t[number])
						}
					}
				}
				usedt = make([]Tetrominoe, 0)
			}
		}
	}
	return false, nil
}

func SolveTetrominoesReversed(arr [][]int, length int, t []Tetrominoe, total int) bool {
	row := -1
	column := -1
	usedt := make([]Tetrominoe, 0)

	if IsValidFinalSquare(t, arr, length, total) {
		return true
	}

	for i := length - 1; i >= 0; i-- {
		for j := length - 1; j >= 0; j-- {
			if (arr)[i][j] == 0 && CheckAdjacentSquares(arr, i, j) {
				row = i
				column = j
				for number := 0; number < len(t); number++ {
					samet := TetrominoeHasBeenUsed(usedt, t[number])
					if !samet && isCorrect(arr, row, column, t[number]) {
						SetSquareData(&arr, t[number], row, column, int(t[number].letter))
						tempt := CreateCopyOfTetrominoes(t, number)
						if SolveTetrominoesReversed(arr, length, tempt, total) {
							return true
						} else {
							SetSquareData(&arr, t[number], row, column, 0)
							usedt = append(usedt, t[number])
						}
					}
				}
				usedt = make([]Tetrominoe, 0)
			}
		}
	}
	return false
}

func SolveTetrominoesCustomOrder(arr [][]int, length int, t []Tetrominoe, total int, orderI []int) (bool, [][]int) {
	usedt := make([]Tetrominoe, 0)

	if IsValidFinalSquare(t, arr, length, total) {
		return true, arr
	}

	for _, i := range orderI {
		for _, j := range orderI {
			row, column := i, j
			if (arr)[row][column] == 0 && CheckAdjacentSquares(arr, row, column) {
				for number := 0; number < len(t); number++ {
					if !TetrominoeHasBeenUsed(usedt, t[number]) && isCorrect(arr, row, column, t[number]) {
						SetSquareData(&arr, t[number], row, column, int(t[number].letter))
						tempt := CreateCopyOfTetrominoes(t, number)
						solved, returnarr := SolveTetrominoesCustomOrder(arr, length, tempt, total, orderI)
						if solved {
							return true, returnarr
						} else {
							SetSquareData(&arr, t[number], row, column, 0)
							usedt = append(usedt, t[number])
						}
					}
				}
				usedt = make([]Tetrominoe, 0)
			}
		}
	}
	return false, nil
}

func SolveTetrominoesHeuristic(arr [][]int, length int, t []Tetrominoe, total, depth, maxDepth int) (bool, [][]int) {
	if depth >= maxDepth {
		return false, nil
	}

	squares := sortPositions(arr, t)
	c := make(chan Result, len(squares))

	for _, square := range squares {
		row, column := square.row, square.col
		if square.adjacent > 0 && square.availability > 0 {
			if (arr)[row][column] == 0 && CheckAdjacentSquares(arr, row, column) {
				for number := 0; number < len(t); number++ {
					if isCorrect(arr, row, column, t[number]) {
						go func(arr [][]int, t []Tetrominoe, number, row, column int) {
							temp := make([][]int, len(arr))
							for i := range arr {
								temp[i] = make([]int, len(arr[i]))
								copy(temp[i], arr[i])
							}
							SetSquareData(&temp, t[number], row, column, int(t[number].letter))
							tempt := CreateCopyOfTetrominoes(t, number)
							solved, returnarr := SolveTetrominoesHeuristic(temp, length, tempt, total, depth+1, maxDepth)
							if solved {
								c <- Result{returnarr, true}
							}
						}(arr, t, number, row, column)
					}
				}
			}
		}
	}

	var result Result
	for i := 0; i < len(squares); i++ {
		select {
		case res := <-c:
			if res.solved && (result.arr == nil || len(res.arr) > len(result.arr)) {
				result = res
			}
		}
	}
	if result.arr != nil {
		return true, result.arr
	}
	return false, nil
}

type Result struct {
	arr    [][]int
	solved bool
}

func TetrominoeHasBeenUsed(usedt []Tetrominoe, t Tetrominoe) bool {
	for _, v := range usedt {
		if SameTetrominoe(v, t) {
			return true
		}
	}
	return false
}

func CreateCopyOfTetrominoes(t []Tetrominoe, number int) []Tetrominoe {
	tempt := make([]Tetrominoe, len(t))
	copy(tempt, t)
	tempt = append(tempt[:number], tempt[number+1:]...)
	return tempt
}

func IsValidFinalSquare(t []Tetrominoe, arr [][]int, length int, total int) bool {
	if len(t) == 0 && CountArrZeros(arr) == ((length*length)-total) {
		return true
	}
	return false
}

func SameTetrominoe(k1, k2 Tetrominoe) bool {
	for i := 0; i < len(k1.squares); i++ {
		if k1.squares[i] != k2.squares[i] {
			return false
		}
	}
	return true
}

func CountArrZeros(arr [][]int) int {
	counter := 0
	for i := range arr {
		for j := range arr {
			if arr[i][j] == 0 {
				counter++
			}
		}
	}
	return counter
}

func isCorrect(arr [][]int, row int, column int, k Tetrominoe) bool {
	for i := 0; i < 4; i++ {
		if row+(k.squares[i].x) < len(arr) && column+(k.squares[i].y) < len(arr) {
			if (arr)[row+(k.squares[i].x)][column+(k.squares[i].y)] != 0 {
				return false
			}
		} else {
			return false
		}
	}
	return true
}

func SetSquareData(arr *[][]int, k Tetrominoe, row, column, value int) {
	for i := 0; i < len(k.squares); i++ {
		(*arr)[row+(k.squares[i].x)][column+(k.squares[i].y)] = value
	}
}

func CheckAdjacentSquares(arr [][]int, i, j int) bool {
	length := len(arr)
	if i == length-1 && j == length-1 {
		return false
	}
	if i < length-1 {
		if arr[i+1][j] == 0 {
			return true
		}
	}
	if j < length-1 {
		if arr[i][j+1] == 0 {
			return true
		}
	}
	return false
}

func ReverseSet(set []int) []int {
	reversed := make([]int, len(set))

	for i, j := 0, len(set)-1; i < len(set); i, j = i+1, j-1 {
		reversed[i] = set[j]
	}

	return reversed
}

func (s Position) Less(other Position) bool {
	return (s.adjacent*1)-(s.centre*1)+(s.availability*1) > (other.adjacent*1)-(other.centre*1)+(other.availability*1)
}

func sortPositions(arr [][]int, t []Tetrominoe) []Position {
	squares := make([]Position, 0)
	for i := 0; i < len(arr); i++ {
		for j := 0; j < len(arr[0]); j++ {
			if arr[i][j] == 0 {
				adjacent := 0
				centre := 0
				availability := 0
				if i < len(arr)-2 && arr[i+2][j] == 0 {
					adjacent++
				}
				if i < len(arr)-1 && arr[i+1][j] == 0 {
					adjacent++
				}
				if j < len(arr[0])-2 && arr[i][j+2] == 0 {
					adjacent++
				}
				if j < len(arr[0])-1 && arr[i][j+1] == 0 {
					adjacent++
				}
				if j < len(arr[0])-1 && i < len(arr)-1 && arr[i+1][j+1] == 0 {
					adjacent++
				}
				if j < len(arr[0])-2 && i < len(arr)-1 && arr[i+1][j+2] == 0 {
					adjacent++
				}
				if j < len(arr[0])-1 && i < len(arr)-2 && arr[i+2][j+1] == 0 {
					adjacent++
				}
				if j < len(arr[0])-2 && i < len(arr)-2 && arr[i+2][j+2] == 0 {
					adjacent++
				}
				for n := 0; n < len(t); n++ {
					if isCorrect(arr, i, j, t[n]) {
						availability++
						if CanClearMultipleRows(arr, t[n], i, j) {
							availability += 2
						}
						SetSquareData(&arr, t[n], i, j, int(t[n].letter))
						adjacent += PredictFutureDifficulty(arr)
						SetSquareData(&arr, t[n], i, j, 0)
					} else {
						adjacent += 500
					}
				}
				centre = int(math.Abs(float64((i - (len(arr) / 2)) + (j - (len(arr) / 2)))))
				adjacent += PredictFutureDifficulty(arr)
				squares = append(squares, Position{i, j, adjacent, centre, availability})
			}
		}
	}
	sort.Slice(squares, func(i, j int) bool {
		return squares[i].Less(squares[j])
	})
	return squares
}

func CanClearMultipleRows(board [][]int, tetrominoe Tetrominoe, row, col int) bool {

	var clearedRows []int
	for i := 0; i < tetrominoe.height; i++ {
		fullRow := true
		for j := 0; j < len(board[0]); j++ {
			if board[row+i][j] == 0 {
				fullRow = false
				break
			}
		}
		if fullRow {
			clearedRows = append(clearedRows, row+i)
		}
	}

	return len(clearedRows) > 1
}

func PredictFutureDifficulty(board [][]int) int {
	difficulty := 0
	for i := 0; i < len(board); i++ {
		rowHoles := 0
		for j := 0; j < len(board[0]); j++ {
			if board[i][j] == 0 {
				rowHoles++
			}
		}
		difficulty += rowHoles * rowHoles
	}
	return difficulty
}
