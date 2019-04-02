package throttle

import (
	"sync"
	"time"
)

type Throttle struct {
	Update       chan struct{} // update channel
	quit         chan struct{} // quit channel
	interval     time.Duration // interval
	isActive     bool          // is active
	sync.RWMutex               // embedded mutex
}

// Create new throttle.
func NewThrottle(interval time.Duration) *Throttle {
	return &Throttle{
		Update:   make(chan struct{}),
		quit:     make(chan struct{}),
		interval: interval,
		isActive: false,
	}
}

// Get tick stat.
func (t *Throttle) getStatus() bool {
	t.RLock()
	defer t.RUnlock()
	return t.isActive
}

// Set tick stat.
func (t *Throttle) setStatus(status bool) {
	t.Lock()
	defer t.Unlock()
	t.isActive = status
}

// Generate tick and send signal to Update channel.
func (t *Throttle) generateTick() {
	t.setStatus(true)
	defer t.setStatus(false)

	ticker := time.NewTicker(t.GetInterval())
	for {
		select {
		case <-ticker.C:
			select {
			case t.Update <- struct{}{}:
			default:
			}
		case <-t.quit:
			ticker.Stop()
			return
		}
	}
}

// Get interval duration.
func (t *Throttle) GetInterval() time.Duration {
	t.RLock()
	defer t.RUnlock()
	return t.interval
}

// Change interval.
func (t *Throttle) ChangeInterval(interval time.Duration) {
	stat := t.getStatus()
	if stat {
		t.Stop()
	}
	t.Lock()
	t.interval = interval
	t.Unlock()
	if stat {
		go t.generateTick()
	}
}

// Start throttle.
func (t *Throttle) Start() {
	t.Stop()
	go t.generateTick()
}

// Stop throttle.
func (t *Throttle) Stop() {
	t.Lock()
	defer t.Unlock()
	select {
	case t.quit <- struct{}{}:
	default:
	}
}
