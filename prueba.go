package main

import "fmt"

func aristar(G []int) {
	n := len(G)
	for i := 0; i < n; i++ {
		fmt.Println("aristando ", i)
		G[i] = 1
	}
}

func main() {
	G := []int{1, 2, 3, 4, 5}
	aristar(G)
	fmt.Println(G)
}
