package fc

import (
	"sync"
)

type WaiterN func(int, func(func())) WaiterN

func WaitN(amount int, f func(func())) WaiterN {

	var wg sync.WaitGroup
	done := func() {
		wg.Done()
	}

	wg.Add(amount)
	f(done)
	wg.Wait()

	var reuse func(int, func(func())) WaiterN
	reuse = func(amount int, g func(func())) WaiterN {

		wg.Add(amount)
		g(done)
		wg.Wait()

		return reuse
	}

	return reuse
}
