package main

import (
	"container/list"
	"fmt"
	"os"

	"github.com/mohae/deepcopy"
)

// Recoder r
type Recoder struct {
	point      [2]int
	pointIndex int
	sudokuList [9][9]interface{}
}

// Sudo Sudo
type Sudo struct {
	guessTimes int
	newPoints  *list.List
	recoder    *list.List
	basePoints [9][2]int
	sudokuList [9][9]interface{}
}

// New New
func New(sudoArr [9][9]int) *Sudo {
	var sudokuList [9][9]interface{}
	newPoints := list.New()
	recoder := list.New()
	for i := 0; i < 81; i++ {
		if sudoArr[i/9][i%9] != 0 {
			sudokuList[i/9][i%9] = sudoArr[i/9][i%9]
			newPoints.PushBack([2]int{i / 9, i % 9})
		} else {
			sudokuList[i/9][i%9] = []int{1, 2, 3, 4, 5, 6, 7, 8, 9}
		}
	}
	basePoints := [9][2]int{
		{0, 0}, {0, 3}, {0, 6}, {3, 0}, {3, 3}, {3, 6}, {6, 0}, {6, 3}, {6, 6},
	}
	return &Sudo{basePoints: basePoints, guessTimes: 0,
		newPoints: newPoints, recoder: recoder, sudokuList: sudokuList}
}

// Calc Calc
func (sudo *Sudo) Calc() {
	sudo.solveSudo()
	for {
		if sudo.checkValue() {
			if sudo.getCount() == 81 {
				break
			}
			point := sudo.getBestPoint()
			sudo.recodeGuess(point, 0)
		} else {
			sudo.reback()
		}
	}
}

func (sudo *Sudo) solveSudo() {
	isRunSame := true
	isRunOne := true
	for isRunSame {
		for isRunOne {
			for sudo.newPoints.Len() > 0 {
				e := sudo.newPoints.Front()
				sudo.newPoints.Remove(e)
				sudo.cutNum(e.Value.([2]int))
			}
			isRunOne = sudo.checkOnePossbile()
		}
		isRunSame = sudo.checkSameNum()
		isRunOne = true
	}
}

func (sudo *Sudo) cutNum(point [2]int) {
	val, err := sudo.sudokuList[point[0]][point[1]].(int)
	if !err {
		return
	}
	//行排除
	for col := 0; col < 9; col++ {
		item, ok := sudo.sudokuList[point[0]][col].([]int)
		if !ok {
			continue
		}
		key := In(val, item)
		if key == -1 {
			continue
		}
		temp := CopySlice(item, key)
		sudo.sudokuList[point[0]][col] = temp
		if len(temp) == 1 {
			sudo.newPoints.PushFront([2]int{point[0], col})
			sudo.sudokuList[point[0]][col] = temp[0]
		}
	}

	//列排除
	for row := 0; row < 9; row++ {
		item, ok := sudo.sudokuList[row][point[1]].([]int)
		if !ok {
			continue
		}
		key := In(val, item)
		if key == -1 {
			continue
		}
		temp := CopySlice(item, key)
		sudo.sudokuList[row][point[1]] = temp
		if len(temp) == 1 {
			sudo.newPoints.PushFront([2]int{row, point[1]})
			sudo.sudokuList[row][point[1]] = temp[0]
		}
	}

	//九宫格排除
	x := point[0] / 3 * 3
	y := point[1] / 3 * 3
	for row := x; row < x+3; row++ {
		for col := y; col < y+3; col++ {
			if _, ok := sudo.sudokuList[row][col].([]int); !ok {
				continue
			}
			key := In(val, sudo.sudokuList[row][col].([]int))
			if key == -1 {
				continue
			}
			temp := CopySlice(sudo.sudokuList[row][col].([]int), key)
			sudo.sudokuList[row][col] = temp
			if len(temp) == 1 {
				sudo.newPoints.PushFront([2]int{row, col})
				sudo.sudokuList[row][col] = temp[0]
			}
		}
	}
}

func (sudo *Sudo) checkOnePossbile() bool {
	for r := 0; r < 9; r++ {
		values := make(map[int][]int, 0)
		for col := 0; col < 9; col++ {
			if item, ok := sudo.sudokuList[r][col].([]int); ok {
				values[col] = item
			}
		}
		for k := 0; k < 9; k++ {
			if _, ok := sudo.sudokuList[r][k].([]int); !ok {
				continue
			}
			for _, val := range sudo.sudokuList[r][k].([]int) {
				sum := sudo.ergodic(values, val)
				if sum == 1 {
					sudo.sudokuList[r][k] = val
					sudo.newPoints.PushFront([2]int{r, k})
					return true
				}
			}
		}
	}

	for c := 0; c < 9; c++ {
		values := make(map[int][]int, 0)
		for row := 0; row < 9; row++ {
			if item, ok := sudo.sudokuList[row][c].([]int); ok {
				values[row] = item
			}
		}
		for r := 0; r < 9; r++ {
			if _, ok := sudo.sudokuList[r][c].([]int); !ok {
				continue
			}
			for _, val := range sudo.sudokuList[r][c].([]int) {
				sum := sudo.ergodic(values, val)
				if sum == 1 {
					sudo.sudokuList[r][c] = val
					sudo.newPoints.PushFront([2]int{r, c})
					return true
				}
			}
		}
	}

	for _, val := range sudo.basePoints {
		values := make(map[int][]int, 0)
		for key := val[0]; key < val[0]+3; key++ {
			for i := val[1]; i < val[1]+3; i++ {
				if item, ok := sudo.sudokuList[key][i].([]int); ok {
					values[key*3+i] = item
				}
			}
		}
		for key := val[0]; key < val[0]+3; key++ {
			for i := val[1]; i < val[1]+3; i++ {
				if _, ok := sudo.sudokuList[key][i].([]int); !ok {
					continue
				}

				for _, v := range sudo.sudokuList[key][i].([]int) {
					sum := sudo.ergodic(values, v)
					if sum == 1 {
						sudo.sudokuList[key][i] = v
						sudo.newPoints.PushFront([2]int{key, i})
						return true
					}
				}
			}
		}
	}
	return false
}

