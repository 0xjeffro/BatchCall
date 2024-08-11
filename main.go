package main

import (
	"fmt"
	"github.com/0xjeffro/BatchCall/batchCall"
	"math/rand"
	"strconv"
	"time"
)

func handler(i int) (string, error) {
	n := rand.Intn(10)

	time.Sleep(time.Duration(n/5) * time.Second)

	if n < 5 {
		return "", fmt.Errorf("error")
	} else {
		res := strconv.Itoa(i)
		return res, nil
	}
}

func main() {
	params := make([]int, 100)
	for i := range params {
		params[i] = i
	}

	var call batchCall.BatchCall[int, string]
	call.Params = params
	call.Op = handler

	results := call.Call(100, 10)

	for idx, r := range results {
		fmt.Println(idx, r)
	}
}
