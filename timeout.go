package awsclient

import "time"

type timeout struct {
	value time.Duration
}

func newTimeout(v time.Duration) *timeout {
	return &timeout{
		value: v,
	}
}

func (t *timeout) Value() time.Duration {
	return t.value
}

func (t *timeout) Set(v time.Duration) {
	t.value = v
}
