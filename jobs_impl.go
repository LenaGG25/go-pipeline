package main

import (
	"fmt"
	"sort"
	"strings"
	"sync"
)

var (
	mu sync.Mutex
	Th = 6
)

func SingleHash(input, output chan any) {
	wg := sync.WaitGroup{}
	for data := range input {
		wg.Add(1)
		go func(data string, output chan any) {
			defer wg.Done()
			handleSingleHash(data, output)
		}(fmt.Sprintf("%v", data), output)
	}
	wg.Wait()
}

func MultiHash(input, output chan any) {
	wg := sync.WaitGroup{}
	for data := range input {
		stringData := fmt.Sprintf("%v", data)
		wg.Add(1)
		go func() {
			defer wg.Done()
			handleMultiHash(stringData, output)
		}()
	}
	wg.Wait()
}

func CombineResults(input, output chan any) {
	var result []string
	for data := range input {
		result = append(result, fmt.Sprintf("%v", data))
	}
	sort.Strings(result)
	output <- strings.Join(result, "_")
}

func handleSingleHash(data string, output chan any) {
	ch := make(chan string)
	go func(ch chan string) {
		ch <- DataSignerCrc32(<-ch)
	}(ch)
	ch <- data

	mu.Lock()
	Md5 := DataSignerMd5(data)
	mu.Unlock()

	Crc32Md5 := DataSignerCrc32(Md5)
	Crc32 := <-ch
	output <- Crc32 + "~" + Crc32Md5
}

func handleMultiHash(data string, output chan any) {
	channels := make([]chan string, Th)
	var result string
	for i := 0; i < Th; i++ {
		channels[i] = make(chan string)
		go func(ch chan string, th int) {
			ch <- DataSignerCrc32(fmt.Sprintf("%v%v", th, data))
		}(channels[i], i)
	}

	for _, channel := range channels {
		result += <-channel
		close(channel)
	}
	output <- result
}
