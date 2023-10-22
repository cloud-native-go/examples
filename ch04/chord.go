package ch04

import "sync"

func Chord(sources ...<-chan int) <-chan []int {
	type input struct { // Used to send inputs
		idx, input int // between goroutines
	}

	dest := make(chan []int) // The output channel

	inputs := make(chan input) // An intermediate channel

	wg := sync.WaitGroup{} // Used to close channels when
	wg.Add(len(sources))   // all sources are closed

	for i, ch := range sources { // Start goroutine for each source
		go func(i int, ch <-chan int) {
			defer wg.Done() // Notify WaitGroup when ch closes

			for n := range ch {
				inputs <- input{i, n} // Transfer input to next goroutine
			}
		}(i, ch)
	}

	go func() {
		wg.Wait()     // Wait for all sources to close
		close(inputs) // then close inputs channel
	}()

	go func() {
		res := make([]int, len(sources))   // Slice for incoming inputs
		sent := make([]bool, len(sources)) // Slice to track sent status
		count := len(sources)              // Counter for channels

		for r := range inputs {
			res[r.idx] = r.input // Update incoming input

			if !sent[r.idx] { // First input from channel?
				sent[r.idx] = true
				count--
			}

			if count == 0 {
				c := make([]int, len(res)) // Copy and send inputs slice
				copy(c, res)
				dest <- c

				count = len(sources) // Reset counter
				clear(sent)          // Clear status tracker
			}
		}

		close(dest)
	}()

	return dest
}
