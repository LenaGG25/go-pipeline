package main

import "sync"

type Executor interface {
	Execute(jobs ...Job)
}

type PipelineExecutor struct{}

func NewPipelineExecutor() Executor {
	return &PipelineExecutor{}
}

func (executor *PipelineExecutor) Execute(jobs ...Job) {
	channels := make([]chan any, len(jobs)+1)
	for i := 0; i < len(jobs)+1; i++ {
		channels[i] = make(chan any)
	}

	wg := sync.WaitGroup{}
	for i, job := range jobs {
		wg.Add(1)
		go func(jobFunc Job, input chan any, output chan any) {
			defer wg.Done()
			defer close(output)

			jobFunc(input, output)
		}(job, channels[i], channels[i+1])
	}
	wg.Wait()
}
