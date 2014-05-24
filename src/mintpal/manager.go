package mintpal

import (
	"sync"
)

var (
	once   sync.Once
	manage *manager
)

type manager struct{}

func Manager() (m *manager) {
	if manage == nil {
		once.Do(func() {
			manage = new(manager)
		})
	}
	m = manage
	return
}
