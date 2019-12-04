// Tool to fetch a full list of DOI from datacite.org API, because as of Fall
// 2019 a full dump is not yet available (https://git.io/Je6bs,
// https://git.io/Je6Dg).
//
// THIS IS THROWAWAY CODE, AS IT IS HOPEFULLY OBSOLETE SOON.
//
// Currently (12/2019) using the "dois" endpoint, from v2 of the datacite API,
// supposedly.
//
// > The current version of the REST API is version 2. If you are using the
// endpoints /works, /members, or /data-centers, you are using version 1.
//
// Various intervals (weekly, daily, hourly, every minute) to mitigate deep
// paging issue and HTTP 502s.
//
// Notes.
//
// a) Hourly slices will fetch 27G, then fail with 502, too.
// b) Every minute is better, smaller chunks, short cursor runs (100s of pages, no 1000s).
//
// Errors encountered: 502, 500, 403, 400, "unexpected EOF" (maybe
// https://stackoverflow.com/q/21147562/89391). Strange error with minute
// interval: "search_after has 3 value(s) but sort has 2."
// (https://is.gd/9b0GF0); only in the 2019-08-02 15:17:00 - 15:17:59 query
// window.
//
// Less informative 500 on https://is.gd/uP0aJ2; 2019-10-07 16:19:00 - 16:19:59.
package main

import (
	"flag"
	"fmt"
	"net/url"
	"os"
	"path"
	"strings"
	"sync"
	"time"

	"github.com/miku/dcdump"
	"github.com/miku/dcdump/atomic"
	"github.com/miku/dcdump/dateutil"
	log "github.com/sirupsen/logrus"
)

var (
	start dateutil.Date = dateutil.Date{Time: dateutil.MustParse("2018-01-01")}
	end   dateutil.Date = dateutil.Date{Time: time.Now().UTC()}

	debug       = flag.Bool("debug", false, "only print intervals then exit")
	prefix      = flag.String("p", "dcdump-", "file prefix for harvested files")
	maxRequests = flag.Int("l", 16777216, "upper limit for number of requests")
	workers     = flag.Int("w", 4, "parallel workers (approximate)")
	interval    = flag.String("i", "d", "[w]eekly, [d]daily, [h]ourly, [e]very minute")
	directory   = flag.String("d", ".", "directory, where to put harvested files")
	showVersion = flag.Bool("version", false, "show version")

	Version   = "dev"
	Buildtime = time.Now().Format("2006-01-02T15:04:05Z")
)

// unrollPages takes a start and end time and will write newline delimited JSON
// into a single file at DIRECTORY/PREFIX-START-END.ndj. If that file already
// exists, we assume we already fetched that particular time window.
func unrollPages(s, e time.Time, directory, prefix string) error {
	filename := path.Join(directory, fmt.Sprintf("%s%s-%s.ndj",
		prefix,
		s.Format("20060102150405"),
		e.Format("20060102150405")))

	// If file exists assume we already harvested it.
	if _, err := os.Stat(filename); err == nil {
		log.Printf("[skip] assuming data already harvested in %s", filename)
		return nil
	}

	// https://api.datacite.org/dois?query=updated:%5B2019-01-01T00:00:00Z+TO+2019-02-01T23:59:59Z%5D&state=findable
	from, to := s.Format(time.RFC3339), e.Format(time.RFC3339) // TODO(martin): This is fine for UTC, but also w/ TZ?
	vs := url.Values{
		"query":        []string{fmt.Sprintf("updated:[%s TO %s]", from, to)},
		"state":        []string{"findable"},
		"page[size]":   []string{"100"},
		"page[cursor]": []string{"1"}, // https://support.datacite.org/docs/pagination#section-cursor
	}
	link := fmt.Sprintf("https://api.datacite.org/dois?%s", vs.Encode())

	// Fetch into temporary file, then move to destination.
	log.Printf("start batch: %s", link)
	fn, err := dcdump.HarvestBatch(link, *maxRequests) // Page through.
	if err != nil {
		return err
	}
	log.Printf("batch done: %s", link)
	return atomic.MoveFile(fn, filename)
}

// hasPrefix returns true, if s starts with prefix, case insensitive.
func hasPrefix(s, prefix string) bool {
	return strings.HasPrefix(
		strings.ToLower(strings.TrimSpace(s)),
		strings.ToLower(strings.TrimSpace(prefix)))
}

func main() {
	flag.Var(&start, "s", "start date for harvest")
	flag.Var(&end, "e", "end date for harvest")

	flag.Parse()

	if *showVersion {
		fmt.Printf("dcdump %s %s\n", Version, Buildtime)
		os.Exit(0)
	}

	sem := make(chan struct{}, *workers) // Have at most ~workers in parallel.
	var wg sync.WaitGroup

	var intervals []dateutil.Interval
	switch {
	case hasPrefix(*interval, "e"):
		intervals = dateutil.EveryMinute(start.Time, end.Time)
	case hasPrefix(*interval, "h"):
		intervals = dateutil.Hourly(start.Time, end.Time)
	case hasPrefix(*interval, "d"):
		intervals = dateutil.Daily(start.Time, end.Time)
	case hasPrefix(*interval, "w"):
		intervals = dateutil.Weekly(start.Time, end.Time)
	case hasPrefix(*interval, "m"):
		intervals = dateutil.Monthly(start.Time, end.Time)
	default:
		log.Fatal("intervals supported: [h]ourly, [d]aily, [w]eekly, [m]onthly and [e]very minute")
	}

	if *debug {
		for _, iv := range intervals {
			fmt.Printf("%s -- %s\n", iv.Start, iv.End)
		}
		log.Printf("%d intervals", len(intervals))
		os.Exit(0)
	}

	log.Printf("attempting to fetch datacite in %d intervals", len(intervals))

	for _, iv := range intervals {
		sem <- struct{}{}
		wg.Add(1)
		go func(iv dateutil.Interval) {
			defer wg.Done()
			if err := unrollPages(iv.Start, iv.End, *directory, *prefix); err != nil {
				log.Warnf("incomplete harvest - maybe rm -f %s*.ndj", path.Join(*directory, *prefix))
				log.Fatal(err)
			}
			<-sem
		}(iv)
	}
	wg.Wait()
	log.Printf("%d date slices succeeded", len(intervals))
}
