package main

import (
	"fmt"

	"github.com/kaliv0/ister/heap"
)

func main() {
	// Top-K elements
	vals := []int{1, 8, 3, 9, 0, 40, 2, 15, 77}
	less := func(a, b int) bool {
		return a > b
	}
	h := heap.FromSlice(less, vals)

	k := 5
	res := make([]int, 0, k)

	for i := 0; i < k; i++ {
		v, _ := h.Pop()
		res = append(res, v)
	}

	fmt.Print(res)
}
