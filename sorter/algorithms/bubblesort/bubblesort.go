package bubblesort

// 冒泡排序的核心理念是什么?那就是相邻两数比较,前面的数比后面的数小的话,就交换位置,
// 每次循环找到该次排序的最小值,然后放到该次循环数组的队尾,因此便利到最后,留的就是最大的数.
// 最佳情况：T(n) = O(n)
// 最差情况：T(n) = O(n2)
// 平均情况：T(n) = O(n2)

// BubbleSort 冒泡排序
func BubbleSort(values []int) {
	flag := true
	for i := 0; i < len(values)-1; i++ {
		flag = true
		for j := 0; j < len(values)-i-1; j++ {
			if values[j] > values[j+1] {
				values[j], values[j+1] = values[j+1], values[j]
				flag = false
			} // end if
		} // end for j = ...
		if flag == true {
			break
		}
	} // end for i = ...
}
