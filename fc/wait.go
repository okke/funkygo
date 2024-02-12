package fc

import (
	"sync"
)

type Done func()
type WaiterN func(int, func(Done)) WaiterN

func WaitN(amount int, f func(Done)) WaiterN {

	var wg sync.WaitGroup
	done := func() {
		wg.Done()
	}

	var waiter WaiterN
	waiter = func(amount int, g func(Done)) WaiterN {

		wg.Add(amount)
		g(done)
		wg.Wait()

		return waiter
	}

	waiter(amount, f)
	return waiter
}

type TaskSubmitter func(func())

type Waiter func(func(TaskSubmitter)) Waiter

func Wait(submitTasks func(TaskSubmitter)) Waiter {

	var wg sync.WaitGroup

	var waiter Waiter
	waiter = func(submitWaiterTasks func(TaskSubmitter)) Waiter {
		submitWaiterTasks(func(bg func()) {
			wg.Add(1)
			go func() {
				defer wg.Done()
				bg()
			}()
		})

		wg.Wait()
		return waiter
	}

	waiter(submitTasks)
	return waiter
}
