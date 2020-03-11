package qsort

// 快速排序是C.R.A.Hoare于1962年提出的一种划分交换排序。它采用了一种分治的策略，
// 通常称其为分治法(Divide-and-ConquerMethod)。

// 该方法的基本思想是：

// 1．先从数列中取出一个数作为基准数。
// 2．分区过程，将比这个数大的数全放到它的右边，小于或等于它的数全放到它的左边。
// 3．再对左右区间重复第二步，直到各区间只有一个数。
// 最佳情况：T(n) = O(nlogn)
// 最差情况：T(n) = O(n2)
// 平均情况：T(n) = O(nlogn)

func quickSort(values []int, left, right int) {
	temp := values[left]
	p := left
	i, j := left, right
	for i <= j {
		for j >= p && values[j] >= temp {
			j--
		}
		if j >= p {
			values[p] = values[j]
			p = j
		}
		if values[i] <= temp && i <= p {
			i++
		}
		if i <= p {
			values[p] = values[i]
			p = i
		}
	}
	values[p] = temp
	if p-left > 1 {
		quickSort(values, left, p-1)
	}
	if right-p > 1 {
		quickSort(values, p+1, right)
	}
}

// QuickSort 快速排序
func QuickSort(values []int) {
	quickSort(values, 0, len(values)-1)
}
