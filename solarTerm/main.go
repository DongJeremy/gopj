package main

import (
	"fmt"
	"strconv"
)

const (
	// D 节气D值
	D float64 = 0.2422
)

var C20 = [...]float64{
	6.11, 20.84, 4.6295, 19.4599, 6.3826, 21.4155, 5.59, 20.888,
	6.318, 21.86, 6.5, 22.2, 7.928, 23.65, 8.35, 23.95,
	8.44, 23.822, 9.098, 24.218, 8.218, 23.08, 7.9, 22.6,
}

var C21 = [...]float64{
	5.4055, 20.12, 3.87, 18.73, 5.63, 20.646, 4.81, 20.1,
	5.52, 21.04, 5.678, 21.37, 7.108, 22.83, 7.5, 23.13,
	7.646, 23.042, 8.318, 23.438, 7.438, 22.36, 7.18, 21.94,
}

var TERM = [...]string{
	"小寒", "大寒", "立春", "雨水", "惊蛰", "春分", "清明", "谷雨",
	"立夏", "小满", "芒种", "夏至", "小暑", "大暑", "立秋", "处暑",
	"白露", "秋分", "寒露", "霜降", "立冬", "小雪", "大雪", "冬至",
}

var termMap = make(map[int]map[string]string)

func main() {
	fmt.Println(getTermName(2020, 2, 17))
}

func getTermName(year, month, date int) string {
	map1 := getYearTermMap(year)
	fmt.Println(map1)
	if map1 == nil {
		return ""
	}
	return map1[getTermKey(month, date)]
}

// 获取指定年份萨拉表
func getYearTermMap(year int) map[string]string {
	// 处理世纪C值
	var c [24]float64
	if year > 1900 && year <= 2000 {
		c = C20
	} else if year > 2000 && year <= 2100 {
		c = C21
	} else {
		panic("不支持的年份")
	}
	// if map1, ok := termMap[year]; ok {
	// 	return map1
	// }

	map1 := termMap[year]
	if map1 == nil {
		y := year % 100
		map1 = make(map[string]string)
		for i := 0; i < 24; i++ {
			//计算节气日期，计算公式：[Y*D+C]-L
			var date int
			if i < 2 || i > 22 {
				date = int(float64(y)*D+c[i]) - int((y-1)/4)
			} else {
				date = int(float64(y)*D+c[i]) - int(y/4)
			}
			// 记录计算结果
			key := getTermKey(i/2+1, date)
			map1[key] = TERM[i]
		}
		// 计算结果添加到节气表
		termMap[year] = map1
	}
	return map1
}

func getTermKey(month, date int) string {
	key := strconv.Itoa(month)
	if month < 10 {
		key = "0" + key
	}
	if date < 10 {
		key += "0"
	}
	key += strconv.Itoa(date)
	return key
}
