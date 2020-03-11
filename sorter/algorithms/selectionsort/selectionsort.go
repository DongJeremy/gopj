package selectionsort

// 首先在未排序序列中找到最小（大）元素，存放到排序序列的起始位置，然后，
// 再从剩余未排序元素中继续寻找最小（大）元素，然后放到已排序序列的末尾。
// 以此类推，直到所有元素均排序完毕。
// 最佳情况：T(n) = O(n2)
// 最差情况：T(n) = O(n2)
// 平均情况：T(n) = O(n2)

func selectionSort(values []int, start int) {
	if start == len(values) {
		return
	}
	minIdx := start
	minVal := values[start]
	for i := start + 1; i < len(values); i++ {
		if values[i] < minVal {
			minIdx, minVal = i, values[i]
		}
	}
	values[start], values[minIdx] = values[minIdx], values[start]
	selectionSort(values, start+1)
}

// SelectionSort 选择排序
func SelectionSort(values []int) {
	selectionSort(values, 0)
}
