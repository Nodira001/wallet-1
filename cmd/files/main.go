package main

import "fmt"

func main() {
	a := aaa(4)

	fmt.Println(a)
}
func aaa(count int) [][]int {

	payments := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11}
	grouped := [][]int{}

	for i := 0; i < len(payments); i++ {

		if i+count > len(payments)-1 {

			grouped = append(grouped, payments[i:])

			break
		}

		grouped = append(grouped, payments[i:i+count])

		i += count - 1
	}

	return grouped
}
