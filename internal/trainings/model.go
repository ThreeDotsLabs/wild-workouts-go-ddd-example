package main

import "time"

func (t Training) canBeCancelled() bool {
	return t.Time.Sub(time.Now()) > time.Hour*24
}
