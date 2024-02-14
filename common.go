package main

import (
	"crypto/md5"
	"fmt"
	"hash/crc32"
	"sync/atomic"
	"time"
)

var dataSignerOverheat = atomic.Int32{}

func OverheatLock() {
	for {
		if !dataSignerOverheat.CompareAndSwap(0, 1) {
			fmt.Println("OverheatLock happend")
			time.Sleep(time.Second)
		} else {
			break
		}
	}
}

func OverheatUnlock() {
	for {
		if !dataSignerOverheat.CompareAndSwap(1, 0) {
			fmt.Println("OverheatUnlock happend")
			time.Sleep(time.Second)
		} else {
			break
		}
	}
}

func DataSignerMd5(data string) string {
	OverheatLock()
	defer OverheatUnlock()
	dataHash := fmt.Sprintf("%x", md5.Sum([]byte(data)))
	time.Sleep(10 * time.Millisecond)
	return dataHash
}

func DataSignerCrc32(data string) string {
	crc32 := crc32.ChecksumIEEE([]byte(data))
	dataHash := fmt.Sprintf("%v", crc32)
	time.Sleep(time.Second)
	return dataHash
}
