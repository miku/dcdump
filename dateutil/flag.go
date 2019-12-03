// Package dateutil provides a custom flag for dates.
package dateutil

import (
	"time"

	"github.com/araddon/dateparse"
	"github.com/jinzhu/now"
)

var (
	// https://english.stackexchange.com/questions/3091/weekly-daily-hourly-minutely
	EveryMinute = makeIntervalFunc(
		func(t time.Time) time.Time { return now.With(t).BeginningOfMinute() },
		func(t time.Time) time.Time { return now.With(t).EndOfMinute() },
	)
	Hourly = makeIntervalFunc(
		func(t time.Time) time.Time { return now.With(t).BeginningOfHour() },
		func(t time.Time) time.Time { return now.With(t).EndOfHour() },
	)
	Daily = makeIntervalFunc(
		func(t time.Time) time.Time { return now.With(t).BeginningOfDay() },
		func(t time.Time) time.Time { return now.With(t).EndOfDay() },
	)
	Weekly = makeIntervalFunc(
		func(t time.Time) time.Time { return now.With(t).BeginningOfWeek() },
		func(t time.Time) time.Time { return now.With(t).EndOfWeek() },
	)
	Monthly = makeIntervalFunc(
		func(t time.Time) time.Time { return now.With(t).BeginningOfMonth() },
		func(t time.Time) time.Time { return now.With(t).EndOfMonth() },
	)
)

type (
	shiftFunc    func(t time.Time) time.Time
	intervalFunc func(s, e time.Time) []Interval
)

// makeIntervalFunc is a helper to create daily, weekly and other intervals.
func makeIntervalFunc(beginningOfX, endOfX shiftFunc) intervalFunc {
	f := func(s, e time.Time) (result []Interval) {
		var (
			l time.Time = s
			r time.Time
		)
		for {
			r = endOfX(l)
			result = append(result, Interval{Start: l, End: r})
			l = beginningOfX(r.Add(1 * time.Second))
			if l.After(e) {
				break
			}
		}
		return result
	}
	return f
}

// Interval groups start and end.
type Interval struct {
	Start time.Time
	End   time.Time
}

// Date can be used to parse command line args into dates.
type Date struct {
	time.Time
}

// String returns a formatted date.
func (d *Date) String() string {
	return d.Format("2006-01-02")
}

// Set parses a value into a date, relatively flexible due to
// araddon/dateparse, 2014-04-26 will work, but oct. 7, 1970, too.
func (d *Date) Set(value string) error {
	t, err := dateparse.ParseStrict(value)
	if err != nil {
		return err
	}
	*d = Date{t}
	return nil
}

// MustParse will panic on an unparsable date string.
func MustParse(value string) time.Time {
	t, err := dateparse.ParseStrict(value)
	if err != nil {
		panic(err)
	}
	return t
}
