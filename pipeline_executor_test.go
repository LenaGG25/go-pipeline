package main

import (
	"sync/atomic"
	"testing"
	"time"
)

func TestPipeline(t *testing.T) {
	testExpectedTime := 100 * time.Millisecond

	ok := true
	received := atomic.Int32{}
	freeFlowJobs := []Job{
		Job(func(input, output chan any) {
			output <- 1
			time.Sleep(10 * time.Millisecond)
			currReceived := received.Load()
			if currReceived == 0 {
				ok = false
			}
		}),
		Job(func(input, output chan any) {
			for range input {
				received.Add(1)
			}
		}),
	}

	start := time.Now()

	executor := NewPipelineExecutor()
	executor.Execute(freeFlowJobs...)

	end := time.Since(start)

	if !ok || received.Load() == 0 {
		t.Errorf(
			"no value free flow - dont collect them\nGot: ok = %v, received = %v \nExpected: ok = %v, received = %v",
			ok, received.Load(),
			false, 1,
		)
	}

	if end > testExpectedTime {
		t.Errorf("time limit\nGot: %s\nExpected: <= %s", end, testExpectedTime)
	}

}

func TestSigner(t *testing.T) {
	testExpectedTime := 3 * time.Second

	testExpected := "1173136728138862632818075107442090076184424490584241521304_1696913515191343735512658979631549563179965036907783101867_27225454331033649287118297354036464389062965355426795162684_29568666068035183841425683795340791879727309630931025356555_3994492081516972096677631278379039212655368881548151736_4958044192186797981418233587017209679042592862002427381542_4958044192186797981418233587017209679042592862002427381542"
	testResult := "NOT_SET"
	inputData := []int{0, 1, 1, 2, 3, 5, 8}

	hashSignPipeline := []Job{
		Job(func(input, output chan any) {
			for _, value := range inputData {
				output <- value
			}
		}),
		SingleHash,
		MultiHash,
		CombineResults,
		Job(func(input, output chan any) {
			testResult = (<-input).(string)
		}),
	}

	start := time.Now()

	executor := NewPipelineExecutor()
	executor.Execute(hashSignPipeline...)

	end := time.Since(start)

	if testExpected != testResult {
		t.Errorf("results not match\nGot: %v\nExpected: %v", testResult, testExpected)
	}

	if end > testExpectedTime {
		t.Errorf("time limit\nGot: %s\nExpected: < %s", end, testExpectedTime)
	}
}

func TestExtra(t *testing.T) {
	testExpectedTime := 400 * time.Millisecond
	var testExpected int32 = (1 + 3 + 4) * 3

	received := atomic.Int32{}

	freeFlowJobs := []Job{
		Job(func(input, output chan any) {
			output <- 1
			output <- 3
			output <- 4
		}),
		Job(func(input, output chan any) {
			for value := range input {
				output <- value.(int) * 3
				time.Sleep(100 * time.Millisecond)
			}
		}),
		Job(func(input, output chan any) {
			for value := range input {
				received.Add(int32(value.(int)))
			}
		}),
	}

	start := time.Now()

	executor := NewPipelineExecutor()
	executor.Execute(freeFlowJobs...)

	end := time.Since(start)

	if received.Load() != testExpected {
		t.Errorf("f3 have not collected inputs\nGot: %v \nExpected: %v", received.Load(), testExpected)
	}

	if end > testExpectedTime {
		t.Errorf("time limit\nGot: %s\nExpected: <= %s", end, testExpectedTime)
	}
}
