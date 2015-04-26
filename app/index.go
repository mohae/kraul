package app

import (
	"fmt"

	"github.com/mohae/contour"
	"github.com/mohae/geomi"
)

// Index indexes the site at the starturl. URLs not in the startUrl domain
// are not crawled
func Index(startUrl string) (string, error) {
	spider, err := geomi.NewSpider(startUrl)
	if err != nil {
		fmt.Printf("Could not create spider: %q", err)
		return "", err
	}
	fetchDelay := contour.GetInt("wait")
	if fetchDelay > 0 {
		spider.SetFetchInterval(int64(fetchDelay))
	}
	msg, err := spider.Crawl(4)
	return msg, err
}
