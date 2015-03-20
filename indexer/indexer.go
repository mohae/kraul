// indexer is a basic html indexer.
package indexer

import (
	"fmt"
	_ "log"
	"net/http"
	"net/url"
	_ "runtime"
	"strconv"
	_ "strings"
	"sync"
	"time"

	"github.com/PuerkitoBio/fetchbot"
	"github.com/PuerkitoBio/goquery"
)

// first get results
// then see about specific parsing
var (
	visited          tracker
	restrictToDomain = true
	stopAfter        time.Duration
)

func init() {
	visited = tracker{urls: map[string]bool{}}
}

func Index(startUrl string) (message string, err error) {
	if startUrl == "" {
		return "", fmt.Errorf("start url expected, none received")
	}

	// only looking at the first; for now
	u, err := url.Parse(startUrl)

	// get a new muxer
	mux := fetchbot.NewMux()

	// Set up error handling
	mux.HandleErrors(fetchbot.HandlerFunc(func(ctx *fetchbot.Context, res *http.Response, err error) {
		fmt.Printf("[ERR] %s %s - %s\n", ctx.Cmd.Method(), ctx.Cmd.URL(), err)
	}))

	// Handle the GET requests
	mux.Response().Method("GET").ContentType("text/html").Handler(fetchbot.HandlerFunc(
		func(ctx *fetchbot.Context, res *http.Response, err error) {
			// Process body for links
			doc, err := goquery.NewDocumentFromResponse(res)
			if err != nil {
				fmt.Printf("[ERR] %s %s - %s\n", ctx.Cmd.Method(), ctx.Cmd.URL(), err)
				return
			}
			// Enqueue all links as HEAD requests
			visited.enqueueLinks(ctx, doc)
		},
	))

	// Handle HEAD requests for html responses coming from the source host - we don't want
	// to index links from other hosts.
	mux.Response().Method("HEAD").Host(u.Host).ContentType("text/html").Handler(fetchbot.HandlerFunc(
		func(ctx *fetchbot.Context, res *http.Response, err error) {
			if _, err := ctx.Q.SendStringGet(ctx.Cmd.URL().String()); err != nil {
				fmt.Printf("[ERR] %s %s - %s\n", ctx.Cmd.Method(), ctx.Cmd.URL(), err)
			}
		}))

	// Create the Fetcher, handle the logging first, then dispatch to the Muxer
	h := logHandler(mux)
	f := fetchbot.New(h)
	// Start processing
	q := f.Start()
	if stopAfter > 0 {
		go func() {
			c := time.After(stopAfter)
			<-c
			q.Close()
		}()
	}
	// Enqueue the startUrl, which is the first entry in the dup map
	visited.mu.Lock()
	visited.urls[startUrl] = false
	visited.mu.Unlock()
	_, err = q.SendStringGet(startUrl)
	if err != nil {
		fmt.Printf("[ERR] GET %s - %s\n", startUrl, err)
	}
	q.Block()

	return indexMessage(len(visited.urls)), nil
}

func handler(ctx *fetchbot.Context, res *http.Response, err error) {
	if err != nil {
		fmt.Printf("error: %s\n", err)
		return
	}
	fmt.Printf("[%d]%s %s\n%s\n", res.StatusCode, ctx.Cmd.Method, ctx.Cmd.URL(), res.Body)
}

func indexMessage(i int) string {
	return strconv.Itoa(i) + " urls processed"
}

type tracker struct {
	mu   sync.Mutex      // prevent race
	urls map[string]bool // url/visited; prevents dups
}

// Add adds a url to the url map
func (t *tracker) Add(url string) (added bool) {
	t.mu.Lock()
	defer t.mu.Unlock()
	_, ok := t.urls[url]
	if ok { // already exists
		return false
	}
	t.urls[url] = false // add, hasn't been processed yet
	return true
}

func (t *tracker) SetProcessed(url string) (processed bool, err error) {
	t.mu.Lock()
	defer t.mu.Unlock()
	processed, ok := t.urls[url]
	if !ok { // doesn't exist yet
		return false, fmt.Errorf("Unable to process %s: not in index", url)
	}
	if processed { // its already processed, return
		return true, nil
	}
	t.urls[url] = true // update to processed
	return true, nil
}

// stopHandler stops the fetcher if the stopurl is reached. Otherwise it dispatches
// the call to the wrapped Handler.
func stopHandler(stopurl string, wrapped fetchbot.Handler) fetchbot.Handler {
	return fetchbot.HandlerFunc(func(ctx *fetchbot.Context, res *http.Response, err error) {
		if ctx.Cmd.URL().String() == stopurl {
			ctx.Q.Close()
			return
		}
		wrapped.Handle(ctx, res, err)
	})
}

// logHandler prints the fetch information and dispatches the call to the wrapped Handler.
func logHandler(wrapped fetchbot.Handler) fetchbot.Handler {
	return fetchbot.HandlerFunc(func(ctx *fetchbot.Context, res *http.Response, err error) {
		if err == nil {
			fmt.Printf("[%d] %s %s - %s\n", res.StatusCode, ctx.Cmd.Method(), ctx.Cmd.URL(), res.Header.Get("Content-Type"))
		}
		wrapped.Handle(ctx, res, err)
	})
}

func (t *tracker) enqueueLinks(ctx *fetchbot.Context, doc *goquery.Document) {
	t.mu.Lock()
	defer t.mu.Unlock()
	doc.Find("a[href]").Each(func(i int, s *goquery.Selection) {
		val, _ := s.Attr("href")
		// Resolve address
		u, err := ctx.Cmd.URL().Parse(val)
		if err != nil {
			fmt.Printf("error: resolve URL %s - %s\n", val, err)
			return
		}
		if !t.urls[u.String()] {
			if _, err := ctx.Q.SendStringHead(u.String()); err != nil {
				fmt.Printf("error: enqueue head %s - %s\n", u, err)
			} else {
				t.urls[u.String()] = false
			}
		}
	})
}
