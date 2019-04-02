package throttle

import (
	"os"
	"testing"
	"time"
)

var (
	timeInterval = time.Second
	th           = new(Throttle)
)

func TestMain(m *testing.M) {
	th = NewThrottle(timeInterval)
	os.Exit(m.Run())
}

func TestGetInterval(t *testing.T) {
	interval := th.GetInterval()
	if interval != timeInterval {
		t.Error("GetInterval error!")
	}
}

func TestStart(t *testing.T) {
	th.Start()
	<-th.Update
}

func TestChangeIntervalBeforeStop(t *testing.T) {
	th.ChangeInterval(time.Second * 2)
	<-th.Update
}

func TestStop(t *testing.T) {
	<-th.Update
	th.Stop()
}

func TestChangeIntervalAfterStop(t *testing.T) {
	th.ChangeInterval(time.Second * 3)
	th.Start()
	<-th.Update
}
