package throttle

import (
	"sync"
	"time"
)

type Throttle struct {
	Update   chan struct{} // update channel
	quit     chan struct{} // quit channel
	interval time.Duration // interval
	isActive bool          // is active
	mu       *sync.Mutex   // mutex for protecting ticker
	lock     *sync.RWMutex // rwmutex for protecting fields
}

// Create new throttle.
func NewThrottle(interval time.Duration) *Throttle {
	return &Throttle{
		Update:   make(chan struct{}),
		quit:     make(chan struct{}),
		interval: interval,
		isActive: false,
		mu:       &sync.Mutex{},
		lock:     &sync.RWMutex{},
	}
}

// Get tick stat.
func (t *Throttle) getStatus() bool {
	t.lock.RLock()
	defer t.lock.RUnlock()
	return t.isActive
}

// Set tick stat.
func (t *Throttle) setStatus(status bool) {
	t.lock.Lock()
	defer t.lock.Unlock()
	t.isActive = status
}

// Generate tick and send signal to Update channel.
func (t *Throttle) generateTick() {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.setStatus(true)

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
	t.lock.RLock()
	defer t.lock.RUnlock()
	return t.interval
}

// Change interval.
func (t *Throttle) ChangeInterval(interval time.Duration) {
	active := t.getStatus()
	if active {
		t.Stop()
	}
	t.lock.Lock()
	t.interval = interval
	t.lock.Unlock()
	if active {
		go t.generateTick()
	}
}

// Start throttle.
func (t *Throttle) Start() {
	if t.getStatus() {
		t.Stop()
	}

	go t.generateTick()
}

// Stop throttle.
// Transmit must be guaranteed.
func (t *Throttle) Stop() {
	active := t.getStatus()
	if active {
		t.setStatus(false)
		t.lock.Lock()
		select {
		case t.quit <- struct{}{}:
		}
		t.lock.Unlock()
	}
}
