package main

import (
	"fmt"
	"goroutine_pool/pool"
	"time"
)

func main() {
	p, err := pool.New(100)
	if err != nil {
		fmt.Println(err)
	}
	for i:= 1; i < 1000; i++ {
		p.Put(&pool.Task{
			Handler:    func(v ...interface{}){fmt.Println(v...)},
			Parameters: []interface{}{i},
		})
	}
	time.Sleep(time.Second)
}

