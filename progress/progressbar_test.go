package progress

import "time"

func ExampleProgressBar() {
	bar := StartProgressBar(20, "counting from 1 to 20")
	for i := 0; i < 20; i++ {
		time.Sleep(1)
		bar.Add(1)
	}
	bar.Done()
}
