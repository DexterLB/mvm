package progress

import (
	"sync"
	"time"
)

func ExampleProgressBar() {
	bar := StartProgressBar(20, "counting from 1 to 20")
	for i := 0; i < 20; i++ {
		time.Sleep(1)
		bar.Add(1)
	}
	bar.Done()
}

func ExampleProgressBar_manyThreads() {
	// the progressbar is thread-safe!
	// this example should produce the same result as the above one

	foos := make(chan string)
	go func() {
		for i := 0; i < 20; i++ {
			time.Sleep(1)
			foos <- "foo"
		}
		close(foos)
	}()

	bar := StartProgressBar(20, "counting foos with 4 parallel threads")

	wg := sync.WaitGroup{}
	for i := 0; i < 4; i++ {
		go func() {
			for _ = range foos {
				bar.Add(1)
			}
			wg.Done()
		}()
	}

	wg.Wait()
	bar.Done()
}
