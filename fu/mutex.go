package fu

import "sync"

func WithMutex(m *sync.Mutex, f func()) {
	m.Lock()
	defer m.Unlock()
	f()
}
