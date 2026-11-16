package main

// Eingebaute Störungen für die Übungen, nur für /notes:
//
//	NOTES_DELAY=300ms        jede Anfrage wartet so lange
//	NOTES_ERROR_RATE=0.2     so viele Anfragen scheitern mit 500

import (
	"errors"
	"math/rand/v2"
	"os"
	"strconv"
	"time"
)

type faults struct {
	delay     time.Duration
	errorRate float64
}

var errFault = errors.New("storage: request timed out")

func loadFaults() (faults, error) {
	var f faults
	var err error
	if v := os.Getenv("NOTES_DELAY"); v != "" {
		if f.delay, err = time.ParseDuration(v); err != nil {
			return f, errors.New("NOTES_DELAY: " + err.Error())
		}
	}
	if v := os.Getenv("NOTES_ERROR_RATE"); v != "" {
		if f.errorRate, err = strconv.ParseFloat(v, 64); err != nil {
			return f, errors.New("NOTES_ERROR_RATE: " + err.Error())
		}
	}
	return f, nil
}

func (f faults) active() bool { return f.delay > 0 || f.errorRate > 0 }

// inject wartet und liefert manchmal einen Fehler.
func (f faults) inject() error {
	time.Sleep(f.delay)
	if f.errorRate > 0 && rand.Float64() < f.errorRate {
		return errFault
	}
	return nil
}
