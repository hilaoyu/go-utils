package utilTime

import (
	"time"
)

type Timer struct {
	*time.Timer
	expiration time.Time
	callback   func()
	expired    bool
}

func NewTimer(duration time.Duration) *Timer {
	return &Timer{
		Timer:      time.NewTimer(duration),
		expiration: time.Now().Add(duration),
	}
}
func AfterFunc(duration time.Duration, f func()) *Timer {
	return &Timer{
		Timer:      time.AfterFunc(duration, f),
		expiration: time.Now().Add(duration),
		callback:   f,
	}
}

func (t *Timer) ExpirationTime() time.Time {
	select {
	case t1 := <-t.C:
		t.expiration = t1
		t.expired = true
		break

	default:
		break
	}
	return t.expiration
}
func (t *Timer) IsExpired() (expired bool) {
	t.ExpirationTime()
	return t.expired
}
func (t *Timer) Reset(duration time.Duration) {
	t.Stop()
	if nil != t.callback {
		t.Timer = time.AfterFunc(duration, t.callback)
	} else {
		t.Timer = time.NewTimer(duration)
	}
	t.expiration = time.Now().Add(duration)
	t.expired = false
	return
}
func (t *Timer) Stop() {
	if nil != t.Timer {
		t.Timer.Stop()
		t.Timer = nil
	}
	if time.Now().Before(t.expiration) {
		t.expiration = time.Now()
	}
	t.expired = true
	return
}
