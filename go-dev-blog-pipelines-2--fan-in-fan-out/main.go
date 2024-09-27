package main

// Example of fan-in, fan-out pattern.
// Source: https://go.dev/blog/pipelines#fan-out-fan-in

import (
	"fmt"
	"sync"
)

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
		// ℹ️ This for loop stops when the `in` channel is closed.
		for n := range in {
			out <- n * n
		}
		close(out)
	}()
	return out
}

func merge(cs ...<-chan int) <-chan int {
	var wg sync.WaitGroup
	out := make(chan int)

	// Start an output goroutine for each input channel in cs.  output
	// copies values from c to out until c is closed, then calls wg.Done.
	output := func(c <-chan int) {
		for n := range c {
			// 🚂 Merge central! Will block until the value is read from the out channel.
			out <- n
		}
		// 🧮 Done decrements the WaitGroup counter by one.
		wg.Done()
	}

	// 🧮 Increment the WaitGroup counter by the number of goroutines we're going to start.
	wg.Add(len(cs))

	// 🎼 Composer.
	for _, c := range cs {
		go output(c)
	}

	// Start a goroutine to close out once all the output goroutines are
	// done.  This must start after the wg.Add call.
	go func() {
		// 🛑 Will block until the WaitGroup counter is zero!
		wg.Wait()
		// ℹ️ Closing the out channel will break the range loop in the main function.
		close(out)
	}()
	// 🟢 Return the out channel.
	return out
}

func main() {
	in := gen(2, 3)

	// Fan-out: Distribute the sq work across two goroutines that both read from in.
	c1 := sq(in)
	c2 := sq(in)

	// Fan-in: Consume the merged output from c1 and c2.
	for n := range merge(c1, c2) {
		fmt.Println(n) // 4 then 9, or 9 then 4
	}
}
