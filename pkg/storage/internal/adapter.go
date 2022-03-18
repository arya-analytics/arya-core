package internal

import "time"

type Adapter interface {
	Healthy() bool
	Acquire()
	Release()
}

type Expiration struct {
	Start    time.Time
	Duration time.Duration
}

func (e Expiration) Remaining() time.Duration {
	return -1 * time.Since(e.Start.Add(e.Duration))
}

func (e *Expiration) Expired() bool {
	return e.Remaining() < 0
}

type Demand struct {
	Max   int
	Min   int
	Value int
}

func (d Demand) Increment() {
	d.Value++
}

func (d Demand) Decrement() {
	d.Value--
}

func (d Demand) Exceeded() bool {
	return d.TooHigh() || d.TooLow()
}

func (d Demand) TooHigh() bool {
	return d.Value > d.Max

}

func (d Demand) TooLow() bool {
	return d.Value < d.Min
}
