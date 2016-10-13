package command

import (
	"fmt"
	"testing"
	"time"
)

func TestMatrixChain(t *testing.T) {
	id := make(chan int, 1)
	next := make(chan bool, 1)
	computeSlice := make([]chan bool, 6)
	for i := 0; i < 5; i++ {
		computeSlice[i] = make(chan bool, 1)
	}
	input := []float64{1, 2, 3, 4}
	output := []string{"test1", "test2"}
	ResultSlice := make(chan []float64, 1)
	for i := 0; i < 5; i++ {
		// Run the function which handle the single slice
		workingSlice := make([]float64, 5)
		for j := 0; j < len(input); j++ {
			workingSlice[j] = input[j]
		}
		go generateSingleSlice(workingSlice, input, output, i, id, next, ResultSlice, computeSlice[i], computeSlice[i+1])
	}
	fmt.Println("pong")
	var supp bool
	computeFinish := &supp
	go func(t *testing.T) {
		time.Sleep(10 * time.Second)
		t.Fatalf("Compute slice time is finish")
	}(t)
	go func(computeFinish *bool) {
		*computeFinish = <-computeSlice[6]
	}(computeFinish)
	go func() {
		for {
			computeSlice[0] <- true
			_ = <-ResultSlice
			_ = <-id
			_ = <-next
		}
	}()

	fmt.Println("pong")
	for {
		if *computeFinish {
			break
		}
	}
}