func (sudo *Sudo) ergodic(list map[int][]int, search int) int {
	sum := 0
	for _, item := range list {
		for _, val := range item {
			if search == val {
				sum++
			}
		}
	}
	return sum
}

func (sudo *Sudo) checkSameNum() bool {
	for _, val := range sudo.basePoints {
		var data []interface{}
		for x := val[0]; x < val[0]+3; x++ {
			for y := val[1]; y < val[1]+3; y++ {
				data = append(data, sudo.sudokuList[x][y])
			}
		}
		for i := 1; i < 10; i++ {
			var result []int
			for blockKey, item := range data {
				list, ok := item.([]int)
				if !ok {
					result = append(result, 0)
					continue
				}
				if retKey := In(i, list); retKey > -1 {
					result = append(result, blockKey+1)
				} else {
					result = append(result, 0)
				}
			}
			var temp []int
			for _, val := range result {
				if val > 0 {
					temp = append(temp, val)
				}
			}
			result = temp
			if rCount := len(result); rCount == 2 || rCount == 3 {
				rows := Map(func(x int) int { return (x - 1) / 3 }, result)
				cols := Map(func(x int) int { return (x - 1) % 3 }, result)
				if len(removeDuplicates(rows)) == 1 {
					row := val[0] + rows[0]
					result = Map(func(x int) int { return val[1] + (x-1)%3 }, result)

					for col := 0; col < 9; col++ {
						if In(col, result) > -1 {
							continue
						}
						if _, ok := sudo.sudokuList[row][col].([]int); !ok {
							continue
						}
						if sudo.replace(row, col, i) {
							return true
						}
					}
				} else if len(removeDuplicates(cols)) == 1 {
					result = Map(func(x int) int { return val[0] + (x-1)/3 }, result)
					col := val[1] + cols[0]

					for row := 0; row < 9; row++ {
						if In(row, result) > -1 {
							continue
						}
						if _, ok := sudo.sudokuList[row][col].([]int); !ok {
							continue
						}

						if sudo.replace(row, col, i) {
							return true
						}
					}
				}

			}
		}
	}
	return false
}

func (sudo *Sudo) replace(row int, col int, search int) bool {
	key := In(search, sudo.sudokuList[row][col].([]int))
	if key == -1 {
		return false
	}
	temp := deepcopy.Copy(sudo.sudokuList[row][col]).([]int)
	temp = append(temp[:key], temp[key+1:]...)
	sudo.sudokuList[row][col] = temp
	if len(temp) == 1 {
		sudo.newPoints.PushFront([2]int{row, col})
		sudo.sudokuList[row][col] = temp[0]
		return true
	}
	return false
}

//得到确定的数字
func (sudo *Sudo) getCount() int {
	sum := 0
	for i := 0; i < 81; i++ {
		if _, ok := sudo.sudokuList[i/9][i%9].(int); ok {
			sum++
		}
	}
	return sum
}

//评分，找到最佳的猜测坐标
func (sudo *Sudo) getBestPoint() [2]int {
	bestScore := 0
	bestPoint := [2]int{0, 0}
	for row := range [9]int{0: 9} {
		for col := range sudo.sudokuList[row] {
			pointScore := sudo.getPointScore(row, col)
			if bestScore < pointScore {
				bestScore = pointScore
				bestPoint = [2]int{row, col}
			}
		}
	}
	return bestPoint
}

//计算某坐标的评分
func (sudo *Sudo) getPointScore(row int, col int) int {
	if _, ok := sudo.sudokuList[row][col].([]int); !ok {
		return 0
	}
	score := 10 - len(sudo.sudokuList[row][col].([]int))
	for _, val := range sudo.sudokuList[row] {
		if _, ok := val.(int); ok {
			score++
		}
	}
	for i := 0; i < 9; i++ {
		if _, ok := sudo.sudokuList[i][col].(int); ok {
			score++
		}
	}
	return score
}

