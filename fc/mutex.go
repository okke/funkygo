package fc

import "sync"

type Mutex func(func())

func NewMutex() Mutex {

	var m sync.Mutex

	return func(f func()) {
		m.Lock()
		defer m.Unlock()
		f()
	}

}
