package main

import (
	"fmt"
	"math/rand"
	"time"
)

func main() {
	ch := make(chan int32, 8)

	n := 10

	t1 := time.Now()
	for i := 0; i < n; i++ {
		go fGo(ch)
	}
	for i := 0; i < n; i++ {
		<-ch
	}
	fmt.Printf("\n%v\n", time.Since(t1))

	t2 := time.Now()
	for i := 0; i < n; i++ {
		fPasGo()
	}
	fmt.Printf("\n%v\n", time.Since(t2))
}

func fGo(ch chan int32) {
	time.Sleep(time.Millisecond)
	ch <- rand.Int31()
}

func fPasGo() int32 {
	time.Sleep(time.Millisecond)
	return rand.Int31()
}
