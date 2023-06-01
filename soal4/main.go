package main

import (
	"fmt"
	"sort"
)
 
func main() {
	ArrBilangan := []float64{4, -7, -5, 3, 3.3, 9, 0, 10, 0.2}

	// order float ascending
	sort.Float64s(ArrBilangan)
	fmt.Println("Sort Float Ascending : ", ArrBilangan)

	// order float decending
	sort.Sort(sort.Reverse(sort.Float64Slice(ArrBilangan)))
	fmt.Println("Sort Float Descending : ", ArrBilangan)
}