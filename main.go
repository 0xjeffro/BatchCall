package main

import (
	"fmt"
	"github.com/0xjeffro/BatchCall.git/batchCall"
	"math/rand"
	"time"
)

func handler(i interface{}) (interface{}, error) {
	n := rand.Intn(10)

	time.Sleep(time.Duration(n/5) * time.Second)

	if n < 5 {
		return nil, fmt.Errorf("error")
	} else {
		return i.(int), nil
	}
}

func main() {
	params := make([]interface{}, 100)
	for i := range params {
		params[i] = i
	}

	results := batchCall.BathCall(params, handler, 100, 2)

	for idx, r := range results {
		fmt.Println(idx, r)
	}
}
