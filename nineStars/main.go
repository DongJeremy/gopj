package main

import (
	"fmt"
)

var starCN = [9]string{"罗候", "土星", "水星", "金星", "太阳", "火星", "计都", "太阴", "木星"}
var starMale = [9]int{1, 2, 3, 4, 5, 6, 7, 8, 0}
var starFeMale = [9]int{6, 5, 0, 8, 7, 2, 1, 4, 3}

func main() {
	for i := 1; i < 91; i++ {
		fmt.Printf("%d %s, ", i, calcStar(i, true))
		if i%9 == 0 {
			fmt.Println()
		}
	}
	for i := 1; i < 91; i++ {
		fmt.Printf("%d %s, ", i, calcStar(i, false))
		if i%9 == 0 {
			fmt.Println()
		}
	}

}

func calcStar(age int, sex bool) string {
	var index int
	value := age % 9
	if sex {
		for i, v := range starMale {
			if v == value {
				index = i
				break
			}
		}
	} else {
		for i, v := range starFeMale {
			if v == value {
				index = i
				break
			}
		}
	}
	return starCN[index]
}
