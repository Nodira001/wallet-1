package wallet

import "sync"

func Regular() int64 {
	sum := int64(0)
	for i := 0; i < 2_000; i++ {
		sum++
	}
	return sum
}
func Concurrently() int64 {
	wg := sync.WaitGroup{}
	mu := sync.Mutex{}
	wg.Add(2)
	sum := int64(0)
	go func() {
		defer wg.Done()
		val:=int64(0)
		for i := 0; i < 1_000; i++ {
			val++
		}
		mu.Lock()
		defer mu.Unlock()
		sum+=val
	}()
	go func() {
		defer wg.Done()
		val:=int64(0)
		for i := 0; i < 1_000; i++ {
			val++
		}
		mu.Lock()
		defer mu.Unlock()
		val+=sum
	}()
	wg.Wait()
	return sum
}
