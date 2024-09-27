package main

// Example of a simple pipeline.
// Source: https://go.dev/blog/pipelines#squaring-numbers

import "fmt"

func gen(nums ...int) <-chan int {
	out := make(chan int)
	go func() {
		for _, n := range nums {
			out <- n
		}
		close(out)
	}()
	return out
}

func sq(in <-chan int) <-chan int {
	out := make(chan int)
	go func() {
		// This for loop stops when the `in` channel is closed!
		for n := range in {
			out <- n * n
		}
		close(out)
	}()
	return out
}

func main() {
	// Set up the pipeline.
	c := gen(2, 3)
	out := sq(c)

	// Consume the output.
	fmt.Println(<-out) // 4
	fmt.Println(<-out) // 9

	// Alternatively, we can compose like this:
	for n := range sq(gen(2, 3)) {
		fmt.Println(n) // 4 then 9
	}

	// Taking it one step further, we can compose like this:
	for n := range sq(sq(gen(2, 3))) {
		fmt.Println(n) // 16 then 81
	}
}
