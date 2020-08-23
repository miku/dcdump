package dcdump

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/sethgrid/pester"
	log "github.com/sirupsen/logrus"
)

// HarvestBatch takes a link (like https://is.gd/0pwu5c) and follows subsequent
// pages, writes everything into a tempfile. Returns path to temporary file and
// an error. Fails, if HTTP status is >= 400; has limited retry capabilities.
func HarvestBatch(link string, maxRequests int, sleep time.Duration) (string, error) {
	var (
		i          int
		client     = pester.New()
		maxRetries = 10 // maxRetries on successful requests but bad HTTP status codes
		retry      = 0
		resp       *http.Response
		err        error
	)
	client.SetRetryOnHTTP429(true)
	f, err := ioutil.TempFile("", "dcdump-")
	if err != nil {
		return "", err
	}
	defer f.Close()
	for {
		retry = 0
		for {
			retry++
			if retry == maxRetries {
				return "", fmt.Errorf("max retries exceeded on %s (last status code was: %d)", link, resp.StatusCode)
			}
			if retry > 1 {
				log.Warnf("[%d] %s", retry, link)
			}
			req, err := http.NewRequest("GET", link, nil)
			if err != nil {
				return "", err
			}
			req.Header.Add("User-Agent", "dcdump")
			// The "unexpected EOF" might be caused by some truncated, compressed
			// content? This is a workaround, it probably should be retried (which
			// pester does not).
			// https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Accept-Encoding#Directives
			// https://stackoverflow.com/questions/21147562/unexpected-eof-using-go-http-client/21160982
			req.Header.Add("Accept-Encoding", "identity")
			resp, err = client.Do(req)
			if err == io.ErrUnexpectedEOF {
				log.Warnf("got %s on %s", err, link)
				time.Sleep(sleep)
				continue
			}
			if err != nil {
				// Retry on any transport error.
				continue
			}
			defer resp.Body.Close()
			if resp.StatusCode < 400 {
				break
			}
			log.Warnf("failed to fetch link with %s, will retry in %s", resp.Status, sleep)
			time.Sleep(sleep)
		}
		tee := io.TeeReader(resp.Body, f)
		// Fabricate newline delimited file.
		if _, err := io.WriteString(f, "\n"); err != nil {
			return "", err
		}
		// Decoding and see, whether we got JSON and if so, parse out some fields.
		dec := json.NewDecoder(tee)
		var dr DOIResponse
		if err := dec.Decode(&dr); err != nil {
			return "", err
		}
		if int64(i) > dr.Meta.TotalPages {
			log.Warnf("more request than pages, next link might be broken: %s", dr.Links.Next)
		}
		if i > maxRequests {
			return "", fmt.Errorf("exceeded maximum number of requests: %d", maxRequests)
		}
		log.Printf("requests=%d, pages=%d, total=%d", i+1, dr.Meta.TotalPages, dr.Meta.Total)
		// https://api.datacite.org/dois?page%5Bnumber%5D=2&page%5Bsize%5D=25
		// https://api.datacite.org/dois?page[number]=2&page[size]=25
		if dr.Meta.TotalPages >= 400 {
			log.Warnf("time slice is over 400 pages wide (%s), might break w/o cursor", link)
		}
		link = dr.Links.Next // Should be independent of whether we use cursor or not.
		if link == "" {
			break
		}
		i++
	}
	return f.Name(), nil
}
