package main

import (
	"log"
	"sync"
	"time"
)

func main() {
	log.Print("main started")
	wg := sync.WaitGroup{}
	wg.Add(2)
	sum := 0
	go func() {
		defer wg.Done()
		for i := 0; i < 1_000; i++ {
			sum++
		}
	}()
	go func() {
		defer wg.Done()
		for i := 0; i < 1_000; i++ {
			sum++
		}
	}()
	wg.Wait()
	log.Print("main finished")
	time.Sleep(time.Second * 5)
	log.Print(sum)
}
