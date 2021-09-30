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
	d := time.Now().Sub(a.firstTime)
	stale = d > (time.Minute * time.Duration(config().FailedTimeout))
	return stale
}

type attempts struct {
	lookup map[string]*attempt
}

func (at *attempts) init() {
	at.lookup = make(map[string]*attempt)
}

func (at *attempts) allowed(user string) bool {
	a, x := at.lookup[user]
	allow := false
	if x {
		if a.isStale() {
			// prevous attemptis too old
			// treat this as new set of attempts
			at.lookup[user] = &attempt{
				config().NumFailedAttempts - 1,
				time.Now(),
			}
			allow = true
		} else {
			// active attempt found
			a.count--
			allow = a.count > 0
		}
	} else {
		// no previous attempts exists for user
		at.lookup[user] = &attempt{
			config().NumFailedAttempts - 1,
			time.Now(),
		}
		allow = true
	}
	return allow
}
