package main

import (
	"time"
)

type attempt struct {
	count     uint
	firstTime time.Time
}

func (a *attempt) isStale() bool {
	stale := true
	//compare time.Sub
	return stale
}

type attempts struct {
	lookup map[string]attempt
}

func (at *attempts) init() {
	at.lookup = make(map[string]attempt)
}