//检查数字有无错误
func (sudo *Sudo) checkValue() bool {
	for row := range sudo.sudokuList {
		nums := make([]int, 0)
		lists := make([][]int, 0)
		for _, val := range sudo.sudokuList[row] {
			if list, ok := val.([]int); ok {
				lists = append(lists, list)
			} else {
				nums = append(nums, val.(int))
			}
		}
		if isRetBol(nums, lists) == false {
			return false
		}
	}

	for i := 0; i < 9; i++ {
		nums := make([]int, 0)
		lists := make([][]int, 0)
		for j := 0; j < 9; j++ {
			if list, ok := sudo.sudokuList[j][i].([]int); ok {
				lists = append(lists, list)
			} else {
				nums = append(nums, sudo.sudokuList[j][i].(int))
			}
		}
		if isRetBol(nums, lists) == false {
			return false
		}
	}

	for _, val := range sudo.basePoints {
		for key := range sudo.sudokuList[val[0] : val[0]+3] {
			nums := make([]int, 0)
			lists := make([][]int, 0)
			for i := val[1]; i < val[1]+3; i++ {
				if list, ok := sudo.sudokuList[key][i].([]int); ok {
					lists = append(lists, list)
				} else {
					nums = append(nums, sudo.sudokuList[key][i].(int))
				}
			}
			if isRetBol(nums, lists) == false {
				return false
			}
		}
	}
	return true
}

//猜测记录
func (sudo *Sudo) recodeGuess(point [2]int, index int) {
	var recoder Recoder
	recoder.point = point
	recoder.pointIndex = index

	recoder.sudokuList = deepcopy.Copy(sudo.sudokuList).([9][9]interface{})
	sudo.recoder.PushBack(recoder)
	sudo.guessTimes++

	//新一轮的排除处理
	item := sudo.sudokuList[point[0]][point[1]].([]int)
	sudo.sudokuList[point[0]][point[1]] = item[index]
	sudo.newPoints.PushFront(point)
	sudo.solveSudo()
}

func (sudo *Sudo) reback() {
	var index int
	var point [2]int
	var recoder Recoder
	for {
		if sudo.recoder.Len() == 0 {
			fmt.Println("sudo is wrong")
			os.Exit(1)
		} else {
			e := sudo.recoder.Back()
			sudo.recoder.Remove(e)
			recoder = e.Value.(Recoder)
			point = recoder.point
			index = recoder.pointIndex + 1
			item := recoder.sudokuList[point[0]][point[1]].([]int)
			if index < len(item) {
				break
			}
		}
	}
	sudo.sudokuList = recoder.sudokuList
	sudo.recodeGuess(point, index)
}

func isRetBol(nums []int, lists [][]int) bool {
	if len(removeDuplicates(nums)) != len(nums) {
		return false
	}
	for _, val := range lists {
		if len(val) == 0 {
			return false
		}
	}
	return true
}

func removeDuplicates(elements []int) []int {
	encountered := map[int]bool{}
	result := []int{}

	for v := range elements {
		if encountered[elements[v]] != true {
			encountered[elements[v]] = true
			result = append(result, elements[v])
		}
	}
	return result
}

// Map Map
func Map(f func(int) int, v []int) (r []int) {
	r = make([]int, len(v))
	for i, value := range v {
		r[i] = f(value)
	}
	return
}

// In In
func In(search int, value []int) int {
	for k, val := range value {
		if search == val {
			return k
		}
	}
	return -1
}

// CopySlice CopySlice
func CopySlice(arr []int, key int) []int {
	var list []int
	for k, val := range arr {
		if k == key {
			continue
		}
		list = append(list, val)
	}
	return list
}

func main() {
	// sudoArr := [9][9]int{
	// 	{0, 0, 0, 0, 0, 2, 0, 5, 0},
	// 	{0, 7, 8, 0, 0, 0, 3, 0, 0},
	// 	{0, 0, 0, 0, 0, 4, 0, 0, 0},
	// 	{5, 0, 0, 0, 0, 0, 0, 0, 0},
	// 	{0, 0, 0, 0, 0, 0, 1, 0, 0},
	// 	{0, 0, 0, 0, 3, 0, 7, 0, 8},
	// 	{2, 0, 0, 0, 0, 0, 0, 4, 0},
	// 	{0, 0, 0, 0, 0, 5, 0, 9, 0},
	// 	{0, 1, 0, 0, 7, 0, 0, 0, 0},
	// }
	sudoArr := [9][9]int{
		{8, 0, 0, 0, 0, 0, 0, 0, 0},
		{0, 0, 3, 6, 0, 0, 0, 0, 0},
		{0, 7, 0, 0, 9, 0, 2, 0, 0},
		{0, 5, 0, 0, 0, 7, 0, 0, 0},
		{0, 0, 0, 0, 4, 5, 7, 0, 0},
		{0, 0, 0, 1, 0, 0, 0, 3, 0},
		{0, 0, 1, 0, 0, 0, 0, 6, 8},
		{0, 0, 8, 5, 0, 0, 0, 1, 0},
		{0, 9, 0, 0, 0, 0, 4, 0, 0},
	}
	data := New(sudoArr)
	data.Calc()
	fmt.Printf("完成，猜测了%d次\n", data.guessTimes)
	for _, item := range data.sudokuList {
		fmt.Println(item)
	}
}
