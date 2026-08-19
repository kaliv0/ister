package main

import (
	"fmt"

	"github.com/kaliv0/ister/list"
)

func main() {
	l := list.Of(1, 2, 3)
	l.Add(4)
	l.Reverse()
	l.RemoveAt(1)
	fmt.Println(l.Get(0))

	for el := range l.All() {
		fmt.Printf("%d ", el)
	}
	fmt.Println()

	l.Set(0, 80)
	fmt.Println(l)

	l.Clear()
	fmt.Print(l.Len())
}
